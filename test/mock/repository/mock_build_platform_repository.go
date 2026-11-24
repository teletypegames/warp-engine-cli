package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockBuildPlatformRepository mocks [domain.BuildPlatformRepositoryInterface].
type MockBuildPlatformRepository struct {
	// AllFunc replaces All.
	AllFunc func() ([]domain.BuildPlatformModel, error)
}

// All calls AllFunc, or returns no platforms.
func (m *MockBuildPlatformRepository) All() ([]domain.BuildPlatformModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return nil, nil
}
