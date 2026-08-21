// Package auth provides the real AuthService: the local OIDC identity provider
// from roadmap 1.3, a Keycloak instance running in the Demo cluster.
//
// The service orchestrates the pinned kubectl against manifests embedded here;
// it does not talk to the Kubernetes API directly. keycloak.yaml and
// freelunch-realm.json are the source of truth for the server and its realm,
// embedded so a released binary carries the definition it was tested with. The
// realm is re-imported on every pod start — dev mode keeps an in-memory
// database, so the committed JSON is what makes the realm durable, not the
// admin console.
package auth

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// Namespace is where every FreeLunch platform component lands, starting here.
const Namespace = "freelunch-system"

// RealmName must match the "realm" field in freelunch-realm.json. A test
// asserts they agree, because Status probes the realm's discovery endpoint by
// this name and a mismatch would report a healthy server as absent.
const RealmName = "freelunch"

// DiscoveryURL is the host-reachable OIDC discovery endpoint. Status performs
// a real HTTP GET against it rather than trusting pod readiness: the pod being
// Ready proves the process started, not that the issuer is right or that the
// ingress routes — and during a pod swap, kubectl reports the rollout done a
// few seconds before Traefik switches endpoints.
const DiscoveryURL = "http://keycloak.localhost:8080/realms/" + RealmName +
	"/.well-known/openid-configuration"

// BinDirEnvVar mirrors the cluster package's variable of the same name. It is
// redeclared rather than imported: concrete service packages never import each
// other, and the constant is the contract with the installer scripts either way.
const BinDirEnvVar = "FREELUNCH_BIN_DIR"

//go:embed keycloak.yaml
var keycloakManifests []byte

//go:embed freelunch-realm.json
var realmJSON []byte

type (
	// Runner executes an external command and returns its combined output.
	// It exists so tests can drive the service without a cluster.
	Runner interface {
		Run(ctx context.Context, name string, args ...string) ([]byte, error)
	}

	// Getter performs the HTTP GET Status uses to probe OIDC discovery. It
	// exists so tests can exercise Status without a server.
	Getter interface {
		Get(url string) (*http.Response, error)
	}

	execRunner struct{}

	authServiceFinal struct {
		sm     managers.ServiceManager
		runner Runner
		getter Getter

		// Resolving kubectl touches the filesystem, so it happens on first
		// real use rather than in Start.
		once       sync.Once
		kubectl    string
		kubectlErr error
	}
)

func (execRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := execCommand(ctx, name, args...)
	return cmd.CombinedOutput()
}

// NewAuthService builds the service against the real system.
func NewAuthService() managers.AuthService {
	return &authServiceFinal{
		runner: execRunner{},
		getter: &http.Client{Timeout: 5 * time.Second},
	}
}

// NewAuthServiceWithDeps builds the service against caller-supplied seams. It
// is what tests use to exercise every path without a cluster or a server.
func NewAuthServiceWithDeps(r Runner, g Getter) managers.AuthService {
	return &authServiceFinal{runner: r, getter: g}
}

// Start does no I/O — locating kubectl touches the filesystem and this runs on
// every invocation including `freelunch --help`, so resolution is deferred to
// first real use behind sync.Once.
func (s *authServiceFinal) Start(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "starting the auth service")
	return nil
}

func (s *authServiceFinal) Close(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "stopping the auth service")
	return nil
}

// Healthy reports whether kubectl can be found. Not called during startup —
// that would pay a filesystem lookup on every command.
func (s *authServiceFinal) Healthy(ctx context.Context) error {
	_, err := s.tool(ctx)
	return err
}

func (s *authServiceFinal) WithServiceManager(sm managers.ServiceManager) managers.AuthService {
	s.sm = sm
	return s
}

func (s *authServiceFinal) ServiceManager() managers.ServiceManager {
	return s.sm
}

