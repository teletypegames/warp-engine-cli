package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockSoftwareCategoryRepository mocks [domain.SoftwareCategoryRepositoryInterface].
type MockSoftwareCategoryRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareCategoryModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(category_id int) (*domain.SoftwareCategoryModel, error)
	// CreateFunc replaces Create.
	CreateFunc func(attributes domain.Attributes) (*domain.SoftwareCategoryModel, error)
	// UpdateFunc replaces Update.
	UpdateFunc func(category_id int, attributes domain.Attributes) (*domain.SoftwareCategoryModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(category_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockSoftwareCategoryRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareCategoryModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.SoftwareCategoryModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty category.
func (m *MockSoftwareCategoryRepository) FindByID(category_id int) (*domain.SoftwareCategoryModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(category_id)
	}
	return &domain.SoftwareCategoryModel{}, nil
}

// Create calls CreateFunc, or returns an empty category.
func (m *MockSoftwareCategoryRepository) Create(attributes domain.Attributes) (*domain.SoftwareCategoryModel, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(attributes)
	}
	return &domain.SoftwareCategoryModel{}, nil
}

// Update calls UpdateFunc, or returns an empty category.
func (m *MockSoftwareCategoryRepository) Update(category_id int, attributes domain.Attributes) (*domain.SoftwareCategoryModel, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(category_id, attributes)
	}
	return &domain.SoftwareCategoryModel{}, nil
}

// Destroy calls DestroyFunc, or returns nil.
func (m *MockSoftwareCategoryRepository) Destroy(category_id int) error {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(category_id)
	}
	return nil
}
