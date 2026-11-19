package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockProfileService mocks [domain.ProfileServiceInterface].
type MockProfileService struct {
	// AllFunc replaces All.
	AllFunc func() (*[]domain.ProfileModel, string, error)
	// ShowFunc replaces Show.
	ShowFunc func(name string) (*domain.ProfileModel, error)
	// SaveFunc replaces Save.
	SaveFunc func(profile domain.ProfileModel, current bool) error
	// UseFunc replaces Use.
	UseFunc func(name string) error
	// DestroyFunc replaces Destroy.
	DestroyFunc func(name string) error
}

// All calls AllFunc, or returns no profiles.
func (m *MockProfileService) All() (*[]domain.ProfileModel, string, error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return &[]domain.ProfileModel{}, "", nil
}

// Show calls ShowFunc, or returns a profile named name.
func (m *MockProfileService) Show(name string) (*domain.ProfileModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(name)
	}
	return &domain.ProfileModel{Name: name}, nil
}

// Save calls SaveFunc, or returns nil.
func (m *MockProfileService) Save(profile domain.ProfileModel, current bool) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(profile, current)
	}
	return nil
}

// Use calls UseFunc, or returns nil.
func (m *MockProfileService) Use(name string) error {
	if m.UseFunc != nil {
		return m.UseFunc(name)
	}
	return nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockProfileService) Destroy(name string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(name)
	}
	return nil
}
