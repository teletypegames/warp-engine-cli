package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockRuntimeRepository mocks [domain.RuntimeRepositoryInterface].
type MockRuntimeRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.RuntimeVersionModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(runtime_id int) (*domain.RuntimeVersionModel, error)
	// InstallFunc replaces Install.
	InstallFunc func(attributes domain.Attributes) (*domain.RuntimeInstallModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockRuntimeRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.RuntimeVersionModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.RuntimeVersionModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty runtime version.
func (m *MockRuntimeRepository) FindByID(runtime_id int) (*domain.RuntimeVersionModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(runtime_id)
	}
	return &domain.RuntimeVersionModel{}, nil
}

// Install calls InstallFunc, or returns an empty install.
func (m *MockRuntimeRepository) Install(attributes domain.Attributes) (*domain.RuntimeInstallModel, error) {
	if m.InstallFunc != nil {
		return m.InstallFunc(attributes)
	}
	return &domain.RuntimeInstallModel{}, nil
}
