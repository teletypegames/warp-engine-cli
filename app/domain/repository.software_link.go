package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SoftwareLinkRepositoryInterface is the Admin API's external links resource.
type SoftwareLinkRepositoryInterface interface {
	// All lists the links of the software software_id (the list has no paging).
	All(software_id int) ([]SoftwareLinkModel, error)
	// FindByID returns the link link_id.
	FindByID(link_id int) (*SoftwareLinkModel, error)
	// Create adds a link to the software software_id from attributes.
	Create(software_id int, attributes Attributes) (*SoftwareLinkModel, error)
	// Update changes the attributes of the link link_id.
	Update(link_id int, attributes Attributes) (*SoftwareLinkModel, error)
	// Destroy deletes the link link_id.
	Destroy(link_id int) error
}

// SoftwareLinkRepository is the [SoftwareLinkRepositoryInterface] over the Admin API.
type SoftwareLinkRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /softwares/{softwareId}/links.
func (r SoftwareLinkRepository) All(software_id int) ([]SoftwareLinkModel, error) {
	page, err := list[SoftwareLinkModel](r.Client, "/softwares/"+strconv.Itoa(software_id)+"/links", nil, ListQuery{})
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

// FindByID gets GET /links/{id}.
func (r SoftwareLinkRepository) FindByID(link_id int) (*SoftwareLinkModel, error) {
	return api.Get[SoftwareLinkModel](r.Client, "/links/"+strconv.Itoa(link_id), nil)
}

// Create posts POST /softwares/{softwareId}/links.
func (r SoftwareLinkRepository) Create(software_id int, attributes Attributes) (*SoftwareLinkModel, error) {
	return api.Send[SoftwareLinkModel](r.Client, http.MethodPost, "/softwares/"+strconv.Itoa(software_id)+"/links", attributes)
}

// Update sends PATCH /links/{id}.
func (r SoftwareLinkRepository) Update(link_id int, attributes Attributes) (*SoftwareLinkModel, error) {
	return api.Send[SoftwareLinkModel](r.Client, http.MethodPatch, "/links/"+strconv.Itoa(link_id), attributes)
}

// Destroy sends DELETE /links/{id}.
func (r SoftwareLinkRepository) Destroy(link_id int) error {
	return api.Delete(r.Client, "/links/"+strconv.Itoa(link_id), nil)
}
