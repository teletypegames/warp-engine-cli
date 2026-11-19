package domain

import "errors"

// SoftwareCategoryServiceInterface manages the software categories, the vocabulary of a software's Category.
type SoftwareCategoryServiceInterface interface {
	// All lists the software categories matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareCategoryModel], error)
	// Show returns the category category_id.
	Show(category_id int) (*SoftwareCategoryModel, error)
	// Create creates a category.
	Create(attributes Attributes) (*SoftwareCategoryModel, error)
	// Update changes the category category_id.
	Update(category_id int, attributes Attributes) (*SoftwareCategoryModel, error)
	// Destroy deletes the category category_id.
	Destroy(category_id int) error
}

// SoftwareCategoryService is the [SoftwareCategoryServiceInterface] over a [SoftwareCategoryRepositoryInterface].
type SoftwareCategoryService struct {
	// SoftwareCategoryRepository talks to the software categories resource.
	SoftwareCategoryRepository SoftwareCategoryRepositoryInterface
}

// All lists the software categories.
func (s SoftwareCategoryService) All(filter Filter, query ListQuery) (*PageModel[SoftwareCategoryModel], error) {
	return s.SoftwareCategoryRepository.All(filter, query)
}

// Show gets the category category_id.
func (s SoftwareCategoryService) Show(category_id int) (*SoftwareCategoryModel, error) {
	return s.SoftwareCategoryRepository.FindByID(category_id)
}

// Create checks that attributes hold key and label — the server would
// refuse without them — and creates the category.
func (s SoftwareCategoryService) Create(attributes Attributes) (*SoftwareCategoryModel, error) {
	for _, key := range []string{"key", "label"} {
		if !attributes.Has(key) || attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("a category needs a " + key)
		}
	}
	return s.SoftwareCategoryRepository.Create(attributes)
}

// Update changes the category category_id; empty attributes are
// [ErrNothingToChange], since there would be nothing to send.
func (s SoftwareCategoryService) Update(category_id int, attributes Attributes) (*SoftwareCategoryModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.SoftwareCategoryRepository.Update(category_id, attributes)
}

// Destroy deletes the category category_id.
func (s SoftwareCategoryService) Destroy(category_id int) error {
	return s.SoftwareCategoryRepository.Destroy(category_id)
}
