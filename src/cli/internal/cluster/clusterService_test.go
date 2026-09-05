package cluster

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// NewManagerForTests builds a container for this package's tests. Only the
// collaborators under test are registered; everything else stays a no-op, which
// is what lets s.sm.LogsService() be called without wiring a real logger.
func NewManagerForTests() (managers.ServiceManager, context.Context) {
	return managers.NewManager(), context.Background()
}

// call records one invocation the service made through the Runner.
type call struct {
	name string
	args []string
}

// fakeRunner stands in for exec. Each entry in results is matched against the
// first argument of the command being run, so a test can fail `cluster list`
// while letting `get nodes` succeed.
type fakeRunner struct {
	calls   []call
	out     map[string][]byte
	fail    map[string]error
	failAll error
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{out: map[string][]byte{}, fail: map[string]error{}}
}

// key identifies a command by its first two arguments, which is enough to tell
// "cluster create" from "cluster delete" from "get nodes".
func key(args []string) string {
	if len(args) >= 2 {
		return args[0] + " " + args[1]
	}
	if len(args) == 1 {
		return args[0]
	}
	return ""
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, call{name: name, args: args})

	if f.failAll != nil {
		return []byte("boom"), f.failAll
	}

	k := key(args)
	if err, ok := f.fail[k]; ok {
		return []byte("boom"), err
	}

	return f.out[k], nil
}

// withTools points the service at a directory containing fake binaries, so the
// lazy toolset resolution succeeds without the real tools being installed.
func withTools(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	for _, name := range []string{"k3d", "kubectl"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"), 0o700); err != nil {
			t.Fatalf("cannot create fake %s: %v", name, err)
		}
	}
	t.Setenv(BinDirEnvVar, dir)
}

// withoutTools points the service at an empty directory, so resolution fails.
func withoutTools(t *testing.T) {
	t.Helper()
	t.Setenv(BinDirEnvVar, t.TempDir())
}

// newServiceForTest returns a registered service and its fake runner.
func newServiceForTest(t *testing.T) (*clusterServiceFinal, *fakeRunner, context.Context) {
	t.Helper()

	// Isolate the airgap image cache. Without this every test would read the
	// developer's real ~/.freelunch/images, so whether Create mounts a volume
	// would depend on whether that machine had run setup — the tests would pass
	// or fail based on the host rather than the code.
	t.Setenv(ImagesDirEnvVar, t.TempDir())

	sm, ctx := NewManagerForTests()
	runner := newFakeRunner()
	sm.WithClusterService(NewClusterServiceWithRunner(runner))

	svc, ok := sm.ClusterService().(*clusterServiceFinal)
	if !ok {
		t.Fatalf("ClusterService() is %T, want *clusterServiceFinal", sm.ClusterService())
	}

	return svc, runner, ctx
}

func TestNewClusterService(t *testing.T) {
	if got := NewClusterService(); got == nil {
		t.Fatal("NewClusterService() = nil")
	}
}

func TestNewClusterServiceWithRunner(t *testing.T) {
	if got := NewClusterServiceWithRunner(newFakeRunner()); got == nil {
		t.Fatal("NewClusterServiceWithRunner() = nil")
	}
}

// Test_clusterServiceFinal_StartDoesNoIO is the rule this service exists to
// respect: Start runs on every invocation, including `freelunch --help`, so it
// must not touch the filesystem. Start is called with no tools present and must
// still succeed — resolution is deferred to first real use.
func Test_clusterServiceFinal_StartDoesNoIO(t *testing.T) {
	withoutTools(t)

	svc, runner, ctx := newServiceForTest(t)

	if err := svc.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	if len(runner.calls) != 0 {
		t.Errorf("Start() ran %d external commands, want 0", len(runner.calls))
	}

	if svc.tools != nil {
		t.Error("Start() resolved the toolset; it must be deferred to first use")
	}
}

