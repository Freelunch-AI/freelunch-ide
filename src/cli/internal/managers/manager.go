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

	// ClusterStatus describes the local Demo cluster at a point in time.
	//
	// It lives here rather than in the cluster package because the interface
	// below names it, and managers cannot import a concrete service package
	// without creating an import cycle.
	ClusterStatus struct {
		// Name of the cluster as k3d knows it.
		Name string
		// Running reports whether the cluster exists and its nodes are up.
		Running bool
		// Nodes are the node names, empty when the cluster is not running.
		Nodes []string
	}

	// ClusterService owns the local Kubernetes cluster described by roadmap
	// 1.2. It orchestrates the pinned k3d and kubectl binaries rather than
	// reimplementing cluster management: the declarative config is the source
	// of truth, and this service only drives it.
	ClusterService interface {
		GenericService
		WithServiceManager(sm ServiceManager) ClusterService
		ServiceManager() ServiceManager

		// Create brings the cluster up from the embedded configuration. It is
		// an error if the cluster already exists — 1.2 specifies "fresh start
		// only", so callers delete and recreate rather than reconciling.
		Create(ctx context.Context) error

		// Delete tears the cluster down. Deleting a cluster that does not
		// exist succeeds, so teardown is safe to run unconditionally.
		Delete(ctx context.Context) error

		// Status reports whether the cluster is running and which nodes it has.
		Status(ctx context.Context) (*ClusterStatus, error)
	}

	// ServiceManager is the container every service holds a reference to, and
	// the only route from one service to another. Registering an implementation
	// returns the manager, so wiring reads as a single chain in main.
	ServiceManager interface {
		GenericService
		WithLogsService(ls LogsService) ServiceManager
		LogsService() LogsService
		WithClusterService(cs ClusterService) ServiceManager
		ClusterService() ClusterService
		WithCommandService(cs CommandService) ServiceManager
		CommandService() CommandService
	}

	serviceManagerFinal struct {
		logsService    LogsService
		clusterService ClusterService
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
		clusterService: NewNoOpsClusterService(),
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

	if err := m.clusterService.Start(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
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

	if err := m.clusterService.Close(ctx); err != nil {
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

	if err := m.clusterService.Healthy(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
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

func (m *serviceManagerFinal) WithClusterService(cs ClusterService) ServiceManager {
	m.clusterService = cs.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) ClusterService() ClusterService {
	return m.clusterService
}

func (m *serviceManagerFinal) WithCommandService(cs CommandService) ServiceManager {
	m.commandService = cs.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) CommandService() CommandService {
	return m.commandService
}
