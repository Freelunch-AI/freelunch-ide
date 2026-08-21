package auth

import (
	"context"
	"os/exec"
)

// execCommand is a seam over exec.CommandContext, kept in its own file so the
// service reads as orchestration rather than process plumbing.
func execCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}
