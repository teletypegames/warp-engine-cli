package domain

import "errors"

// ExtraServiceInterface reads and replaces the host's additional resources
// of a software (hosted engines only). A software is named by a reference
// (its ID or exact name), a resource by its association name.
type ExtraServiceInterface interface {
	// Show returns what the resource association holds for the software reference.
	Show(reference string, association string) (*ExtraModel, error)
	// Update replaces what the resource association holds for the software
	// reference with the "value" of attributes.
	Update(reference string, association string, attributes Attributes) (*ExtraModel, error)
}

// ExtraService is the [ExtraServiceInterface] over an
// [ExtraRepositoryInterface] and, to resolve software references, a
// [SoftwareRepositoryInterface].
type ExtraService struct {
	// ExtraRepository talks to the extras resource.
	ExtraRepository ExtraRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// Show checks the association, resolves reference and gets the resource.
func (s ExtraService) Show(reference string, association string) (*ExtraModel, error) {
	if association == "" {
		return nil, errors.New("no association given")
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.ExtraRepository.Find(software_id, association)
}

// Update checks the association and that attributes hold a "value" — the
// server would refuse without one — resolves reference and replaces the
// resource.
func (s ExtraService) Update(reference string, association string, attributes Attributes) (*ExtraModel, error) {
	if association == "" {
		return nil, errors.New("no association given")
	}
	if !attributes.Has("value") {
		return nil, errors.New("an extra update needs a value: --value JSON or --data '{\"value\": …}'")
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.ExtraRepository.Update(software_id, association, attributes)
}
