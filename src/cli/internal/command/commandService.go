// Package command provides the real CommandService: the freelunch command tree.
//
// cobra is an implementation detail of this package. The managers.CommandService
// interface exposes only Execute(ctx, args), and every command body is a plain
// method taking a context and already-parsed values, so the behaviour of a
// command can be tested without constructing a *cobra.Command.
package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/buildinfo"
	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// ExitError lets a command choose the process exit code. main unwraps it; any
// other error exits 1.
type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit status %d", e.Code)
	}
	return e.Err.Error()
}

func (e ExitError) Unwrap() error { return e.Err }

type (
	commandServiceFinal struct {
		sm managers.ServiceManager
		// Program output and command diagnostics. Held as fields rather than
		// taken by the constructor so tests can redirect them; production always
		// uses the process streams.
		out    io.Writer
		errOut io.Writer
	}
)

// NewCommandService builds the service writing to the process streams.
func NewCommandService() managers.CommandService {
	return &commandServiceFinal{
		out:    os.Stdout,
		errOut: os.Stderr,
	}
}

// Start does no work yet. The command tree is built per Execute call, which
// keeps invocations independent and costs nothing to skip.
func (s *commandServiceFinal) Start(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "starting the command service")
	return nil
}

func (s *commandServiceFinal) Close(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "stopping the command service")
	return nil
}

func (s *commandServiceFinal) Healthy(_ context.Context) error {
	return nil
}

func (s *commandServiceFinal) WithServiceManager(sm managers.ServiceManager) managers.CommandService {
	s.sm = sm
	return s
}

func (s *commandServiceFinal) ServiceManager() managers.ServiceManager {
	return s.sm
}

func (s *commandServiceFinal) Execute(ctx context.Context, args []string) error {
	root := s.newRootCommand()
	root.SetArgs(args)
	return root.ExecuteContext(ctx)
}

// newRootCommand assembles the command tree. Unexported on purpose: a
// *cobra.Command must never leave this package.
func (s *commandServiceFinal) newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "freelunch",
		Short: "FreeLunch platform CLI",
		Long: "freelunch bootstraps and inspects a FreeLunch platform installation.\n\n" +
			"init creates a customer monorepo with the canonical FreeLunch structure;\n" +
			"install creates the local Demo environment on Docker (cluster, auth\n" +
			"service, secrets store), uninstall tears it down, and status reports what\n" +
			"is running. Every change to a running system goes through GitOps, not\n" +
			"through this CLI.",
		SilenceUsage:  true,
		SilenceErrors: false,
		// Print help instead of an unhelpful "unknown command" for a bare invocation.
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	root.SetOut(s.out)
	root.SetErr(s.errOut)
	root.CompletionOptions.HiddenDefaultCmd = true

	root.AddCommand(s.newInitCommand())
	root.AddCommand(s.newInstallCommand())
	root.AddCommand(s.newUninstallCommand())
	root.AddCommand(s.newStatusCommand())
	root.AddCommand(s.newVersionCommand())

	return root
}

func (s *commandServiceFinal) newVersionCommand() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the freelunch version and build metadata",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return s.RunVersion(cmd.Context(), cmd.OutOrStdout(), asJSON)
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "emit build metadata as JSON")

	return cmd
}

// RunVersion is the cobra-free body of `freelunch version`.
func (s *commandServiceFinal) RunVersion(ctx context.Context, out io.Writer, asJSON bool) error {
	s.sm.LogsService().Debug(ctx, "running the version command")

	info := buildinfo.Get()

	if asJSON {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	_, err := fmt.Fprintf(out,
		"freelunch %s\n  commit:   %s\n  built:    %s\n  go:       %s\n  platform: %s\n",
		info.Version, info.Commit, info.Date, info.GoVersion, info.Platform)
	return err
}

// newInitCommand builds `freelunch init`.
func (s *commandServiceFinal) newInitCommand() *cobra.Command {
	var product string

	cmd := &cobra.Command{
		Use:   "init <directory>",
		Short: "Create a customer monorepo with the canonical FreeLunch structure",
		Long: "init creates a new directory with the canonical customer-monorepo\n" +
			"structure — platform/, products/<product>/services/ and workflows/, and\n" +
			".github/workflows/ — stamps the platform version, and runs `git init` so\n" +
			"the result is ready to commit and push.\n\n" +
			"The products directory starts as the example_product placeholder; pass\n" +
			"--product to name it now, or rename it when the first real product\n" +
			"exists. git must be installed.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.RunInit(cmd.Context(), cmd.OutOrStdout(), args[0], product)
		},
	}

	cmd.Flags().StringVar(&product, "product", "example_product",
		"name for the products/<product> directory")

	return cmd
}

// RunInit is the cobra-free body of `freelunch init`.
func (s *commandServiceFinal) RunInit(ctx context.Context, out io.Writer, dir, product string) error {
	s.sm.LogsService().Debug(ctx, "running the init command")

	if err := s.sm.ScaffoldService().Init(ctx, dir, product); err != nil {
		return err
	}

	_, err := fmt.Fprintf(out,
		"Monorepo created at %s (product: %s).\n"+
			"Next: commit and push it, then run `freelunch install` for the local environment.\n",
		dir, product)
	return err
}

