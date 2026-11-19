package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockReleaseService mocks [domain.ReleaseServiceInterface].
type MockReleaseService struct {
	// AllFunc replaces All.
	AllFunc func(reference string, filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ReleaseModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(id string) (*domain.ReleaseModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(reference string, attributes domain.Attributes) (*domain.ReleaseModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(id string, attributes domain.Attributes) (*domain.ReleaseModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(id string) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockReleaseService) All(reference string, filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ReleaseModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(reference, filter, query)
	}
	return &domain.PageModel[domain.ReleaseModel]{}, nil
}

// Show calls ShowFunc, or returns an empty release.
func (m *MockReleaseService) Show(id string) (*domain.ReleaseModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(id)
	}
	return &domain.ReleaseModel{}, nil
}

// Create calls CreateFunc, or returns an empty release.
func (m *MockReleaseService) Create(reference string, attributes domain.Attributes) (*domain.ReleaseModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(reference, attributes)
	}
	return &domain.ReleaseModel{}, nil
}

// Update calls UpdateFunc, or returns an empty release.
func (m *MockReleaseService) Update(id string, attributes domain.Attributes) (*domain.ReleaseModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, attributes)
	}
	return &domain.ReleaseModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockReleaseService) Destroy(id string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(id)
	}
	return nil
}
