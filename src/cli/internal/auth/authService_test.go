package auth

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// NewManagerForTests builds a container for this package's tests. By convention
// it is duplicated per package rather than imported from another one.
func NewManagerForTests() (managers.ServiceManager, context.Context) {
	return managers.NewManager(), context.Background()
}

// call records one invocation the service made through the Runner.
type call struct {
	name string
	args []string
}

// fakeRunner stands in for exec. Each entry in fail is matched against the
// first two arguments, so a test can fail `apply -f` while letting
// `rollout restart` succeed.
type fakeRunner struct {
	calls   []call
	out     map[string][]byte
	fail    map[string]error
	failAll error
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{out: map[string][]byte{}, fail: map[string]error{}}
}

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
	// The configmap dry-run is invoked with -n first; identify calls by any
	// two-token window so the fail/out maps stay simple.
	for k2 := range f.fail {
		if strings.Contains(strings.Join(args, " "), k2) {
			return []byte("boom"), f.fail[k2]
		}
	}
	for k2 := range f.out {
		if strings.Contains(strings.Join(args, " "), k2) {
			return f.out[k2], nil
		}
	}
	_ = k
	return []byte("ok"), nil
}

// fakeGetter stands in for the HTTP client Status probes discovery with.
type fakeGetter struct {
	resp *http.Response
	err  error
}

func (f *fakeGetter) Get(_ string) (*http.Response, error) {
	return f.resp, f.err
}

// jsonResponse fabricates the response a fakeGetter hands to Status. The
// bodyclose findings on its call sites are false positives — Status, the code
// under test, is what closes these bodies — hence the nolint directives there.
func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// withTools points the bin dir at a temp dir containing a kubectl stub.
func withTools(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "kubectl"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("cannot write kubectl stub: %v", err)
	}
	t.Setenv(BinDirEnvVar, dir)
}

// withoutTools points the bin dir somewhere empty.
func withoutTools(t *testing.T) {
	t.Helper()
	t.Setenv(BinDirEnvVar, t.TempDir())
}

func newServiceForTest(t *testing.T, g Getter) (*authServiceFinal, *fakeRunner, context.Context) {
	t.Helper()

	sm, ctx := NewManagerForTests()
	runner := newFakeRunner()
	if g == nil {
		g = &fakeGetter{err: errors.New("no server")}
	}
	sm.WithAuthService(NewAuthServiceWithDeps(runner, g))

	svc, ok := sm.AuthService().(*authServiceFinal)
	if !ok {
		t.Fatalf("AuthService() is %T, want *authServiceFinal", sm.AuthService())
	}

	return svc, runner, ctx
}

func TestNewAuthService(t *testing.T) {
	if got := NewAuthService(); got == nil {
		t.Fatal("NewAuthService() = nil")
	}
}

func Test_authServiceFinal_Start(t *testing.T) {
	withTools(t)
	svc, runner, ctx := newServiceForTest(t, nil)

	if err := svc.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	// Start must not touch the system — the framework rule this package
	// inherits from ClusterService.
	if len(runner.calls) != 0 {
		t.Errorf("Start() ran %d external commands, want 0", len(runner.calls))
	}
}

func Test_authServiceFinal_Close(t *testing.T) {
	withTools(t)
	svc, _, ctx := newServiceForTest(t, nil)

	if err := svc.Close(ctx); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}

