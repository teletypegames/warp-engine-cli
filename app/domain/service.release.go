package domain

import "errors"

// ReleaseServiceInterface manages the releases of the softwares. A software
// is named by a reference (its ID or exact name), a release by its ID.
type ReleaseServiceInterface interface {
	// All lists the releases of the software reference matching filter, paged by query.
	All(reference string, filter Filter, query ListQuery) (*PageModel[ReleaseModel], error)
	// Show returns the release id with its assets.
	Show(id string) (*ReleaseModel, error)
	// Create adds a release to the software reference.
	Create(reference string, attributes Attributes) (*ReleaseModel, error)
	// Update changes the release id.
	Update(id string, attributes Attributes) (*ReleaseModel, error)
	// Destroy deletes the release id with its assets.
	Destroy(id string) error
}

// ReleaseService is the [ReleaseServiceInterface] over a
// [ReleaseRepositoryInterface] and, to resolve software references, a
// [SoftwareRepositoryInterface].
type ReleaseService struct {
	// ReleaseRepository talks to the releases resource.
	ReleaseRepository ReleaseRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// All resolves reference and lists its releases.
func (s ReleaseService) All(reference string, filter Filter, query ListQuery) (*PageModel[ReleaseModel], error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.ReleaseRepository.All(software_id, filter, query)
}

// Show parses id and gets the release.
func (s ReleaseService) Show(id string) (*ReleaseModel, error) {
	release_id, err := ParseID("release", id)
	if err != nil {
		return nil, err
	}
	return s.ReleaseRepository.FindByID(release_id)
}

// Create checks that attributes hold a version — the server would refuse
// without one; the channel follows from it — resolves reference and adds
// the release.
func (s ReleaseService) Create(reference string, attributes Attributes) (*ReleaseModel, error) {
	if attributes["version"] == nil || attributes["version"] == "" {
		return nil, errors.New("a release needs a version")
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.ReleaseRepository.Create(software_id, attributes)
}

// Update parses id and changes the release; empty attributes are
// [ErrNothingToChange].
func (s ReleaseService) Update(id string, attributes Attributes) (*ReleaseModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	release_id, err := ParseID("release", id)
	if err != nil {
		return nil, err
	}
	return s.ReleaseRepository.Update(release_id, attributes)
}

// Destroy parses id and deletes the release.
func (s ReleaseService) Destroy(id string) error {
	release_id, err := ParseID("release", id)
	if err != nil {
		return err
	}
	return s.ReleaseRepository.Destroy(release_id)
}
