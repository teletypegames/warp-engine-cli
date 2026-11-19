package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSyncPairingService mocks [domain.SyncPairingServiceInterface].
type MockSyncPairingService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncPairingModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(sync_pairing_id int) (*domain.SyncPairingModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.SyncPairingModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(sync_pairing_id int, attributes domain.Attributes) (*domain.SyncPairingModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(sync_pairing_id int) error
	// TokenFunc replaces Token.
	TokenFunc func(sync_pairing_id int) (*domain.SyncPairingModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockSyncPairingService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncPairingModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SyncPairingModel]{}, nil
}

// Show calls ShowFunc, or returns an empty record.
func (m *MockSyncPairingService) Show(sync_pairing_id int) (*domain.SyncPairingModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(sync_pairing_id)
	}
	return &domain.SyncPairingModel{}, nil
}

// Create calls CreateFunc, or returns an empty record.
func (m *MockSyncPairingService) Create(attributes domain.Attributes) (*domain.SyncPairingModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.SyncPairingModel{}, nil
}

// Update calls UpdateFunc, or returns an empty record.
func (m *MockSyncPairingService) Update(sync_pairing_id int, attributes domain.Attributes) (*domain.SyncPairingModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(sync_pairing_id, attributes)
	}
	return &domain.SyncPairingModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSyncPairingService) Destroy(sync_pairing_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(sync_pairing_id)
	}
	return nil
}

// Token calls TokenFunc, or returns an empty record.
func (m *MockSyncPairingService) Token(sync_pairing_id int) (*domain.SyncPairingModel, error) {
	if m.TokenFunc != nil {
		return m.TokenFunc(sync_pairing_id)
	}
	return &domain.SyncPairingModel{}, nil
}