func Test_authServiceFinal_Healthy(t *testing.T) {
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
			svc, _, ctx := newServiceForTest(t, nil)

			if err := svc.Healthy(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Healthy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authServiceFinal_WithServiceManager(t *testing.T) {
	sm, _ := NewManagerForTests()
	s := NewAuthService()

	got := s.WithServiceManager(sm)
	if got.ServiceManager() != sm {
		t.Errorf("ServiceManager() = %v, want %v", got.ServiceManager(), sm)
	}
}

func Test_authServiceFinal_Install(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		fail      map[string]error
		wantErr   bool
		wantCalls int
	}{
		{
			// apply manifests, configmap dry-run, apply configmap, rollout restart.
			name:      "success applies manifests, realm and restart",
			setup:     withTools,
			wantCalls: 5, // namespace, manifests, configmap render, configmap apply, restart
		},
		{
			name:    "tools missing",
			setup:   withoutTools,
			wantErr: true,
		},
		{
			name:    "apply failure surfaces",
			setup:   withTools,
			fail:    map[string]error{"apply -f": errors.New("connection refused")},
			wantErr: true,
		},
		{
			name:    "rollout restart failure surfaces",
			setup:   withTools,
			fail:    map[string]error{"rollout restart": errors.New("not found")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			svc, runner, ctx := newServiceForTest(t, nil)
			for k, v := range tt.fail {
				runner.fail[k] = v
			}

			err := svc.Install(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Install() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantCalls > 0 && len(runner.calls) != tt.wantCalls {
				t.Errorf("Install() made %d calls, want %d: %v", len(runner.calls), tt.wantCalls, runner.calls)
			}
		})
	}
}

// Test_authServiceFinal_InstallWritesManifests pins that the embedded files
// are materialised and passed by path, and cleaned up afterwards.
func Test_authServiceFinal_InstallWritesManifests(t *testing.T) {
	withTools(t)
	svc, runner, ctx := newServiceForTest(t, nil)

	if err := svc.Install(ctx); err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	// The shared namespace first, so `install --only auth` works on a fresh
	// cluster; the component's own manifests second.
	first, second := runner.calls[0], runner.calls[1]
	if first.args[0] != "apply" || first.args[1] != "-f" ||
		!strings.HasSuffix(first.args[2], "namespace.yaml") {
		t.Errorf("first call = %v, want apply -f <path>/namespace.yaml", first.args)
	}
	if second.args[0] != "apply" || second.args[1] != "-f" ||
		!strings.HasSuffix(second.args[2], "keycloak.yaml") {
		t.Errorf("second call = %v, want apply -f <path>/keycloak.yaml", second.args)
	}
	if _, err := os.Stat(second.args[2]); !os.IsNotExist(err) {
		t.Errorf("manifest %s still exists after Install(); want it removed", second.args[2])
	}
}

// Test_componentManifestDoesNotOwnNamespace guards the split: the namespace
// is shared and lives only in namespace.yaml. If it crept back into
// keycloak.yaml, Delete would remove it and take the secrets store down.
func Test_componentManifestDoesNotOwnNamespace(t *testing.T) {
	if strings.Contains(string(keycloakManifests), "kind: Namespace") {
		t.Error("keycloak.yaml declares a Namespace; it belongs in namespace.yaml only")
	}
	if !strings.Contains(string(namespaceManifest), "kind: Namespace") ||
		!strings.Contains(string(namespaceManifest), "name: "+Namespace) {
		t.Errorf("namespace.yaml must declare the %s Namespace", Namespace)
	}
}

// Test_namespaceManifestMatchesSecrets guards the deliberate duplication of
// namespace.yaml across component packages: go:embed cannot reach across
// packages and they may not import each other, so each carries a copy, and
// this test is what keeps the copies identical.
func Test_namespaceManifestMatchesSecrets(t *testing.T) {
	other, err := os.ReadFile(filepath.Join("..", "secrets", "namespace.yaml"))
	if err != nil {
		t.Fatalf("reading the secrets copy: %v", err)
	}
	if !bytes.Equal(other, namespaceManifest) {
		t.Error("auth/namespace.yaml and secrets/namespace.yaml differ; edit them together")
	}
}

func Test_authServiceFinal_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		fail    map[string]error
		wantErr bool
	}{
		{name: "success", setup: withTools},
		{name: "tools missing", setup: withoutTools, wantErr: true},
		{
			name:    "delete failure surfaces",
			setup:   withTools,
			fail:    map[string]error{"delete -f": errors.New("boom")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			svc, runner, ctx := newServiceForTest(t, nil)
			for k, v := range tt.fail {
				runner.fail[k] = v
			}

			err := svc.Delete(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.name == "success" {
				got := runner.calls[0].args
				if got[0] != "delete" || got[len(got)-1] != "--ignore-not-found" {
					t.Errorf("Delete() args = %v, want delete ... --ignore-not-found", got)
				}
				// The namespace is shared: deleting it would take the secrets
				// store down with the auth service.
				for _, c := range runner.calls {
					if strings.Contains(strings.Join(c.args, " "), "namespace.yaml") {
						t.Errorf("Delete() touched namespace.yaml: %v", c.args)
					}
				}
			}
		})
	}
}

func Test_authServiceFinal_Status(t *testing.T) {
	tests := []struct {
		name   string
		getter Getter
		want   managers.AuthStatus
	}{
		{
			name: "ready with issuer",
			getter: &fakeGetter{resp: jsonResponse(200, //nolint:bodyclose // closed by Status
				`{"issuer":"http://keycloak.localhost:8080/realms/freelunch"}`)},
			want: managers.AuthStatus{
				Ready:     true,
				Realm:     RealmName,
				IssuerURL: "http://keycloak.localhost:8080/realms/freelunch",
			},
		},
		{
			// Traefik answers before the pod does during a swap.
			name:   "503 is not ready, not an error",
			getter: &fakeGetter{resp: jsonResponse(503, "Service Unavailable")}, //nolint:bodyclose // closed by Status
			want:   managers.AuthStatus{},
		},
		{
			name:   "unreachable is not ready, not an error",
			getter: &fakeGetter{err: errors.New("connection refused")},
			want:   managers.AuthStatus{},
		},
		{
			name:   "garbage body is not ready, not an error",
			getter: &fakeGetter{resp: jsonResponse(200, "<html>not json</html>")}, //nolint:bodyclose // closed by Status
			want:   managers.AuthStatus{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTools(t)
			svc, _, ctx := newServiceForTest(t, tt.getter)

			got, err := svc.Status(ctx)
			if err != nil {
				t.Fatalf("Status() error = %v, want nil", err)
			}
			if *got != tt.want {
				t.Errorf("Status() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

// Test_realmNameMatchesJSON guards the constant against the committed realm.
// Status probes discovery by RealmName; if the JSON drifted, a healthy server
// would be reported as absent.
func Test_realmNameMatchesJSON(t *testing.T) {
	if !strings.Contains(string(realmJSON), `"realm": "`+RealmName+`"`) {
		t.Errorf("freelunch-realm.json does not declare realm %q", RealmName)
	}
}

// Test_realmDefinesThreePersonas pins decision D5: three Personas, with the
// grants and the hotfix permission modelled separately.
func Test_realmDefinesThreePersonas(t *testing.T) {
	for _, g := range []string{"platform-admin", "platform-engineer", "developer"} {
		if !strings.Contains(string(realmJSON), `"name": "`+g+`"`) {
			t.Errorf("realm is missing Persona group %q", g)
		}
	}
	for _, g := range []string{"developer-tech-lead", "platform-tech-lead"} {
		if !strings.Contains(string(realmJSON), `"name": "`+g+`"`) {
			t.Errorf("realm is missing temporary-grant group %q", g)
		}
	}
	if !strings.Contains(string(realmJSON), `"name": "hotfix"`) {
		t.Error("realm is missing the hotfix realm role")
	}
}

// Test_manifestsPinImage guards against a floating image tag sneaking in.
func Test_manifestsPinImage(t *testing.T) {
	if strings.Contains(string(keycloakManifests), ":latest") {
		t.Error("keycloak.yaml uses a :latest tag; pin a version")
	}
	if !strings.Contains(string(keycloakManifests), "quay.io/keycloak/keycloak:") {
		t.Error("keycloak.yaml does not reference a pinned keycloak image")
	}
}
