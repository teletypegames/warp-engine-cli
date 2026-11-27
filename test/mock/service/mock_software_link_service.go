package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareLinkService mocks [domain.SoftwareLinkServiceInterface].
type MockSoftwareLinkService struct {
	// AllFunc replaces All.
	AllFunc func(reference string) ([]domain.SoftwareLinkModel, error)
	// ShowFunc replaces Show.
	ShowFunc func(id string) (*domain.SoftwareLinkModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(reference string, attributes domain.Attributes) (*domain.SoftwareLinkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(id string, attributes domain.Attributes) (*domain.SoftwareLinkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(id string) error
}

// All calls AllFunc, or returns no links.
func (m *MockSoftwareLinkService) All(reference string) ([]domain.SoftwareLinkModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(reference)
	}
	return nil, nil
}

// Show calls ShowFunc, or returns an empty link.
func (m *MockSoftwareLinkService) Show(id string) (*domain.SoftwareLinkModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(id)
	}
	return &domain.SoftwareLinkModel{}, nil
}

// Create calls CreateFunc, or returns an empty link.
func (m *MockSoftwareLinkService) Create(reference string, attributes domain.Attributes) (*domain.SoftwareLinkModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(reference, attributes)
	}
	return &domain.SoftwareLinkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty link.
func (m *MockSoftwareLinkService) Update(id string, attributes domain.Attributes) (*domain.SoftwareLinkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, attributes)
	}
	return &domain.SoftwareLinkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareLinkService) Destroy(id string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(id)
	}
	return nil
}
