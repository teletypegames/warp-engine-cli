package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockBuildRepository mocks [domain.BuildRepositoryInterface].
type MockBuildRepository struct {
	// UploadFunc replaces Upload.
	UploadFunc func(target domain.BuildTargetModel, file_path string, sha256 string) (*domain.BuildUploadModel, error)
	// PublishFunc replaces Publish.
	PublishFunc func(target domain.BuildTargetModel) (*domain.BuildPublishModel, error)
}

// Upload calls UploadFunc, or returns an empty upload.
func (m *MockBuildRepository) Upload(target domain.BuildTargetModel, file_path string, sha256 string) (*domain.BuildUploadModel, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(target, file_path, sha256)
	}
	return &domain.BuildUploadModel{}, nil
}

// Publish calls PublishFunc, or returns an empty answer.
func (m *MockBuildRepository) Publish(target domain.BuildTargetModel) (*domain.BuildPublishModel, error) {
	if m.PublishFunc != nil {
		return m.PublishFunc(target)
	}
	return &domain.BuildPublishModel{}, nil
}
