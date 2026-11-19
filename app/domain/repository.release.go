package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// ReleaseRepositoryInterface is the Admin API's releases resource.
type ReleaseRepositoryInterface interface {
	// All lists the releases of the software software_id matching filter
	// (commitSha, channel), paged by query.
	All(software_id int, filter Filter, query ListQuery) (*PageModel[ReleaseModel], error)
	// FindByID returns the release release_id with its assets.
	FindByID(release_id int) (*ReleaseModel, error)
	// Create adds a release to the software software_id from attributes.
	Create(software_id int, attributes Attributes) (*ReleaseModel, error)
	// Update changes the attributes of the release release_id.
	Update(release_id int, attributes Attributes) (*ReleaseModel, error)
	// Destroy deletes the release release_id with its assets.
	Destroy(release_id int) error
}

// ReleaseRepository is the [ReleaseRepositoryInterface] over the Admin API.
type ReleaseRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /softwares/{softwareId}/releases.
func (r ReleaseRepository) All(software_id int, filter Filter, query ListQuery) (*PageModel[ReleaseModel], error) {
	return list[ReleaseModel](r.Client, "/softwares/"+strconv.Itoa(software_id)+"/releases", filter, query)
}

// FindByID gets GET /releases/{id}.
func (r ReleaseRepository) FindByID(release_id int) (*ReleaseModel, error) {
	return api.Get[ReleaseModel](r.Client, "/releases/"+strconv.Itoa(release_id), nil)
}

// Create posts POST /softwares/{softwareId}/releases.
func (r ReleaseRepository) Create(software_id int, attributes Attributes) (*ReleaseModel, error) {
	return api.Send[ReleaseModel](r.Client, http.MethodPost, "/softwares/"+strconv.Itoa(software_id)+"/releases", attributes)
}

// Update sends PATCH /releases/{id}.
func (r ReleaseRepository) Update(release_id int, attributes Attributes) (*ReleaseModel, error) {
	return api.Send[ReleaseModel](r.Client, http.MethodPatch, "/releases/"+strconv.Itoa(release_id), attributes)
}

// Destroy sends DELETE /releases/{id}.
func (r ReleaseRepository) Destroy(release_id int) error {
	return api.Delete(r.Client, "/releases/"+strconv.Itoa(release_id), nil)
}
