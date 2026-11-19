package domain

import (
	"fmt"
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SoftwareImageRepositoryInterface is the Admin API's software images resource.
type SoftwareImageRepositoryInterface interface {
	// All lists the images of the software software_id in order (the list
	// has no paging).
	All(software_id int) ([]SoftwareImageModel, error)
	// FindByID returns the software image image_id.
	FindByID(image_id int) (*SoftwareImageModel, error)
	// Create binds an image of the host's image library (attributes hold its
	// imageId) to the software software_id.
	Create(software_id int, attributes Attributes) (*SoftwareImageModel, error)
	// Upload sends the file at file_path into the image library and binds it
	// to the software software_id; attributes may hold isDefault and position.
	Upload(software_id int, file_path string, attributes Attributes) (*SoftwareImageModel, error)
	// Update changes isDefault or position of the software image image_id.
	Update(image_id int, attributes Attributes) (*SoftwareImageModel, error)
	// Destroy unbinds the software image image_id from its software.
	Destroy(image_id int) error
}

// SoftwareImageRepository is the [SoftwareImageRepositoryInterface] over the Admin API.
type SoftwareImageRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /softwares/{softwareId}/images.
func (r SoftwareImageRepository) All(software_id int) ([]SoftwareImageModel, error) {
	page, err := list[SoftwareImageModel](r.Client, "/softwares/"+strconv.Itoa(software_id)+"/images", nil, ListQuery{})
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

// FindByID gets GET /software_images/{id}.
func (r SoftwareImageRepository) FindByID(image_id int) (*SoftwareImageModel, error) {
	return api.Get[SoftwareImageModel](r.Client, "/software_images/"+strconv.Itoa(image_id), nil)
}

// Create posts POST /softwares/{softwareId}/images.
func (r SoftwareImageRepository) Create(software_id int, attributes Attributes) (*SoftwareImageModel, error) {
	return api.Send[SoftwareImageModel](r.Client, http.MethodPost, "/softwares/"+strconv.Itoa(software_id)+"/images", attributes)
}

// Upload posts POST /softwares/{softwareId}/images/upload as
// multipart/form-data: the file as the part "file", every attribute as a
// plain part with its value as text.
func (r SoftwareImageRepository) Upload(software_id int, file_path string, attributes Attributes) (*SoftwareImageModel, error) {
	fields := map[string]string{}
	for key, value := range attributes {
		if value != nil {
			fields[key] = fmt.Sprint(value)
		}
	}
	var envelope api.Envelope[SoftwareImageModel]
	if err := r.Client.Upload("/softwares/"+strconv.Itoa(software_id)+"/images/upload", fields, "file", file_path, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

// Update sends PATCH /software_images/{id}.
func (r SoftwareImageRepository) Update(image_id int, attributes Attributes) (*SoftwareImageModel, error) {
	return api.Send[SoftwareImageModel](r.Client, http.MethodPatch, "/software_images/"+strconv.Itoa(image_id), attributes)
}

// Destroy sends DELETE /software_images/{id}.
func (r SoftwareImageRepository) Destroy(image_id int) error {
	return api.Delete(r.Client, "/software_images/"+strconv.Itoa(image_id), nil)
}
