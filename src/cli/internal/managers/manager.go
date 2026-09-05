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

	// ClusterNode is one node of the Demo cluster and whether it is Ready.
	ClusterNode struct {
		// Name is the node name as Kubernetes reports it.
		Name string
		// Ready reports whether the node's Ready condition is True. A node whose
		// container has died, or that is still joining, is listed with Ready
		// false rather than omitted: it exists, it is just not usable.
		Ready bool
	}

	// ClusterStatus describes the local Demo cluster at a point in time.
	//
	// Existence and readiness are reported separately. A cluster can be known
	// to k3d while its API is still coming up, or answer the API with nodes that
	// are NotReady; neither is "absent" and neither is "running". Callers must
	// not infer one from the other.
	//
	// It lives here rather than in the cluster package because the interface
	// below names it, and managers cannot import a concrete service package
	// without creating an import cycle.
	ClusterStatus struct {
		// Name of the cluster as k3d knows it.
		Name string
		// Exists reports whether k3d knows a cluster of this name.
		Exists bool
		// Running reports whether the cluster exists, its API answers, and every
		// node is Ready. It is never true while Exists is false.
		Running bool
		// Nodes lists every node the API reports, Ready or not. Empty when the
		// cluster is absent or its API is not answering.
		Nodes []ClusterNode
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

	// AuthStatus describes the local auth service at a point in time.
	//
	// It lives here rather than in the auth package because the interface
	// below names it, and managers cannot import a concrete service package
	// without creating an import cycle.
	AuthStatus struct {
		// Ready reports whether the server is up and serving OIDC discovery.
		Ready bool
		// Realm is the realm name, empty when the realm is not imported yet.
		Realm string
		// IssuerURL is the issuer the server advertises, which is the value
		// that actually breaks when hostname configuration is wrong.
		IssuerURL string
	}

	// AuthService owns the local OIDC identity provider described by roadmap
	// 1.3. It orchestrates kubectl against manifests embedded in the binary;
	// it does not talk to the Kubernetes API directly.
	AuthService interface {
		GenericService
		WithServiceManager(sm ServiceManager) AuthService
		ServiceManager() ServiceManager

		// Install applies the auth service and its realm to the running
		// cluster. It is an error if no cluster is running.
		Install(ctx context.Context) error

		// Delete removes it. Deleting what is not there succeeds.
		Delete(ctx context.Context) error

		// Status reports readiness and the advertised issuer.
		Status(ctx context.Context) (*AuthStatus, error)
	}

	// SecretsStatus describes the local secrets store at a point in time.
	//
	// It lives here rather than in the secrets package because the interface
	// below names it, and managers cannot import a concrete service package
	// without creating an import cycle.
	SecretsStatus struct {
		// Up reports whether the pod runs and the store answers `bao status` as
		// initialized. It says nothing about whether the store can serve
		// secrets; see Ready.
		Up bool
		// Sealed is reported separately because a running-but-sealed store is
		// the failure that looks most like success from the outside.
		Sealed bool
		// Engine is what is mounted at the KV path, e.g. "secret/ (kv v2)",
		// empty when nothing is. Recorded because the v1/v2 distinction is what
		// breaks 2.1.
		Engine string
		// Ready reports whether the store is usable as roadmap 1.4 promises:
		// Up, not Sealed, and a KV v2 engine mounted at secret/. Dev mode
		// mounts that engine on every start, so Up without Ready means the
		// deployment was changed out from under us — waiting will not fix it.
		Ready bool
	}

	// SecretsService owns the local secrets store described by roadmap 1.4.
	// It orchestrates kubectl against manifests embedded in the binary.
	SecretsService interface {
		GenericService
		WithServiceManager(sm ServiceManager) SecretsService
		ServiceManager() ServiceManager

		// Install applies the store to the running cluster and waits for it
		// to answer. It is an error if no cluster is running.
		Install(ctx context.Context) error

		// Delete removes it. Deleting what is not there succeeds.
		Delete(ctx context.Context) error

		// Status reports readiness, seal state and the mounted engine.
		Status(ctx context.Context) (*SecretsStatus, error)

		// PutSecret writes one key/value pair at a logical path, e.g.
		// ("example_service", "api-key", "..."). Re-runnable by design: dev mode
		// loses its contents on restart, so seeding is part of install, not a
		// one-off.
		PutSecret(ctx context.Context, path, key, value string) error
	}

	// ScaffoldService owns customer-monorepo bootstrapping (roadmap 1.1).
	// The canonical template is embedded in the binary, so a bare binary can
	// init with nothing else installed.
	ScaffoldService interface {
		GenericService
		WithServiceManager(sm ServiceManager) ScaffoldService
		ServiceManager() ServiceManager

		// Init creates dir with the canonical monorepo structure, renames the
		// example_product placeholder to product, stamps the platform version,
		// and runs `git init`. It is an error if dir already exists, and a
		// failed init removes what it created rather than leaving a partial
		// repository behind.
		Init(ctx context.Context, dir, product string) error
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
		WithAuthService(as AuthService) ServiceManager
		AuthService() AuthService
		WithSecretsService(ss SecretsService) ServiceManager
		SecretsService() SecretsService
		WithScaffoldService(sc ScaffoldService) ServiceManager
		ScaffoldService() ScaffoldService
		WithCommandService(cs CommandService) ServiceManager
		CommandService() CommandService
	}

	serviceManagerFinal struct {
		logsService     LogsService
		clusterService  ClusterService
		authService     AuthService
		secretsService  SecretsService
		scaffoldService ScaffoldService
		commandService  CommandService
	}
)

// NewManager returns a manager with every slot filled by a no-op. Registering a
// real implementation replaces one; whatever is left alone stays safely inert,
// so a caller — a test especially — only wires what it needs and calls into the
// rest without nil checks.
func NewManager() ServiceManager {
	return &serviceManagerFinal{
		logsService:     NewNoOpsLogsService(),
		clusterService:  NewNoOpsClusterService(),
		authService:     NewNoOpsAuthService(),
		secretsService:  NewNoOpsSecretsService(),
		scaffoldService: NewNoOpsScaffoldService(),
		commandService:  NewNoOpsCommandService(),
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

	if err := m.authService.Start(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	if err := m.secretsService.Start(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	if err := m.scaffoldService.Start(ctx); err != nil {
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

	if err := m.scaffoldService.Close(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	if err := m.secretsService.Close(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	if err := m.authService.Close(ctx); err != nil {
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

	if err := m.authService.Healthy(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	if err := m.secretsService.Healthy(ctx); err != nil {
		m.logsService.Error(ctx, err.Error())
		return err
	}

	if err := m.scaffoldService.Healthy(ctx); err != nil {
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

func (m *serviceManagerFinal) WithAuthService(as AuthService) ServiceManager {
	m.authService = as.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) AuthService() AuthService {
	return m.authService
}

func (m *serviceManagerFinal) WithSecretsService(ss SecretsService) ServiceManager {
	m.secretsService = ss.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) SecretsService() SecretsService {
	return m.secretsService
}

func (m *serviceManagerFinal) WithScaffoldService(sc ScaffoldService) ServiceManager {
	m.scaffoldService = sc.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) ScaffoldService() ScaffoldService {
	return m.scaffoldService
}

func (m *serviceManagerFinal) WithCommandService(cs CommandService) ServiceManager {
	m.commandService = cs.WithServiceManager(m)
	return m
}

func (m *serviceManagerFinal) CommandService() CommandService {
	return m.commandService
}
