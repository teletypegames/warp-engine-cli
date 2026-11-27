package domain

import (
	"net/http"
	"net/url"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// LocalProviderRepositoryInterface is the Admin API's local_providers
// resource and the local CI of a software. Only WarpEngine Local has them;
// any other server answers 404.
type LocalProviderRepositoryInterface interface {
	// All lists the providers with their state.
	All() (*PageModel[LocalProviderModel], error)
	// Enable queues starting the provider key (and what it needs).
	Enable(key string) (*LocalTaskModel, error)
	// Disable queues stopping the provider key (and what needs it).
	Disable(key string) (*LocalTaskModel, error)
	// Destroy queues removing the containers of the provider key; with
	// purge its data and generated sign-ins too.
	Destroy(key string, purge bool) (*LocalTaskModel, error)
	// Credentials returns the generated sign-in of the provider key.
	Credentials(key string) (*LocalCredentialsModel, error)
	// SetupCI queues setting the software software_id up on the local git
	// server and CI.
	SetupCI(software_id int) (*LocalTaskModel, error)
	// TeardownCI stops the local CI of the software software_id.
	TeardownCI(software_id int) error
}

// LocalProviderRepository is the [LocalProviderRepositoryInterface] over
// the Admin API.
type LocalProviderRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /local_providers. The list has no paging.
func (r LocalProviderRepository) All() (*PageModel[LocalProviderModel], error) {
	return list[LocalProviderModel](r.Client, "/local_providers", nil, ListQuery{})
}

// Enable posts POST /local_providers/{key}/enable.
func (r LocalProviderRepository) Enable(key string) (*LocalTaskModel, error) {
	return api.Send[LocalTaskModel](r.Client, http.MethodPost, "/local_providers/"+url.PathEscape(key)+"/enable", nil)
}

// Disable posts POST /local_providers/{key}/disable.
func (r LocalProviderRepository) Disable(key string) (*LocalTaskModel, error) {
	return api.Send[LocalTaskModel](r.Client, http.MethodPost, "/local_providers/"+url.PathEscape(key)+"/disable", nil)
}

// Destroy sends DELETE /local_providers/{key}, with purge=true when purge
// is set, and returns the queued task.
func (r LocalProviderRepository) Destroy(key string, purge bool) (*LocalTaskModel, error) {
	var query url.Values
	if purge {
		query = url.Values{"purge": {"true"}}
	}
	var envelope api.Envelope[LocalTaskModel]
	if err := r.Client.Admin(http.MethodDelete, "/local_providers/"+url.PathEscape(key), query, nil, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

// Credentials gets GET /local_providers/{key}/credentials.
func (r LocalProviderRepository) Credentials(key string) (*LocalCredentialsModel, error) {
	return api.Get[LocalCredentialsModel](r.Client, "/local_providers/"+url.PathEscape(key)+"/credentials", nil)
}

// SetupCI posts POST /softwares/{softwareId}/local_ci.
func (r LocalProviderRepository) SetupCI(software_id int) (*LocalTaskModel, error) {
	return api.Send[LocalTaskModel](r.Client, http.MethodPost, "/softwares/"+strconv.Itoa(software_id)+"/local_ci", nil)
}

// TeardownCI sends DELETE /softwares/{softwareId}/local_ci.
func (r LocalProviderRepository) TeardownCI(software_id int) error {
	return api.Delete(r.Client, "/softwares/"+strconv.Itoa(software_id)+"/local_ci", nil)
}
