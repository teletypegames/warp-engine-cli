package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockExternalLinkRepository mocks [domain.ExternalLinkRepositoryInterface].
type MockExternalLinkRepository struct {
	// AllFunc replaces All.
	AllFunc func(software_id int) ([]domain.ExternalLinkModel, error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(link_id int) (*domain.ExternalLinkModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(software_id int, attributes domain.Attributes) (*domain.ExternalLinkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(link_id int, attributes domain.Attributes) (*domain.ExternalLinkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(link_id int) error
}

// All calls AllFunc, or returns no links.
func (m *MockExternalLinkRepository) All(software_id int) ([]domain.ExternalLinkModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(software_id)
	}
	return nil, nil
}

// FindByID calls FindByIDFunc, or returns an empty link.
func (m *MockExternalLinkRepository) FindByID(link_id int) (*domain.ExternalLinkModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(link_id)
	}
	return &domain.ExternalLinkModel{}, nil
}

// Create calls CreateFunc, or returns an empty link.
func (m *MockExternalLinkRepository) Create(software_id int, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(software_id, attributes)
	}
	return &domain.ExternalLinkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty link.
func (m *MockExternalLinkRepository) Update(link_id int, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(link_id, attributes)
	}
	return &domain.ExternalLinkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockExternalLinkRepository) Destroy(link_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(link_id)
	}
	return nil
}
