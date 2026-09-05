package scaffold

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/buildinfo"
	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// NewManagerForTests builds a container for this package's tests. By convention
// it is duplicated per package rather than imported from another one.
func NewManagerForTests() (managers.ServiceManager, context.Context) {
	return managers.NewManager(), context.Background()
}

type call struct {
	name string
	args []string
}

type fakeRunner struct {
	calls   []call
	failAll error
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, call{name: name, args: args})
	if f.failAll != nil {
		return []byte("boom"), f.failAll
	}
	return []byte("ok"), nil
}

func newServiceForTest(t *testing.T) (*scaffoldServiceFinal, *fakeRunner, context.Context) {
	t.Helper()

	sm, ctx := NewManagerForTests()
	runner := &fakeRunner{}
	sm.WithScaffoldService(NewScaffoldServiceWithRunner(runner))

	svc, ok := sm.ScaffoldService().(*scaffoldServiceFinal)
	if !ok {
		t.Fatalf("ScaffoldService() is %T, want *scaffoldServiceFinal", sm.ScaffoldService())
	}

	return svc, runner, ctx
}

func TestNewScaffoldService(t *testing.T) {
	if got := NewScaffoldService(); got == nil {
		t.Fatal("NewScaffoldService() = nil")
	}
}

func Test_scaffoldServiceFinal_Start(t *testing.T) {
	svc, runner, ctx := newServiceForTest(t)

	if err := svc.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	if len(runner.calls) != 0 {
		t.Errorf("Start() ran %d external commands, want 0", len(runner.calls))
	}
}

func Test_scaffoldServiceFinal_Close(t *testing.T) {
	svc, _, ctx := newServiceForTest(t)

	if err := svc.Close(ctx); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}

func Test_scaffoldServiceFinal_Healthy(t *testing.T) {
	svc, _, ctx := newServiceForTest(t)

	if err := svc.Healthy(ctx); err != nil {
		t.Fatalf("Healthy() error = %v, want nil", err)
	}
}

func Test_scaffoldServiceFinal_WithServiceManager(t *testing.T) {
	sm, _ := NewManagerForTests()
	s := NewScaffoldService()

	got := s.WithServiceManager(sm)
	if got.ServiceManager() != sm {
		t.Errorf("ServiceManager() = %v, want %v", got.ServiceManager(), sm)
	}
}

func Test_scaffoldServiceFinal_Init(t *testing.T) {
	tests := []struct {
		name    string
		product string
		wantErr bool
		// wantDirs are checked relative to the created repo.
		wantDirs []string
	}{
		{
			name:    "default placeholder product",
			product: "example_product",
			wantDirs: []string{
				"platform",
				".github/workflows",
				"products/example_product/services",
				"products/example_product/workflows",
			},
		},
		{
			name:    "renamed product",
			product: "shop",
			wantDirs: []string{
				"products/shop/services",
				"products/shop/workflows",
			},
		},
		{name: "empty product name", product: "", wantErr: true},
		{name: "path separator in product", product: "a/b", wantErr: true},
		{name: "dotdot in product", product: "..evil", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, runner, ctx := newServiceForTest(t)
			dir := filepath.Join(t.TempDir(), "my-company")

			err := svc.Init(ctx, dir, tt.product)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
					t.Errorf("failed Init() left %s behind", dir)
				}
				return
			}

			for _, want := range tt.wantDirs {
				if _, statErr := os.Stat(filepath.Join(dir, want)); statErr != nil {
					t.Errorf("missing %s: %v", want, statErr)
				}
			}

			// The placeholder must be gone when renamed.
			if tt.product != Placeholder {
				if _, statErr := os.Stat(filepath.Join(dir, "products", Placeholder)); !os.IsNotExist(statErr) {
					t.Errorf("placeholder products/%s still exists after rename", Placeholder)
				}
			}

			// git init must have been asked for, on the right directory.
			if len(runner.calls) != 1 {
				t.Fatalf("Init() ran %d commands, want 1 (git init)", len(runner.calls))
			}
			got := runner.calls[0]
			if got.name != "git" || got.args[0] != "init" || got.args[1] != dir {
				t.Errorf("command = %s %v, want git init %s", got.name, got.args, dir)
			}
		})
	}
}

