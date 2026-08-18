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
			"Setup commands scaffold a customer monorepo and install the platform;\n" +
			"inspection commands report the state of workloads, environments, and\n" +
			"pipelines. Every change to a running system goes through GitOps, not\n" +
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
