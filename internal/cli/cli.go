// Package cli assembles the freelunch command tree.
//
// Commands are wired here rather than in package main so that the whole surface
// is exercisable from tests with an explicit argv, stdout, and stderr.
package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
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

// NewRootCommand builds the freelunch command tree writing to out and errOut.
func NewRootCommand(out, errOut io.Writer) *cobra.Command {
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

	root.SetOut(out)
	root.SetErr(errOut)
	root.CompletionOptions.HiddenDefaultCmd = true

	root.AddCommand(newVersionCommand())

	return root
}

// Execute runs the freelunch command tree with the given arguments.
func Execute(ctx context.Context, args []string, out, errOut io.Writer) error {
	root := NewRootCommand(out, errOut)
	root.SetArgs(args)
	return root.ExecuteContext(ctx)
}
