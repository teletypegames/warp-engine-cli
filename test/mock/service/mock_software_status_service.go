package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareStatusService mocks [domain.SoftwareStatusServiceInterface].
type MockSoftwareStatusService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareStatusModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(status_id int) (*domain.SoftwareStatusModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.SoftwareStatusModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(status_id int, attributes domain.Attributes) (*domain.SoftwareStatusModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(status_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockSoftwareStatusService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareStatusModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SoftwareStatusModel]{}, nil
}

// Show calls ShowFunc, or returns an empty status.
func (m *MockSoftwareStatusService) Show(status_id int) (*domain.SoftwareStatusModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(status_id)
	}
	return &domain.SoftwareStatusModel{}, nil
}

// Create calls CreateFunc, or returns an empty status.
func (m *MockSoftwareStatusService) Create(attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.SoftwareStatusModel{}, nil
}

// Update calls UpdateFunc, or returns an empty status.
func (m *MockSoftwareStatusService) Update(status_id int, attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(status_id, attributes)
	}
	return &domain.SoftwareStatusModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareStatusService) Destroy(status_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(status_id)
	}
	return nil
}
