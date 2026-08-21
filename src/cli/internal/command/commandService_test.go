package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
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

// fakeAuthService is a hand-written AuthService for the command tests.
type fakeAuthService struct {
	sm managers.ServiceManager

	installErr error
	deleteErr  error
	status     *managers.AuthStatus
	statusErr  error

	installCalls int
	deleteCalls  int
	statusCalls  int
}

func (f *fakeAuthService) Start(_ context.Context) error   { return nil }
func (f *fakeAuthService) Close(_ context.Context) error   { return nil }
func (f *fakeAuthService) Healthy(_ context.Context) error { return nil }

func (f *fakeAuthService) WithServiceManager(sm managers.ServiceManager) managers.AuthService {
	f.sm = sm
	return f
}

func (f *fakeAuthService) ServiceManager() managers.ServiceManager { return f.sm }

func (f *fakeAuthService) Install(_ context.Context) error {
	f.installCalls++
	return f.installErr
}

func (f *fakeAuthService) Delete(_ context.Context) error {
	f.deleteCalls++
	return f.deleteErr
}

func (f *fakeAuthService) Status(_ context.Context) (*managers.AuthStatus, error) {
	f.statusCalls++
	return f.status, f.statusErr
}

// newTestServiceWithCluster wires a fake ClusterService alongside the real
// command service, so the cluster commands can be exercised without a cluster.
// The auth slot stays a no-op unless the test registers a fake of its own.
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

			err := s.RunInstall(ctx, &buf, nil, nil)
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

	err := s.RunInstall(ctx, &bytes.Buffer{}, nil, nil)
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
			wantOut: []string{"Demo environment deleted"},
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
			wantOut: []string{"freelunch", "running", "k3d-freelunch-server-0", "k3d-freelunch-agent-0", "Auth service is not ready"},
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
			wantOut:    []string{"Demo environment deleted"},
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

	for _, want := range []string{"init", "install", "uninstall", "status", "version"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("root help missing command %q; got:\n%s", want, out.String())
		}
	}
}

// newTestServiceWithBoth wires fake cluster and auth services.
func newTestServiceWithBoth(
	t *testing.T, cs *fakeClusterService, as *fakeAuthService,
) (s *commandServiceFinal, out *bytes.Buffer, ctx context.Context) {
	t.Helper()

	sm, ctx := NewManagerForTests()
	sm.WithClusterService(cs)
	sm.WithAuthService(as)
	s = sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)

	out = &bytes.Buffer{}
	s.out = out
	s.errOut = &bytes.Buffer{}

	return s, out, ctx
}

func Test_selectComponents(t *testing.T) {
	tests := []struct {
		name    string
		only    []string
		skip    []string
		want    []string
		wantErr bool
	}{
		{name: "default is everything in order", want: []string{"cluster", "auth", "secrets"}},
		{name: "only cluster", only: []string{"cluster"}, want: []string{"cluster"}},
		{name: "only auth", only: []string{"auth"}, want: []string{"auth"}},
		{name: "skip auth and secrets", skip: []string{"auth", "secrets"}, want: []string{"cluster"}},
		{name: "skip beats only", only: []string{"auth"}, skip: []string{"auth"}, want: nil},
		{name: "unknown only is an error", only: []string{"vault"}, wantErr: true},
		{name: "unknown skip is an error", skip: []string{"clutser"}, wantErr: true},
		{
			// --only out of order still installs in dependency order.
			name: "order is canonical regardless of flag order",
			only: []string{"secrets", "auth", "cluster"},
			want: []string{"cluster", "auth", "secrets"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectComponents(tt.only, tt.skip)
			if (err != nil) != tt.wantErr {
				t.Fatalf("selectComponents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectComponents(%v, %v) = %v, want %v", tt.only, tt.skip, got, tt.want)
			}
		})
	}
}

