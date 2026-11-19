package domain

import (
	"net/http"
	"net/url"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// ExtraRepositoryInterface is the Admin API's extras resource: the host's
// additional resources of a software. Only hosted engines that declare some
// have it; elsewhere every request answers 404.
type ExtraRepositoryInterface interface {
	// Find returns what the resource association holds for the software
	// software_id.
	Find(software_id int, association string) (*ExtraModel, error)
	// Update replaces what the resource association holds for the software
	// software_id; attributes hold the "value".
	Update(software_id int, association string, attributes Attributes) (*ExtraModel, error)
}

// ExtraRepository is the [ExtraRepositoryInterface] over the Admin API.
type ExtraRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// Find gets GET /softwares/{softwareId}/extras/{association}.
func (r ExtraRepository) Find(software_id int, association string) (*ExtraModel, error) {
	return api.Get[ExtraModel](r.Client, extraPath(software_id, association), nil)
}

// Update sends PATCH /softwares/{softwareId}/extras/{association}.
func (r ExtraRepository) Update(software_id int, association string, attributes Attributes) (*ExtraModel, error) {
	return api.Send[ExtraModel](r.Client, http.MethodPatch, extraPath(software_id, association), attributes)
}

// extraPath is the path of the resource association of the software software_id.
func extraPath(software_id int, association string) string {
	return "/softwares/" + strconv.Itoa(software_id) + "/extras/" + url.PathEscape(association)
}
