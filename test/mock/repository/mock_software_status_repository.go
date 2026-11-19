package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareStatusRepository mocks [domain.SoftwareStatusRepositoryInterface].
type MockSoftwareStatusRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareStatusModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(status_id int) (*domain.SoftwareStatusModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.SoftwareStatusModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(status_id int, attributes domain.Attributes) (*domain.SoftwareStatusModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(status_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockSoftwareStatusRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareStatusModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SoftwareStatusModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty status.
func (m *MockSoftwareStatusRepository) FindByID(status_id int) (*domain.SoftwareStatusModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(status_id)
	}
	return &domain.SoftwareStatusModel{}, nil
}

// Create calls CreateFunc, or returns an empty status.
func (m *MockSoftwareStatusRepository) Create(attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.SoftwareStatusModel{}, nil
}

// Update calls UpdateFunc, or returns an empty status.
func (m *MockSoftwareStatusRepository) Update(status_id int, attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(status_id, attributes)
	}
	return &domain.SoftwareStatusModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareStatusRepository) Destroy(status_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(status_id)
	}
	return nil
}
