package domain

// CiRepositoryServiceInterface manages the CI repositories the engine knows and
// their runs.
type CiRepositoryServiceInterface interface {
	// All lists the CI repositories matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[CiRepositoryModel], error)
	// Show returns the CI repository ci_repository_id.
	Show(ci_repository_id int) (*CiRepositoryModel, error)
	// Update changes the platform or the software of the CI repository
	// ci_repository_id; softwareId may be given as a software reference.
	Update(ci_repository_id int, attributes Attributes) (*CiRepositoryModel, error)
	// Sync reads the repositories from the CI server.
	Sync() (*CiRepositorySyncModel, error)
	// Runs lists one page of the runs of the CI repository ci_repository_id.
	Runs(ci_repository_id int, page int) (*PageModel[PipelineRunModel], error)
	// Trigger starts a run of the CI repository ci_repository_id on branch (empty:
	// the server's default, main).
	Trigger(ci_repository_id int, branch string) (*PipelineRunModel, error)
}

// CiRepositoryService is the [CiRepositoryServiceInterface] over a
// [CiRepositoryRepositoryInterface]; it resolves software references through
// the [SoftwareRepositoryInterface].
type CiRepositoryService struct {
	// CiRepositoryRepository talks to the CI repositories resource.
	CiRepositoryRepository CiRepositoryRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
	// BuildPlatformService checks the platform
	// against the server's; nil checks nothing.
	BuildPlatformService BuildPlatformServiceInterface
}

// All lists the CI repositories.
func (s CiRepositoryService) All(filter Filter, query ListQuery) (*PageModel[CiRepositoryModel], error) {
	return s.CiRepositoryRepository.All(filter, query)
}

// Show gets the ci_repository.
func (s CiRepositoryService) Show(ci_repository_id int) (*CiRepositoryModel, error) {
	return s.CiRepositoryRepository.FindByID(ci_repository_id)
}

// Update changes the ci_repository. A platform must be one of the server's or
// [CiRepositoryUnknownPlatform]. A softwareId given as a string is a
// software reference (ID or exact name) and is resolved to its ID; nil
// detaches the software. Empty attributes are [ErrNothingToChange].
func (s CiRepositoryService) Update(ci_repository_id int, attributes Attributes) (*CiRepositoryModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	if err := checkPlatformAttribute(s.BuildPlatformService, attributes, CiRepositoryUnknownPlatform); err != nil {
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
	return s.CiRepositoryRepository.Update(ci_repository_id, attributes)
}

// Sync syncs the CI repositories with the CI server.
func (s CiRepositoryService) Sync() (*CiRepositorySyncModel, error) {
	return s.CiRepositoryRepository.Sync()
}

// Runs lists the runs of the ci_repository.
func (s CiRepositoryService) Runs(ci_repository_id int, page int) (*PageModel[PipelineRunModel], error) {
	return s.CiRepositoryRepository.Runs(ci_repository_id, page)
}

// Trigger starts a run; an empty branch sends no body field, so the server
// picks its default.
func (s CiRepositoryService) Trigger(ci_repository_id int, branch string) (*PipelineRunModel, error) {
	attributes := Attributes{}
	if branch != "" {
		attributes["branch"] = branch
	}
	return s.CiRepositoryRepository.Trigger(ci_repository_id, attributes)
}
