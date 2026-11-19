package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SoftwareCategoryRepositoryInterface is the Admin API's software categories resource.
type SoftwareCategoryRepositoryInterface interface {
	// All lists the software categories matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareCategoryModel], error)
	// FindByID returns the category category_id.
	FindByID(category_id int) (*SoftwareCategoryModel, error)
	// Create creates a category from attributes.
	Create(attributes Attributes) (*SoftwareCategoryModel, error)
	// Update changes the attributes of the category category_id.
	Update(category_id int, attributes Attributes) (*SoftwareCategoryModel, error)
	// Destroy deletes the category category_id.
	Destroy(category_id int) error
}

// SoftwareCategoryRepository is the [SoftwareCategoryRepositoryInterface] over the Admin API.
type SoftwareCategoryRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /software_categories.
func (r SoftwareCategoryRepository) All(filter Filter, query ListQuery) (*PageModel[SoftwareCategoryModel], error) {
	return list[SoftwareCategoryModel](r.Client, "/software_categories", filter, query)
}

// FindByID gets GET /software_categories/{id}.
func (r SoftwareCategoryRepository) FindByID(category_id int) (*SoftwareCategoryModel, error) {
	return api.Get[SoftwareCategoryModel](r.Client, "/software_categories/"+strconv.Itoa(category_id), nil)
}

// Create posts POST /software_categories.
func (r SoftwareCategoryRepository) Create(attributes Attributes) (*SoftwareCategoryModel, error) {
	return api.Send[SoftwareCategoryModel](r.Client, http.MethodPost, "/software_categories", attributes)
}

// Update sends PATCH /software_categories/{id}.
func (r SoftwareCategoryRepository) Update(category_id int, attributes Attributes) (*SoftwareCategoryModel, error) {
	return api.Send[SoftwareCategoryModel](r.Client, http.MethodPatch, "/software_categories/"+strconv.Itoa(category_id), attributes)
}

// Destroy sends DELETE /software_categories/{id}.
func (r SoftwareCategoryRepository) Destroy(category_id int) error {
	return api.Delete(r.Client, "/software_categories/"+strconv.Itoa(category_id), nil)
}
