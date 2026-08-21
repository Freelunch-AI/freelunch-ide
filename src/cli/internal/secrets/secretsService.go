// Package secrets provides the real SecretsService: the local secrets store
// from roadmap 1.4, an OpenBao instance running in the Demo cluster.
//
// OpenBao rather than HashiCorp Vault (decision D1): Vault has been BUSL-1.1
// since August 2023, which restricts embedded commercial use; OpenBao is the
// Linux Foundation fork of the last MPL-2.0 codebase with the same API
// surface. Note the renamed CLI — `bao`, with BAO_* environment variables.
//
// The service orchestrates the pinned kubectl against openbao.yaml embedded
// here; the store's own CLI is reached via `kubectl exec` into the pod, so no
// bao binary is ever installed on the host. Dev mode is in-memory and loses
// everything on restart, which is why PutSecret is re-runnable and seeding
// belongs to install.
package secrets

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// Namespace is where every FreeLunch platform component lands.
const Namespace = "freelunch-system"

// DevRootToken is the dev-mode root token. It is deliberately a constant and
// deliberately not a secret: this store runs in dev mode, holds nothing real,
// and roadmap 2.1's SecretStore needs a known token to authenticate with.
const DevRootToken = "freelunch-dev-root"

// Mount is the KV engine mount path. Dev mode creates it as KV version 2,
// which means the HTTP API reads secrets at Mount + "data/" + path — the
// indirection external-secrets-operator must be configured for in 2.1.
const Mount = "secret"

// BinDirEnvVar mirrors the other service packages' variable of the same name.
// Redeclared rather than imported: concrete service packages never import each
// other.
const BinDirEnvVar = "FREELUNCH_BIN_DIR"

// podSelector finds the store's pod for exec. Only running pods qualify —
// during a pod swap the terminated one is still listed, and exec into it
// fails with a confusing "completed pod" error.
const podSelector = "app.kubernetes.io/name=openbao"

//go:embed openbao.yaml
var openbaoManifests []byte

type (
	// Runner executes an external command and returns its combined output.
	// It exists so tests can drive the service without a cluster.
	Runner interface {
		Run(ctx context.Context, name string, args ...string) ([]byte, error)
	}

	execRunner struct{}

	secretsServiceFinal struct {
		sm     managers.ServiceManager
		runner Runner

		// Resolving kubectl touches the filesystem, so it happens on first
		// real use rather than in Start.
		once       sync.Once
		kubectl    string
		kubectlErr error
	}
)

func (execRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// NewSecretsService builds the service against the real system.
func NewSecretsService() managers.SecretsService {
	return &secretsServiceFinal{runner: execRunner{}}
}

// NewSecretsServiceWithRunner builds the service against a caller-supplied
// Runner. It is the seam tests use to exercise every path without a cluster.
func NewSecretsServiceWithRunner(r Runner) managers.SecretsService {
	return &secretsServiceFinal{runner: r}
}

// Start does no I/O; see the note on the cluster service's Start.
func (s *secretsServiceFinal) Start(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "starting the secrets service")
	return nil
}

func (s *secretsServiceFinal) Close(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "stopping the secrets service")
	return nil
}

// Healthy reports whether kubectl can be found. Not called during startup.
func (s *secretsServiceFinal) Healthy(ctx context.Context) error {
	_, err := s.tool(ctx)
	return err
}

func (s *secretsServiceFinal) WithServiceManager(sm managers.ServiceManager) managers.SecretsService {
	s.sm = sm
	return s
}

func (s *secretsServiceFinal) ServiceManager() managers.ServiceManager {
	return s.sm
}

// tool resolves the pinned kubectl once, on first use. Never a bare "kubectl".
func (s *secretsServiceFinal) tool(ctx context.Context) (string, error) {
	s.once.Do(func() {
		s.sm.LogsService().Debug(ctx, "resolving the pinned kubectl")

		dir := os.Getenv(BinDirEnvVar)
		if dir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				s.kubectlErr = fmt.Errorf("cannot determine home directory: %w", err)
				return
			}
			dir = filepath.Join(home, ".freelunch", "bin")
		}

		path := filepath.Join(dir, "kubectl")
		if _, err := os.Stat(path); err != nil {
			s.kubectlErr = fmt.Errorf(
				"pinned kubectl not found: %s is missing; run `pixi run task setup:cluster-tools`: %w",
				path, err)
			return
		}
		s.kubectl = path
	})

	return s.kubectl, s.kubectlErr
}

// writeManifests materialises the embedded definition so kubectl can read it.
func writeManifests() (path string, cleanup func(), err error) {
	dir, err := os.MkdirTemp("", "freelunch-secrets-")
	if err != nil {
		return "", func() {}, fmt.Errorf("cannot create temp dir for secrets manifests: %w", err)
	}

	cleanup = func() { _ = os.RemoveAll(dir) }

	path = filepath.Join(dir, "openbao.yaml")
	if writeErr := os.WriteFile(path, openbaoManifests, 0o600); writeErr != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("cannot write openbao.yaml: %w", writeErr)
	}

	return path, cleanup, nil
}

