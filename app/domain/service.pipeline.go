package domain

// PipelineServiceInterface manages the CI pipelines the engine knows and
// their runs.
type PipelineServiceInterface interface {
	// All lists the pipelines matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[PipelineModel], error)
	// Show returns the pipeline pipeline_id.
	Show(pipeline_id int) (*PipelineModel, error)
	// Update changes the platform or the software of the pipeline
	// pipeline_id; softwareId may be given as a software reference.
	Update(pipeline_id int, attributes Attributes) (*PipelineModel, error)
	// Sync reads the repositories from the CI server.
	Sync() (*PipelineSyncModel, error)
	// Runs lists one page of the runs of the pipeline pipeline_id.
	Runs(pipeline_id int, page int) (*PageModel[PipelineRunModel], error)
	// Trigger starts a run of the pipeline pipeline_id on branch (empty:
	// the server's default, main).
	Trigger(pipeline_id int, branch string) (*PipelineRunModel, error)
}

// PipelineService is the [PipelineServiceInterface] over a
// [PipelineRepositoryInterface]; it resolves software references through
// the [SoftwareRepositoryInterface].
type PipelineService struct {
	// PipelineRepository talks to the pipelines resource.
	PipelineRepository PipelineRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
	// BuildPlatformService checks the platform
	// against the server's; nil checks nothing.
	BuildPlatformService BuildPlatformServiceInterface
}

// All lists the pipelines.
func (s PipelineService) All(filter Filter, query ListQuery) (*PageModel[PipelineModel], error) {
	return s.PipelineRepository.All(filter, query)
}

// Show gets the pipeline.
func (s PipelineService) Show(pipeline_id int) (*PipelineModel, error) {
	return s.PipelineRepository.FindByID(pipeline_id)
}

// Update changes the pipeline. A platform must be one of the server's or
// [PipelineUnknownPlatform]. A softwareId given as a string is a
// software reference (ID or exact name) and is resolved to its ID; nil
// detaches the software. Empty attributes are [ErrNothingToChange].
func (s PipelineService) Update(pipeline_id int, attributes Attributes) (*PipelineModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	if err := checkPlatformAttribute(s.BuildPlatformService, attributes, PipelineUnknownPlatform); err != nil {
		return nil, err
	}
	if reference, ok := attributes["softwareId"].(string); ok {
		software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
		if err != nil {
			return nil, err
		}
		resolved := Attributes{}
		for key, value := range attributes {
			resolved[key] = value
		}
		resolved["softwareId"] = software_id
		attributes = resolved
	}
	return s.PipelineRepository.Update(pipeline_id, attributes)
}

// Sync syncs the pipelines with the CI server.
func (s PipelineService) Sync() (*PipelineSyncModel, error) {
	return s.PipelineRepository.Sync()
}

// Runs lists the runs of the pipeline.
func (s PipelineService) Runs(pipeline_id int, page int) (*PageModel[PipelineRunModel], error) {
	return s.PipelineRepository.Runs(pipeline_id, page)
}

// Trigger starts a run; an empty branch sends no body field, so the server
// picks its default.
func (s PipelineService) Trigger(pipeline_id int, branch string) (*PipelineRunModel, error) {
	attributes := Attributes{}
	if branch != "" {
		attributes["branch"] = branch
	}
	return s.PipelineRepository.Trigger(pipeline_id, attributes)
}
