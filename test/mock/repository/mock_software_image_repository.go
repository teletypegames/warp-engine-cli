package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareImageRepository mocks [domain.SoftwareImageRepositoryInterface].
type MockSoftwareImageRepository struct {
	// AllFunc replaces All.
	AllFunc func(software_id int) ([]domain.SoftwareImageModel, error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(image_id int) (*domain.SoftwareImageModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(software_id int, attributes domain.Attributes) (*domain.SoftwareImageModel, error)
	// UploadFunc replaces Upload.
	UploadFunc func(software_id int, file_path string, attributes domain.Attributes) (*domain.SoftwareImageModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(image_id int, attributes domain.Attributes) (*domain.SoftwareImageModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(image_id int) error
}

// All calls AllFunc, or returns no images.
func (m *MockSoftwareImageRepository) All(software_id int) ([]domain.SoftwareImageModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(software_id)
	}
	return nil, nil
}

// FindByID calls FindByIDFunc, or returns an empty image.
func (m *MockSoftwareImageRepository) FindByID(image_id int) (*domain.SoftwareImageModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(image_id)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Create calls CreateFunc, or returns an empty image.
func (m *MockSoftwareImageRepository) Create(software_id int, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(software_id, attributes)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Upload calls UploadFunc, or returns an empty image.
func (m *MockSoftwareImageRepository) Upload(software_id int, file_path string, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(software_id, file_path, attributes)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Update calls UpdateFunc, or returns an empty image.
func (m *MockSoftwareImageRepository) Update(image_id int, attributes domain.Attributes) (*domain.SoftwareImageModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(image_id, attributes)
	}
	return &domain.SoftwareImageModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareImageRepository) Destroy(image_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(image_id)
	}
	return nil
}
