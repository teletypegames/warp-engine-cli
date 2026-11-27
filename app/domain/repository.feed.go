package domain

import (
	"net/url"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// FeedRepositoryInterface is the Admin API's feeds resource.
type FeedRepositoryInterface interface {
	// All lists the configured feeds.
	All() ([]FeedModel, error)
	// Check returns what the feed key offers and leaves out.
	Check(key string) (*FeedCheckModel, error)
}

// FeedRepository is the [FeedRepositoryInterface] over the Admin API.
type FeedRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /feeds.
func (r FeedRepository) All() ([]FeedModel, error) {
	items, _, err := api.List[FeedModel](r.Client, "/feeds", nil, false)
	return items, err
}

// Check gets GET /feeds/{key}/check.
func (r FeedRepository) Check(key string) (*FeedCheckModel, error) {
	return api.Get[FeedCheckModel](r.Client, "/feeds/"+url.PathEscape(key)+"/check", nil)
}