// tool resolves the pinned kubectl once, on first use. Deliberately never a
// bare "kubectl": the pinned binary is not on PATH, and resolving via PATH
// could pick up an unrelated version.
func (s *authServiceFinal) tool(ctx context.Context) (string, error) {
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

// writeManifests materialises the embedded definitions so kubectl can read
// them. The caller invokes cleanup to remove the directory.
func writeManifests() (dir string, cleanup func(), err error) {
	dir, err = os.MkdirTemp("", "freelunch-auth-")
	if err != nil {
		return "", func() {}, fmt.Errorf("cannot create temp dir for auth manifests: %w", err)
	}

	cleanup = func() { _ = os.RemoveAll(dir) }

	for name, data := range map[string][]byte{
		"keycloak.yaml":        keycloakManifests,
		"freelunch-realm.json": realmJSON,
	} {
		if writeErr := os.WriteFile(filepath.Join(dir, name), data, 0o600); writeErr != nil {
			cleanup()
			return "", func() {}, fmt.Errorf("cannot write %s: %w", name, writeErr)
		}
	}

	return dir, cleanup, nil
}

// Install applies the auth service and its realm to the running cluster.
//
// The realm ConfigMap is applied first because the Deployment mounts it — the
// other order works too (Kubernetes reconciles), but this way a watching human
// never sees the pod crash-loop on a missing mount.
func (s *authServiceFinal) Install(ctx context.Context) error {
	kubectl, err := s.tool(ctx)
	if err != nil {
		return err
	}

	dir, cleanup, err := writeManifests()
	if err != nil {
		return err
	}
	defer cleanup()

	s.sm.LogsService().Debug(ctx, "installing the auth service")

	// The namespace must exist before the ConfigMap can land in it, and both
	// live in keycloak.yaml — so apply the manifests, then the ConfigMap, then
	// restart to pick up realm changes on re-install.
	out, err := s.runner.Run(ctx, kubectl, "apply", "-f", filepath.Join(dir, "keycloak.yaml"))
	if err != nil {
		return fmt.Errorf("kubectl apply failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	out, err = s.runner.Run(ctx, kubectl,
		"-n", Namespace, "create", "configmap", "keycloak-realm",
		"--from-file=freelunch-realm.json="+filepath.Join(dir, "freelunch-realm.json"),
		"--dry-run=client", "-o", "yaml")
	if err != nil {
		return fmt.Errorf("kubectl create configmap failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	tmpCM := filepath.Join(dir, "realm-configmap.yaml")
	if err = os.WriteFile(tmpCM, out, 0o600); err != nil {
		return fmt.Errorf("cannot write realm configmap: %w", err)
	}

	out, err = s.runner.Run(ctx, kubectl, "apply", "-f", tmpCM)
	if err != nil {
		return fmt.Errorf("kubectl apply configmap failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	// A changed realm JSON does not restart the pod on its own — the import
	// runs at startup. Rollout restart makes re-install pick changes up; on a
	// fresh install it is a no-op restart of a pod still coming up.
	out, err = s.runner.Run(ctx, kubectl,
		"-n", Namespace, "rollout", "restart", "deployment/keycloak")
	if err != nil {
		return fmt.Errorf("kubectl rollout restart failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// Delete removes the auth service. Deleting what is not there succeeds, so
// teardown can run unconditionally.
func (s *authServiceFinal) Delete(ctx context.Context) error {
	kubectl, err := s.tool(ctx)
	if err != nil {
		return err
	}

	dir, cleanup, err := writeManifests()
	if err != nil {
		return err
	}
	defer cleanup()

	s.sm.LogsService().Debug(ctx, "deleting the auth service")

	out, err := s.runner.Run(ctx, kubectl,
		"delete", "-f", filepath.Join(dir, "keycloak.yaml"), "--ignore-not-found")
	if err != nil {
		return fmt.Errorf("kubectl delete failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// Status reports whether the auth service answers OIDC discovery from the
// host, and with which issuer.
//
// An unreachable or not-yet-ready server is a normal answer, not an error:
// callers ask precisely because they do not know. Only being unable to ask —
// no kubectl, say — would be an error, and Status needs no kubectl at all.
func (s *authServiceFinal) Status(_ context.Context) (*managers.AuthStatus, error) {
	status := &managers.AuthStatus{}

	resp, err := s.getter.Get(DiscoveryURL)
	if err != nil {
		return status, nil
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return status, nil
	}

	var doc struct {
		Issuer string `json:"issuer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return status, nil
	}

	status.Ready = true
	status.Realm = RealmName
	status.IssuerURL = doc.Issuer

	return status, nil
}
