package domain

import "strconv"

// SoftwareForkServiceInterface moderates the forks other engines announced of
// this engine's softwares.
type SoftwareForkServiceInterface interface {
	// All lists the announced forks matching filter, paged by query; a
	// "software" key in filter (an ID or an exact name) becomes softwareId.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareForkModel], error)
	// Show returns the announced fork software_fork_id.
	Show(software_fork_id int) (*SoftwareForkModel, error)
	// Update changes the status of the announced fork software_fork_id.
	Update(software_fork_id int, attributes Attributes) (*SoftwareForkModel, error)
	// Destroy deletes the announced fork software_fork_id.
	Destroy(software_fork_id int) error
}

// SoftwareForkService is the [SoftwareForkServiceInterface] over a
// [SoftwareForkRepositoryInterface]; software references are resolved through a
// [SoftwareRepositoryInterface].
type SoftwareForkService struct {
	// SoftwareForkRepository talks to the forks resource.
	SoftwareForkRepository SoftwareForkRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// All resolves the "software" filter to a softwareId, if given, and lists
// the announced forks.
func (s SoftwareForkService) All(filter Filter, query ListQuery) (*PageModel[SoftwareForkModel], error) {
	if reference, ok := filter["software"]; ok {
		software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
		if err != nil {
			return nil, err
		}
		resolved := Filter{}
		for key, value := range filter {
			resolved[key] = value
		}
		delete(resolved, "software")
		resolved["softwareId"] = strconv.Itoa(software_id)
		filter = resolved
	}
	return s.SoftwareForkRepository.All(filter, query)
}

// Show gets the announced fork software_fork_id.
func (s SoftwareForkService) Show(software_fork_id int) (*SoftwareForkModel, error) {
	return s.SoftwareForkRepository.FindByID(software_fork_id)
}

// Update changes the announced fork software_fork_id; empty attributes are
// [ErrNothingToChange].
func (s SoftwareForkService) Update(software_fork_id int, attributes Attributes) (*SoftwareForkModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.SoftwareForkRepository.Update(software_fork_id, attributes)
}

// Destroy deletes the announced fork software_fork_id.
func (s SoftwareForkService) Destroy(software_fork_id int) error {
	return s.SoftwareForkRepository.Destroy(software_fork_id)
}
