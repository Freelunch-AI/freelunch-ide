package managers

import "context"

type (
	noOpsSecretsService struct {
		sm ServiceManager
	}
)

// NewNoOpsSecretsService returns a SecretsService that touches nothing and
// succeeds. It is the default in NewManager, which lets a test exercise the
// container's lifecycle — and every command that does not touch secrets —
// without a cluster or a store.
func NewNoOpsSecretsService() SecretsService {
	return &noOpsSecretsService{}
}

func (n *noOpsSecretsService) Start(_ context.Context) error {
	return nil
}

func (n *noOpsSecretsService) Close(_ context.Context) error {
	return nil
}

func (n *noOpsSecretsService) Healthy(_ context.Context) error {
	return nil
}

func (n *noOpsSecretsService) WithServiceManager(sm ServiceManager) SecretsService {
	n.sm = sm
	return n
}

func (n *noOpsSecretsService) ServiceManager() ServiceManager {
	return n.sm
}

func (n *noOpsSecretsService) Install(_ context.Context) error {
	return nil
}

func (n *noOpsSecretsService) Delete(_ context.Context) error {
	return nil
}

func (n *noOpsSecretsService) Status(_ context.Context) (*SecretsStatus, error) {
	return &SecretsStatus{}, nil
}

func (n *noOpsSecretsService) PutSecret(_ context.Context, _, _, _ string) error {
	return nil
}
