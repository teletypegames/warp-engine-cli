package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockForkLinkRepository mocks [domain.ForkLinkRepositoryInterface].
type MockForkLinkRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ForkLinkModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(fork_link_id int) (*domain.ForkLinkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(fork_link_id int, attributes domain.Attributes) (*domain.ForkLinkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(fork_link_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockForkLinkRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ForkLinkModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.ForkLinkModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty fork link.
func (m *MockForkLinkRepository) FindByID(fork_link_id int) (*domain.ForkLinkModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(fork_link_id)
	}
	return &domain.ForkLinkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty fork link.
func (m *MockForkLinkRepository) Update(fork_link_id int, attributes domain.Attributes) (*domain.ForkLinkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(fork_link_id, attributes)
	}
	return &domain.ForkLinkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockForkLinkRepository) Destroy(fork_link_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(fork_link_id)
	}
	return nil
}
