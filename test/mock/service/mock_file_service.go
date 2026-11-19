package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockFileService mocks [domain.FileServiceInterface].
type MockFileService struct {
	// AllFunc replaces All.
	AllFunc func(dir string) (*domain.PageModel[domain.FileEntryModel], error)
	// UploadFunc replaces Upload.
	UploadFunc func(file_path string, dir string) (*domain.FilePathModel, error)
	// RenameFunc replaces Rename.
	RenameFunc func(file_path string, new_name string) (*domain.FilePathModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(file_path string) error
	// CreateDirectoryFunc replaces CreateDirectory.
	CreateDirectoryFunc func(dir_path string) (*domain.FilePathModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockFileService) All(dir string) (*domain.PageModel[domain.FileEntryModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(dir)
	}
	return &domain.PageModel[domain.FileEntryModel]{}, nil
}

// Upload calls UploadFunc, or returns an empty path.
func (m *MockFileService) Upload(file_path string, dir string) (*domain.FilePathModel, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(file_path, dir)
	}
	return &domain.FilePathModel{}, nil
}

// Rename calls RenameFunc, or returns an empty path.
func (m *MockFileService) Rename(file_path string, new_name string) (*domain.FilePathModel, error) {
	if m.RenameFunc != nil {
		return m.RenameFunc(file_path, new_name)
	}
	return &domain.FilePathModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockFileService) Destroy(file_path string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(file_path)
	}
	return nil
}

// CreateDirectory calls CreateDirectoryFunc, or returns an empty path.
func (m *MockFileService) CreateDirectory(dir_path string) (*domain.FilePathModel, error) {
	if m.CreateDirectoryFunc != nil {
		return m.CreateDirectoryFunc(dir_path)
	}
	return &domain.FilePathModel{}, nil
}
