package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// PlatformLinkRepositoryInterface is the Admin API's platform links resource.
type PlatformLinkRepositoryInterface interface {
	// All lists the platform links matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[PlatformLinkModel], error)
	// FindByID returns the platform link platform_link_id.
	FindByID(platform_link_id int) (*PlatformLinkModel, error)
	// Create creates a platform link from attributes.
	Create(attributes Attributes) (*PlatformLinkModel, error)
	// Update changes the attributes of the platform link platform_link_id.
	Update(platform_link_id int, attributes Attributes) (*PlatformLinkModel, error)
	// Destroy deletes the platform link platform_link_id.
	Destroy(platform_link_id int) error
}

// PlatformLinkRepository is the [PlatformLinkRepositoryInterface] over the Admin API.
type PlatformLinkRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /platform_links.
func (r PlatformLinkRepository) All(filter Filter, query ListQuery) (*PageModel[PlatformLinkModel], error) {
	return list[PlatformLinkModel](r.Client, "/platform_links", filter, query)
}

// FindByID gets GET /platform_links/{id}.
func (r PlatformLinkRepository) FindByID(platform_link_id int) (*PlatformLinkModel, error) {
	return api.Get[PlatformLinkModel](r.Client, "/platform_links/"+strconv.Itoa(platform_link_id), nil)
}

// Create posts POST /platform_links.
func (r PlatformLinkRepository) Create(attributes Attributes) (*PlatformLinkModel, error) {
	return api.Send[PlatformLinkModel](r.Client, http.MethodPost, "/platform_links", attributes)
}

// Update sends PATCH /platform_links/{id}.
func (r PlatformLinkRepository) Update(platform_link_id int, attributes Attributes) (*PlatformLinkModel, error) {
	return api.Send[PlatformLinkModel](r.Client, http.MethodPatch, "/platform_links/"+strconv.Itoa(platform_link_id), attributes)
}

// Destroy sends DELETE /platform_links/{id}; the server keeps the link, marked deleted.
func (r PlatformLinkRepository) Destroy(platform_link_id int) error {
	return api.Delete(r.Client, "/platform_links/"+strconv.Itoa(platform_link_id), nil)
}
