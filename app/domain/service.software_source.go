package domain

import "errors"

// SoftwareSourceServiceInterface manages the source repositories of the
// softwares. A software is named by a reference (its ID or exact name).
type SoftwareSourceServiceInterface interface {
	// Show returns the source of the software reference.
	Show(reference string) (*SoftwareSourceModel, error)
	// Set sets the source of the software reference from attributes (url,
	// ref).
	Set(reference string, attributes Attributes) (*SoftwareSourceModel, error)
	// Destroy removes the source of the software reference.
	Destroy(reference string) error
}

// SoftwareSourceService is the [SoftwareSourceServiceInterface] over a
// [SoftwareSourceRepositoryInterface] and, to resolve software references,
// a [SoftwareRepositoryInterface].
type SoftwareSourceService struct {
	// SoftwareSourceRepository talks to the source of a software.
	SoftwareSourceRepository SoftwareSourceRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// Show resolves reference and gets its source.
func (s SoftwareSourceService) Show(reference string) (*SoftwareSourceModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareSourceRepository.Find(software_id)
}

// Set checks that attributes hold a url — the server would refuse without
// it, as the source is replaced as a whole —, resolves reference and sets
// its source.
func (s SoftwareSourceService) Set(reference string, attributes Attributes) (*SoftwareSourceModel, error) {
	if attributes["url"] == nil || attributes["url"] == "" {
		return nil, errors.New("a source needs a url")
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareSourceRepository.Save(software_id, attributes)
}

// Destroy resolves reference and removes its source.
func (s SoftwareSourceService) Destroy(reference string) error {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return err
	}
	return s.SoftwareSourceRepository.Destroy(software_id)
}
