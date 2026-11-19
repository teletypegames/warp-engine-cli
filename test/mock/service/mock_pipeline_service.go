package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockPipelineService mocks [domain.PipelineServiceInterface].
type MockPipelineService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PipelineModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(pipeline_id int) (*domain.PipelineModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(pipeline_id int, attributes domain.Attributes) (*domain.PipelineModel, error)
	// SyncFunc replaces Sync.
	SyncFunc func() (*domain.PipelineSyncModel, error)
	// RunsFunc replaces Runs.
	RunsFunc func(pipeline_id int, page int) (*domain.PageModel[domain.PipelineRunModel], error)
	// TriggerFunc replaces Trigger.
	TriggerFunc func(pipeline_id int, branch string) (*domain.PipelineRunModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockPipelineService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PipelineModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.PipelineModel]{}, nil
}

// Show calls ShowFunc, or returns an empty pipeline.
func (m *MockPipelineService) Show(pipeline_id int) (*domain.PipelineModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(pipeline_id)
	}
	return &domain.PipelineModel{}, nil
}

// Update calls UpdateFunc, or returns an empty pipeline.
func (m *MockPipelineService) Update(pipeline_id int, attributes domain.Attributes) (*domain.PipelineModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(pipeline_id, attributes)
	}
	return &domain.PipelineModel{}, nil
}

// Sync calls SyncFunc, or returns an empty sync.
func (m *MockPipelineService) Sync() (*domain.PipelineSyncModel, error) {
	if m.SyncFunc != nil {
		return m.SyncFunc()
	}
	return &domain.PipelineSyncModel{}, nil
}

// Runs calls RunsFunc, or returns an empty page.
func (m *MockPipelineService) Runs(pipeline_id int, page int) (*domain.PageModel[domain.PipelineRunModel], error) {
	if m.RunsFunc != nil {
		return m.RunsFunc(pipeline_id, page)
	}
	return &domain.PageModel[domain.PipelineRunModel]{}, nil
}

// Trigger calls TriggerFunc, or returns an empty run.
func (m *MockPipelineService) Trigger(pipeline_id int, branch string) (*domain.PipelineRunModel, error) {
	if m.TriggerFunc != nil {
		return m.TriggerFunc(pipeline_id, branch)
	}
	return &domain.PipelineRunModel{}, nil
}
