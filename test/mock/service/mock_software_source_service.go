package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareSourceService mocks [domain.SoftwareSourceServiceInterface].
type MockSoftwareSourceService struct {
	// ShowFunc replaces Show.
	ShowFunc func(reference string) (*domain.SoftwareSourceModel, error)
	// SetFunc replaces Set.
	SetFunc func(reference string, attributes domain.Attributes) (*domain.SoftwareSourceModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(reference string) error
}

// Show calls ShowFunc, or returns an empty source.
func (m *MockSoftwareSourceService) Show(reference string) (*domain.SoftwareSourceModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(reference)
	}
	return &domain.SoftwareSourceModel{}, nil
}

// Set calls SetFunc, or returns an empty source.
func (m *MockSoftwareSourceService) Set(reference string, attributes domain.Attributes) (*domain.SoftwareSourceModel, error) {
	if m.SetFunc != nil {
		return m.SetFunc(reference, attributes)
	}
	return &domain.SoftwareSourceModel{}, nil
}

// Destroy calls DestroyFunc, or does nothing.
func (m *MockSoftwareSourceService) Destroy(reference string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(reference)
	}
	return nil
}
