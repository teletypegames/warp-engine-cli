package domain

import (
	"net/http"
	"net/url"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// LocalComponentRepositoryInterface is the Admin API's local_components
// resource and the local CI of a software. Only WarpEngine Local has them;
// any other server answers 404.
type LocalComponentRepositoryInterface interface {
	// All lists the components with their state.
	All() (*PageModel[LocalComponentModel], error)
	// Enable queues starting the component key (and what it needs).
	Enable(key string) (*LocalOperationModel, error)
	// Disable queues stopping the component key (and what needs it).
	Disable(key string) (*LocalOperationModel, error)
	// Destroy queues removing the containers of the component key; with
	// purge its data and generated sign-ins too.
	Destroy(key string, purge bool) (*LocalOperationModel, error)
	// Credentials returns the generated sign-in of the component key.
	Credentials(key string) (*LocalCredentialsModel, error)
	// SetupCI queues setting the software software_id up on the local git
	// server and CI.
	SetupCI(software_id int) (*LocalOperationModel, error)
	// TeardownCI stops the local CI of the software software_id.
	TeardownCI(software_id int) error
}

// LocalComponentRepository is the [LocalComponentRepositoryInterface] over
// the Admin API.
type LocalComponentRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /local_components. The list has no paging.
func (r LocalComponentRepository) All() (*PageModel[LocalComponentModel], error) {
	return list[LocalComponentModel](r.Client, "/local_components", nil, ListQuery{})
}

// Enable posts POST /local_components/{key}/enable.
func (r LocalComponentRepository) Enable(key string) (*LocalOperationModel, error) {
	return api.Send[LocalOperationModel](r.Client, http.MethodPost, "/local_components/"+url.PathEscape(key)+"/enable", nil)
}

// Disable posts POST /local_components/{key}/disable.
func (r LocalComponentRepository) Disable(key string) (*LocalOperationModel, error) {
	return api.Send[LocalOperationModel](r.Client, http.MethodPost, "/local_components/"+url.PathEscape(key)+"/disable", nil)
}

// Destroy sends DELETE /local_components/{key}, with purge=true when purge
// is set, and returns the queued operation.
func (r LocalComponentRepository) Destroy(key string, purge bool) (*LocalOperationModel, error) {
	var query url.Values
	if purge {
		query = url.Values{"purge": {"true"}}
	}
	var envelope api.Envelope[LocalOperationModel]
	if err := r.Client.Admin(http.MethodDelete, "/local_components/"+url.PathEscape(key), query, nil, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

// Credentials gets GET /local_components/{key}/credentials.
func (r LocalComponentRepository) Credentials(key string) (*LocalCredentialsModel, error) {
	return api.Get[LocalCredentialsModel](r.Client, "/local_components/"+url.PathEscape(key)+"/credentials", nil)
}

// SetupCI posts POST /softwares/{softwareId}/local_ci.
func (r LocalComponentRepository) SetupCI(software_id int) (*LocalOperationModel, error) {
	return api.Send[LocalOperationModel](r.Client, http.MethodPost, "/softwares/"+strconv.Itoa(software_id)+"/local_ci", nil)
}

// TeardownCI sends DELETE /softwares/{softwareId}/local_ci.
func (r LocalComponentRepository) TeardownCI(software_id int) error {
	return api.Delete(r.Client, "/softwares/"+strconv.Itoa(software_id)+"/local_ci", nil)
}
