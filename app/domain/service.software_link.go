package domain

import "errors"

// SoftwareLinkServiceInterface manages the plain links of the
// softwares. A software is named by a reference (its ID or exact name), a
// link by its ID.
type SoftwareLinkServiceInterface interface {
	// All lists the links of the software reference.
	All(reference string) ([]SoftwareLinkModel, error)
	// Show returns the link id.
	Show(id string) (*SoftwareLinkModel, error)
	// Create adds a link to the software reference.
	Create(reference string, attributes Attributes) (*SoftwareLinkModel, error)
	// Update changes the link id.
	Update(id string, attributes Attributes) (*SoftwareLinkModel, error)
	// Destroy deletes the link id.
	Destroy(id string) error
}

// SoftwareLinkService is the [SoftwareLinkServiceInterface] over an
// [SoftwareLinkRepositoryInterface] and, to resolve software references, a
// [SoftwareRepositoryInterface].
type SoftwareLinkService struct {
	// SoftwareLinkRepository talks to the links resource.
	SoftwareLinkRepository SoftwareLinkRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// All resolves reference and lists its links.
func (s SoftwareLinkService) All(reference string) ([]SoftwareLinkModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareLinkRepository.All(software_id)
}

// Show parses id and gets the link.
func (s SoftwareLinkService) Show(id string) (*SoftwareLinkModel, error) {
	link_id, err := ParseID("link", id)
	if err != nil {
		return nil, err
	}
	return s.SoftwareLinkRepository.FindByID(link_id)
}

// Create checks that attributes hold a label and a url — the server would
// refuse without them — resolves reference and adds the link.
func (s SoftwareLinkService) Create(reference string, attributes Attributes) (*SoftwareLinkModel, error) {
	for _, key := range []string{"label", "url"} {
		if attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("a link needs a " + key)
		}
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareLinkRepository.Create(software_id, attributes)
}

// Update parses id and changes the link; empty attributes are
// [ErrNothingToChange].
func (s SoftwareLinkService) Update(id string, attributes Attributes) (*SoftwareLinkModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	link_id, err := ParseID("link", id)
	if err != nil {
		return nil, err
	}
	return s.SoftwareLinkRepository.Update(link_id, attributes)
}

// Destroy parses id and deletes the link.
func (s SoftwareLinkService) Destroy(id string) error {
	link_id, err := ParseID("link", id)
	if err != nil {
		return err
	}
	return s.SoftwareLinkRepository.Destroy(link_id)
}