func Test_clusterServiceFinal_Close(t *testing.T) {
	svc, _, ctx := newServiceForTest(t)

	if err := svc.Close(ctx); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func Test_clusterServiceFinal_Healthy(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		wantErr bool
	}{
		{name: "tools present", setup: withTools, wantErr: false},
		{name: "tools missing", setup: withoutTools, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			svc, _, ctx := newServiceForTest(t)

			if err := svc.Healthy(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Healthy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clusterServiceFinal_WithServiceManager(t *testing.T) {
	sm, _ := NewManagerForTests()
	svc, ok := NewClusterServiceWithRunner(newFakeRunner()).(*clusterServiceFinal)
	if !ok {
		t.Fatal("constructor did not return *clusterServiceFinal")
	}

	got := svc.WithServiceManager(sm)

	if !reflect.DeepEqual(got, managers.ClusterService(svc)) {
		t.Errorf("WithServiceManager() = %v, want the service itself", got)
	}
	if !reflect.DeepEqual(svc.ServiceManager(), sm) {
		t.Error("ServiceManager() did not return the registered manager")
	}
}

// Test_clusterServiceFinal_toolsetResolvedOnce covers the sync.Once guard: the
// filesystem is inspected on first use and not again.
func Test_clusterServiceFinal_toolsetResolvedOnce(t *testing.T) {
	withTools(t)
	svc, _, ctx := newServiceForTest(t)

	first, err := svc.toolset(ctx)
	if err != nil {
		t.Fatalf("toolset() error = %v, want nil", err)
	}

	// Remove the tools; a second call must return the cached result rather
	// than re-inspecting the filesystem and failing.
	if err := os.RemoveAll(os.Getenv(BinDirEnvVar)); err != nil {
		t.Fatalf("cannot remove fake tools: %v", err)
	}

	second, err := svc.toolset(ctx)
	if err != nil {
		t.Fatalf("second toolset() error = %v, want the cached success", err)
	}
	if first != second {
		t.Error("toolset() resolved twice; sync.Once is not guarding it")
	}
}

func Test_clusterServiceFinal_toolsetMissingIsErrNoTools(t *testing.T) {
	withoutTools(t)
	svc, _, ctx := newServiceForTest(t)

	_, err := svc.toolset(ctx)
	if !errors.Is(err, ErrNoTools) {
		t.Errorf("toolset() error = %v, want it to wrap ErrNoTools", err)
	}
}

func Test_clusterServiceFinal_Create(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T)
		failWith error
		wantErr  bool
	}{
		{name: "success", setup: withTools, wantErr: false},
		{name: "k3d fails", setup: withTools, failWith: errors.New("exit 1"), wantErr: true},
		{name: "tools missing", setup: withoutTools, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			svc, runner, ctx := newServiceForTest(t)
			if tt.failWith != nil {
				runner.fail["cluster create"] = tt.failWith
			}

			err := svc.Create(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.name != "success" {
				return
			}

			if len(runner.calls) != 1 {
				t.Fatalf("Create() made %d calls, want 1", len(runner.calls))
			}
			got := runner.calls[0]
			if filepath.Base(got.name) != "k3d" {
				t.Errorf("Create() ran %q, want the pinned k3d", got.name)
			}
			if got.args[0] != "cluster" || got.args[1] != "create" {
				t.Errorf("Create() args = %v, want cluster create", got.args)
			}
			// The config must be materialised from the embedded copy and
			// passed by path, not assumed to exist in the working directory.
			if got.args[2] != "--config" || !strings.HasSuffix(got.args[3], "k3d-cluster.yaml") {
				t.Errorf("Create() args = %v, want --config <path>/k3d-cluster.yaml", got.args)
			}
		})
	}
}

// Test_clusterServiceFinal_CreateCleansUpConfig checks the temp file does not
// outlive the call.
func Test_clusterServiceFinal_CreateCleansUpConfig(t *testing.T) {
	withTools(t)
	svc, runner, ctx := newServiceForTest(t)

	if err := svc.Create(ctx); err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	configPath := runner.calls[0].args[3]
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Errorf("config %s still exists after Create(); want it removed", configPath)
	}
}

func Test_clusterServiceFinal_Delete(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T)
		failWith error
		wantErr  bool
	}{
		{name: "success", setup: withTools, wantErr: false},
		{name: "k3d fails", setup: withTools, failWith: errors.New("exit 1"), wantErr: true},
		{name: "tools missing", setup: withoutTools, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			svc, runner, ctx := newServiceForTest(t)
			if tt.failWith != nil {
				runner.fail["cluster delete"] = tt.failWith
			}

			err := svc.Delete(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.name == "success" {
				got := runner.calls[0].args
				if got[0] != "cluster" || got[1] != "delete" || got[2] != ClusterName {
					t.Errorf("Delete() args = %v, want cluster delete %s", got, ClusterName)
				}
			}
		})
	}
}

