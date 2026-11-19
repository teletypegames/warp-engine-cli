package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// ReleaseAssetRepositoryInterface is the Admin API's release assets resource.
type ReleaseAssetRepositoryInterface interface {
	// All lists the assets of the release release_id (the list has no paging).
	All(release_id int) ([]ReleaseAssetModel, error)
	// FindByID returns the asset asset_id.
	FindByID(asset_id int) (*ReleaseAssetModel, error)
	// Create adds an asset to the release release_id from attributes.
	Create(release_id int, attributes Attributes) (*ReleaseAssetModel, error)
	// Update changes the attributes of the asset asset_id.
	Update(asset_id int, attributes Attributes) (*ReleaseAssetModel, error)
	// Destroy deletes the asset asset_id.
	Destroy(asset_id int) error
}

// ReleaseAssetRepository is the [ReleaseAssetRepositoryInterface] over the Admin API.
type ReleaseAssetRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /releases/{releaseId}/assets.
func (r ReleaseAssetRepository) All(release_id int) ([]ReleaseAssetModel, error) {
	page, err := list[ReleaseAssetModel](r.Client, "/releases/"+strconv.Itoa(release_id)+"/assets", nil, ListQuery{})
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

// FindByID gets GET /assets/{id}.
func (r ReleaseAssetRepository) FindByID(asset_id int) (*ReleaseAssetModel, error) {
	return api.Get[ReleaseAssetModel](r.Client, "/assets/"+strconv.Itoa(asset_id), nil)
}

// Create posts POST /releases/{releaseId}/assets.
func (r ReleaseAssetRepository) Create(release_id int, attributes Attributes) (*ReleaseAssetModel, error) {
	return api.Send[ReleaseAssetModel](r.Client, http.MethodPost, "/releases/"+strconv.Itoa(release_id)+"/assets", attributes)
}

// Update sends PATCH /assets/{id}.
func (r ReleaseAssetRepository) Update(asset_id int, attributes Attributes) (*ReleaseAssetModel, error) {
	return api.Send[ReleaseAssetModel](r.Client, http.MethodPatch, "/assets/"+strconv.Itoa(asset_id), attributes)
}

// Destroy sends DELETE /assets/{id}.
func (r ReleaseAssetRepository) Destroy(asset_id int) error {
	return api.Delete(r.Client, "/assets/"+strconv.Itoa(asset_id), nil)
}