func Test_commandServiceFinal_RunInstall_components(t *testing.T) {
	tests := []struct {
		name        string
		only, skip  []string
		auth        *fakeAuthService
		cluster     *fakeClusterService
		wantErr     bool
		wantCreate  int
		wantInstall int
		wantOut     []string
		wantNotOut  []string
	}{
		{
			name:        "default installs cluster then auth",
			cluster:     &fakeClusterService{},
			auth:        &fakeAuthService{},
			wantCreate:  1,
			wantInstall: 1,
			wantOut:     []string{"Cluster created", "Auth service installed", "Done."},
		},
		{
			name:        "--only cluster leaves auth untouched",
			only:        []string{"cluster"},
			cluster:     &fakeClusterService{},
			auth:        &fakeAuthService{},
			wantCreate:  1,
			wantInstall: 0,
			wantNotOut:  []string{"Auth service"},
		},
		{
			name:        "--skip auth leaves auth untouched",
			skip:        []string{"auth"},
			cluster:     &fakeClusterService{},
			auth:        &fakeAuthService{},
			wantCreate:  1,
			wantInstall: 0,
		},
		{
			name:        "--only auth skips the cluster",
			only:        []string{"auth"},
			cluster:     &fakeClusterService{},
			auth:        &fakeAuthService{},
			wantCreate:  0,
			wantInstall: 1,
		},
		{
			name:        "cluster failure stops before auth",
			cluster:     &fakeClusterService{createErr: errors.New("docker down")},
			auth:        &fakeAuthService{},
			wantErr:     true,
			wantCreate:  1,
			wantInstall: 0,
		},
		{
			name:        "auth failure surfaces after cluster",
			cluster:     &fakeClusterService{},
			auth:        &fakeAuthService{installErr: errors.New("no kubectl")},
			wantErr:     true,
			wantCreate:  1,
			wantInstall: 1,
		},
		{
			name:    "unknown component fails before any work",
			only:    []string{"nonsense"},
			cluster: &fakeClusterService{},
			auth:    &fakeAuthService{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, out, ctx := newTestServiceWithBoth(t, tt.cluster, tt.auth)
			var buf bytes.Buffer

			err := s.RunInstall(ctx, &buf, tt.only, tt.skip)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.cluster.createCalls != tt.wantCreate {
				t.Errorf("Create() called %d times, want %d", tt.cluster.createCalls, tt.wantCreate)
			}
			if tt.auth.installCalls != tt.wantInstall {
				t.Errorf("Install() called %d times, want %d", tt.auth.installCalls, tt.wantInstall)
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
			_ = out
		})
	}
}

func Test_commandServiceFinal_RunUninstall_deletesAuthFirst(t *testing.T) {
	cs := &fakeClusterService{}
	as := &fakeAuthService{}
	s, _, ctx := newTestServiceWithBoth(t, cs, as)
	var buf bytes.Buffer

	if err := s.RunUninstall(ctx, &buf); err != nil {
		t.Fatalf("RunUninstall() error = %v", err)
	}
	if as.deleteCalls != 1 || cs.deleteCalls != 1 {
		t.Errorf("deletes = auth %d, cluster %d; want 1 and 1", as.deleteCalls, cs.deleteCalls)
	}
}

func Test_commandServiceFinal_RunStatus_reportsAuth(t *testing.T) {
	tests := []struct {
		name    string
		auth    *fakeAuthService
		wantOut []string
	}{
		{
			name: "ready auth prints realm and issuer",
			auth: &fakeAuthService{status: &managers.AuthStatus{
				Ready:     true,
				Realm:     "freelunch",
				IssuerURL: "http://keycloak.localhost:8080/realms/freelunch",
			}},
			wantOut: []string{"realm \"freelunch\"", "issuer http://keycloak.localhost:8080/realms/freelunch"},
		},
		{
			name:    "auth not ready says so",
			auth:    &fakeAuthService{status: &managers.AuthStatus{}},
			wantOut: []string{"Auth service is not ready"},
		},
		{
			// The interface permits nil; it must not panic.
			name:    "nil auth status treated as not ready",
			auth:    &fakeAuthService{status: nil},
			wantOut: []string{"Auth service is not ready"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &fakeClusterService{status: &managers.ClusterStatus{
				Name: "freelunch", Running: true, Nodes: []string{"n1"},
			}}
			s, _, ctx := newTestServiceWithBoth(t, cs, tt.auth)
			var buf bytes.Buffer

			if err := s.RunStatus(ctx, &buf); err != nil {
				t.Fatalf("RunStatus() error = %v", err)
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, buf.String())
				}
			}
		})
	}
}

// fakeSecretsService is a hand-written SecretsService for the command tests.
type fakeSecretsService struct {
	sm managers.ServiceManager

	installErr error
	deleteErr  error
	putErr     error
	status     *managers.SecretsStatus
	statusErr  error

	installCalls int
	deleteCalls  int
	putCalls     []string
	statusCalls  int
}

