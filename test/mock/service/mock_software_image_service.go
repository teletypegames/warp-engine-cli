package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareImageService mocks [domain.SoftwareImageServiceInterface].
type MockSoftwareImageService struct {
	// AllFunc replaces All.
	AllFunc func(reference string) ([]domain.SoftwareImageModel, error)
	// ShowFunc replaces Show.
	ShowFunc func(id string) (*domain.SoftwareImageModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(reference string, attributes domain.Attributes) (*domain.SoftwareImageModel, error)
	// UploadFunc replaces Upload.
	UploadFunc func(reference string, file_path string, attributes domain.Attributes) (*domain.SoftwareImageModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(id string, attributes domain.Attributes) (*domain.SoftwareImageModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(id string) error
}

// All calls AllFunc, or returns no images.
func (m *MockSoftwareImageService) All(reference string) ([]domain.SoftwareImageModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(reference)
	}
	return nil, nil
}

// Show calls ShowFunc, or returns an empty image.
func (m *MockSoftwareImageService) Show(id string) (*domain.SoftwareImageModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(id)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Create calls CreateFunc, or returns an empty image.
func (m *MockSoftwareImageService) Create(reference string, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(reference, attributes)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Upload calls UploadFunc, or returns an empty image.
func (m *MockSoftwareImageService) Upload(reference string, file_path string, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(reference, file_path, attributes)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Update calls UpdateFunc, or returns an empty image.
func (m *MockSoftwareImageService) Update(id string, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, attributes)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareImageService) Destroy(id string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(id)
	}
	return nil
}
