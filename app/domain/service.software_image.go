package domain

import "errors"

// SoftwareImageServiceInterface manages the images of the softwares. A
// software is named by a reference (its ID or exact name), a software image
// by its ID.
type SoftwareImageServiceInterface interface {
	// All lists the images of the software reference in order.
	All(reference string) ([]SoftwareImageModel, error)
	// Show returns the software image id.
	Show(id string) (*SoftwareImageModel, error)
	// Create binds an image of the host's image library to the software reference.
	Create(reference string, attributes Attributes) (*SoftwareImageModel, error)
	// Upload sends the file at file_path into the image library and binds it
	// to the software reference.
	Upload(reference string, file_path string, attributes Attributes) (*SoftwareImageModel, error)
	// Update makes the software image id the default, or moves it.
	Update(id string, attributes Attributes) (*SoftwareImageModel, error)
	// Destroy unbinds the software image id from its software.
	Destroy(id string) error
}

// SoftwareImageService is the [SoftwareImageServiceInterface] over a
// [SoftwareImageRepositoryInterface] and, to resolve software references, a
// [SoftwareRepositoryInterface].
type SoftwareImageService struct {
	// SoftwareImageRepository talks to the software images resource.
	SoftwareImageRepository SoftwareImageRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// All resolves reference and lists its images.
func (s SoftwareImageService) All(reference string) ([]SoftwareImageModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareImageRepository.All(software_id)
}

// Show parses id and gets the software image.
func (s SoftwareImageService) Show(id string) (*SoftwareImageModel, error) {
	image_id, err := ParseID("software image", id)
	if err != nil {
		return nil, err
	}
	return s.SoftwareImageRepository.FindByID(image_id)
}

// Create checks that attributes hold an imageId — the server would refuse
// without one — resolves reference and binds the image.
func (s SoftwareImageService) Create(reference string, attributes Attributes) (*SoftwareImageModel, error) {
	if attributes["imageId"] == nil {
		return nil, errors.New("binding an image needs the imageId of an image in the library")
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareImageRepository.Create(software_id, attributes)
}

// Upload checks that file_path is given, resolves reference and uploads
// the file; attributes may hold isDefault and position.
func (s SoftwareImageService) Upload(reference string, file_path string, attributes Attributes) (*SoftwareImageModel, error) {
	if file_path == "" {
		return nil, errors.New("no image file given")
	}
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SoftwareImageRepository.Upload(software_id, file_path, attributes)
}

// Update parses id and changes the software image; empty attributes are
// [ErrNothingToChange].
func (s SoftwareImageService) Update(id string, attributes Attributes) (*SoftwareImageModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	image_id, err := ParseID("software image", id)
	if err != nil {
		return nil, err
	}
	return s.SoftwareImageRepository.Update(image_id, attributes)
}

// Destroy parses id and unbinds the software image.
func (s SoftwareImageService) Destroy(id string) error {
	image_id, err := ParseID("software image", id)
	if err != nil {
		return err
	}
	return s.SoftwareImageRepository.Destroy(image_id)
}