func (f *fakeSecretsService) Start(_ context.Context) error   { return nil }
func (f *fakeSecretsService) Close(_ context.Context) error   { return nil }
func (f *fakeSecretsService) Healthy(_ context.Context) error { return nil }

func (f *fakeSecretsService) WithServiceManager(sm managers.ServiceManager) managers.SecretsService {
	f.sm = sm
	return f
}

func (f *fakeSecretsService) ServiceManager() managers.ServiceManager { return f.sm }

func (f *fakeSecretsService) Install(_ context.Context) error {
	f.installCalls++
	return f.installErr
}

func (f *fakeSecretsService) Delete(_ context.Context) error {
	f.deleteCalls++
	return f.deleteErr
}

func (f *fakeSecretsService) Status(_ context.Context) (*managers.SecretsStatus, error) {
	f.statusCalls++
	return f.status, f.statusErr
}

func (f *fakeSecretsService) PutSecret(_ context.Context, path, key, value string) error {
	f.putCalls = append(f.putCalls, path+"/"+key+"="+value)
	return f.putErr
}

func Test_commandServiceFinal_RunInstall_secrets(t *testing.T) {
	tests := []struct {
		name        string
		only        []string
		secrets     *fakeSecretsService
		wantErr     bool
		wantInstall int
		wantPuts    int
	}{
		{
			name:        "secrets install seeds the demo credential",
			only:        []string{"secrets"},
			secrets:     &fakeSecretsService{},
			wantInstall: 1,
			wantPuts:    1,
		},
		{
			name:        "seed failure surfaces",
			only:        []string{"secrets"},
			secrets:     &fakeSecretsService{putErr: errors.New("no pod")},
			wantErr:     true,
			wantInstall: 1,
			wantPuts:    1,
		},
		{
			name:        "install failure stops before seeding",
			only:        []string{"secrets"},
			secrets:     &fakeSecretsService{installErr: errors.New("timed out")},
			wantErr:     true,
			wantInstall: 1,
			wantPuts:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm, ctx := NewManagerForTests()
			sm.WithClusterService(&fakeClusterService{})
			sm.WithSecretsService(tt.secrets)
			s := sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)
			var buf bytes.Buffer

			err := s.RunInstall(ctx, &buf, tt.only, nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.secrets.installCalls != tt.wantInstall {
				t.Errorf("Install() called %d times, want %d", tt.secrets.installCalls, tt.wantInstall)
			}
			if len(tt.secrets.putCalls) != tt.wantPuts {
				t.Errorf("PutSecret() called %d times, want %d", len(tt.secrets.putCalls), tt.wantPuts)
			}
			if tt.wantPuts == 1 && !tt.wantErr {
				want := "example_service/api-key=example-api-key-value"
				if tt.secrets.putCalls[0] != want {
					t.Errorf("seeded %q, want %q", tt.secrets.putCalls[0], want)
				}
			}
		})
	}
}

func Test_commandServiceFinal_RunUninstall_deletesSecretsFirst(t *testing.T) {
	cs := &fakeClusterService{}
	as := &fakeAuthService{}
	ss := &fakeSecretsService{}
	sm, ctx := NewManagerForTests()
	sm.WithClusterService(cs)
	sm.WithAuthService(as)
	sm.WithSecretsService(ss)
	s := sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)

	if err := s.RunUninstall(ctx, &bytes.Buffer{}); err != nil {
		t.Fatalf("RunUninstall() error = %v", err)
	}
	if ss.deleteCalls != 1 || as.deleteCalls != 1 || cs.deleteCalls != 1 {
		t.Errorf("deletes = secrets %d, auth %d, cluster %d; want 1,1,1",
			ss.deleteCalls, as.deleteCalls, cs.deleteCalls)
	}
}