// Test_scaffoldServiceFinal_InitStampsVersion pins the platform version stamp.
func Test_scaffoldServiceFinal_InitStampsVersion(t *testing.T) {
	svc, _, ctx := newServiceForTest(t)
	dir := filepath.Join(t.TempDir(), "repo")

	if err := svc.Init(ctx, dir, "example_product"); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "platform", "freelunch.yaml"))
	if err != nil {
		t.Fatalf("cannot read freelunch.yaml: %v", err)
	}
	// Line-wise, not substring: the dev version is "0.0.0-dev", of which the
	// unstamped "version: 0.0.0" is a prefix — a Contains check would pass
	// against an unstamped file and fail against a correctly stamped one.
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == versionLine {
			t.Errorf("freelunch.yaml still carries the literal line %q; the version was not stamped", versionLine)
		}
	}
	want := "version: " + buildinfo.Get().Version
	if !strings.Contains(string(data), want) {
		t.Errorf("freelunch.yaml missing %q; got:\n%s", want, string(data))
	}
}

func Test_scaffoldServiceFinal_InitRefusesExistingDir(t *testing.T) {
	svc, _, ctx := newServiceForTest(t)
	dir := t.TempDir() // exists by construction

	err := svc.Init(ctx, dir, "example_product")
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Init() error = %v, want an already-exists refusal", err)
	}
}

// Test_scaffoldServiceFinal_InitCleansUpOnGitFailure pins the atomicity
// contract: a failed init leaves nothing behind.
func Test_scaffoldServiceFinal_InitCleansUpOnGitFailure(t *testing.T) {
	svc, runner, ctx := newServiceForTest(t)
	runner.failAll = errors.New("git: command not found")
	dir := filepath.Join(t.TempDir(), "repo")

	err := svc.Init(ctx, dir, "example_product")
	if err == nil || !strings.Contains(err.Error(), "git init") {
		t.Fatalf("Init() error = %v, want a git init failure", err)
	}
	if _, statErr := os.Stat(dir); !os.IsNotExist(statErr) {
		t.Errorf("failed Init() left %s behind", dir)
	}
}

// Test_embeddedTemplateMatchesRepoTemplate guards the two template copies
// against drift. The embedded copy is what ships; templates/monorepo is what
// the GitHub template publishing uses. They must stay identical.
func Test_embeddedTemplateMatchesRepoTemplate(t *testing.T) {
	repoTemplate := filepath.Join("..", "..", "..", "..", "templates", "monorepo")
	if _, err := os.Stat(repoTemplate); err != nil {
		t.Skipf("repo template not present (release archive?): %v", err)
	}

	// Collect relative path → content for both sides.
	repoFiles := map[string]string{}
	err := filepath.Walk(repoTemplate, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return walkErr
		}
		rel, _ := filepath.Rel(repoTemplate, path)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		repoFiles[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("cannot walk repo template: %v", err)
	}

	embeddedCount := 0
	for rel, want := range repoFiles {
		data, readErr := templateFS.ReadFile("template/" + rel)
		if readErr != nil {
			t.Errorf("embedded template is missing %s — re-copy templates/monorepo into internal/scaffold/template", rel)
			continue
		}
		if string(data) != want {
			t.Errorf("embedded template %s differs from templates/monorepo — re-copy to fix", rel)
		}
		embeddedCount++
	}
	if embeddedCount != len(repoFiles) {
		t.Errorf("compared %d files, repo template has %d", embeddedCount, len(repoFiles))
	}
}

// Test_templateCarriesVersionLine guards the stamp target: if the template's
// freelunch.yaml stops saying `version: 0.0.0`, Init silently stops stamping.
func Test_templateCarriesVersionLine(t *testing.T) {
	data, err := templateFS.ReadFile("template/platform/freelunch.yaml")
	if err != nil {
		t.Fatalf("embedded template has no platform/freelunch.yaml: %v", err)
	}
	if !strings.Contains(string(data), versionLine) {
		t.Errorf("template freelunch.yaml does not contain %q; Init's stamp has no target", versionLine)
	}
}