// newInstallCommand builds `freelunch install`.
//
// The name follows roadmap.md:370 and the 1.2 story, which both use install for
// provisioning. It is deliberately not `start`: k3d distinguishes creating a
// cluster from starting one that already exists, and reserving `start` for the
// latter keeps both words meaning what they ordinarily mean.
func (s *commandServiceFinal) newInstallCommand() *cobra.Command {
	var only []string
	var skip []string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Create the local Demo environment: cluster, auth service, secrets store",
		Long: "install creates the local Demo Kubernetes cluster from the pinned k3d\n" +
			"configuration that ships inside this binary, then installs the platform\n" +
			"components into it — the auth service (roadmap 1.3) and the secrets\n" +
			"store (roadmap 1.4).\n\n" +
			"Docker must be running and the pinned tools must already be installed.\n" +
			"Cluster creation is fresh-start only: run uninstall first to recreate a\n" +
			"cluster that already exists. Components re-apply cleanly.\n\n" +
			"Use --only to install a subset, or --skip to leave components out:\n" +
			"  freelunch install --only cluster\n" +
			"  freelunch install --skip auth,secrets",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return s.RunInstall(cmd.Context(), cmd.OutOrStdout(), only, skip)
		},
	}

	cmd.Flags().StringSliceVar(&only, "only", nil, "install only these components (cluster, auth, secrets)")
	cmd.Flags().StringSliceVar(&skip, "skip", nil, "skip these components (cluster, auth, secrets)")

	return cmd
}

// newUninstallCommand builds `freelunch uninstall`.
func (s *commandServiceFinal) newUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Delete the local Demo cluster",
		Long: "uninstall tears the local Demo cluster down, removing its containers,\n" +
			"volumes and kubeconfig entry.\n\n" +
			"Deleting a cluster that is not there succeeds, so this is safe to run\n" +
			"unconditionally.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return s.RunUninstall(cmd.Context(), cmd.OutOrStdout())
		},
	}
}

// newStatusCommand builds `freelunch status`.
func (s *commandServiceFinal) newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Report whether the local Demo cluster is running",
		Long: "status reports whether the local Demo cluster exists and which nodes it\n" +
			"has.\n\n" +
			"A cluster that is absent is a normal answer rather than a failure, so this\n" +
			"exits zero either way; only being unable to ask is an error.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return s.RunStatus(cmd.Context(), cmd.OutOrStdout())
		},
	}
}

// installComponents is the ordered component list install works through.
// Order matters: auth deploys into the cluster, so the cluster goes first.
var installComponents = []string{"cluster", "auth", "secrets"}

// selectComponents resolves --only/--skip into the subset to act on, in
// canonical order. Unknown names are an error rather than silently ignored —
// a typo in --skip that installs the very thing the user excluded is worse
// than a failed command.
func selectComponents(only, skip []string) ([]string, error) {
	known := map[string]bool{}
	for _, c := range installComponents {
		known[c] = true
	}
	for _, name := range append(append([]string{}, only...), skip...) {
		if !known[name] {
			return nil, fmt.Errorf("unknown component %q (known: %s)",
				name, strings.Join(installComponents, ", "))
		}
	}

	selected := map[string]bool{}
	if len(only) > 0 {
		for _, name := range only {
			selected[name] = true
		}
	} else {
		for _, c := range installComponents {
			selected[c] = true
		}
	}
	for _, name := range skip {
		delete(selected, name)
	}

	var result []string
	for _, c := range installComponents {
		if selected[c] {
			result = append(result, c)
		}
	}
	return result, nil
}

