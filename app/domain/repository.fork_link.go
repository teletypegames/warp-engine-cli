package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// ForkLinkRepositoryInterface is the Admin API's fork_links resource: the
// forks other engines announced.
type ForkLinkRepositoryInterface interface {
	// All lists the fork links matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[ForkLinkModel], error)
	// FindByID returns the fork link fork_link_id.
	FindByID(fork_link_id int) (*ForkLinkModel, error)
	// Update changes the attributes (the status) of the fork link
	// fork_link_id.
	Update(fork_link_id int, attributes Attributes) (*ForkLinkModel, error)
	// Destroy deletes the fork link fork_link_id.
	Destroy(fork_link_id int) error
}

// ForkLinkRepository is the [ForkLinkRepositoryInterface] over the Admin API.
type ForkLinkRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /fork_links.
func (r ForkLinkRepository) All(filter Filter, query ListQuery) (*PageModel[ForkLinkModel], error) {
	return list[ForkLinkModel](r.Client, "/fork_links", filter, query)
}

// FindByID gets GET /fork_links/{id}.
func (r ForkLinkRepository) FindByID(fork_link_id int) (*ForkLinkModel, error) {
	return api.Get[ForkLinkModel](r.Client, "/fork_links/"+strconv.Itoa(fork_link_id), nil)
}

// Update sends PATCH /fork_links/{id}.
func (r ForkLinkRepository) Update(fork_link_id int, attributes Attributes) (*ForkLinkModel, error) {
	return api.Send[ForkLinkModel](r.Client, http.MethodPatch, "/fork_links/"+strconv.Itoa(fork_link_id), attributes)
}

// Destroy sends DELETE /fork_links/{id}.
func (r ForkLinkRepository) Destroy(fork_link_id int) error {
	return api.Delete(r.Client, "/fork_links/"+strconv.Itoa(fork_link_id), nil)
}
