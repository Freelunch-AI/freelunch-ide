package secrets

import (
	"context"
	"errors"
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

// fakeRunner stands in for exec. Entries in out/fail are substring-matched
// against the joined argument list, so a test can fail `kv put` while letting
// `get pod` succeed.
type fakeRunner struct {
	calls   []call
	out     map[string][]byte
	fail    map[string]error
	failAll error
}

func newFakeRunner() *fakeRunner {
	r := &fakeRunner{out: map[string][]byte{}, fail: map[string]error{}}
	// Defaults that make the happy path work: a running pod exists, status
	// reports unsealed, and the KV v2 engine is mounted.
	r.out["get pod"] = []byte("openbao-abc123")
	r.out["status -format=json"] = []byte(`{"initialized": true, "sealed": false}`)
	r.out["secrets list"] = []byte(`{"secret/": {"type": "kv", "options": {"version": "2"}}}`)
	return r
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, call{name: name, args: args})

	if f.failAll != nil {
		return []byte("boom"), f.failAll
	}
	joined := strings.Join(args, " ")
	for k, v := range f.fail {
		if strings.Contains(joined, k) {
			// A failing command still has output — kubectl exec returns the
			// bao CLI's body together with its non-zero exit — so a test may
			// pair fail with out for the same key.
			if body, ok := f.out[k]; ok {
				return body, v
			}
			return []byte("boom"), v
		}
	}
	for k, v := range f.out {
		if strings.Contains(joined, k) {
			return v, nil
		}
	}
	return []byte("ok"), nil
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

func withoutTools(t *testing.T) {
	t.Helper()
	t.Setenv(BinDirEnvVar, t.TempDir())
}

func newServiceForTest(t *testing.T) (*secretsServiceFinal, *fakeRunner, context.Context) {
	t.Helper()

	sm, ctx := NewManagerForTests()
	runner := newFakeRunner()
	sm.WithSecretsService(NewSecretsServiceWithRunner(runner))

	svc, ok := sm.SecretsService().(*secretsServiceFinal)
	if !ok {
		t.Fatalf("SecretsService() is %T, want *secretsServiceFinal", sm.SecretsService())
	}

	return svc, runner, ctx
}

func TestNewSecretsService(t *testing.T) {
	if got := NewSecretsService(); got == nil {
		t.Fatal("NewSecretsService() = nil")
	}
}

func Test_secretsServiceFinal_Start(t *testing.T) {
	withTools(t)
	svc, runner, ctx := newServiceForTest(t)

	if err := svc.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	if len(runner.calls) != 0 {
		t.Errorf("Start() ran %d external commands, want 0", len(runner.calls))
	}
}

func Test_secretsServiceFinal_Close(t *testing.T) {
	withTools(t)
	svc, _, ctx := newServiceForTest(t)

	if err := svc.Close(ctx); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}

func Test_secretsServiceFinal_Healthy(t *testing.T) {
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

func Test_secretsServiceFinal_WithServiceManager(t *testing.T) {
	sm, _ := NewManagerForTests()
	s := NewSecretsService()

	got := s.WithServiceManager(sm)
	if got.ServiceManager() != sm {
		t.Errorf("ServiceManager() = %v, want %v", got.ServiceManager(), sm)
	}
}

func Test_secretsServiceFinal_Install(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T)
		fail      map[string]error
		wantErr   bool
		wantCalls int
	}{
		{
			// apply, rollout status.
			name:      "success applies and waits",
			setup:     withTools,
			wantCalls: 2,
		},
		{name: "tools missing", setup: withoutTools, wantErr: true},
		{
			name:    "apply failure surfaces",
			setup:   withTools,
			fail:    map[string]error{"apply -f": errors.New("connection refused")},
			wantErr: true,
		},
		{
			name:    "rollout timeout surfaces",
			setup:   withTools,
			fail:    map[string]error{"rollout status": errors.New("timed out")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			svc, runner, ctx := newServiceForTest(t)
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

// Test_secretsServiceFinal_InstallWritesManifests pins that the embedded file
// is materialised, passed by path, and cleaned up afterwards.
func Test_secretsServiceFinal_InstallWritesManifests(t *testing.T) {
	withTools(t)
	svc, runner, ctx := newServiceForTest(t)

	if err := svc.Install(ctx); err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	first := runner.calls[0]
	if first.args[0] != "apply" || !strings.HasSuffix(first.args[2], "openbao.yaml") {
		t.Errorf("first call = %v, want apply -f <path>/openbao.yaml", first.args)
	}
	if _, err := os.Stat(first.args[2]); !os.IsNotExist(err) {
		t.Errorf("manifest %s still exists after Install(); want it removed", first.args[2])
	}
}

func Test_secretsServiceFinal_Delete(t *testing.T) {
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
			svc, runner, ctx := newServiceForTest(t)
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
			}
		})
	}
}

