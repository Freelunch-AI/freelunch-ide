package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func run(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	var out, errOut bytes.Buffer
	err = Execute(context.Background(), args, &out, &errOut)
	return out.String(), errOut.String(), err
}

func TestBareInvocationPrintsHelp(t *testing.T) {
	stdout, _, err := run(t)
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}
	if !strings.Contains(stdout, "freelunch") {
		t.Errorf("help output missing command name:\n%s", stdout)
	}
}

func TestVersionCommandText(t *testing.T) {
	stdout, _, err := run(t, "version")
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}
	for _, want := range []string{"freelunch ", "commit:", "go:", "platform:"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("version output missing %q:\n%s", want, stdout)
		}
	}
}

func TestVersionCommandJSON(t *testing.T) {
	stdout, _, err := run(t, "version", "--json")
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	var got map[string]any
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("version --json emitted invalid JSON: %v\n%s", err, stdout)
	}
	for _, key := range []string{"version", "commit", "date", "goVersion", "platform"} {
		if _, ok := got[key]; !ok {
			t.Errorf("version --json missing key %q, got %v", key, got)
		}
	}
}

func TestUnknownCommandFails(t *testing.T) {
	if _, _, err := run(t, "definitely-not-a-command"); err == nil {
		t.Fatal("Execute() error = nil, want an error for an unknown command")
	}
}

func TestVersionRejectsPositionalArgs(t *testing.T) {
	if _, _, err := run(t, "version", "extra"); err == nil {
		t.Fatal("Execute() error = nil, want an error for unexpected arguments")
	}
}
