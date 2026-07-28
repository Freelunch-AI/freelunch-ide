// Command freelunch is the FreeLunch platform CLI.
//
// Setup and inspection commands are defined in internal/cli. See docs/roadmap.md
// section 7.3 for the command surface planned for the Demo.
package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/Freelunch-AI/freelunch-ide/internal/cli"
)

func main() {
	// Exit codes are returned by run so its deferred cleanup runs first;
	// os.Exit would skip it.
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := cli.Execute(ctx, os.Args[1:], os.Stdout, os.Stderr)
	if err == nil {
		return 0
	}

	// Cobra has already reported the error to stderr.
	var exit cli.ExitError
	if errors.As(err, &exit) {
		return exit.Code
	}
	return 1
}
