package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockFileRepository mocks [domain.FileRepositoryInterface].
type MockFileRepository struct {
	// AllFunc replaces All.
	AllFunc func(dir string) (*domain.PageModel[domain.FileEntryModel], error)
	// UploadFunc replaces Upload.
	UploadFunc func(dir string, file_path string) (*domain.FilePathModel, error)
	// RenameFunc replaces Rename.
	RenameFunc func(path string, new_name string) (*domain.FilePathModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(path string) error
	// CreateDirectoryFunc replaces CreateDirectory.
	CreateDirectoryFunc func(dir string, name string) (*domain.FilePathModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockFileRepository) All(dir string) (*domain.PageModel[domain.FileEntryModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(dir)
	}
	return &domain.PageModel[domain.FileEntryModel]{}, nil
}

// Upload calls UploadFunc, or returns an empty path.
func (m *MockFileRepository) Upload(dir string, file_path string) (*domain.FilePathModel, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(dir, file_path)
	}
	return &domain.FilePathModel{}, nil
}

// Rename calls RenameFunc, or returns an empty path.
func (m *MockFileRepository) Rename(path string, new_name string) (*domain.FilePathModel, error) {
	if m.RenameFunc != nil {
		return m.RenameFunc(path, new_name)
	}
	return &domain.FilePathModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockFileRepository) Destroy(path string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(path)
	}
	return nil
}

// CreateDirectory calls CreateDirectoryFunc, or returns an empty path.
func (m *MockFileRepository) CreateDirectory(dir string, name string) (*domain.FilePathModel, error) {
	if m.CreateDirectoryFunc != nil {
		return m.CreateDirectoryFunc(dir, name)
	}
	return &domain.FilePathModel{}, nil
}