// pod returns the name of the running store pod. Callers exec into it — the
// bao CLI ships inside the image, so nothing is installed on the host.
func (s *secretsServiceFinal) pod(ctx context.Context, kubectl string) (string, error) {
	out, err := s.runner.Run(ctx, kubectl,
		"-n", Namespace, "get", "pod",
		"-l", podSelector,
		"--field-selector=status.phase=Running",
		"-o", "jsonpath={.items[0].metadata.name}")
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return "", fmt.Errorf("no running secrets store pod: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

// bao runs one bao CLI command inside the store pod, authenticated as the
// dev root.
func (s *secretsServiceFinal) bao(ctx context.Context, kubectl, pod string, args ...string) ([]byte, error) {
	full := append([]string{
		"-n", Namespace, "exec", pod, "--",
		"env", "BAO_ADDR=http://127.0.0.1:8200", "BAO_TOKEN=" + DevRootToken,
		"bao",
	}, args...)
	return s.runner.Run(ctx, kubectl, full...)
}

// Install applies the store to the running cluster and waits for it to come
// up. Unlike the auth service — whose Keycloak takes ~60s and is therefore
// left to come up in the background — OpenBao answers within seconds, and
// install-time seeding needs a running pod, so waiting here is cheap and buys
// a store that is usable the moment install returns.
func (s *secretsServiceFinal) Install(ctx context.Context) error {
	kubectl, err := s.tool(ctx)
	if err != nil {
		return err
	}

	path, cleanup, err := writeManifests()
	if err != nil {
		return err
	}
	defer cleanup()

	s.sm.LogsService().Debug(ctx, "installing the secrets store")

	out, err := s.runner.Run(ctx, kubectl, "apply", "-f", path)
	if err != nil {
		return fmt.Errorf("kubectl apply failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	out, err = s.runner.Run(ctx, kubectl,
		"-n", Namespace, "rollout", "status", "deployment/openbao", "--timeout=180s")
	if err != nil {
		return fmt.Errorf("secrets store did not come up: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// Delete removes the store. Deleting what is not there succeeds.
func (s *secretsServiceFinal) Delete(ctx context.Context) error {
	kubectl, err := s.tool(ctx)
	if err != nil {
		return err
	}

	path, cleanup, err := writeManifests()
	if err != nil {
		return err
	}
	defer cleanup()

	s.sm.LogsService().Debug(ctx, "deleting the secrets store")

	out, err := s.runner.Run(ctx, kubectl, "delete", "-f", path, "--ignore-not-found")
	if err != nil {
		return fmt.Errorf("kubectl delete failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// Status reports readiness, seal state and the mounted engine.
//
// An absent or unready store is a normal answer, not an error. Sealed is
// reported separately from Ready because a sealed store still has a Ready pod
// — it is the failure that looks most like success from the outside. Dev mode
// never seals, so Sealed=true here means someone changed the deployment out
// from under us, which is exactly when it matters.
func (s *secretsServiceFinal) Status(ctx context.Context) (*managers.SecretsStatus, error) {
	kubectl, err := s.tool(ctx)
	if err != nil {
		return nil, err
	}

	status := &managers.SecretsStatus{}

	pod, err := s.pod(ctx, kubectl)
	if err != nil {
		return status, nil
	}

	out, err := s.bao(ctx, kubectl, pod, "status", "-format=json")
	if err != nil {
		return status, nil
	}

	var st struct {
		Initialized bool `json:"initialized"`
		Sealed      bool `json:"sealed"`
	}
	if jsonErr := json.Unmarshal(out, &st); jsonErr != nil {
		return status, nil
	}

	status.Sealed = st.Sealed
	status.Ready = st.Initialized && !st.Sealed

	if !status.Ready {
		return status, nil
	}

	out, err = s.bao(ctx, kubectl, pod, "secrets", "list", "-format=json")
	if err != nil {
		return status, nil
	}

	var mounts map[string]struct {
		Type    string            `json:"type"`
		Options map[string]string `json:"options"`
	}
	if jsonErr := json.Unmarshal(out, &mounts); jsonErr != nil {
		return status, nil
	}

	if m, ok := mounts[Mount+"/"]; ok {
		status.Engine = fmt.Sprintf("%s/ (%s v%s)", Mount, m.Type, m.Options["version"])
	}

	return status, nil
}

// PutSecret writes one key/value pair at a logical path under the KV mount.
// The CLI accepts the logical path (secret/<path>); the HTTP API will serve
// it back at secret/data/<path> — KV v2's rewrite, which 2.1 must configure
// for.
func (s *secretsServiceFinal) PutSecret(ctx context.Context, path, key, value string) error {
	kubectl, err := s.tool(ctx)
	if err != nil {
		return err
	}

	pod, err := s.pod(ctx, kubectl)
	if err != nil {
		return err
	}

	s.sm.LogsService().Debug(ctx, "seeding secret "+Mount+"/"+path)

	out, err := s.bao(ctx, kubectl, pod, "kv", "put", Mount+"/"+path, key+"="+value)
	if err != nil {
		return fmt.Errorf("bao kv put failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}
