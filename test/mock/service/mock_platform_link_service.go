package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockPlatformLinkService mocks [domain.PlatformLinkServiceInterface].
type MockPlatformLinkService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PlatformLinkModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(platform_link_id int) (*domain.PlatformLinkModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.PlatformLinkModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(platform_link_id int, attributes domain.Attributes) (*domain.PlatformLinkModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(platform_link_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockPlatformLinkService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PlatformLinkModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.PlatformLinkModel]{}, nil
}

// Show calls ShowFunc, or returns an empty platform link.
func (m *MockPlatformLinkService) Show(platform_link_id int) (*domain.PlatformLinkModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(platform_link_id)
	}
	return &domain.PlatformLinkModel{}, nil
}

// Create calls CreateFunc, or returns an empty platform link.
func (m *MockPlatformLinkService) Create(attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.PlatformLinkModel{}, nil
}

// Update calls UpdateFunc, or returns an empty platform link.
func (m *MockPlatformLinkService) Update(platform_link_id int, attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(platform_link_id, attributes)
	}
	return &domain.PlatformLinkModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockPlatformLinkService) Destroy(platform_link_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(platform_link_id)
	}
	return nil
}
