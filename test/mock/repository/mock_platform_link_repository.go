package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockPlatformLinkRepository mocks [domain.PlatformLinkRepositoryInterface].
type MockPlatformLinkRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PlatformLinkModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(platform_link_id int) (*domain.PlatformLinkModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.PlatformLinkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(platform_link_id int, attributes domain.Attributes) (*domain.PlatformLinkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(platform_link_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockPlatformLinkRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PlatformLinkModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.PlatformLinkModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty platform link.
func (m *MockPlatformLinkRepository) FindByID(platform_link_id int) (*domain.PlatformLinkModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(platform_link_id)
	}
	return &domain.PlatformLinkModel{}, nil
}

// Create calls CreateFunc, or returns an empty platform link.
func (m *MockPlatformLinkRepository) Create(attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.PlatformLinkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty platform link.
func (m *MockPlatformLinkRepository) Update(platform_link_id int, attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(platform_link_id, attributes)
	}
	return &domain.PlatformLinkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockPlatformLinkRepository) Destroy(platform_link_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(platform_link_id)
	}
	return nil
}
