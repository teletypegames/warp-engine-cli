package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockExternalLinkService mocks [domain.ExternalLinkServiceInterface].
type MockExternalLinkService struct {
	// AllFunc replaces All.
	AllFunc func(reference string) ([]domain.ExternalLinkModel, error)
	// ShowFunc replaces Show.
	ShowFunc func(id string) (*domain.ExternalLinkModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(reference string, attributes domain.Attributes) (*domain.ExternalLinkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(id string, attributes domain.Attributes) (*domain.ExternalLinkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(id string) error
}

// All calls AllFunc, or returns no links.
func (m *MockExternalLinkService) All(reference string) ([]domain.ExternalLinkModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(reference)
	}
	return nil, nil
}

// Show calls ShowFunc, or returns an empty link.
func (m *MockExternalLinkService) Show(id string) (*domain.ExternalLinkModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(id)
	}
	return &domain.ExternalLinkModel{}, nil
}

// Create calls CreateFunc, or returns an empty link.
func (m *MockExternalLinkService) Create(reference string, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(reference, attributes)
	}
	return &domain.ExternalLinkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty link.
func (m *MockExternalLinkService) Update(id string, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, attributes)
	}
	return &domain.ExternalLinkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockExternalLinkService) Destroy(id string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(id)
	}
	return nil
}
