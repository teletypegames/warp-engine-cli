package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SyncPairingRepositoryInterface is the Admin API's sync_pairings resource.
type SyncPairingRepositoryInterface interface {
	// All lists the caller's pairings, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SyncPairingModel], error)
	// FindByID returns the pairing sync_pairing_id.
	FindByID(sync_pairing_id int) (*SyncPairingModel, error)
	// Create pairs: adds an inbound peer from attributes and issues its token.
	Create(attributes Attributes) (*SyncPairingModel, error)
	// Update renames the pairing sync_pairing_id.
	Update(sync_pairing_id int, attributes Attributes) (*SyncPairingModel, error)
	// Destroy unpairs: deletes the pairing sync_pairing_id and revokes its token.
	Destroy(sync_pairing_id int) error
	// Token issues a new token for the pairing sync_pairing_id; the old one
	// stops working.
	Token(sync_pairing_id int) (*SyncPairingModel, error)
}

// SyncPairingRepository is the [SyncPairingRepositoryInterface] over the
// Admin API.
type SyncPairingRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /sync_pairings.
func (r SyncPairingRepository) All(filter Filter, query ListQuery) (*PageModel[SyncPairingModel], error) {
	return list[SyncPairingModel](r.Client, "/sync_pairings", filter, query)
}

// FindByID gets GET /sync_pairings/{id}.
func (r SyncPairingRepository) FindByID(sync_pairing_id int) (*SyncPairingModel, error) {
	return api.Get[SyncPairingModel](r.Client, syncPairingPath(sync_pairing_id), nil)
}

// Create posts POST /sync_pairings.
func (r SyncPairingRepository) Create(attributes Attributes) (*SyncPairingModel, error) {
	return api.Send[SyncPairingModel](r.Client, http.MethodPost, "/sync_pairings", attributes)
}

// Update sends PATCH /sync_pairings/{id}.
func (r SyncPairingRepository) Update(sync_pairing_id int, attributes Attributes) (*SyncPairingModel, error) {
	return api.Send[SyncPairingModel](r.Client, http.MethodPatch, syncPairingPath(sync_pairing_id), attributes)
}

// Destroy sends DELETE /sync_pairings/{id}.
func (r SyncPairingRepository) Destroy(sync_pairing_id int) error {
	return api.Delete(r.Client, syncPairingPath(sync_pairing_id), nil)
}

// Token posts POST /sync_pairings/{id}/token.
func (r SyncPairingRepository) Token(sync_pairing_id int) (*SyncPairingModel, error) {
	return api.Send[SyncPairingModel](r.Client, http.MethodPost, syncPairingPath(sync_pairing_id)+"/token", nil)
}

// syncPairingPath is the path of the pairing sync_pairing_id.
func syncPairingPath(sync_pairing_id int) string {
	return "/sync_pairings/" + strconv.Itoa(sync_pairing_id)
}