func Test_clusterServiceFinal_Status(t *testing.T) {
	ready := func(name string) managers.ClusterNode { return managers.ClusterNode{Name: name, Ready: true} }
	notReady := func(name string) managers.ClusterNode { return managers.ClusterNode{Name: name} }

	tests := []struct {
		name        string
		setup       func(t *testing.T)
		nodesOutput string
		failList    bool
		failNodes   bool
		wantExists  bool
		wantRunning bool
		wantNodes   []managers.ClusterNode
		wantErr     bool
	}{
		{
			name:        "running with two Ready nodes",
			setup:       withTools,
			nodesOutput: "k3d-freelunch-server-0\tTrue\nk3d-freelunch-agent-0\tTrue\n",
			wantExists:  true,
			wantRunning: true,
			wantNodes:   []managers.ClusterNode{ready("k3d-freelunch-server-0"), ready("k3d-freelunch-agent-0")},
		},
		{
			// The case the -o name implementation got wrong: the node exists,
			// so it was counted, but it is not usable.
			name:        "a NotReady node exists but does not make the cluster running",
			setup:       withTools,
			nodesOutput: "k3d-freelunch-server-0\tTrue\nk3d-freelunch-agent-0\tFalse\n",
			wantExists:  true,
			wantRunning: false,
			wantNodes:   []managers.ClusterNode{ready("k3d-freelunch-server-0"), notReady("k3d-freelunch-agent-0")},
		},
		{
			// "Unknown" is what the API reports once the kubelet stops
			// heartbeating (container killed, Docker paused); a node still
			// joining has no Ready condition at all. Neither is Ready.
			name:        "Unknown and missing Ready conditions are NotReady",
			setup:       withTools,
			nodesOutput: "k3d-freelunch-server-0\tUnknown\nk3d-freelunch-agent-0\t\n",
			wantExists:  true,
			wantRunning: false,
			wantNodes:   []managers.ClusterNode{notReady("k3d-freelunch-server-0"), notReady("k3d-freelunch-agent-0")},
		},
		{
			name:        "an API answering with no nodes is not running",
			setup:       withTools,
			nodesOutput: "",
			wantExists:  true,
			wantRunning: false,
		},
		{
			name:        "cluster absent is not an error",
			setup:       withTools,
			failList:    true,
			wantExists:  false,
			wantRunning: false,
		},
		{
			name:        "api not answering yet is not an error, and the cluster still exists",
			setup:       withTools,
			failNodes:   true,
			wantExists:  true,
			wantRunning: false,
		},
		{
			name:    "tools missing",
			setup:   withoutTools,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			svc, runner, ctx := newServiceForTest(t)

			runner.out["get nodes"] = []byte(tt.nodesOutput)
			if tt.failList {
				runner.fail["cluster list"] = errors.New("no such cluster")
			}
			if tt.failNodes {
				runner.fail["get nodes"] = errors.New("connection refused")
			}

			got, err := svc.Status(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Status() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if got.Name != ClusterName {
				t.Errorf("Status().Name = %q, want %q", got.Name, ClusterName)
			}
			if got.Exists != tt.wantExists {
				t.Errorf("Status().Exists = %v, want %v", got.Exists, tt.wantExists)
			}
			if got.Running != tt.wantRunning {
				t.Errorf("Status().Running = %v, want %v", got.Running, tt.wantRunning)
			}
			if !reflect.DeepEqual(got.Nodes, tt.wantNodes) {
				t.Errorf("Status().Nodes = %v, want %v", got.Nodes, tt.wantNodes)
			}
		})
	}
}

