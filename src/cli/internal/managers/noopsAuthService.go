package managers

import "context"

type (
	noOpsAuthService struct {
		sm ServiceManager
	}
)

// NewNoOpsAuthService returns an AuthService that touches nothing and
// succeeds. It is the default in NewManager, which lets a test exercise the
// container's lifecycle — and every command that does not touch auth — without
// a cluster or an identity provider.
func NewNoOpsAuthService() AuthService {
	return &noOpsAuthService{}
}

func (n *noOpsAuthService) Start(_ context.Context) error {
	return nil
}

func (n *noOpsAuthService) Close(_ context.Context) error {
	return nil
}

func (n *noOpsAuthService) Healthy(_ context.Context) error {
	return nil
}

func (n *noOpsAuthService) WithServiceManager(sm ServiceManager) AuthService {
	n.sm = sm
	return n
}

func (n *noOpsAuthService) ServiceManager() ServiceManager {
	return n.sm
}

func (n *noOpsAuthService) Install(_ context.Context) error {
	return nil
}

func (n *noOpsAuthService) Delete(_ context.Context) error {
	return nil
}

func (n *noOpsAuthService) Status(_ context.Context) (*AuthStatus, error) {
	return &AuthStatus{}, nil
}
