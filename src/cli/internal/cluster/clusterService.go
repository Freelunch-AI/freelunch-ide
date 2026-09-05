// Package cluster provides the real ClusterService: the local Kubernetes
// cluster from roadmap 1.2.
//
// The service orchestrates the pinned k3d and kubectl binaries; it does not
// reimplement cluster management. k3d-cluster.yaml is the source of truth for
// the cluster's shape and is embedded here so a released binary carries the
// same definition the repository was tested with.
package cluster

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Freelunch-AI/freelunch-ide/src/cli/internal/managers"
)

// ClusterName is the cluster's name as k3d knows it.
//
// It must match metadata.name in k3d-cluster.yaml. A test asserts they agree,
// because a mismatch would make Delete and Status silently address a cluster
// that does not exist.
const ClusterName = "freelunch"

// BinDirEnvVar overrides the directory the pinned tools are read from. It is
// the same variable the installer scripts honour, so redirecting it moves both
// halves together.
const BinDirEnvVar = "FREELUNCH_BIN_DIR"

// ImagesDirEnvVar overrides the directory the k3s airgap image bundle is read
// from, mirroring BinDirEnvVar. The installer script honours the same variable.
const ImagesDirEnvVar = "FREELUNCH_IMAGES_DIR"

// nodeImagesPath is where k3s looks for image tarballs to import at startup.
// Anything mounted here is loaded into the node's containerd before workloads
// are scheduled, which is what lets a cluster come up with no network.
const nodeImagesPath = "/var/lib/rancher/k3s/agent/images"

// ErrNoTools is returned when the pinned binaries cannot be located. It is
// exported so a caller can tell "setup was never run" apart from a genuine
// cluster failure, via errors.Is.
var ErrNoTools = errors.New("pinned cluster tools not found")

// k3dConfig is the declarative cluster definition, embedded so the shipped
// binary does not depend on the repository being present.
//
//go:embed k3d-cluster.yaml
var k3dConfig []byte

type (
	// Runner executes an external command and returns its combined output.
	// It exists so tests can drive the service without Docker or a cluster.
	Runner interface {
		Run(ctx context.Context, name string, args ...string) ([]byte, error)
	}

	execRunner struct{}

	// toolset holds the resolved absolute paths of the pinned binaries.
	toolset struct {
		k3d     string
		kubectl string
	}

	clusterServiceFinal struct {
		sm     managers.ServiceManager
		runner Runner

		// Resolving the tools touches the filesystem, so it happens on first
		// real use rather than in Start. See the note on Start below.
		once     sync.Once
		tools    *toolset
		toolsErr error
	}
)

func (execRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// NewClusterService builds the service against the real system.
func NewClusterService() managers.ClusterService {
	return &clusterServiceFinal{runner: execRunner{}}
}

// NewClusterServiceWithRunner builds the service against a caller-supplied
// Runner. It is the seam tests use to exercise every path without creating a
// cluster.
func NewClusterServiceWithRunner(r Runner) managers.ClusterService {
	return &clusterServiceFinal{runner: r}
}

// Start does no I/O. Locating the k3d and kubectl binaries touches the
// filesystem, and this runs on every invocation including `freelunch --help`,
// so resolution is deferred to first real use behind sync.Once.
func (s *clusterServiceFinal) Start(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "starting the cluster service")
	return nil
}

func (s *clusterServiceFinal) Close(ctx context.Context) error {
	s.sm.LogsService().Debug(ctx, "stopping the cluster service")
	return nil
}

// Healthy reports whether the pinned tools can be found. It is not called
// during startup — that would pay filesystem lookups on every command — and
// exists for an explicit diagnostic command to drive.
func (s *clusterServiceFinal) Healthy(ctx context.Context) error {
	_, err := s.toolset(ctx)
	return err
}

func (s *clusterServiceFinal) WithServiceManager(sm managers.ServiceManager) managers.ClusterService {
	s.sm = sm
	return s
}

func (s *clusterServiceFinal) ServiceManager() managers.ServiceManager {
	return s.sm
}

// binDir is where the installer scripts place the pinned tools. They are
// deliberately not on the user's PATH, so a bare "k3d" must never be executed —
// that could pick up an unrelated version.
func binDir() (string, error) {
	if dir := os.Getenv(BinDirEnvVar); dir != "" {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	return filepath.Join(home, ".freelunch", "bin"), nil
}

// imagesDir is where install-airgap-images.sh caches the k3s image bundle.
func imagesDir() (string, error) {
	if dir := os.Getenv(ImagesDirEnvVar); dir != "" {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	return filepath.Join(home, ".freelunch", "images"), nil
}

// airgapVolume returns the k3d --volume argument that mounts the cached image
// bundle into every node, or "" when no cache is present.
//
// A missing cache is deliberately not an error: creating a cluster with a
// working network still succeeds without it, and failing here would break the
// path most contributors use. The caller logs which mode it chose, because
// "airgap worked" and "the network quietly carried it" look identical
// otherwise.
func airgapVolume() (string, error) {
	dir, err := imagesDir()
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		// Not cached, or unreadable. Either way there is nothing to mount.
		return "", nil
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "k3s-airgap-images-") {
			// Mounted on servers and agents alike: an agent schedules pods too,
			// so it needs the same images the server does.
			return fmt.Sprintf("%s:%s@server:*;agent:*", dir, nodeImagesPath), nil
		}
	}

	return "", nil
}

