package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSyncPeerRepository mocks [domain.SyncPeerRepositoryInterface].
type MockSyncPeerRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncPeerModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(sync_peer_id int) (*domain.SyncPeerModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.SyncPeerModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(sync_peer_id int, attributes domain.Attributes) (*domain.SyncPeerModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(sync_peer_id int) error
	// CheckFunc replaces Check.
	CheckFunc func(sync_peer_id int) (*domain.SyncCheckModel, error)
	// InboundTokenFunc replaces InboundToken.
	InboundTokenFunc func(sync_peer_id int) (*domain.SyncInboundTokenModel, error)
	// PushFunc replaces Push.
	PushFunc func(sync_peer_id int, attributes domain.Attributes) (*domain.SyncPushModel, error)
	// RunFunc replaces Run.
	RunFunc func(attributes domain.Attributes) (*domain.SyncPushModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockSyncPeerRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncPeerModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SyncPeerModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty record.
func (m *MockSyncPeerRepository) FindByID(sync_peer_id int) (*domain.SyncPeerModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(sync_peer_id)
	}
	return &domain.SyncPeerModel{}, nil
}

// Create calls CreateFunc, or returns an empty record.
func (m *MockSyncPeerRepository) Create(attributes domain.Attributes) (*domain.SyncPeerModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.SyncPeerModel{}, nil
}

// Update calls UpdateFunc, or returns an empty record.
func (m *MockSyncPeerRepository) Update(sync_peer_id int, attributes domain.Attributes) (*domain.SyncPeerModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(sync_peer_id, attributes)
	}
	return &domain.SyncPeerModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSyncPeerRepository) Destroy(sync_peer_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(sync_peer_id)
	}
	return nil
}

// Check calls CheckFunc, or returns an empty record.
func (m *MockSyncPeerRepository) Check(sync_peer_id int) (*domain.SyncCheckModel, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(sync_peer_id)
	}
	return &domain.SyncCheckModel{}, nil
}

// InboundToken calls InboundTokenFunc, or returns an empty record.
func (m *MockSyncPeerRepository) InboundToken(sync_peer_id int) (*domain.SyncInboundTokenModel, error) {
	if m.InboundTokenFunc != nil {
		return m.InboundTokenFunc(sync_peer_id)
	}
	return &domain.SyncInboundTokenModel{}, nil
}

// Push calls PushFunc, or returns an empty record.
func (m *MockSyncPeerRepository) Push(sync_peer_id int, attributes domain.Attributes) (*domain.SyncPushModel, error) {
	if m.PushFunc != nil {
		return m.PushFunc(sync_peer_id, attributes)
	}
	return &domain.SyncPushModel{}, nil
}

// Run calls RunFunc, or returns an empty record.
func (m *MockSyncPeerRepository) Run(attributes domain.Attributes) (*domain.SyncPushModel, error) {
	if m.RunFunc != nil {
		return m.RunFunc(attributes)
	}
	return &domain.SyncPushModel{}, nil
}
