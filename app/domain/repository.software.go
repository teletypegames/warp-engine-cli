package domain

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SoftwareRepositoryInterface is the Admin API's softwares resource.
type SoftwareRepositoryInterface interface {
	// All lists the softwares matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareModel], error)
	// FindByID returns the software software_id with the includes
	// (see [SoftwareIncludes]) embedded.
	FindByID(software_id int, include []string) (*SoftwareModel, error)
	// FindByName returns the software whose name is exactly name, or an
	// error with the status 404.
	FindByName(name string) (*SoftwareModel, error)
	// Create creates a software from attributes.
	Create(attributes Attributes) (*SoftwareModel, error)
	// Update changes the attributes of the software software_id.
	Update(software_id int, attributes Attributes) (*SoftwareModel, error)
	// Destroy deletes the software software_id.
	Destroy(software_id int) error
	// Scaffold creates a software together with its forge repository and CI repository.
	Scaffold(attributes Attributes) (*ScaffoldModel, error)
}

// SoftwareRepository is the [SoftwareRepositoryInterface] over the Admin API.
type SoftwareRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /softwares.
func (r SoftwareRepository) All(filter Filter, query ListQuery) (*PageModel[SoftwareModel], error) {
	return list[SoftwareModel](r.Client, "/softwares", filter, query)
}

// FindByID gets GET /softwares/{id}?include=….
func (r SoftwareRepository) FindByID(software_id int, include []string) (*SoftwareModel, error) {
	var query url.Values
	if len(include) > 0 {
		query = url.Values{"include": {strings.Join(include, ",")}}
	}
	return api.Get[SoftwareModel](r.Client, "/softwares/"+strconv.Itoa(software_id), query)
}

// FindByName lists the softwares whose name contains name and returns the
// one whose name is name; none is an [*api.Error] with the status 404.
func (r SoftwareRepository) FindByName(name string) (*SoftwareModel, error) {
	page, err := r.All(Filter{"name": name}, ListQuery{All: true})
	if err != nil {
		return nil, err
	}
	for _, software := range page.Items {
		if software.Name == name {
			return &software, nil
		}
	}
	return nil, notFound("/softwares", "no software named "+strconv.Quote(name))
}

// Create posts POST /softwares.
func (r SoftwareRepository) Create(attributes Attributes) (*SoftwareModel, error) {
	return api.Send[SoftwareModel](r.Client, http.MethodPost, "/softwares", attributes)
}

// Update sends PATCH /softwares/{id}.
func (r SoftwareRepository) Update(software_id int, attributes Attributes) (*SoftwareModel, error) {
	return api.Send[SoftwareModel](r.Client, http.MethodPatch, "/softwares/"+strconv.Itoa(software_id), attributes)
}

// Destroy sends DELETE /softwares/{id}.
func (r SoftwareRepository) Destroy(software_id int) error {
	return api.Delete(r.Client, "/softwares/"+strconv.Itoa(software_id), nil)
}

// Scaffold posts POST /softwares/scaffold.
func (r SoftwareRepository) Scaffold(attributes Attributes) (*ScaffoldModel, error) {
	return api.Send[ScaffoldModel](r.Client, http.MethodPost, "/softwares/scaffold", attributes)
}
