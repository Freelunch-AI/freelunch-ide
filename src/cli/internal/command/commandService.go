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
			"install creates the local Demo cluster on Docker, uninstall tears it down,\n" +
			"and status reports what is running. Every change to a running system goes\n" +
			"through GitOps, not through this CLI.",
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

// newInstallCommand builds `freelunch install`.
//
// The name follows roadmap.md:370 and the 1.2 story, which both use install for
// provisioning. It is deliberately not `start`: k3d distinguishes creating a
// cluster from starting one that already exists, and reserving `start` for the
// latter keeps both words meaning what they ordinarily mean.
func (s *commandServiceFinal) newInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Create the local Demo cluster",
		Long: "install creates the local Demo Kubernetes cluster from the pinned k3d\n" +
			"configuration that ships inside this binary.\n\n" +
			"Docker must be running and the pinned tools must already be installed.\n" +
			"Creation is fresh-start only: run uninstall first to recreate a cluster\n" +
			"that already exists.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return s.RunInstall(cmd.Context(), cmd.OutOrStdout())
		},
	}
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

// RunInstall is the cobra-free body of `freelunch install`.
func (s *commandServiceFinal) RunInstall(ctx context.Context, out io.Writer) error {
	s.sm.LogsService().Debug(ctx, "running the install command")

	// Creation takes roughly a minute, so say something before blocking. This is
	// program output rather than a diagnostic: it is the command reporting on the
	// work it was asked to do.
	if _, err := fmt.Fprintln(out, "Creating the local Demo cluster; this takes about a minute."); err != nil {
		return err
	}

	if err := s.sm.ClusterService().Create(ctx); err != nil {
		return err
	}

	_, err := fmt.Fprintln(out, "Cluster created. Run `freelunch status` to inspect it.")
	return err
}

// RunUninstall is the cobra-free body of `freelunch uninstall`.
func (s *commandServiceFinal) RunUninstall(ctx context.Context, out io.Writer) error {
	s.sm.LogsService().Debug(ctx, "running the uninstall command")

	if err := s.sm.ClusterService().Delete(ctx); err != nil {
		return err
	}

	_, err := fmt.Fprintln(out, "Cluster deleted.")
	return err
}

// RunStatus is the cobra-free body of `freelunch status`.
func (s *commandServiceFinal) RunStatus(ctx context.Context, out io.Writer) error {
	s.sm.LogsService().Debug(ctx, "running the status command")

	status, err := s.sm.ClusterService().Status(ctx)
	if err != nil {
		return err
	}

	// The interface permits a nil status, so treat it the same as a cluster that
	// is not there rather than dereferencing it.
	if status == nil || !status.Running {
		_, err = fmt.Fprintln(out, "Cluster is not running. Run `freelunch install` to create it.")
		return err
	}

	if _, err = fmt.Fprintf(out, "Cluster %q is running, with %d node(s):\n",
		status.Name, len(status.Nodes)); err != nil {
		return err
	}

	for _, node := range status.Nodes {
		if _, err = fmt.Fprintf(out, "  %s\n", node); err != nil {
			return err
		}
	}

	return nil
}
