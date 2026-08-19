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
	tests := []struct {
		name        string
		setup       func(t *testing.T)
		nodesOutput string
		failList    bool
		failNodes   bool
		wantRunning bool
		wantNodes   []string
		wantErr     bool
	}{
		{
			name:        "running with two nodes",
			setup:       withTools,
			nodesOutput: "node/k3d-freelunch-server-0\nnode/k3d-freelunch-agent-0\n",
			wantRunning: true,
			wantNodes:   []string{"k3d-freelunch-server-0", "k3d-freelunch-agent-0"},
		},
		{
			name:        "cluster absent is not an error",
			setup:       withTools,
			failList:    true,
			wantRunning: false,
		},
		{
			name:        "api not answering yet is not an error",
			setup:       withTools,
			failNodes:   true,
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
			if got.Running != tt.wantRunning {
				t.Errorf("Status().Running = %v, want %v", got.Running, tt.wantRunning)
			}
			if tt.wantNodes != nil && !reflect.DeepEqual(got.Nodes, tt.wantNodes) {
				t.Errorf("Status().Nodes = %v, want %v", got.Nodes, tt.wantNodes)
			}
		})
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
