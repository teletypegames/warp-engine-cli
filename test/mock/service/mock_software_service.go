package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareService mocks [domain.SoftwareServiceInterface].
type MockSoftwareService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(reference string, include []string) (*domain.SoftwareModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.SoftwareModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(reference string, attributes domain.Attributes) (*domain.SoftwareModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(reference string) error
	// ScaffoldFunc replaces Scaffold.
	ScaffoldFunc func(attributes domain.Attributes) (*domain.ScaffoldModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockSoftwareService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SoftwareModel]{}, nil
}

// Show calls ShowFunc, or returns an empty software.
func (m *MockSoftwareService) Show(reference string, include []string) (*domain.SoftwareModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(reference, include)
	}
	return &domain.SoftwareModel{}, nil
}

// Create calls CreateFunc, or returns an empty software.
func (m *MockSoftwareService) Create(attributes domain.Attributes) (*domain.SoftwareModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.SoftwareModel{}, nil
}

// Update calls UpdateFunc, or returns an empty software.
func (m *MockSoftwareService) Update(reference string, attributes domain.Attributes) (*domain.SoftwareModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(reference, attributes)
	}
	return &domain.SoftwareModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareService) Destroy(reference string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(reference)
	}
	return nil
}

// Scaffold calls ScaffoldFunc, or returns an empty scaffold.
func (m *MockSoftwareService) Scaffold(attributes domain.Attributes) (*domain.ScaffoldModel, error) {
	if m.ScaffoldFunc != nil {
		return m.ScaffoldFunc(attributes)
	}
	return &domain.ScaffoldModel{}, nil
}
