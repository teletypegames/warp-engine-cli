package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SoftwareForkRepositoryInterface is the Admin API's forks resource: the
// forks other engines announced.
type SoftwareForkRepositoryInterface interface {
	// All lists the announced forks matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareForkModel], error)
	// FindByID returns the announced fork software_fork_id.
	FindByID(software_fork_id int) (*SoftwareForkModel, error)
	// Update changes the attributes (the status) of the announced fork
	// software_fork_id.
	Update(software_fork_id int, attributes Attributes) (*SoftwareForkModel, error)
	// Destroy deletes the announced fork software_fork_id.
	Destroy(software_fork_id int) error
}

// SoftwareForkRepository is the [SoftwareForkRepositoryInterface] over the Admin API.
type SoftwareForkRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /forks.
func (r SoftwareForkRepository) All(filter Filter, query ListQuery) (*PageModel[SoftwareForkModel], error) {
	return list[SoftwareForkModel](r.Client, "/forks", filter, query)
}

// FindByID gets GET /forks/{id}.
func (r SoftwareForkRepository) FindByID(software_fork_id int) (*SoftwareForkModel, error) {
	return api.Get[SoftwareForkModel](r.Client, "/forks/"+strconv.Itoa(software_fork_id), nil)
}

// Update sends PATCH /forks/{id}.
func (r SoftwareForkRepository) Update(software_fork_id int, attributes Attributes) (*SoftwareForkModel, error) {
	return api.Send[SoftwareForkModel](r.Client, http.MethodPatch, "/forks/"+strconv.Itoa(software_fork_id), attributes)
}

// Destroy sends DELETE /forks/{id}.
func (r SoftwareForkRepository) Destroy(software_fork_id int) error {
	return api.Delete(r.Client, "/forks/"+strconv.Itoa(software_fork_id), nil)
}
