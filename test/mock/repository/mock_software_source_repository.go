package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareSourceRepository mocks [domain.SoftwareSourceRepositoryInterface].
type MockSoftwareSourceRepository struct {
	// FindFunc replaces Find.
	FindFunc func(software_id int) (*domain.SoftwareSourceModel, error)
	// SaveFunc replaces Save.
	SaveFunc func(software_id int, attributes domain.Attributes) (*domain.SoftwareSourceModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(software_id int) error
}

// Find calls FindFunc, or returns an empty source.
func (m *MockSoftwareSourceRepository) Find(software_id int) (*domain.SoftwareSourceModel, error) {
	if m.FindFunc != nil {
		return m.FindFunc(software_id)
	}
	return &domain.SoftwareSourceModel{}, nil
}

// Save calls SaveFunc, or returns an empty source.
func (m *MockSoftwareSourceRepository) Save(software_id int, attributes domain.Attributes) (*domain.SoftwareSourceModel, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(software_id, attributes)
	}
	return &domain.SoftwareSourceModel{}, nil
}

// Destroy calls DestroyFunc, or does nothing.
func (m *MockSoftwareSourceRepository) Destroy(software_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(software_id)
	}
	return nil
}
