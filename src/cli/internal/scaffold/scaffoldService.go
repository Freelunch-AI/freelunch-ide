// Package scaffold provides the real ScaffoldService: customer-monorepo
// bootstrapping from roadmap 1.1.
//
// The canonical template is embedded from the package-local template/ copy —
// a released binary carries the exact structure it was tested with, and a
// bare binary can init with nothing else installed. The repository-level
// templates/monorepo (published as the GitHub template and shipped in release
// archives) is the same content; a test fails if the two drift.
//
// git is the one external dependency, and deliberately from PATH rather than
// the pinned toolchain: the story ends "ready to commit and push to GitHub",
// which presumes the customer's own git, and pinning a vcs binary would buy
// nothing.
package scaffold

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/buildinfo"
	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// Placeholder is the product directory name inside the template. The repo
// convention (never <angle> tokens — they are illegal in Windows filenames)
// makes it a normal, buildable name; Init swaps it for the real product name.
const Placeholder = "example_product"

// versionLine is the line Init stamps in platform/freelunch.yaml. A test
// asserts the template actually contains it, because a silent mismatch would
// leave customers with a monorepo claiming platform version 0.0.0.
const versionLine = "version: 0.0.0"

// The all: prefix is load-bearing: without it, embed skips dotfiles, and the
// .gitkeep files that hold the empty directories open would vanish from the
// scaffolded repo.
//
//go:embed all:template
var templateFS embed.FS

type (
	// Runner executes an external command and returns its combined output.
	// It exists so tests can drive the service without a git installation.
	Runner interface {
		Run(ctx context.Context, name string, args ...string) ([]byte, error)
	}

	execRunner struct{}

	scaffoldServiceFinal struct {
		sm     managers.ServiceManager
		runner Runner
	}
)

func (execRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// NewScaffoldService builds the service against the real system.
func NewScaffoldService() managers.ScaffoldService {
	return &scaffoldServiceFinal{runner: execRunner{}}
}

// NewScaffoldServiceWithRunner builds the service against a caller-supplied
// Runner. It is the seam tests use to exercise Init without git.
func NewScaffoldServiceWithRunner(r Runner) managers.ScaffoldService {
	return &scaffoldServiceFinal{runner: r}
}

// Start does no I/O; the template is embedded and needs no resolution.
func (s *scaffoldServiceFinal) Start(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "starting the scaffold service")
	return nil
}

func (s *scaffoldServiceFinal) Close(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "stopping the scaffold service")
	return nil
}

func (s *scaffoldServiceFinal) Healthy(_ context.Context) error {
	return nil
}

func (s *scaffoldServiceFinal) WithServiceManager(sm managers.ServiceManager) managers.ScaffoldService {
	s.sm = sm
	return s
}

func (s *scaffoldServiceFinal) ServiceManager() managers.ServiceManager {
	return s.sm
}

// validName rejects product names that would escape the products/ directory
// or fail on another contributor's OS. The rules are deliberately tight —
// a directory name, not a display name.
func validName(name string) error {
	if name == "" {
		return fmt.Errorf("product name is empty")
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) || strings.Contains(name, "..") {
		return fmt.Errorf("product name %q contains characters that are not portable across filesystems", name)
	}
	return nil
}

// Init creates dir with the canonical monorepo structure. See the interface
// documentation in managers for the contract; one addition worth knowing: on
// any failure the created directory is removed, so a failed init never leaves
// a half-scaffolded repository for the next attempt to trip over.
func (s *scaffoldServiceFinal) Init(ctx context.Context, dir, product string) error {
	if err := validName(product); err != nil {
		return err
	}

	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s already exists; init only creates new repositories", dir)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot inspect %s: %w", dir, err)
	}

	s.sm.LogsService().Debug(ctx, "scaffolding monorepo at "+dir)

	// Everything below is undone on failure.
	fail := func(err error) error {
		_ = os.RemoveAll(dir)
		return err
	}

	err := fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel := strings.TrimPrefix(path, "template")
		rel = strings.TrimPrefix(rel, "/")
		// The placeholder is a path segment, never a substring of file
		// content — the template holds no templated text by design.
		rel = strings.ReplaceAll(rel, Placeholder, product)
		target := filepath.Join(dir, filepath.FromSlash(rel))

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, readErr := templateFS.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		// Stamp the platform version the moment it passes through, so the
		// scaffolded repo records which release created it. buildinfo is a
		// plain leaf package, not a service — importing it directly is the
		// framework-sanctioned exception.
		if rel == filepath.ToSlash(filepath.Join("platform", "freelunch.yaml")) {
			data = bytes.Replace(data,
				[]byte(versionLine),
				[]byte("version: "+buildinfo.Get().Version), 1)
		}

		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		return fail(fmt.Errorf("cannot write monorepo structure: %w", err))
	}

	// `git init <dir>` takes the target as an argument, so no working-directory
	// plumbing is needed. The story requires a repository, not just a tree —
	// a missing git fails the init rather than quietly skipping the step.
	out, err := s.runner.Run(ctx, "git", "init", dir)
	if err != nil {
		return fail(fmt.Errorf(
			"monorepo written but `git init` failed (is git installed?): %w: %s",
			err, strings.TrimSpace(string(out))))
	}

	return nil
}
