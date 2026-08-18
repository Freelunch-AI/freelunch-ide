package managers

import "context"

type (
	noOpsLogsService struct {
		sm ServiceManager
	}
)

// NewNoOpsLogsService returns a LogsService that discards everything. It is the
// default in NewManager so a caller that never registers a real logger can
// still call s.sm.LogsService().Info(...) without a nil check.
func NewNoOpsLogsService() LogsService {
	return &noOpsLogsService{}
}

func (n *noOpsLogsService) Start(_ context.Context) error {
	return nil
}

func (n *noOpsLogsService) Close(_ context.Context) error {
	return nil
}

func (n *noOpsLogsService) Healthy(_ context.Context) error {
	return nil
}

func (n *noOpsLogsService) WithServiceManager(sm ServiceManager) LogsService {
	n.sm = sm
	return n
}

func (n *noOpsLogsService) ServiceManager() ServiceManager {
	return n.sm
}

func (n *noOpsLogsService) Info(_ context.Context, _ string)  {}
func (n *noOpsLogsService) Warn(_ context.Context, _ string)  {}
func (n *noOpsLogsService) Error(_ context.Context, _ string) {}
func (n *noOpsLogsService) Debug(_ context.Context, _ string) {}
