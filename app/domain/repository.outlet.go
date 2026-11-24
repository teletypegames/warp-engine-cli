package domain

import (
	"net/url"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// OutletRepositoryInterface is the Admin API's outlets resource.
type OutletRepositoryInterface interface {
	// All lists the configured outlets.
	All() ([]OutletModel, error)
	// Check returns what the outlet key offers and leaves out.
	Check(key string) (*OutletCheckModel, error)
}

// OutletRepository is the [OutletRepositoryInterface] over the Admin API.
type OutletRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /outlets.
func (r OutletRepository) All() ([]OutletModel, error) {
	items, _, err := api.List[OutletModel](r.Client, "/outlets", nil, false)
	return items, err
}

// Check gets GET /outlets/{key}/check.
func (r OutletRepository) Check(key string) (*OutletCheckModel, error) {
	return api.Get[OutletCheckModel](r.Client, "/outlets/"+url.PathEscape(key)+"/check", nil)
}
