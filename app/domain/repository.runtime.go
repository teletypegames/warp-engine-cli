package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// RuntimeRepositoryInterface is the Admin API's runtimes resource.
type RuntimeRepositoryInterface interface {
	// All lists the installed runtime versions matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[RuntimeVersionModel], error)
	// FindByID returns the runtime version runtime_id.
	FindByID(runtime_id int) (*RuntimeVersionModel, error)
	// Install installs a downloaded runtime from attributes (name, version,
	// dryRun).
	Install(attributes Attributes) (*RuntimeInstallModel, error)
}

// RuntimeRepository is the [RuntimeRepositoryInterface] over the Admin API.
type RuntimeRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /runtimes.
func (r RuntimeRepository) All(filter Filter, query ListQuery) (*PageModel[RuntimeVersionModel], error) {
	return list[RuntimeVersionModel](r.Client, "/runtimes", filter, query)
}

// FindByID gets GET /runtimes/{id}.
func (r RuntimeRepository) FindByID(runtime_id int) (*RuntimeVersionModel, error) {
	return api.Get[RuntimeVersionModel](r.Client, "/runtimes/"+strconv.Itoa(runtime_id), nil)
}

// Install posts POST /runtimes/install.
func (r RuntimeRepository) Install(attributes Attributes) (*RuntimeInstallModel, error) {
	return api.Send[RuntimeInstallModel](r.Client, http.MethodPost, "/runtimes/install", attributes)
}
