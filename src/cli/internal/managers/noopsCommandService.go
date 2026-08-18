package managers

import "context"

type (
	noOpsCommandService struct {
		sm ServiceManager
	}
)

// NewNoOpsCommandService returns a CommandService that runs nothing and
// succeeds. It is the default in NewManager, which lets a test exercise the
// container's lifecycle without building a real command tree.
func NewNoOpsCommandService() CommandService {
	return &noOpsCommandService{}
}

func (n *noOpsCommandService) Start(_ context.Context) error {
	return nil
}

func (n *noOpsCommandService) Close(_ context.Context) error {
	return nil
}

func (n *noOpsCommandService) Healthy(_ context.Context) error {
	return nil
}

func (n *noOpsCommandService) WithServiceManager(sm ServiceManager) CommandService {
	n.sm = sm
	return n
}

func (n *noOpsCommandService) ServiceManager() ServiceManager {
	return n.sm
}

func (n *noOpsCommandService) Execute(_ context.Context, _ []string) error {
	return nil
}
