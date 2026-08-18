// Package managers defines the service interfaces that make up the freelunch
// CLI, and the ServiceManager that owns them.
//
// The ServiceManager is an inversion-of-control container: no service holds a
// reference to another service. A service that needs a peer reaches it through
// its own manager at call time — s.sm.LogsService() — so an implementation can
// be swapped without its consumers changing, and a unit test can register only
// the collaborators it actually cares about.
//
// Unlike a long-running server, the CLI runs in three phases: Start wires and
// initialises every service and returns, CommandService.Execute does the work,
// and Close tears everything down. Start therefore MUST stay cheap — no network
// calls, no file I/O, nothing that would be paid again on `freelunch --help`.
// A service needing an expensive resource acquires it lazily on first real use.
package managers

import "context"

type (
	// GenericService is the lifecycle contract every service shares.
	GenericService interface {
		Start(ctx context.Context) error
		Close(ctx context.Context) error
		Healthy(ctx context.Context) error
	}

	// LogsService emits diagnostics. Every service reaches it through the
	// manager, so a caller that registers no logger still logs safely into the
	// no-op.
	LogsService interface {
		GenericService
		WithServiceManager(sm ServiceManager) LogsService
		ServiceManager() ServiceManager
		Info(ctx context.Context, s string)
		Warn(ctx context.Context, s string)
		Error(ctx context.Context, s string)
		Debug(ctx context.Context, s string)
	}

	// CommandService owns the command tree. The cobra dependency stops at its
	// implementation: no *cobra.Command appears in this package, and no other
	// service ever sees one.
	CommandService interface {
		GenericService
		WithServiceManager(sm ServiceManager) CommandService
		ServiceManager() ServiceManager

		// Execute parses args and runs the selected command, returning that
		// command's error. main derives the process exit code from it.
		Execute(ctx context.Context, args []string) error
	}

	// ServiceManager is the container every service holds a reference to, and
	// the only route from one service to another. Registering an implementation
	// returns the manager, so wiring reads as a single chain in main.
	ServiceManager interface {
		GenericService
		WithLogsService(ls LogsService) ServiceManager
		LogsService() LogsService
		WithCommandService(cs CommandService) ServiceManager
		CommandService() CommandService
	}

	serviceManagerFinal struct {
		logsService    LogsService
		commandService CommandService
	}
)

// NewManager returns a manager with every slot filled by a no-op. Registering a
// real implementation replaces one; whatever is left alone stays safely inert,
// so a caller — a test especially — only wires what it needs and calls into the
// rest without nil checks.
func NewManager() ServiceManager {
	return &serviceManagerFinal{
		logsService:    NewNoOpsLogsService(),
		commandService: NewNoOpsCommandService(),
	}
}

// Start brings services up in dependency order. LogsService is first so every
// later failure can be reported through it. Start does not block: the CLI's
// work happens in CommandService.Execute, after this returns.
func (m *serviceManagerFinal) Start(ctx context.Context) error {
	if err := m.logsService.Start(ctx); err != nil {
		return err
	}

	if err := m.commandService.Start(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	return nil
}

// Close tears services down in reverse order. LogsService is last so everything
// else can still report while shutting down.
func (m *serviceManagerFinal) Close(ctx context.Context) error {
	if err := m.commandService.Close(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	if err := m.logsService.Close(ctx); err != nil {
		return err
	}

	return nil
}

// Healthy reports whether every service is usable. The CLI does not call this
// during startup — a per-invocation health check would just be latency. It
// exists for an explicit diagnostic command to drive.
func (m *serviceManagerFinal) Healthy(ctx context.Context) error {
	if err := m.logsService.Healthy(ctx); err != nil {
		return err
	}

	if err := m.commandService.Healthy(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (m *serviceManagerFinal) WithLogsService(ls LogsService) ServiceManager {
	m.logsService = ls.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) LogsService() LogsService {
	return m.logsService
}

func (m *serviceManagerFinal) WithCommandService(cs CommandService) ServiceManager {
	m.commandService = cs.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) CommandService() CommandService {
	return m.commandService
}