// Test_clusterServiceFinal_Status_asksForReadiness pins the shape of the node
// query. `kubectl get nodes -o name` lists NotReady nodes indistinguishably
// from Ready ones, so the query must ask for the Ready condition explicitly.
func Test_clusterServiceFinal_Status_asksForReadiness(t *testing.T) {
	withTools(t)
	svc, runner, ctx := newServiceForTest(t)

	if _, err := svc.Status(ctx); err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	for _, c := range runner.calls {
		if key(c.args) != "get nodes" {
			continue
		}
		joined := strings.Join(c.args, " ")
		if !strings.Contains(joined, "-o jsonpath=") || !strings.Contains(joined, `@.type=="Ready"`) {
			t.Fatalf("get nodes must select the Ready condition via jsonpath; got %q", joined)
		}
		return
	}
	t.Fatal("Status() never ran `get nodes`")
}

func Test_parseNodeReadiness(t *testing.T) {
	got := parseNodeReadiness([]byte("\n  a\tTrue \n\nb\tFalse\nc\n\t\n"))
	want := []managers.ClusterNode{{Name: "a", Ready: true}, {Name: "b"}, {Name: "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseNodeReadiness() = %v, want %v", got, want)
	}
	if got := parseNodeReadiness(nil); got != nil {
		t.Errorf("parseNodeReadiness(nil) = %v, want nil", got)
	}
}

// Test_embeddedConfigMatchesClusterName guards the one piece of duplication in
// this package: ClusterName is a Go constant while the cluster's real name
// lives in the embedded YAML. If they drift, Delete and Status would address a
// cluster that does not exist.
func Test_embeddedConfigMatchesClusterName(t *testing.T) {
	if len(k3dConfig) == 0 {
		t.Fatal("embedded k3d-cluster.yaml is empty")
	}

	want := "name: " + ClusterName
	if !strings.Contains(string(k3dConfig), want) {
		t.Errorf("embedded config does not contain %q; ClusterName and the YAML have drifted", want)
	}
}

// withAirgapCache points the images dir at a temp dir holding a stub bundle,
// and returns that dir. The contents are never read — k3s consumes them inside
// the node, so all this code does is decide whether to mount.
func withAirgapCache(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	stub := filepath.Join(dir, "k3s-airgap-images-arm64.tar.gz")
	if err := os.WriteFile(stub, []byte("not a real tarball"), 0o600); err != nil {
		t.Fatalf("cannot write stub bundle: %v", err)
	}
	t.Setenv(ImagesDirEnvVar, dir)

	return dir
}

