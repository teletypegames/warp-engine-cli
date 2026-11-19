package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SoftwareStatusRepositoryInterface is the Admin API's software statuses resource.
type SoftwareStatusRepositoryInterface interface {
	// All lists the software statuses matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareStatusModel], error)
	// FindByID returns the status status_id.
	FindByID(status_id int) (*SoftwareStatusModel, error)
	// Create creates a status from attributes.
	Create(attributes Attributes) (*SoftwareStatusModel, error)
	// Update changes the attributes of the status status_id.
	Update(status_id int, attributes Attributes) (*SoftwareStatusModel, error)
	// Destroy deletes the status status_id.
	Destroy(status_id int) error
}

// SoftwareStatusRepository is the [SoftwareStatusRepositoryInterface] over the Admin API.
type SoftwareStatusRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /software_statuses.
func (r SoftwareStatusRepository) All(filter Filter, query ListQuery) (*PageModel[SoftwareStatusModel], error) {
	return list[SoftwareStatusModel](r.Client, "/software_statuses", filter, query)
}

// FindByID gets GET /software_statuses/{id}.
func (r SoftwareStatusRepository) FindByID(status_id int) (*SoftwareStatusModel, error) {
	return api.Get[SoftwareStatusModel](r.Client, "/software_statuses/"+strconv.Itoa(status_id), nil)
}

// Create posts POST /software_statuses.
func (r SoftwareStatusRepository) Create(attributes Attributes) (*SoftwareStatusModel, error) {
	return api.Send[SoftwareStatusModel](r.Client, http.MethodPost, "/software_statuses", attributes)
}

// Update sends PATCH /software_statuses/{id}.
func (r SoftwareStatusRepository) Update(status_id int, attributes Attributes) (*SoftwareStatusModel, error) {
	return api.Send[SoftwareStatusModel](r.Client, http.MethodPatch, "/software_statuses/"+strconv.Itoa(status_id), attributes)
}

// Destroy sends DELETE /software_statuses/{id}.
func (r SoftwareStatusRepository) Destroy(status_id int) error {
	return api.Delete(r.Client, "/software_statuses/"+strconv.Itoa(status_id), nil)
}
