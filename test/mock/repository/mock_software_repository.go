package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareRepository mocks [domain.SoftwareRepositoryInterface].
type MockSoftwareRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(software_id int, include []string) (*domain.SoftwareModel, error)
	// FindByNameFunc replaces FindByName.
	FindByNameFunc func(name string) (*domain.SoftwareModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.SoftwareModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(software_id int, attributes domain.Attributes) (*domain.SoftwareModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(software_id int) error
	// ScaffoldFunc replaces Scaffold.
	ScaffoldFunc func(attributes domain.Attributes) (*domain.ScaffoldModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockSoftwareRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SoftwareModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty software.
func (m *MockSoftwareRepository) FindByID(software_id int, include []string) (*domain.SoftwareModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(software_id, include)
	}
	return &domain.SoftwareModel{}, nil
}

// FindByName calls FindByNameFunc, or returns an empty software.
func (m *MockSoftwareRepository) FindByName(name string) (*domain.SoftwareModel, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(name)
	}
	return &domain.SoftwareModel{}, nil
}

// Create calls CreateFunc, or returns an empty software.
func (m *MockSoftwareRepository) Create(attributes domain.Attributes) (*domain.SoftwareModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.SoftwareModel{}, nil
}

// Update calls UpdateFunc, or returns an empty software.
func (m *MockSoftwareRepository) Update(software_id int, attributes domain.Attributes) (*domain.SoftwareModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(software_id, attributes)
	}
	return &domain.SoftwareModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareRepository) Destroy(software_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(software_id)
	}
	return nil
}

// Scaffold calls ScaffoldFunc, or returns an empty scaffold.
func (m *MockSoftwareRepository) Scaffold(attributes domain.Attributes) (*domain.ScaffoldModel, error) {
	if m.ScaffoldFunc != nil {
		return m.ScaffoldFunc(attributes)
	}
	return &domain.ScaffoldModel{}, nil
}
