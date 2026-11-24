package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockBuildService mocks [domain.BuildServiceInterface].
type MockBuildService struct {
	// UploadFunc replaces Upload.
	UploadFunc func(target domain.BuildTargetModel, file_paths []string) ([]domain.BuildUploadModel, error)
	// PublishFunc replaces Publish.
	PublishFunc func(target domain.BuildTargetModel) (*domain.BuildPublishModel, error)
}

// Upload calls UploadFunc, or returns no uploads.
func (m *MockBuildService) Upload(target domain.BuildTargetModel, file_paths []string) ([]domain.BuildUploadModel, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(target, file_paths)
	}
	return nil, nil
}

// Publish calls PublishFunc, or returns an empty answer.
func (m *MockBuildService) Publish(target domain.BuildTargetModel) (*domain.BuildPublishModel, error) {
	if m.PublishFunc != nil {
		return m.PublishFunc(target)
	}
	return &domain.BuildPublishModel{}, nil
}
