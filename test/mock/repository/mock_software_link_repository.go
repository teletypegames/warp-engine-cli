package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareLinkRepository mocks [domain.SoftwareLinkRepositoryInterface].
type MockSoftwareLinkRepository struct {
	// AllFunc replaces All.
	AllFunc func(software_id int) ([]domain.SoftwareLinkModel, error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(link_id int) (*domain.SoftwareLinkModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(software_id int, attributes domain.Attributes) (*domain.SoftwareLinkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(link_id int, attributes domain.Attributes) (*domain.SoftwareLinkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(link_id int) error
}

// All calls AllFunc, or returns no links.
func (m *MockSoftwareLinkRepository) All(software_id int) ([]domain.SoftwareLinkModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(software_id)
	}
	return nil, nil
}

// FindByID calls FindByIDFunc, or returns an empty link.
func (m *MockSoftwareLinkRepository) FindByID(link_id int) (*domain.SoftwareLinkModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(link_id)
	}
	return &domain.SoftwareLinkModel{}, nil
}

// Create calls CreateFunc, or returns an empty link.
func (m *MockSoftwareLinkRepository) Create(software_id int, attributes domain.Attributes) (*domain.SoftwareLinkModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(software_id, attributes)
	}
	return &domain.SoftwareLinkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty link.
func (m *MockSoftwareLinkRepository) Update(link_id int, attributes domain.Attributes) (*domain.SoftwareLinkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(link_id, attributes)
	}
	return &domain.SoftwareLinkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareLinkRepository) Destroy(link_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(link_id)
	}
	return nil
}
