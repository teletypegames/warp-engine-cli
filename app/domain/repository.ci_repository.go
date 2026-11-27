package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// CiRepositoryRepositoryInterface is the Admin API's CI repositories resource.
type CiRepositoryRepositoryInterface interface {
	// All lists the CI repositories matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[CiRepositoryModel], error)
	// FindByID returns the CI repository ci_repository_id.
	FindByID(ci_repository_id int) (*CiRepositoryModel, error)
	// Update changes the platform or software of the CI repository ci_repository_id.
	Update(ci_repository_id int, attributes Attributes) (*CiRepositoryModel, error)
	// Sync reads the repositories from the CI server.
	Sync() (*CiRepositorySyncModel, error)
	// Runs lists the runs of the CI repository ci_repository_id, newest first: page
	// page of the CI server's listing (0 for the first).
	Runs(ci_repository_id int, page int) (*PageModel[PipelineRunModel], error)
	// Trigger starts a run of the CI repository ci_repository_id with attributes
	// (branch).
	Trigger(ci_repository_id int, attributes Attributes) (*PipelineRunModel, error)
}

// CiRepositoryRepository is the [CiRepositoryRepositoryInterface] over the Admin API.
type CiRepositoryRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /CI repositories.
func (r CiRepositoryRepository) All(filter Filter, query ListQuery) (*PageModel[CiRepositoryModel], error) {
	return list[CiRepositoryModel](r.Client, "/ci_repositories", filter, query)
}

// FindByID gets GET /ci_repositories/{id}.
func (r CiRepositoryRepository) FindByID(ci_repository_id int) (*CiRepositoryModel, error) {
	return api.Get[CiRepositoryModel](r.Client, "/ci_repositories/"+strconv.Itoa(ci_repository_id), nil)
}

// Update sends PATCH /ci_repositories/{id}.
func (r CiRepositoryRepository) Update(ci_repository_id int, attributes Attributes) (*CiRepositoryModel, error) {
	return api.Send[CiRepositoryModel](r.Client, http.MethodPatch, "/ci_repositories/"+strconv.Itoa(ci_repository_id), attributes)
}

// Sync posts POST /ci_repositories/sync.
func (r CiRepositoryRepository) Sync() (*CiRepositorySyncModel, error) {
	return api.Send[CiRepositorySyncModel](r.Client, http.MethodPost, "/ci_repositories/sync", nil)
}

// Runs gets GET /ci_repositories/{id}/runs?page=…. The list has no paging block:
// the page answers with Page set to page.
func (r CiRepositoryRepository) Runs(ci_repository_id int, page int) (*PageModel[PipelineRunModel], error) {
	result, err := list[PipelineRunModel](r.Client, "/ci_repositories/"+strconv.Itoa(ci_repository_id)+"/runs", nil, ListQuery{Page: page})
	if err != nil {
		return nil, err
	}
	result.Page = page
	return result, nil
}

// Trigger posts POST /ci_repositories/{id}/trigger.
func (r CiRepositoryRepository) Trigger(ci_repository_id int, attributes Attributes) (*PipelineRunModel, error) {
	return api.Send[PipelineRunModel](r.Client, http.MethodPost, "/ci_repositories/"+strconv.Itoa(ci_repository_id)+"/trigger", attributes)
}
