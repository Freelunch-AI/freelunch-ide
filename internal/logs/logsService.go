// Package logs provides the real LogsService: diagnostics written to stderr.
//
// Diagnostics are not program output. Anything a command produces for a user or
// a pipe goes to stdout through CommandService; this package only ever writes to
// stderr, so redirecting `freelunch ... > file` never captures log noise.
package logs

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/Freelunch-AI/freelunch-ide/internal/managers"
)

// LevelEnvVar names the environment variable that overrides the log level.
const LevelEnvVar = "FREELUNCH_LOG_LEVEL"

// DefaultLevel keeps a normal invocation quiet. A CLI runs to completion on
// every call, so logging each service start at info would print noise before
// the output of something as small as `freelunch version`.
const DefaultLevel = slog.LevelWarn

type (
	logsServiceFinal struct {
		sm     managers.ServiceManager
		logger *slog.Logger
		level  *slog.LevelVar
	}
)

// NewLogsService builds the service already able to log, so a failure occurring
// before Start still has somewhere to go.
func NewLogsService() managers.LogsService {
	level := &slog.LevelVar{}
	level.Set(DefaultLevel)

	return &logsServiceFinal{
		level:  level,
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})),
	}
}

// Start applies the configured level. Reading an environment variable is the
// only work done here — Start must stay cheap enough to pay on every invocation.
func (s *logsServiceFinal) Start(ctx context.Context) error {
	s.level.Set(ParseLevel(os.Getenv(LevelEnvVar)))
	s.Debug(ctx, "starting the logs service")
	return nil
}

func (s *logsServiceFinal) Close(ctx context.Context) error {
	s.Debug(ctx, "stopping the logs service")
	return nil
}

func (s *logsServiceFinal) Healthy(_ context.Context) error {
	return nil
}

func (s *logsServiceFinal) WithServiceManager(sm managers.ServiceManager) managers.LogsService {
	s.sm = sm
	return s
}

func (s *logsServiceFinal) ServiceManager() managers.ServiceManager {
	return s.sm
}

func (s *logsServiceFinal) Info(ctx context.Context, msg string) {
	s.logger.InfoContext(ctx, msg)
}

func (s *logsServiceFinal) Warn(ctx context.Context, msg string) {
	s.logger.WarnContext(ctx, msg)
}

func (s *logsServiceFinal) Error(ctx context.Context, msg string) {
	s.logger.ErrorContext(ctx, msg)
}

func (s *logsServiceFinal) Debug(ctx context.Context, msg string) {
	s.logger.DebugContext(ctx, msg)
}

// ParseLevel maps a level name to a slog level, falling back to DefaultLevel for
// anything unrecognised. An unreadable value must not make the CLI fail to
// start, so this never returns an error.
func ParseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return DefaultLevel
	}
}
