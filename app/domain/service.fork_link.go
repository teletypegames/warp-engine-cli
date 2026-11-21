package domain

import "strconv"

// ForkLinkServiceInterface moderates the forks other engines announced of
// this engine's softwares.
type ForkLinkServiceInterface interface {
	// All lists the fork links matching filter, paged by query; a
	// "software" key in filter (an ID or an exact name) becomes softwareId.
	All(filter Filter, query ListQuery) (*PageModel[ForkLinkModel], error)
	// Show returns the fork link fork_link_id.
	Show(fork_link_id int) (*ForkLinkModel, error)
	// Update changes the status of the fork link fork_link_id.
	Update(fork_link_id int, attributes Attributes) (*ForkLinkModel, error)
	// Destroy deletes the fork link fork_link_id.
	Destroy(fork_link_id int) error
}

// ForkLinkService is the [ForkLinkServiceInterface] over a
// [ForkLinkRepositoryInterface]; software references are resolved through a
// [SoftwareRepositoryInterface].
type ForkLinkService struct {
	// ForkLinkRepository talks to the fork_links resource.
	ForkLinkRepository ForkLinkRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// All resolves the "software" filter to a softwareId, if given, and lists
// the fork links.
func (s ForkLinkService) All(filter Filter, query ListQuery) (*PageModel[ForkLinkModel], error) {
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
	return s.ForkLinkRepository.All(filter, query)
}

// Show gets the fork link fork_link_id.
func (s ForkLinkService) Show(fork_link_id int) (*ForkLinkModel, error) {
	return s.ForkLinkRepository.FindByID(fork_link_id)
}

// Update changes the fork link fork_link_id; empty attributes are
// [ErrNothingToChange].
func (s ForkLinkService) Update(fork_link_id int, attributes Attributes) (*ForkLinkModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.ForkLinkRepository.Update(fork_link_id, attributes)
}

// Destroy deletes the fork link fork_link_id.
func (s ForkLinkService) Destroy(fork_link_id int) error {
	return s.ForkLinkRepository.Destroy(fork_link_id)
}
