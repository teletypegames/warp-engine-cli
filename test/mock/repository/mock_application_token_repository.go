package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockApplicationTokenRepository mocks [domain.ApplicationTokenRepositoryInterface].
type MockApplicationTokenRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ApplicationTokenModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(application_token_id int) (*domain.ApplicationTokenModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.ApplicationTokenModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(application_token_id int, attributes domain.Attributes) (*domain.ApplicationTokenModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(application_token_id int) error
	// RotateFunc replaces Rotate.
	RotateFunc func(application_token_id int) (*domain.ApplicationTokenModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockApplicationTokenRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ApplicationTokenModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.ApplicationTokenModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty record.
func (m *MockApplicationTokenRepository) FindByID(application_token_id int) (*domain.ApplicationTokenModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(application_token_id)
	}
	return &domain.ApplicationTokenModel{}, nil
}

// Create calls CreateFunc, or returns an empty record.
func (m *MockApplicationTokenRepository) Create(attributes domain.Attributes) (*domain.ApplicationTokenModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.ApplicationTokenModel{}, nil
}

// Update calls UpdateFunc, or returns an empty record.
func (m *MockApplicationTokenRepository) Update(application_token_id int, attributes domain.Attributes) (*domain.ApplicationTokenModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(application_token_id, attributes)
	}
	return &domain.ApplicationTokenModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockApplicationTokenRepository) Destroy(application_token_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(application_token_id)
	}
	return nil
}

// Rotate calls RotateFunc, or returns an empty record.
func (m *MockApplicationTokenRepository) Rotate(application_token_id int) (*domain.ApplicationTokenModel, error) {
	if m.RotateFunc != nil {
		return m.RotateFunc(application_token_id)
	}
	return &domain.ApplicationTokenModel{}, nil
}
