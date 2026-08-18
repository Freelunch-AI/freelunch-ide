// Command freelunch is the FreeLunch platform CLI.
//
// Every capability is a service registered on the ServiceManager (see
// internal/managers). A run loads the container, starts its services, executes
// one command, tears down, and exits — there is no daemon. See docs/roadmap.md
// section 7.3 for the command surface planned for the Demo.
package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/command"
	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/logs"
	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

func main() {
	// Exit codes are returned by run so its cleanup happens first; os.Exit
	// would skip it.
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sm := managers.NewManager().
		WithLogsService(logs.NewLogsService()).
		WithCommandService(command.NewCommandService())

	if err := sm.Start(ctx); err != nil {
		// Start has already reported the failure through LogsService.
		return 1
	}

	err := sm.CommandService().Execute(ctx, os.Args[1:])

	// Teardown runs whether or not the command succeeded, but must not mask the
	// command's own exit code.
	if closeErr := sm.Close(ctx); closeErr != nil && err == nil {
		return 1
	}

	if err == nil {
		return 0
	}

	// Cobra has already reported the error to stderr.
	var exit command.ExitError
	if errors.As(err, &exit) {
		return exit.Code
	}
	return 1
}
