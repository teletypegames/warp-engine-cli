package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockReleaseAssetService mocks [domain.ReleaseAssetServiceInterface].
type MockReleaseAssetService struct {
	// AllFunc replaces All.
	AllFunc func(release string) ([]domain.ReleaseAssetModel, error)
	// ShowFunc replaces Show.
	ShowFunc func(id string) (*domain.ReleaseAssetModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(release string, attributes domain.Attributes) (*domain.ReleaseAssetModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(id string, attributes domain.Attributes) (*domain.ReleaseAssetModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(id string) error
}

// All calls AllFunc, or returns no assets.
func (m *MockReleaseAssetService) All(release string) ([]domain.ReleaseAssetModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(release)
	}
	return nil, nil
}

// Show calls ShowFunc, or returns an empty asset.
func (m *MockReleaseAssetService) Show(id string) (*domain.ReleaseAssetModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(id)
	}
	return &domain.ReleaseAssetModel{}, nil
}

// Create calls CreateFunc, or returns an empty asset.
func (m *MockReleaseAssetService) Create(release string, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(release, attributes)
	}
	return &domain.ReleaseAssetModel{}, nil
}

// Update calls UpdateFunc, or returns an empty asset.
func (m *MockReleaseAssetService) Update(id string, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, attributes)
	}
	return &domain.ReleaseAssetModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockReleaseAssetService) Destroy(id string) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(id)
	}
	return nil
}
