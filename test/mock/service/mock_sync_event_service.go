package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSyncEventService mocks [domain.SyncEventServiceInterface].
type MockSyncEventService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncEventModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(sync_event_id int) (*domain.SyncEventModel, error)
	// SoftwareFunc replaces Software.
	SoftwareFunc func(reference string, direction string) ([]domain.SyncEventModel, error)
	// StatusFunc replaces Status.
	StatusFunc func(reference string) ([]domain.SyncStatusModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockSyncEventService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncEventModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SyncEventModel]{}, nil
}

// Show calls ShowFunc, or returns an empty record.
func (m *MockSyncEventService) Show(sync_event_id int) (*domain.SyncEventModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(sync_event_id)
	}
	return &domain.SyncEventModel{}, nil
}

// Software calls SoftwareFunc, or returns no events.
func (m *MockSyncEventService) Software(reference string, direction string) ([]domain.SyncEventModel, error) {
	if m.SoftwareFunc != nil {
		return m.SoftwareFunc(reference, direction)
	}
	return nil, nil
}

// Status calls StatusFunc, or returns no states.
func (m *MockSyncEventService) Status(reference string) ([]domain.SyncStatusModel, error) {
	if m.StatusFunc != nil {
		return m.StatusFunc(reference)
	}
	return nil, nil
}
