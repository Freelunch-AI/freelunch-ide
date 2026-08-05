package logs

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/Freelunch-AI/freelunch-ide/internal/managers"
)

// NewManagerForTests builds a container for this package's tests. By convention
// it is duplicated per package rather than imported from another one.
func NewManagerForTests() (managers.ServiceManager, context.Context) {
	return managers.NewManager(), context.Background()
}

// newCapturingService returns a service writing to buf at the given level, so
// the emit methods can be asserted on rather than merely called.
func newCapturingService(buf *bytes.Buffer, lvl slog.Level) *logsServiceFinal {
	level := &slog.LevelVar{}
	level.Set(lvl)
	return &logsServiceFinal{
		level:  level,
		logger: slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: level})),
	}
}

func TestNewLogsService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogsService(); got == nil {
				t.Error("NewLogsService() = nil")
			}
		})
	}
}

func Test_logsServiceFinal_Start(t *testing.T) {
	sm, ctx := NewManagerForTests()
	s := sm.WithLogsService(NewLogsService()).LogsService()
	tests := []struct {
		name    string
		s       *logsServiceFinal
		wantErr bool
	}{
		{name: "success", s: s.(*logsServiceFinal), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Start(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test_logsServiceFinal_StartReadsLevelEnv covers the only configuration Start
// performs.
func Test_logsServiceFinal_StartReadsLevelEnv(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		want  slog.Level
		unset bool
	}{
		{name: "debug", env: "debug", want: slog.LevelDebug},
		{name: "info", env: "INFO", want: slog.LevelInfo},
		{name: "error", env: "error", want: slog.LevelError},
		{name: "garbage falls back", env: "not-a-level", want: DefaultLevel},
		{name: "unset falls back", unset: true, want: DefaultLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.unset {
				t.Setenv(LevelEnvVar, tt.env)
			}
			sm, ctx := NewManagerForTests()
			s := sm.WithLogsService(NewLogsService()).LogsService().(*logsServiceFinal)

			if err := s.Start(ctx); err != nil {
				t.Fatalf("Start() error = %v", err)
			}
			if got := s.level.Level(); got != tt.want {
				t.Errorf("level = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logsServiceFinal_Close(t *testing.T) {
	sm, ctx := NewManagerForTests()
	s := sm.WithLogsService(NewLogsService()).LogsService()
	tests := []struct {
		name    string
		s       *logsServiceFinal
		wantErr bool
	}{
		{name: "success", s: s.(*logsServiceFinal), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Close(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logsServiceFinal_Healthy(t *testing.T) {
	sm, ctx := NewManagerForTests()
	s := sm.WithLogsService(NewLogsService()).LogsService()
	tests := []struct {
		name    string
		s       *logsServiceFinal
		wantErr bool
	}{
		{name: "success", s: s.(*logsServiceFinal), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Healthy(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Healthy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logsServiceFinal_WithServiceManager(t *testing.T) {
	sm, _ := NewManagerForTests()
	s := NewLogsService()
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

func Test_logsServiceFinal_emit(t *testing.T) {
	tests := []struct {
		name    string
		level   slog.Level
		emit    func(s *logsServiceFinal, ctx context.Context, msg string)
		msg     string
		wantOut bool
	}{
		{
			name:  "info is emitted at info level",
			level: slog.LevelInfo,
			emit:  (*logsServiceFinal).Info,
			msg:   "an info line", wantOut: true,
		},
		{
			name:  "warn is emitted at default level",
			level: DefaultLevel,
			emit:  (*logsServiceFinal).Warn,
			msg:   "a warn line", wantOut: true,
		},
		{
			name:  "error is emitted at default level",
			level: DefaultLevel,
			emit:  (*logsServiceFinal).Error,
			msg:   "an error line", wantOut: true,
		},
		{
			name:  "debug is emitted at debug level",
			level: slog.LevelDebug,
			emit:  (*logsServiceFinal).Debug,
			msg:   "a debug line", wantOut: true,
		},
		{
			// The reason DefaultLevel is warn: a plain invocation stays quiet.
			name:  "info is suppressed at default level",
			level: DefaultLevel,
			emit:  (*logsServiceFinal).Info,
			msg:   "should not appear", wantOut: false,
		},
		{
			name:  "debug is suppressed at default level",
			level: DefaultLevel,
			emit:  (*logsServiceFinal).Debug,
			msg:   "should not appear either", wantOut: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			s := newCapturingService(&buf, tt.level)

			tt.emit(s, context.Background(), tt.msg)

			if got := strings.Contains(buf.String(), tt.msg); got != tt.wantOut {
				t.Errorf("output contains %q = %v, want %v (got %q)", tt.msg, got, tt.wantOut, buf.String())
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want slog.Level
	}{
		{name: "debug", arg: "debug", want: slog.LevelDebug},
		{name: "info", arg: "info", want: slog.LevelInfo},
		{name: "warn", arg: "warn", want: slog.LevelWarn},
		{name: "warning alias", arg: "warning", want: slog.LevelWarn},
		{name: "error", arg: "error", want: slog.LevelError},
		{name: "uppercase", arg: "ERROR", want: slog.LevelError},
		{name: "surrounding space", arg: "  info  ", want: slog.LevelInfo},
		{name: "empty", arg: "", want: DefaultLevel},
		{name: "unknown", arg: "verbose", want: DefaultLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLevel(tt.arg); got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}
