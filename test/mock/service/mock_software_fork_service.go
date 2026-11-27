package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareForkService mocks [domain.SoftwareForkServiceInterface].
type MockSoftwareForkService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareForkModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(software_fork_id int) (*domain.SoftwareForkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(software_fork_id int, attributes domain.Attributes) (*domain.SoftwareForkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(software_fork_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockSoftwareForkService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareForkModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SoftwareForkModel]{}, nil
}

// Show calls ShowFunc, or returns an empty announced fork.
func (m *MockSoftwareForkService) Show(software_fork_id int) (*domain.SoftwareForkModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(software_fork_id)
	}
	return &domain.SoftwareForkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty announced fork.
func (m *MockSoftwareForkService) Update(software_fork_id int, attributes domain.Attributes) (*domain.SoftwareForkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(software_fork_id, attributes)
	}
	return &domain.SoftwareForkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareForkService) Destroy(software_fork_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(software_fork_id)
	}
	return nil
}
