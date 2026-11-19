package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockExtraService mocks [domain.ExtraServiceInterface].
type MockExtraService struct {
	// ShowFunc replaces Show.
	ShowFunc func(reference string, association string) (*domain.ExtraModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(reference string, association string, attributes domain.Attributes) (*domain.ExtraModel, error)
}

// Show calls ShowFunc, or returns an empty extra.
func (m *MockExtraService) Show(reference string, association string) (*domain.ExtraModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(reference, association)
	}
	return &domain.ExtraModel{}, nil
}

// Update calls UpdateFunc, or returns an empty extra.
func (m *MockExtraService) Update(reference string, association string, attributes domain.Attributes) (*domain.ExtraModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(reference, association, attributes)
	}
	return &domain.ExtraModel{}, nil
}
