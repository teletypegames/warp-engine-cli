package domain

import "errors"

// SoftwareServiceInterface manages the softwares of the catalog. A software
// is named by a reference: its numeric ID or its exact name.
type SoftwareServiceInterface interface {
	// All lists the softwares matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareModel], error)
	// Show returns the software reference with include embedded.
	Show(reference string, include []string) (*SoftwareModel, error)
	// Create creates a software.
	Create(attributes Attributes) (*SoftwareModel, error)
	// Update changes the software reference.
	Update(reference string, attributes Attributes) (*SoftwareModel, error)
	// Destroy deletes the software reference.
	Destroy(reference string) error
	// Scaffold creates a software with its forge repository and pipeline.
	Scaffold(attributes Attributes) (*ScaffoldModel, error)
}

// SoftwareService is the [SoftwareServiceInterface] over a [SoftwareRepositoryInterface].
type SoftwareService struct {
	// SoftwareRepository talks to the softwares resource.
	SoftwareRepository SoftwareRepositoryInterface
}

// All lists the softwares.
func (s SoftwareService) All(filter Filter, query ListQuery) (*PageModel[SoftwareModel], error) {
	return s.SoftwareRepository.All(filter, query)
}

// Show resolves reference and gets the software with include.
func (s SoftwareService) Show(reference string, include []string) (*SoftwareModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareRepository.FindByID(software_id, include)
}

// Create checks that attributes hold name, title and platform — the server
// would refuse without them — and creates the software.
func (s SoftwareService) Create(attributes Attributes) (*SoftwareModel, error) {
	for _, key := range []string{"name", "title", "platform"} {
		if !attributes.Has(key) || attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("a software needs a " + key)
		}
	}
	return s.SoftwareRepository.Create(attributes)
}

// Update resolves reference and changes it; empty attributes are an error,
// since there would be nothing to send.
func (s SoftwareService) Update(reference string, attributes Attributes) (*SoftwareModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareRepository.Update(software_id, attributes)
}

// Destroy resolves reference and deletes it.
func (s SoftwareService) Destroy(reference string) error {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return err
	}
	return s.SoftwareRepository.Destroy(software_id)
}

// Scaffold checks that attributes hold name, title and platform and sets
// the software up.
func (s SoftwareService) Scaffold(attributes Attributes) (*ScaffoldModel, error) {
	for _, key := range []string{"name", "title", "platform"} {
		if !attributes.Has(key) || attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("a scaffold needs a " + key)
		}
	}
	return s.SoftwareRepository.Scaffold(attributes)
}
