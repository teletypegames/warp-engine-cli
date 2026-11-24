package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockMakefileService mocks [domain.MakefileServiceInterface].
type MockMakefileService struct {
	// CheckFunc replaces Check.
	CheckFunc func(directory string, platform string, app string) (*domain.MakefileReportModel, error)
}

// Check calls CheckFunc, or returns an empty report.
func (m *MockMakefileService) Check(directory string, platform string, app string) (*domain.MakefileReportModel, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(directory, platform, app)
	}
	return &domain.MakefileReportModel{}, nil
}
