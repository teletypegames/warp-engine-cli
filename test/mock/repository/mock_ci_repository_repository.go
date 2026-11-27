package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockCiRepositoryRepository mocks [domain.CiRepositoryRepositoryInterface].
type MockCiRepositoryRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.CiRepositoryModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(ci_repository_id int) (*domain.CiRepositoryModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(ci_repository_id int, attributes domain.Attributes) (*domain.CiRepositoryModel, error)
	// SyncFunc replaces Sync.
	SyncFunc func() (*domain.CiRepositorySyncModel, error)
	// RunsFunc replaces Runs.
	RunsFunc func(ci_repository_id int, page int) (*domain.PageModel[domain.PipelineRunModel], error)
	// TriggerFunc replaces Trigger.
	TriggerFunc func(ci_repository_id int, attributes domain.Attributes) (*domain.PipelineRunModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockCiRepositoryRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.CiRepositoryModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.CiRepositoryModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty ci_repository.
func (m *MockCiRepositoryRepository) FindByID(ci_repository_id int) (*domain.CiRepositoryModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ci_repository_id)
	}
	return &domain.CiRepositoryModel{}, nil
}

// Update calls UpdateFunc, or returns an empty ci_repository.
func (m *MockCiRepositoryRepository) Update(ci_repository_id int, attributes domain.Attributes) (*domain.CiRepositoryModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ci_repository_id, attributes)
	}
	return &domain.CiRepositoryModel{}, nil
}

// Sync calls SyncFunc, or returns an empty sync.
func (m *MockCiRepositoryRepository) Sync() (*domain.CiRepositorySyncModel, error) {
	if m.SyncFunc != nil {
		return m.SyncFunc()
	}
	return &domain.CiRepositorySyncModel{}, nil
}

// Runs calls RunsFunc, or returns an empty page.
func (m *MockCiRepositoryRepository) Runs(ci_repository_id int, page int) (*domain.PageModel[domain.PipelineRunModel], error) {
	if m.RunsFunc != nil {
		return m.RunsFunc(ci_repository_id, page)
	}
	return &domain.PageModel[domain.PipelineRunModel]{}, nil
}

// Trigger calls TriggerFunc, or returns an empty run.
func (m *MockCiRepositoryRepository) Trigger(ci_repository_id int, attributes domain.Attributes) (*domain.PipelineRunModel, error) {
	if m.TriggerFunc != nil {
		return m.TriggerFunc(ci_repository_id, attributes)
	}
	return &domain.PipelineRunModel{}, nil
}
