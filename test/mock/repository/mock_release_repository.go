package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockReleaseRepository mocks [domain.ReleaseRepositoryInterface].
type MockReleaseRepository struct {
	// AllFunc replaces All.
	AllFunc func(software_id int, filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ReleaseModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(release_id int) (*domain.ReleaseModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(software_id int, attributes domain.Attributes) (*domain.ReleaseModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(release_id int, attributes domain.Attributes) (*domain.ReleaseModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(release_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockReleaseRepository) All(software_id int, filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ReleaseModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(software_id, filter, query)
	}
	return &domain.PageModel[domain.ReleaseModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty release.
func (m *MockReleaseRepository) FindByID(release_id int) (*domain.ReleaseModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(release_id)
	}
	return &domain.ReleaseModel{}, nil
}

// Create calls CreateFunc, or returns an empty release.
func (m *MockReleaseRepository) Create(software_id int, attributes domain.Attributes) (*domain.ReleaseModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(software_id, attributes)
	}
	return &domain.ReleaseModel{}, nil
}

// Update calls UpdateFunc, or returns an empty release.
func (m *MockReleaseRepository) Update(release_id int, attributes domain.Attributes) (*domain.ReleaseModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(release_id, attributes)
	}
	return &domain.ReleaseModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockReleaseRepository) Destroy(release_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(release_id)
	}
	return nil
}
