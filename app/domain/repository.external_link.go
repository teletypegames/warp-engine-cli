package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// ExternalLinkRepositoryInterface is the Admin API's external links resource.
type ExternalLinkRepositoryInterface interface {
	// All lists the links of the software software_id (the list has no paging).
	All(software_id int) ([]ExternalLinkModel, error)
	// FindByID returns the link link_id.
	FindByID(link_id int) (*ExternalLinkModel, error)
	// Create adds a link to the software software_id from attributes.
	Create(software_id int, attributes Attributes) (*ExternalLinkModel, error)
	// Update changes the attributes of the link link_id.
	Update(link_id int, attributes Attributes) (*ExternalLinkModel, error)
	// Destroy deletes the link link_id.
	Destroy(link_id int) error
}

// ExternalLinkRepository is the [ExternalLinkRepositoryInterface] over the Admin API.
type ExternalLinkRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /softwares/{softwareId}/external_links.
func (r ExternalLinkRepository) All(software_id int) ([]ExternalLinkModel, error) {
	page, err := list[ExternalLinkModel](r.Client, "/softwares/"+strconv.Itoa(software_id)+"/external_links", nil, ListQuery{})
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

// FindByID gets GET /external_links/{id}.
func (r ExternalLinkRepository) FindByID(link_id int) (*ExternalLinkModel, error) {
	return api.Get[ExternalLinkModel](r.Client, "/external_links/"+strconv.Itoa(link_id), nil)
}

// Create posts POST /softwares/{softwareId}/external_links.
func (r ExternalLinkRepository) Create(software_id int, attributes Attributes) (*ExternalLinkModel, error) {
	return api.Send[ExternalLinkModel](r.Client, http.MethodPost, "/softwares/"+strconv.Itoa(software_id)+"/external_links", attributes)
}

// Update sends PATCH /external_links/{id}.
func (r ExternalLinkRepository) Update(link_id int, attributes Attributes) (*ExternalLinkModel, error) {
	return api.Send[ExternalLinkModel](r.Client, http.MethodPatch, "/external_links/"+strconv.Itoa(link_id), attributes)
}

// Destroy sends DELETE /external_links/{id}.
func (r ExternalLinkRepository) Destroy(link_id int) error {
	return api.Delete(r.Client, "/external_links/"+strconv.Itoa(link_id), nil)
}
