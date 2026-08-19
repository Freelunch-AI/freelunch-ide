package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// NewManagerForTests builds a container for this package's tests. By convention
// it is duplicated per package rather than imported from another one. Only the
// command service is registered — LogsService stays a no-op, which is what lets
// s.sm.LogsService() be called here without wiring a real logger.
func NewManagerForTests() (managers.ServiceManager, context.Context) {
	return managers.NewManager(), context.Background()
}

// newTestService registers a real command service with its output captured.
func newTestService(t *testing.T) (s *commandServiceFinal, out, errOut *bytes.Buffer, ctx context.Context) {
	t.Helper()

	sm, ctx := NewManagerForTests()
	s = sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)

	out, errOut = &bytes.Buffer{}, &bytes.Buffer{}
	s.out = out
	s.errOut = errOut

	return s, out, errOut, ctx
}

func TestNewCommandService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommandService(); got == nil {
				t.Error("NewCommandService() = nil")
			}
		})
	}
}

func Test_commandServiceFinal_Start(t *testing.T) {
	s, _, _, ctx := newTestService(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Start(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_commandServiceFinal_Close(t *testing.T) {
	s, _, _, ctx := newTestService(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Close(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_commandServiceFinal_Healthy(t *testing.T) {
	s, _, _, ctx := newTestService(t)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "success", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Healthy(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Healthy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_commandServiceFinal_WithServiceManager(t *testing.T) {
	sm, _ := NewManagerForTests()
	s := NewCommandService()
	tests := []struct {
		name string
		want managers.ServiceManager
	}{
		{name: "success", want: sm},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.WithServiceManager(tt.want)
			if got.ServiceManager() != tt.want {
				t.Errorf("ServiceManager() = %v, want %v", got.ServiceManager(), tt.want)
			}
		})
	}
}

func Test_commandServiceFinal_Execute(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantOut    []string
		wantNotOut []string
	}{
		{
			name:    "version prints build metadata",
			args:    []string{"version"},
			wantErr: false,
			wantOut: []string{"freelunch", "commit:", "built:", "go:", "platform:"},
		},
		{
			name:    "bare invocation prints help",
			args:    []string{},
			wantErr: false,
			wantOut: []string{"freelunch", "Available Commands", "version"},
		},
		{
			name:    "unknown command fails",
			args:    []string{"no-such-command"},
			wantErr: true,
		},
		{
			name:    "version rejects arguments",
			args:    []string{"version", "extra"},
			wantErr: true,
		},
		{
			// The default completion command is hidden deliberately.
			name:       "completion is hidden from help",
			args:       []string{},
			wantErr:    false,
			wantNotOut: []string{"completion"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, out, _, ctx := newTestService(t)

			err := s.Execute(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}

			for _, want := range tt.wantOut {
				if !strings.Contains(out.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, out.String())
				}
			}
			for _, notWant := range tt.wantNotOut {
				if strings.Contains(out.String(), notWant) {
					t.Errorf("output unexpectedly contains %q; got:\n%s", notWant, out.String())
				}
			}
		})
	}
}

// Test_commandServiceFinal_Execute_versionJSON asserts the machine-readable form
// is real JSON, since it is what a script would parse.
func Test_commandServiceFinal_Execute_versionJSON(t *testing.T) {
	s, out, _, ctx := newTestService(t)

	if err := s.Execute(ctx, []string{"version", "--json"}); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v; got:\n%s", err, out.String())
	}

	for _, key := range []string{"version", "commit", "date", "goVersion", "platform"} {
		if _, ok := got[key]; !ok {
			t.Errorf("JSON output missing key %q; got %v", key, got)
		}
	}
}

func Test_commandServiceFinal_RunVersion(t *testing.T) {
	tests := []struct {
		name    string
		asJSON  bool
		wantErr bool
		check   func(t *testing.T, out string)
	}{
		{
			name:   "text form",
			asJSON: false,
			check: func(t *testing.T, out string) {
				if !strings.HasPrefix(out, "freelunch ") {
					t.Errorf("text output = %q, want it to start with %q", out, "freelunch ")
				}
			},
		},
		{
			name:   "json form",
			asJSON: true,
			check: func(t *testing.T, out string) {
				var v map[string]any
				if err := json.Unmarshal([]byte(out), &v); err != nil {
					t.Errorf("json output invalid: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _, _, ctx := newTestService(t)
			var buf bytes.Buffer

			if err := s.RunVersion(ctx, &buf, tt.asJSON); (err != nil) != tt.wantErr {
				t.Fatalf("RunVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.check(t, buf.String())
		})
	}
}

func TestExitError(t *testing.T) {
	inner := errors.New("underlying failure")

	tests := []struct {
		name        string
		err         ExitError
		wantMessage string
		wantUnwrap  error
	}{
		{
			name:        "wraps an error",
			err:         ExitError{Code: 3, Err: inner},
			wantMessage: "underlying failure",
			wantUnwrap:  inner,
		},
		{
			name:        "bare code",
			err:         ExitError{Code: 2},
			wantMessage: "exit status 2",
			wantUnwrap:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMessage {
				t.Errorf("Error() = %q, want %q", got, tt.wantMessage)
			}
			if got := tt.err.Unwrap(); !errors.Is(got, tt.wantUnwrap) {
				t.Errorf("Unwrap() = %v, want %v", got, tt.wantUnwrap)
			}
		})
	}
}

// TestExitError_errorsAs is the behaviour main relies on to pick an exit code.
func TestExitError_errorsAs(t *testing.T) {
	err := error(ExitError{Code: 7, Err: errors.New("boom")})

	var exit ExitError
	if !errors.As(err, &exit) {
		t.Fatal("errors.As() = false, want true")
	}
	if exit.Code != 7 {
		t.Errorf("Code = %d, want 7", exit.Code)
	}
}

// fakeClusterService is a hand-written ClusterService for the cluster command
// tests. The command bodies are what is under test here, so it records calls and
// returns canned answers rather than driving k3d or Docker.
type fakeClusterService struct {
	sm managers.ServiceManager

	createErr error
	deleteErr error
	status    *managers.ClusterStatus
	statusErr error

	createCalls int
	deleteCalls int
	statusCalls int
}

func (f *fakeClusterService) Start(_ context.Context) error   { return nil }
func (f *fakeClusterService) Close(_ context.Context) error   { return nil }
func (f *fakeClusterService) Healthy(_ context.Context) error { return nil }

func (f *fakeClusterService) WithServiceManager(sm managers.ServiceManager) managers.ClusterService {
	f.sm = sm
	return f
}

func (f *fakeClusterService) ServiceManager() managers.ServiceManager { return f.sm }

func (f *fakeClusterService) Create(_ context.Context) error {
	f.createCalls++
	return f.createErr
}

func (f *fakeClusterService) Delete(_ context.Context) error {
	f.deleteCalls++
	return f.deleteErr
}

func (f *fakeClusterService) Status(_ context.Context) (*managers.ClusterStatus, error) {
	f.statusCalls++
	return f.status, f.statusErr
}

// newTestServiceWithCluster wires a fake ClusterService alongside the real
// command service, so the cluster commands can be exercised without a cluster.
func newTestServiceWithCluster(
	t *testing.T, cs *fakeClusterService,
) (s *commandServiceFinal, out, errOut *bytes.Buffer, ctx context.Context) {
	t.Helper()

	sm, ctx := NewManagerForTests()
	sm.WithClusterService(cs)
	s = sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)

	out, errOut = &bytes.Buffer{}, &bytes.Buffer{}
	s.out = out
	s.errOut = errOut

	return s, out, errOut, ctx
}

func Test_commandServiceFinal_RunInstall(t *testing.T) {
	createErr := errors.New("k3d cluster create failed")

	tests := []struct {
		name    string
		cluster *fakeClusterService
		wantErr bool
		wantOut []string
	}{
		{
			name:    "creates the cluster",
			cluster: &fakeClusterService{},
			wantErr: false,
			wantOut: []string{"Cluster created"},
		},
		{
			name:    "surfaces a creation failure",
			cluster: &fakeClusterService{createErr: createErr},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _, _, ctx := newTestServiceWithCluster(t, tt.cluster)
			var buf bytes.Buffer

			err := s.RunInstall(ctx, &buf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.cluster.createCalls != 1 {
				t.Errorf("Create() called %d times, want 1", tt.cluster.createCalls)
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, buf.String())
				}
			}
		})
	}
}

// Test_commandServiceFinal_RunInstall_wrapsClusterError asserts the underlying
// error is not flattened into a string, so a caller can still match on it.
func Test_commandServiceFinal_RunInstall_wrapsClusterError(t *testing.T) {
	createErr := errors.New("no tools")
	s, _, _, ctx := newTestServiceWithCluster(t, &fakeClusterService{createErr: createErr})

	err := s.RunInstall(ctx, &bytes.Buffer{})
	if !errors.Is(err, createErr) {
		t.Errorf("RunInstall() error = %v, want it to wrap %v", err, createErr)
	}
}

func Test_commandServiceFinal_RunUninstall(t *testing.T) {
	deleteErr := errors.New("k3d cluster delete failed")

	tests := []struct {
		name    string
		cluster *fakeClusterService
		wantErr bool
		wantOut []string
	}{
		{
			name:    "deletes the cluster",
			cluster: &fakeClusterService{},
			wantErr: false,
			wantOut: []string{"Cluster deleted"},
		},
		{
			name:    "surfaces a deletion failure",
			cluster: &fakeClusterService{deleteErr: deleteErr},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _, _, ctx := newTestServiceWithCluster(t, tt.cluster)
			var buf bytes.Buffer

			err := s.RunUninstall(ctx, &buf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunUninstall() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.cluster.deleteCalls != 1 {
				t.Errorf("Delete() called %d times, want 1", tt.cluster.deleteCalls)
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, buf.String())
				}
			}
		})
	}
}

func Test_commandServiceFinal_RunStatus(t *testing.T) {
	statusErr := errors.New("pinned cluster tools not found")

	tests := []struct {
		name       string
		cluster    *fakeClusterService
		wantErr    bool
		wantOut    []string
		wantNotOut []string
	}{
		{
			name: "reports a running cluster and its nodes",
			cluster: &fakeClusterService{status: &managers.ClusterStatus{
				Name:    "freelunch",
				Running: true,
				Nodes:   []string{"k3d-freelunch-server-0", "k3d-freelunch-agent-0"},
			}},
			wantOut: []string{"freelunch", "running", "k3d-freelunch-server-0", "k3d-freelunch-agent-0"},
		},
		{
			name: "reports a cluster that is not running",
			cluster: &fakeClusterService{status: &managers.ClusterStatus{
				Name: "freelunch",
			}},
			wantOut:    []string{"not running", "freelunch install"},
			wantNotOut: []string{"k3d-freelunch"},
		},
		{
			// The interface permits a nil status; it must not panic.
			name:    "tolerates a nil status",
			cluster: &fakeClusterService{status: nil},
			wantOut: []string{"not running"},
		},
		{
			name:    "surfaces a status failure",
			cluster: &fakeClusterService{statusErr: statusErr},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _, _, ctx := newTestServiceWithCluster(t, tt.cluster)
			var buf bytes.Buffer

			err := s.RunStatus(ctx, &buf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.cluster.statusCalls != 1 {
				t.Errorf("Status() called %d times, want 1", tt.cluster.statusCalls)
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, buf.String())
				}
			}
			for _, notWant := range tt.wantNotOut {
				if strings.Contains(buf.String(), notWant) {
					t.Errorf("output unexpectedly contains %q; got:\n%s", notWant, buf.String())
				}
			}
		})
	}
}

