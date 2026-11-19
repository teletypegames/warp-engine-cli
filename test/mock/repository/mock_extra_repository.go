package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockExtraRepository mocks [domain.ExtraRepositoryInterface].
type MockExtraRepository struct {
	// FindFunc replaces Find.
	FindFunc func(software_id int, association string) (*domain.ExtraModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(software_id int, association string, attributes domain.Attributes) (*domain.ExtraModel, error)
}

// Find calls FindFunc, or returns an empty extra.
func (m *MockExtraRepository) Find(software_id int, association string) (*domain.ExtraModel, error) {
	if m.FindFunc != nil {
		return m.FindFunc(software_id, association)
	}
	return &domain.ExtraModel{}, nil
}

// Update calls UpdateFunc, or returns an empty extra.
func (m *MockExtraRepository) Update(software_id int, association string, attributes domain.Attributes) (*domain.ExtraModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(software_id, association, attributes)
	}
	return &domain.ExtraModel{}, nil
}