// toolset resolves the pinned binaries once, on first use.
func (s *clusterServiceFinal) toolset(ctx context.Context) (*toolset, error) {
	s.once.Do(func() {
		s.sm.LogsService().Debug(ctx, "resolving the pinned cluster tools")

		dir, err := binDir()
		if err != nil {
			s.toolsErr = err
			return
		}

		resolved := &toolset{
			k3d:     filepath.Join(dir, "k3d"),
			kubectl: filepath.Join(dir, "kubectl"),
		}

		for _, path := range []string{resolved.k3d, resolved.kubectl} {
			if _, statErr := os.Stat(path); statErr != nil {
				s.toolsErr = fmt.Errorf(
					"%w: %s is missing; run `pixi run task setup:cluster-tools`: %w",
					ErrNoTools, path, statErr)
				return
			}
		}

		s.tools = resolved
	})

	return s.tools, s.toolsErr
}

// writeConfig materialises the embedded cluster definition so k3d can read it.
// The caller invokes cleanup to remove the file.
func writeConfig() (path string, cleanup func(), err error) {
	dir, err := os.MkdirTemp("", "freelunch-cluster-")
	if err != nil {
		return "", func() {}, fmt.Errorf("cannot create temp dir for cluster config: %w", err)
	}

	cleanup = func() { _ = os.RemoveAll(dir) }

	path = filepath.Join(dir, "k3d-cluster.yaml")
	if writeErr := os.WriteFile(path, k3dConfig, 0o600); writeErr != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("cannot write cluster config: %w", writeErr)
	}

	return path, cleanup, nil
}

// Create brings the cluster up from the embedded configuration.
func (s *clusterServiceFinal) Create(ctx context.Context) error {
	tools, err := s.toolset(ctx)
	if err != nil {
		return err
	}

	path, cleanup, err := writeConfig()
	if err != nil {
		return err
	}
	defer cleanup()

	s.sm.LogsService().Debug(ctx, "creating cluster "+ClusterName)

	args := []string{"cluster", "create", "--config", path}

	volume, err := airgapVolume()
	if err != nil {
		return err
	}
	if volume != "" {
		args = append(args, "--volume", volume)
		s.sm.LogsService().Debug(ctx, "mounting cached k3s images: "+volume)
	} else {
		s.sm.LogsService().Debug(ctx,
			"no cached k3s images; the cluster will pull them from the network")
	}

	out, err := s.runner.Run(ctx, tools.k3d, args...)
	if err != nil {
		return fmt.Errorf("k3d cluster create failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// Delete tears the cluster down. Deleting a cluster that is already absent is
// not an error, so teardown can run unconditionally.
func (s *clusterServiceFinal) Delete(ctx context.Context) error {
	tools, err := s.toolset(ctx)
	if err != nil {
		return err
	}

	s.sm.LogsService().Debug(ctx, "deleting cluster "+ClusterName)

	out, err := s.runner.Run(ctx, tools.k3d, "cluster", "delete", ClusterName)
	if err != nil {
		return fmt.Errorf("k3d cluster delete failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// nodeReadinessTemplate is the kubectl jsonpath that prints one line per node:
// the node name, a tab, and the status of its Ready condition ("True",
// "False" or "Unknown"). A node that has no Ready condition yet prints an empty
// second field. `-o name` was not enough here: it lists NotReady nodes exactly
// like healthy ones, so existence would have been reported as readiness.
const nodeReadinessTemplate = `{range .items[*]}{.metadata.name}{"\t"}` +
	`{range .status.conditions[?(@.type=="Ready")]}{.status}{end}{"\n"}{end}`

// parseNodeReadiness turns nodeReadinessTemplate output into nodes. Only a
// Ready condition of exactly "True" counts as Ready; "False", "Unknown" (the
// kubelet stopped reporting) and a missing condition are all NotReady.
func parseNodeReadiness(out []byte) []managers.ClusterNode {
	var nodes []managers.ClusterNode
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		name, ready, _ := strings.Cut(strings.TrimSpace(line), "\t")
		if name == "" {
			continue
		}
		nodes = append(nodes, managers.ClusterNode{
			Name:  name,
			Ready: strings.TrimSpace(ready) == "True",
		})
	}
	return nodes
}

// Status reports whether the cluster exists, whether its API answers, and the
// readiness of each node. Running is true only when every node is Ready.
//
// A missing cluster, or one whose API is not answering, is a normal answer
// rather than an error: callers ask this precisely because they do not know
// yet, and `status` must stay useful while a cluster is still coming up.
func (s *clusterServiceFinal) Status(ctx context.Context) (*managers.ClusterStatus, error) {
	tools, err := s.toolset(ctx)
	if err != nil {
		return nil, err
	}

	status := &managers.ClusterStatus{Name: ClusterName}

	if _, listErr := s.runner.Run(ctx, tools.k3d, "cluster", "list", ClusterName); listErr != nil {
		return status, nil
	}
	status.Exists = true

	out, err := s.runner.Run(ctx, tools.kubectl, "get", "nodes", "-o", "jsonpath="+nodeReadinessTemplate)
	if err != nil {
		return status, nil
	}

	status.Nodes = parseNodeReadiness(out)

	status.Running = len(status.Nodes) > 0
	for _, node := range status.Nodes {
		if !node.Ready {
			status.Running = false
			break
		}
	}

	return status, nil
}