// Test_commandServiceFinal_Execute_clusterCommands covers the wiring end to end:
// that each command is registered, reaches the ClusterService, and rejects args.
func Test_commandServiceFinal_Execute_clusterCommands(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		cluster     *fakeClusterService
		wantErr     bool
		wantOut     []string
		wantCreate  int
		wantDelete  int
		wantStatusN int
	}{
		{
			name:       "install creates the cluster",
			args:       []string{"install"},
			cluster:    &fakeClusterService{},
			wantOut:    []string{"Cluster created"},
			wantCreate: 1,
		},
		{
			name:       "uninstall deletes the cluster",
			args:       []string{"uninstall"},
			cluster:    &fakeClusterService{},
			wantOut:    []string{"Cluster deleted"},
			wantDelete: 1,
		},
		{
			name:        "status reports the cluster",
			args:        []string{"status"},
			cluster:     &fakeClusterService{status: &managers.ClusterStatus{Name: "freelunch"}},
			wantOut:     []string{"not running"},
			wantStatusN: 1,
		},
		{
			name:    "install rejects arguments",
			args:    []string{"install", "extra"},
			cluster: &fakeClusterService{},
			wantErr: true,
		},
		{
			name:    "uninstall rejects arguments",
			args:    []string{"uninstall", "extra"},
			cluster: &fakeClusterService{},
			wantErr: true,
		},
		{
			name:    "status rejects arguments",
			args:    []string{"status", "extra"},
			cluster: &fakeClusterService{},
			wantErr: true,
		},
		{
			name:       "install propagates a cluster failure",
			args:       []string{"install"},
			cluster:    &fakeClusterService{createErr: errors.New("docker is not running")},
			wantErr:    true,
			wantCreate: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, out, _, ctx := newTestServiceWithCluster(t, tt.cluster)

			err := s.Execute(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if tt.cluster.createCalls != tt.wantCreate {
				t.Errorf("Create() called %d times, want %d", tt.cluster.createCalls, tt.wantCreate)
			}
			if tt.cluster.deleteCalls != tt.wantDelete {
				t.Errorf("Delete() called %d times, want %d", tt.cluster.deleteCalls, tt.wantDelete)
			}
			if tt.cluster.statusCalls != tt.wantStatusN {
				t.Errorf("Status() called %d times, want %d", tt.cluster.statusCalls, tt.wantStatusN)
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(out.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, out.String())
				}
			}
		})
	}
}

// Test_commandServiceFinal_rootHelpAdvertisesRealCommands guards the defect
// Phase 6 closes: the root help promised a monorepo scaffold that does not
// exist, while the commands that do exist went unmentioned.
func Test_commandServiceFinal_rootHelpAdvertisesRealCommands(t *testing.T) {
	s, out, _, ctx := newTestServiceWithCluster(t, &fakeClusterService{})

	if err := s.Execute(ctx, []string{}); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	for _, want := range []string{"install", "uninstall", "status", "version"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("root help missing command %q; got:\n%s", want, out.String())
		}
	}
	// Nothing scaffolds a monorepo yet; `freelunch init` is roadmap 1.1's
	// remaining half.
	if strings.Contains(out.String(), "scaffold") {
		t.Errorf("root help still advertises scaffolding; got:\n%s", out.String())
	}
}