func Test_airgapVolume(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) string
		wantMount bool
	}{
		{
			name:      "cached bundle is mounted",
			setup:     withAirgapCache,
			wantMount: true,
		},
		{
			name: "empty cache directory mounts nothing",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				t.Setenv(ImagesDirEnvVar, dir)
				return dir
			},
			wantMount: false,
		},
		{
			name: "missing cache directory is not an error",
			setup: func(t *testing.T) string {
				dir := filepath.Join(t.TempDir(), "never-created")
				t.Setenv(ImagesDirEnvVar, dir)
				return dir
			},
			wantMount: false,
		},
		{
			name: "unrelated files are ignored",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "README"), []byte("x"), 0o600); err != nil {
					t.Fatalf("cannot write file: %v", err)
				}
				t.Setenv(ImagesDirEnvVar, dir)
				return dir
			},
			wantMount: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup(t)

			got, err := airgapVolume()
			if err != nil {
				t.Fatalf("airgapVolume() error = %v, want nil", err)
			}

			if !tt.wantMount {
				if got != "" {
					t.Errorf("airgapVolume() = %q, want empty", got)
				}
				return
			}

			// The destination is the only path k3s scans at startup; getting it
			// wrong yields a silently network-dependent cluster.
			want := dir + ":" + nodeImagesPath + "@server:*;agent:*"
			if got != want {
				t.Errorf("airgapVolume() = %q, want %q", got, want)
			}
		})
	}
}

// Test_clusterServiceFinal_CreateMountsAirgapCache is the unit-level half of
// Phase 7: a cached bundle must reach k3d as a --volume, on both node roles.
func Test_clusterServiceFinal_CreateMountsAirgapCache(t *testing.T) {
	withTools(t)
	dir := withAirgapCache(t)

	svc, runner, ctx := newServiceForTest(t)
	// newServiceForTest isolates the images dir; re-point it at the populated
	// one, since it runs after the setup helper above.
	t.Setenv(ImagesDirEnvVar, dir)

	if err := svc.Create(ctx); err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	args := runner.calls[0].args
	idx := -1
	for i, a := range args {
		if a == "--volume" {
			idx = i
			break
		}
	}
	if idx == -1 {
		t.Fatalf("Create() args = %v, want a --volume mounting the image cache", args)
	}
	if idx+1 >= len(args) {
		t.Fatalf("Create() args = %v, --volume has no value", args)
	}

	want := dir + ":" + nodeImagesPath + "@server:*;agent:*"
	if args[idx+1] != want {
		t.Errorf("Create() volume = %q, want %q", args[idx+1], want)
	}
}

// Test_clusterServiceFinal_CreateWithoutAirgapCache pins the fallback: no
// cache must still create a cluster, because that is the path a contributor
// who has not run setup takes.
func Test_clusterServiceFinal_CreateWithoutAirgapCache(t *testing.T) {
	withTools(t)
	svc, runner, ctx := newServiceForTest(t)

	if err := svc.Create(ctx); err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	for _, a := range runner.calls[0].args {
		if a == "--volume" {
			t.Errorf("Create() args = %v, want no --volume without a cache", runner.calls[0].args)
		}
	}
}

// Test_k3sVersionMatchesClusterConfig guards the one mismatch that would make
// an airgap look like it works.
//
// versions.env pins the k3s release whose image bundle we download; the
// embedded k3d config pins the k3s image the cluster actually boots. If they
// drift, the cache is populated with images the cluster never asks for, and it
// silently falls back to pulling from the network — indistinguishable from
// success until someone turns the network off.
//
// The two spellings differ by necessity: Docker tags cannot contain '+', so
// the release tag v1.35.5+k3s1 appears in the image as v1.35.5-k3s1.
func Test_k3sVersionMatchesClusterConfig(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "toolchain", "scripts", "versions.env"))
	if err != nil {
		t.Fatalf("cannot read versions.env: %v", err)
	}

	var pinned string
	for _, line := range strings.Split(string(raw), "\n") {
		if after, ok := strings.CutPrefix(strings.TrimSpace(line), "K3S_VERSION="); ok {
			pinned = after
			break
		}
	}
	if pinned == "" {
		t.Fatal("versions.env has no K3S_VERSION; the airgap bundle cannot be pinned")
	}

	wantTag := "rancher/k3s:" + strings.ReplaceAll(pinned, "+", "-")
	if !strings.Contains(string(k3dConfig), wantTag) {
		t.Errorf("k3d-cluster.yaml does not use image %q (K3S_VERSION=%s in versions.env)",
			wantTag, pinned)
	}
}
