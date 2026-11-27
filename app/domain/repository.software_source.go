package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SoftwareSourceRepositoryInterface is the Admin API's source of a software:
// one record per software, read, set and removed through the software.
type SoftwareSourceRepositoryInterface interface {
	// Find returns the source of the software software_id.
	Find(software_id int) (*SoftwareSourceModel, error)
	// Save sets the source of the software software_id from attributes (url,
	// ref), creating or replacing it.
	Save(software_id int, attributes Attributes) (*SoftwareSourceModel, error)
	// Destroy removes the source of the software software_id.
	Destroy(software_id int) error
}

// SoftwareSourceRepository is the [SoftwareSourceRepositoryInterface] over
// the Admin API.
type SoftwareSourceRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// Find gets GET /softwares/{softwareId}/source.
func (r SoftwareSourceRepository) Find(software_id int) (*SoftwareSourceModel, error) {
	return api.Get[SoftwareSourceModel](r.Client, softwareSourcePath(software_id), nil)
}

// Save sends PUT /softwares/{softwareId}/source.
func (r SoftwareSourceRepository) Save(software_id int, attributes Attributes) (*SoftwareSourceModel, error) {
	return api.Send[SoftwareSourceModel](r.Client, http.MethodPut, softwareSourcePath(software_id), attributes)
}

// Destroy sends DELETE /softwares/{softwareId}/source.
func (r SoftwareSourceRepository) Destroy(software_id int) error {
	return api.Delete(r.Client, softwareSourcePath(software_id), nil)
}

// softwareSourcePath is the path of the source of the software software_id.
func softwareSourcePath(software_id int) string {
	return "/softwares/" + strconv.Itoa(software_id) + "/source"
}
