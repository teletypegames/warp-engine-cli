package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockProjectRepository mocks [domain.ProjectRepositoryInterface].
type MockProjectRepository struct {
	// MarkerFunc replaces Marker.
	MarkerFunc func(directory string) (*domain.ProjectMarkerModel, error)
}

// Marker calls MarkerFunc, or returns an empty marker.
func (m *MockProjectRepository) Marker(directory string) (*domain.ProjectMarkerModel, error) {
	if m.MarkerFunc != nil {
		return m.MarkerFunc(directory)
	}
	return &domain.ProjectMarkerModel{}, nil
}
