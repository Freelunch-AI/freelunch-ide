package managers

import "context"

type (
	noOpsClusterService struct {
		sm ServiceManager
	}
)

// NewNoOpsClusterService returns a ClusterService that touches no cluster and
// succeeds. It is the default in NewManager, which lets a test exercise the
// container's lifecycle — and every command that does not touch the cluster —
// without Docker being installed.
func NewNoOpsClusterService() ClusterService {
	return &noOpsClusterService{}
}

func (n *noOpsClusterService) Start(_ context.Context) error {
	return nil
}

func (n *noOpsClusterService) Close(_ context.Context) error {
	return nil
}

func (n *noOpsClusterService) Healthy(_ context.Context) error {
	return nil
}

func (n *noOpsClusterService) WithServiceManager(sm ServiceManager) ClusterService {
	n.sm = sm
	return n
}

func (n *noOpsClusterService) ServiceManager() ServiceManager {
	return n.sm
}

func (n *noOpsClusterService) Create(_ context.Context) error {
	return nil
}

func (n *noOpsClusterService) Delete(_ context.Context) error {
	return nil
}

func (n *noOpsClusterService) Status(_ context.Context) (*ClusterStatus, error) {
	return &ClusterStatus{}, nil
}
