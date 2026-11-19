package domain

import (
	"net/url"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SyncEventRepositoryInterface is the Admin API's sync_events resource and
// the sync views of a software.
type SyncEventRepositoryInterface interface {
	// All lists the sync events matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SyncEventModel], error)
	// FindByID returns the sync event sync_event_id.
	FindByID(sync_event_id int) (*SyncEventModel, error)
	// AllOfSoftware lists the sync events of the software software_id,
	// in direction ("" for both).
	AllOfSoftware(software_id int, direction string) ([]SyncEventModel, error)
	// Status returns where the software software_id stands with every peer.
	Status(software_id int) ([]SyncStatusModel, error)
}

// SyncEventRepository is the [SyncEventRepositoryInterface] over the Admin API.
type SyncEventRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /sync_events.
func (r SyncEventRepository) All(filter Filter, query ListQuery) (*PageModel[SyncEventModel], error) {
	return list[SyncEventModel](r.Client, "/sync_events", filter, query)
}

// FindByID gets GET /sync_events/{id}.
func (r SyncEventRepository) FindByID(sync_event_id int) (*SyncEventModel, error) {
	return api.Get[SyncEventModel](r.Client, "/sync_events/"+strconv.Itoa(sync_event_id), nil)
}

// AllOfSoftware gets GET /softwares/{softwareId}/sync_events?direction=….
func (r SyncEventRepository) AllOfSoftware(software_id int, direction string) ([]SyncEventModel, error) {
	query := url.Values{}
	if direction != "" {
		query.Set("direction", direction)
	}
	items, _, err := api.List[SyncEventModel](r.Client, "/softwares/"+strconv.Itoa(software_id)+"/sync_events", query, false)
	return items, err
}

// Status gets GET /softwares/{softwareId}/sync_status.
func (r SyncEventRepository) Status(software_id int) ([]SyncStatusModel, error) {
	items, _, err := api.List[SyncStatusModel](r.Client, "/softwares/"+strconv.Itoa(software_id)+"/sync_status", nil, false)
	return items, err
}
