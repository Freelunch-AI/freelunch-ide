package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Freelunch-AI/freelunch-ide/internal/managers"
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
