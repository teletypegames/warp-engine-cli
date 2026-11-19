package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSyncEventRepository mocks [domain.SyncEventRepositoryInterface].
type MockSyncEventRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncEventModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(sync_event_id int) (*domain.SyncEventModel, error)
	// AllOfSoftwareFunc replaces AllOfSoftware.
	AllOfSoftwareFunc func(software_id int, direction string) ([]domain.SyncEventModel, error)
	// StatusFunc replaces Status.
	StatusFunc func(software_id int) ([]domain.SyncStatusModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockSyncEventRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncEventModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SyncEventModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty record.
func (m *MockSyncEventRepository) FindByID(sync_event_id int) (*domain.SyncEventModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(sync_event_id)
	}
	return &domain.SyncEventModel{}, nil
}

// AllOfSoftware calls AllOfSoftwareFunc, or returns no events.
func (m *MockSyncEventRepository) AllOfSoftware(software_id int, direction string) ([]domain.SyncEventModel, error) {
	if m.AllOfSoftwareFunc != nil {
		return m.AllOfSoftwareFunc(software_id, direction)
	}
	return nil, nil
}

// Status calls StatusFunc, or returns no states.
func (m *MockSyncEventRepository) Status(software_id int) ([]domain.SyncStatusModel, error) {
	if m.StatusFunc != nil {
		return m.StatusFunc(software_id)
	}
	return nil, nil
}