func Test_commandServiceFinal_RunStatus_reportsSecrets(t *testing.T) {
	tests := []struct {
		name    string
		secrets *fakeSecretsService
		wantOut []string
	}{
		{
			name:    "ready store prints the engine",
			secrets: &fakeSecretsService{status: &managers.SecretsStatus{Ready: true, Engine: "secret/ (kv v2)"}},
			wantOut: []string{"Secrets store is ready: secret/ (kv v2)"},
		},
		{
			name:    "sealed store gets its own loud line",
			secrets: &fakeSecretsService{status: &managers.SecretsStatus{Sealed: true}},
			wantOut: []string{"SEALED"},
		},
		{
			name:    "absent store says not ready",
			secrets: &fakeSecretsService{status: &managers.SecretsStatus{}},
			wantOut: []string{"Secrets store is not ready"},
		},
		{
			name:    "nil status treated as not ready",
			secrets: &fakeSecretsService{status: nil},
			wantOut: []string{"Secrets store is not ready"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm, ctx := NewManagerForTests()
			sm.WithClusterService(&fakeClusterService{status: &managers.ClusterStatus{
				Name: "freelunch", Running: true, Nodes: []string{"n1"},
			}})
			sm.WithAuthService(&fakeAuthService{status: &managers.AuthStatus{
				Ready: true, Realm: "freelunch", IssuerURL: "http://x",
			}})
			sm.WithSecretsService(tt.secrets)
			s := sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)
			var buf bytes.Buffer

			if err := s.RunStatus(ctx, &buf); err != nil {
				t.Fatalf("RunStatus() error = %v", err)
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, buf.String())
				}
			}
		})
	}
}

// fakeScaffoldService is a hand-written ScaffoldService for the init tests.
type fakeScaffoldService struct {
	sm managers.ServiceManager

	initErr   error
	initCalls []string
}

func (f *fakeScaffoldService) Start(_ context.Context) error   { return nil }
func (f *fakeScaffoldService) Close(_ context.Context) error   { return nil }
func (f *fakeScaffoldService) Healthy(_ context.Context) error { return nil }

func (f *fakeScaffoldService) WithServiceManager(sm managers.ServiceManager) managers.ScaffoldService {
	f.sm = sm
	return f
}

func (f *fakeScaffoldService) ServiceManager() managers.ServiceManager { return f.sm }

func (f *fakeScaffoldService) Init(_ context.Context, dir, product string) error {
	f.initCalls = append(f.initCalls, dir+"|"+product)
	return f.initErr
}

func Test_commandServiceFinal_RunInit(t *testing.T) {
	tests := []struct {
		name     string
		scaffold *fakeScaffoldService
		wantErr  bool
		wantOut  []string
	}{
		{
			name:     "success reports the location and next step",
			scaffold: &fakeScaffoldService{},
			wantOut:  []string{"Monorepo created at my-company", "product: shop", "freelunch install"},
		},
		{
			name:     "scaffold failure surfaces",
			scaffold: &fakeScaffoldService{initErr: errors.New("already exists")},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm, ctx := NewManagerForTests()
			sm.WithScaffoldService(tt.scaffold)
			s := sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)
			var buf bytes.Buffer

			err := s.RunInit(ctx, &buf, "my-company", "shop")
			if (err != nil) != tt.wantErr {
				t.Fatalf("RunInit() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.scaffold.initCalls) != 1 || tt.scaffold.initCalls[0] != "my-company|shop" {
				t.Errorf("Init() calls = %v, want [my-company|shop]", tt.scaffold.initCalls)
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("output missing %q; got:\n%s", want, buf.String())
				}
			}
		})
	}
}

// Test_commandServiceFinal_Execute_init covers the wired command end to end.
func Test_commandServiceFinal_Execute_init(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		wantCall string
	}{
		{
			name:     "default product is the placeholder",
			args:     []string{"init", "my-company"},
			wantCall: "my-company|example_product",
		},
		{
			name:     "--product renames",
			args:     []string{"init", "my-company", "--product", "shop"},
			wantCall: "my-company|shop",
		},
		{name: "no directory is an error", args: []string{"init"}, wantErr: true},
		{name: "two directories is an error", args: []string{"init", "a", "b"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scaffold := &fakeScaffoldService{}
			sm, ctx := NewManagerForTests()
			sm.WithScaffoldService(scaffold)
			s := sm.WithCommandService(NewCommandService()).CommandService().(*commandServiceFinal)
			s.out, s.errOut = &bytes.Buffer{}, &bytes.Buffer{}

			err := s.Execute(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if tt.wantCall != "" {
				if len(scaffold.initCalls) != 1 || scaffold.initCalls[0] != tt.wantCall {
					t.Errorf("Init() calls = %v, want [%s]", scaffold.initCalls, tt.wantCall)
				}
			}
		})
	}
}
