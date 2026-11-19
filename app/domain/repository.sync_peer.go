package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// SyncPeerRepositoryInterface is the Admin API's sync_peers resource.
type SyncPeerRepositoryInterface interface {
	// All lists the sync peers matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SyncPeerModel], error)
	// FindByID returns the sync peer sync_peer_id.
	FindByID(sync_peer_id int) (*SyncPeerModel, error)
	// Create creates a sync peer from attributes.
	Create(attributes Attributes) (*SyncPeerModel, error)
	// Update changes the attributes of the sync peer sync_peer_id.
	Update(sync_peer_id int, attributes Attributes) (*SyncPeerModel, error)
	// Destroy deletes the sync peer sync_peer_id and revokes its inbound token.
	Destroy(sync_peer_id int) error
	// Check asks whether the sync peer sync_peer_id accepts this engine.
	Check(sync_peer_id int) (*SyncCheckModel, error)
	// InboundToken issues a new inbound token for the sync peer
	// sync_peer_id; the old one stops working.
	InboundToken(sync_peer_id int) (*SyncInboundTokenModel, error)
	// Push queues the releases the sync peer sync_peer_id may be missing,
	// chosen by attributes (the push criteria and dryRun).
	Push(sync_peer_id int, attributes Attributes) (*SyncPushModel, error)
	// Run queues the pushes to several peers, chosen by attributes (the
	// push criteria, dryRun and peerIds).
	Run(attributes Attributes) (*SyncPushModel, error)
}

// SyncPeerRepository is the [SyncPeerRepositoryInterface] over the Admin API.
type SyncPeerRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /sync_peers.
func (r SyncPeerRepository) All(filter Filter, query ListQuery) (*PageModel[SyncPeerModel], error) {
	return list[SyncPeerModel](r.Client, "/sync_peers", filter, query)
}

// FindByID gets GET /sync_peers/{id}.
func (r SyncPeerRepository) FindByID(sync_peer_id int) (*SyncPeerModel, error) {
	return api.Get[SyncPeerModel](r.Client, syncPeerPath(sync_peer_id), nil)
}

// Create posts POST /sync_peers.
func (r SyncPeerRepository) Create(attributes Attributes) (*SyncPeerModel, error) {
	return api.Send[SyncPeerModel](r.Client, http.MethodPost, "/sync_peers", attributes)
}

// Update sends PATCH /sync_peers/{id}.
func (r SyncPeerRepository) Update(sync_peer_id int, attributes Attributes) (*SyncPeerModel, error) {
	return api.Send[SyncPeerModel](r.Client, http.MethodPatch, syncPeerPath(sync_peer_id), attributes)
}

// Destroy sends DELETE /sync_peers/{id}.
func (r SyncPeerRepository) Destroy(sync_peer_id int) error {
	return api.Delete(r.Client, syncPeerPath(sync_peer_id), nil)
}

// Check posts POST /sync_peers/{id}/check.
func (r SyncPeerRepository) Check(sync_peer_id int) (*SyncCheckModel, error) {
	return api.Send[SyncCheckModel](r.Client, http.MethodPost, syncPeerPath(sync_peer_id)+"/check", nil)
}

// InboundToken posts POST /sync_peers/{id}/inbound_token.
func (r SyncPeerRepository) InboundToken(sync_peer_id int) (*SyncInboundTokenModel, error) {
	return api.Send[SyncInboundTokenModel](r.Client, http.MethodPost, syncPeerPath(sync_peer_id)+"/inbound_token", nil)
}

// Push posts POST /sync_peers/{id}/push with attributes.
func (r SyncPeerRepository) Push(sync_peer_id int, attributes Attributes) (*SyncPushModel, error) {
	return api.Send[SyncPushModel](r.Client, http.MethodPost, syncPeerPath(sync_peer_id)+"/push", attributes)
}

// Run posts POST /sync_peers/run with attributes.
func (r SyncPeerRepository) Run(attributes Attributes) (*SyncPushModel, error) {
	return api.Send[SyncPushModel](r.Client, http.MethodPost, "/sync_peers/run", attributes)
}

// syncPeerPath is the path of the sync peer sync_peer_id.
func syncPeerPath(sync_peer_id int) string {
	return "/sync_peers/" + strconv.Itoa(sync_peer_id)
}
