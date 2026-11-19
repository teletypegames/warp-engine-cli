package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockReleaseAssetRepository mocks [domain.ReleaseAssetRepositoryInterface].
type MockReleaseAssetRepository struct {
	// AllFunc replaces All.
	AllFunc func(release_id int) ([]domain.ReleaseAssetModel, error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(asset_id int) (*domain.ReleaseAssetModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(release_id int, attributes domain.Attributes) (*domain.ReleaseAssetModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(asset_id int, attributes domain.Attributes) (*domain.ReleaseAssetModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(asset_id int) error
}

// All calls AllFunc, or returns no assets.
func (m *MockReleaseAssetRepository) All(release_id int) ([]domain.ReleaseAssetModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc(release_id)
	}
	return nil, nil
}

// FindByID calls FindByIDFunc, or returns an empty asset.
func (m *MockReleaseAssetRepository) FindByID(asset_id int) (*domain.ReleaseAssetModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(asset_id)
	}
	return &domain.ReleaseAssetModel{}, nil
}

// Create calls CreateFunc, or returns an empty asset.
func (m *MockReleaseAssetRepository) Create(release_id int, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(release_id, attributes)
	}
	return &domain.ReleaseAssetModel{}, nil
}

// Update calls UpdateFunc, or returns an empty asset.
func (m *MockReleaseAssetRepository) Update(asset_id int, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(asset_id, attributes)
	}
	return &domain.ReleaseAssetModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockReleaseAssetRepository) Destroy(asset_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(asset_id)
	}
	return nil
}
