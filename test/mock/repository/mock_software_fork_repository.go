package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareForkRepository mocks [domain.SoftwareForkRepositoryInterface].
type MockSoftwareForkRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareForkModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(software_fork_id int) (*domain.SoftwareForkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(software_fork_id int, attributes domain.Attributes) (*domain.SoftwareForkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(software_fork_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockSoftwareForkRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareForkModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SoftwareForkModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty announced fork.
func (m *MockSoftwareForkRepository) FindByID(software_fork_id int) (*domain.SoftwareForkModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(software_fork_id)
	}
	return &domain.SoftwareForkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty announced fork.
func (m *MockSoftwareForkRepository) Update(software_fork_id int, attributes domain.Attributes) (*domain.SoftwareForkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(software_fork_id, attributes)
	}
	return &domain.SoftwareForkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareForkRepository) Destroy(software_fork_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(software_fork_id)
	}
	return nil
}
