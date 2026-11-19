package domain

import "errors"

// ExternalLinkServiceInterface manages the external links of the
// softwares. A software is named by a reference (its ID or exact name), a
// link by its ID.
type ExternalLinkServiceInterface interface {
	// All lists the links of the software reference.
	All(reference string) ([]ExternalLinkModel, error)
	// Show returns the link id.
	Show(id string) (*ExternalLinkModel, error)
	// Create adds a link to the software reference.
	Create(reference string, attributes Attributes) (*ExternalLinkModel, error)
	// Update changes the link id.
	Update(id string, attributes Attributes) (*ExternalLinkModel, error)
	// Destroy deletes the link id.
	Destroy(id string) error
}

// ExternalLinkService is the [ExternalLinkServiceInterface] over an
// [ExternalLinkRepositoryInterface] and, to resolve software references, a
// [SoftwareRepositoryInterface].
type ExternalLinkService struct {
	// ExternalLinkRepository talks to the external links resource.
	ExternalLinkRepository ExternalLinkRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// All resolves reference and lists its links.
func (s ExternalLinkService) All(reference string) ([]ExternalLinkModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.ExternalLinkRepository.All(software_id)
}

// Show parses id and gets the link.
func (s ExternalLinkService) Show(id string) (*ExternalLinkModel, error) {
	link_id, err := ParseID("link", id)
	if err != nil {
		return nil, err
	}
	return s.ExternalLinkRepository.FindByID(link_id)
}

// Create checks that attributes hold a label and a url — the server would
// refuse without them — resolves reference and adds the link.
func (s ExternalLinkService) Create(reference string, attributes Attributes) (*ExternalLinkModel, error) {
	for _, key := range []string{"label", "url"} {
		if attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("a link needs a " + key)
		}
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.ExternalLinkRepository.Create(software_id, attributes)
}

// Update parses id and changes the link; empty attributes are
// [ErrNothingToChange].
func (s ExternalLinkService) Update(id string, attributes Attributes) (*ExternalLinkModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	link_id, err := ParseID("link", id)
	if err != nil {
		return nil, err
	}
	return s.ExternalLinkRepository.Update(link_id, attributes)
}

// Destroy parses id and deletes the link.
func (s ExternalLinkService) Destroy(id string) error {
	link_id, err := ParseID("link", id)
	if err != nil {
		return err
	}
	return s.ExternalLinkRepository.Destroy(link_id)
}