// RunInstall is the cobra-free body of `freelunch install`.
func (s *commandServiceFinal) RunInstall(ctx context.Context, out io.Writer, only, skip []string) error {
	s.sm.LogsService().Debug(ctx, "running the install command")

	components, err := selectComponents(only, skip)
	if err != nil {
		return err
	}

	for _, component := range components {
		switch component {
		case "cluster":
			// Creation takes roughly a minute, so say something before
			// blocking. Program output, not a diagnostic: the command is
			// reporting on the work it was asked to do.
			if _, err := fmt.Fprintln(out, "Creating the local Demo cluster; this takes about a minute."); err != nil {
				return err
			}
			if err := s.sm.ClusterService().Create(ctx); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(out, "Cluster created."); err != nil {
				return err
			}
		case "auth":
			if _, err := fmt.Fprintln(out, "Installing the auth service."); err != nil {
				return err
			}
			if err := s.sm.AuthService().Install(ctx); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(out, "Auth service installed; it takes a moment to come up."); err != nil {
				return err
			}
		case "secrets":
			if _, err := fmt.Fprintln(out, "Installing the secrets store."); err != nil {
				return err
			}
			if err := s.sm.SecretsService().Install(ctx); err != nil {
				return err
			}
			// Seed the demo credential from the 1.4 story (the spec names the
			// workload "my-service"; the seed follows the repo's example_*
			// placeholder convention). Dev mode loses contents on restart, so
			// this runs on every install rather than once — which is also what
			// makes re-install after a pod crash produce a working demo.
			if err := s.sm.SecretsService().PutSecret(ctx,
				"example_service", "api-key", "example-api-key-value"); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(out, "Secrets store installed and seeded."); err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintln(out, "Done. Run `freelunch status` to inspect the environment.")
	return err
}

// RunUninstall is the cobra-free body of `freelunch uninstall`.
func (s *commandServiceFinal) RunUninstall(ctx context.Context, out io.Writer) error {
	s.sm.LogsService().Debug(ctx, "running the uninstall command")

	// Reverse of install order. Deleting components before the cluster is
	// mostly symbolic — deleting the cluster removes everything in it — but it
	// keeps the teardown meaningful when the cluster itself is kept in future
	// variants, and every delete tolerates absence.
	if err := s.sm.SecretsService().Delete(ctx); err != nil {
		return err
	}

	if err := s.sm.AuthService().Delete(ctx); err != nil {
		return err
	}

	if err := s.sm.ClusterService().Delete(ctx); err != nil {
		return err
	}

	_, err := fmt.Fprintln(out, "Demo environment deleted.")
	return err
}

// RunStatus is the cobra-free body of `freelunch status`.
func (s *commandServiceFinal) RunStatus(ctx context.Context, out io.Writer) error {
	s.sm.LogsService().Debug(ctx, "running the status command")

	status, err := s.sm.ClusterService().Status(ctx)
	if err != nil {
		return err
	}

	// Absent, API not answering, nodes NotReady and running are four different
	// situations with four different remedies, so they get four different
	// messages. The interface permits a nil status; treat it as absent rather
	// than dereferencing it.
	switch {
	case status == nil || !status.Exists:
		_, err = fmt.Fprintln(out, "Cluster is not running. Run `freelunch install` to create it.")
		return err
	case len(status.Nodes) == 0:
		_, err = fmt.Fprintf(out, "Cluster %q exists but its API is not answering yet; check again shortly.\n", status.Name)
		return err
	case !status.Running:
		if _, err = fmt.Fprintf(out, "Cluster %q exists but is not ready: %d of %d node(s) Ready.\n",
			status.Name, readyNodeCount(status.Nodes), len(status.Nodes)); err != nil {
			return err
		}
		if err := writeNodes(out, status.Nodes); err != nil {
			return err
		}
		// The cluster is disposable by design ("fresh start only"), so the
		// remedy is recreation, not repair.
		_, err = fmt.Fprintln(out, "Run `freelunch uninstall` then `freelunch install` to recreate it.")
		return err
	}

	if _, err = fmt.Fprintf(out, "Cluster %q is running, with %d node(s):\n",
		status.Name, len(status.Nodes)); err != nil {
		return err
	}

	if err := writeNodes(out, status.Nodes); err != nil {
		return err
	}

	authStatus, err := s.sm.AuthService().Status(ctx)
	if err != nil {
		return err
	}

	if authStatus == nil || !authStatus.Ready {
		_, err = fmt.Fprintln(out, "Auth service is not ready. It takes ~60s after install; check again shortly.")
		return err
	}

	// The issuer is printed because it is the value that breaks first when the
	// hostname configuration is wrong — visible here, it is diagnosable at a
	// glance instead of surfacing as a mysterious 401 in a consumer.
	if _, err = fmt.Fprintf(out, "Auth service is ready: realm %q, issuer %s\n",
		authStatus.Realm, authStatus.IssuerURL); err != nil {
		return err
	}

	secretsStatus, err := s.sm.SecretsService().Status(ctx)
	if err != nil {
		return err
	}

	switch {
	case secretsStatus == nil || (!secretsStatus.Ready && !secretsStatus.Sealed):
		_, err = fmt.Fprintln(out, "Secrets store is not ready.")
	case secretsStatus.Sealed:
		// Dev mode never seals, so this means the deployment was changed out
		// from under us — worth its own line, since the pod looks healthy.
		_, err = fmt.Fprintln(out, "Secrets store is SEALED — this dev-mode store should never seal; inspect the deployment.")
	default:
		// The engine is printed because the kv v1/v2 distinction is what 2.1
		// must configure for; visible here, drift is caught at a glance.
		_, err = fmt.Fprintf(out, "Secrets store is ready: %s\n", secretsStatus.Engine)
	}
	return err
}

// readyNodeCount counts the nodes whose Ready condition is True.
func readyNodeCount(nodes []managers.ClusterNode) int {
	n := 0
	for _, node := range nodes {
		if node.Ready {
			n++
		}
	}
	return n
}

// writeNodes prints one line per node with its readiness, in the wording
// `kubectl get nodes` uses, so the output reads the same as the tool people
// reach for next.
func writeNodes(out io.Writer, nodes []managers.ClusterNode) error {
	for _, node := range nodes {
		state := "Ready"
		if !node.Ready {
			state = "NotReady"
		}
		if _, err := fmt.Fprintf(out, "  %s  %s\n", node.Name, state); err != nil {
			return err
		}
	}
	return nil
}