func Test_secretsServiceFinal_Status(t *testing.T) {
	// kubectl exec passes the bao CLI's non-zero exit on as an error, and the
	// runner returns combined output, so this is what a sealed store's
	// `bao status` really looks like to the service.
	sealedOutput := []byte(`{"initialized": true, "sealed": true}` + "\ncommand terminated with exit code 2\n")
	sealedErr := errors.New("exit status 2")

	tests := []struct {
		name    string
		mut     func(r *fakeRunner)
		want    managers.SecretsStatus
		wantErr bool
	}{
		{
			name: "ready unsealed with the kv v2 engine",
			mut:  func(_ *fakeRunner) {},
			want: managers.SecretsStatus{Up: true, Ready: true, Engine: "secret/ (kv v2)"},
		},
		{
			name: "no running pod is not up, not an error",
			mut: func(r *fakeRunner) {
				r.fail["get pod"] = errors.New("no items")
			},
			want: managers.SecretsStatus{},
		},
		{
			// The failure that looks most like success: pod Ready, store sealed.
			// bao exits 2 here, so the JSON must be read despite the error.
			name: "sealed store is reported up and sealed, not ready",
			mut: func(r *fakeRunner) {
				r.out["status -format=json"] = sealedOutput
				r.fail["status -format=json"] = sealedErr
			},
			want: managers.SecretsStatus{Up: true, Sealed: true},
		},
		{
			name: "status failing with no JSON is not up, not an error",
			mut: func(r *fakeRunner) {
				r.out["status -format=json"] = []byte("connection refused")
				r.fail["status -format=json"] = errors.New("exit status 1")
			},
			want: managers.SecretsStatus{},
		},
		{
			name: "garbage status is not up, not an error",
			mut: func(r *fakeRunner) {
				r.out["status -format=json"] = []byte("not json at all")
			},
			want: managers.SecretsStatus{},
		},
		{
			name: "uninitialized store is not up",
			mut: func(r *fakeRunner) {
				r.out["status -format=json"] = []byte(`{"initialized": false, "sealed": true}`)
			},
			want: managers.SecretsStatus{Sealed: true},
		},
		{
			// Mohit's case: the store answers but secret/ is gone. Up, so the
			// user is told to inspect rather than wait; not Ready, because the
			// seed and 2.1's consumer have nowhere to read from.
			name: "missing engine is up but not ready",
			mut: func(r *fakeRunner) {
				r.out["secrets list"] = []byte(`{"sys/": {"type": "system"}}`)
			},
			want: managers.SecretsStatus{Up: true},
		},
		{
			// KV v1 at secret/ serves a different HTTP path; 2.1 would read
			// nothing. The engine is reported so the message can name it.
			name: "kv v1 engine is up but not ready",
			mut: func(r *fakeRunner) {
				r.out["secrets list"] = []byte(`{"secret/": {"type": "kv", "options": {"version": "1"}}}`)
			},
			want: managers.SecretsStatus{Up: true, Engine: "secret/ (kv v1)"},
		},
		{
			name: "engine listing failure on an up store is an error, not a verdict",
			mut: func(r *fakeRunner) {
				r.fail["secrets list"] = errors.New("permission denied")
			},
			wantErr: true,
		},
		{
			name: "unparseable engine listing on an up store is an error, not a verdict",
			mut: func(r *fakeRunner) {
				r.out["secrets list"] = []byte("not json")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTools(t)
			svc, runner, ctx := newServiceForTest(t)
			tt.mut(runner)

			got, err := svc.Status(ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Status() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if *got != tt.want {
				t.Errorf("Status() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func Test_unmarshalJSONObject(t *testing.T) {
	var v struct {
		A int `json:"a"`
	}
	if err := unmarshalJSONObject([]byte("noise {\"a\": 1}\ncommand terminated with exit code 2\n"), &v); err != nil || v.A != 1 {
		t.Errorf("unmarshalJSONObject() = %+v, %v; want A=1, nil", v, err)
	}
	if err := unmarshalJSONObject([]byte("no object here"), &v); err == nil {
		t.Error("unmarshalJSONObject() expected an error for output without a JSON object")
	}
	if err := unmarshalJSONObject([]byte("{not json}"), &v); err == nil {
		t.Error("unmarshalJSONObject() expected an error for a malformed object")
	}
}

func Test_secretsServiceFinal_PutSecret(t *testing.T) {
	tests := []struct {
		name    string
		fail    map[string]error
		wantErr bool
	}{
		{name: "success"},
		{
			name:    "kv put failure surfaces",
			fail:    map[string]error{"kv put": errors.New("permission denied")},
			wantErr: true,
		},
		{
			name:    "no pod surfaces as error",
			fail:    map[string]error{"get pod": errors.New("no items")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTools(t)
			svc, runner, ctx := newServiceForTest(t)
			for k, v := range tt.fail {
				runner.fail[k] = v
			}

			err := svc.PutSecret(ctx, "example_service", "api-key", "v")
			if (err != nil) != tt.wantErr {
				t.Fatalf("PutSecret() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// The write must go through the logical CLI path under the
				// mount — KV v2's data/ rewrite belongs to the API, never here.
				last := runner.calls[len(runner.calls)-1].args
				joined := strings.Join(last, " ")
				if !strings.Contains(joined, "kv put secret/example_service api-key=v") {
					t.Errorf("PutSecret() args = %q, want kv put secret/example_service", joined)
				}
			}
		})
	}
}

// Test_manifestsPinImage guards against a floating image tag sneaking in.
func Test_manifestsPinImage(t *testing.T) {
	if strings.Contains(string(openbaoManifests), ":latest") {
		t.Error("openbao.yaml uses a :latest tag; pin a version")
	}
	if !strings.Contains(string(openbaoManifests), "ghcr.io/openbao/openbao:") {
		t.Error("openbao.yaml does not reference a pinned openbao image")
	}
}

// Test_manifestsBindAllInterfaces guards the dev-mode listener address. Dev
// mode binds 127.0.0.1 by default, which would leave the Service routing to a
// port nothing listens on for out-of-pod callers — including 2.1's operator.
func Test_manifestsBindAllInterfaces(t *testing.T) {
	if !strings.Contains(string(openbaoManifests), "-dev-listen-address=0.0.0.0:8200") {
		t.Error("openbao.yaml does not set -dev-listen-address=0.0.0.0:8200")
	}
}
