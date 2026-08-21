package managers

import "context"

type (
	noOpsScaffoldService struct {
		sm ServiceManager
	}
)

// NewNoOpsScaffoldService returns a ScaffoldService that touches nothing and
// succeeds. It is the default in NewManager, which lets a test exercise the
// container's lifecycle — and every command that does not scaffold — without
// touching the filesystem.
func NewNoOpsScaffoldService() ScaffoldService {
	return &noOpsScaffoldService{}
}

func (n *noOpsScaffoldService) Start(_ context.Context) error {
	return nil
}

func (n *noOpsScaffoldService) Close(_ context.Context) error {
	return nil
}

func (n *noOpsScaffoldService) Healthy(_ context.Context) error {
	return nil
}

func (n *noOpsScaffoldService) WithServiceManager(sm ServiceManager) ScaffoldService {
	n.sm = sm
	return n
}

func (n *noOpsScaffoldService) ServiceManager() ServiceManager {
	return n.sm
}

func (n *noOpsScaffoldService) Init(_ context.Context, _, _ string) error {
	return nil
}
