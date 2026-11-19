package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockRuntimeService mocks [domain.RuntimeServiceInterface].
type MockRuntimeService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.RuntimeVersionModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(runtime_id int) (*domain.RuntimeVersionModel, error)
	// InstallFunc replaces Install.
	InstallFunc func(name string, version string, dry_run bool) (*domain.RuntimeInstallModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockRuntimeService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.RuntimeVersionModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.RuntimeVersionModel]{}, nil
}

// Show calls ShowFunc, or returns an empty runtime version.
func (m *MockRuntimeService) Show(runtime_id int) (*domain.RuntimeVersionModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(runtime_id)
	}
	return &domain.RuntimeVersionModel{}, nil
}

// Install calls InstallFunc, or returns an empty install.
func (m *MockRuntimeService) Install(name string, version string, dry_run bool) (*domain.RuntimeInstallModel, error) {
	if m.InstallFunc != nil {
		return m.InstallFunc(name, version, dry_run)
	}
	return &domain.RuntimeInstallModel{}, nil
}
