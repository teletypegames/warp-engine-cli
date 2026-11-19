package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// PipelineRepositoryInterface is the Admin API's pipelines resource.
type PipelineRepositoryInterface interface {
	// All lists the pipelines matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[PipelineModel], error)
	// FindByID returns the pipeline pipeline_id.
	FindByID(pipeline_id int) (*PipelineModel, error)
	// Update changes the platform or software of the pipeline pipeline_id.
	Update(pipeline_id int, attributes Attributes) (*PipelineModel, error)
	// Sync reads the repositories from the CI server.
	Sync() (*PipelineSyncModel, error)
	// Runs lists the runs of the pipeline pipeline_id, newest first: page
	// page of the CI server's listing (0 for the first).
	Runs(pipeline_id int, page int) (*PageModel[PipelineRunModel], error)
	// Trigger starts a run of the pipeline pipeline_id with attributes
	// (branch).
	Trigger(pipeline_id int, attributes Attributes) (*PipelineRunModel, error)
}

// PipelineRepository is the [PipelineRepositoryInterface] over the Admin API.
type PipelineRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /pipelines.
func (r PipelineRepository) All(filter Filter, query ListQuery) (*PageModel[PipelineModel], error) {
	return list[PipelineModel](r.Client, "/pipelines", filter, query)
}

// FindByID gets GET /pipelines/{id}.
func (r PipelineRepository) FindByID(pipeline_id int) (*PipelineModel, error) {
	return api.Get[PipelineModel](r.Client, "/pipelines/"+strconv.Itoa(pipeline_id), nil)
}

// Update sends PATCH /pipelines/{id}.
func (r PipelineRepository) Update(pipeline_id int, attributes Attributes) (*PipelineModel, error) {
	return api.Send[PipelineModel](r.Client, http.MethodPatch, "/pipelines/"+strconv.Itoa(pipeline_id), attributes)
}

// Sync posts POST /pipelines/sync.
func (r PipelineRepository) Sync() (*PipelineSyncModel, error) {
	return api.Send[PipelineSyncModel](r.Client, http.MethodPost, "/pipelines/sync", nil)
}

// Runs gets GET /pipelines/{id}/runs?page=…. The list has no paging block:
// the page answers with Page set to page.
func (r PipelineRepository) Runs(pipeline_id int, page int) (*PageModel[PipelineRunModel], error) {
	result, err := list[PipelineRunModel](r.Client, "/pipelines/"+strconv.Itoa(pipeline_id)+"/runs", nil, ListQuery{Page: page})
	if err != nil {
		return nil, err
	}
	result.Page = page
	return result, nil
}

// Trigger posts POST /pipelines/{id}/trigger.
func (r PipelineRepository) Trigger(pipeline_id int, attributes Attributes) (*PipelineRunModel, error) {
	return api.Send[PipelineRunModel](r.Client, http.MethodPost, "/pipelines/"+strconv.Itoa(pipeline_id)+"/trigger", attributes)
}
