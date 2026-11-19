package domain

import "errors"

// SyncPeerServiceInterface manages the sync peers: the other engines this
// one pushes releases to or takes them from.
type SyncPeerServiceInterface interface {
	// All lists the sync peers matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SyncPeerModel], error)
	// Show returns the sync peer sync_peer_id.
	Show(sync_peer_id int) (*SyncPeerModel, error)
	// Create creates a sync peer.
	Create(attributes Attributes) (*SyncPeerModel, error)
	// Update changes the sync peer sync_peer_id.
	Update(sync_peer_id int, attributes Attributes) (*SyncPeerModel, error)
	// Destroy deletes the sync peer sync_peer_id.
	Destroy(sync_peer_id int) error
	// Check asks whether the sync peer sync_peer_id accepts this engine.
	Check(sync_peer_id int) (*SyncCheckModel, error)
	// InboundToken issues a new inbound token for the sync peer sync_peer_id.
	InboundToken(sync_peer_id int) (*SyncInboundTokenModel, error)
	// Push queues the releases the sync peer sync_peer_id may be missing.
	Push(sync_peer_id int, attributes Attributes) (*SyncPushModel, error)
	// Run queues the pushes to every outbound peer, or to the peers in
	// attributes["peerIds"].
	Run(attributes Attributes) (*SyncPushModel, error)
}

// SyncPeerService is the [SyncPeerServiceInterface] over a
// [SyncPeerRepositoryInterface].
type SyncPeerService struct {
	// SyncPeerRepository talks to the sync_peers resource.
	SyncPeerRepository SyncPeerRepositoryInterface
}

// All lists the sync peers.
func (s SyncPeerService) All(filter Filter, query ListQuery) (*PageModel[SyncPeerModel], error) {
	return s.SyncPeerRepository.All(filter, query)
}

// Show gets the sync peer sync_peer_id.
func (s SyncPeerService) Show(sync_peer_id int) (*SyncPeerModel, error) {
	return s.SyncPeerRepository.FindByID(sync_peer_id)
}

// Create checks that attributes hold a label — the server would refuse
// without one — and creates the sync peer.
func (s SyncPeerService) Create(attributes Attributes) (*SyncPeerModel, error) {
	if label, _ := attributes["label"].(string); label == "" {
		return nil, errors.New("a sync peer needs a label")
	}
	return s.SyncPeerRepository.Create(attributes)
}

// Update changes the sync peer sync_peer_id; empty attributes are
// [ErrNothingToChange].
func (s SyncPeerService) Update(sync_peer_id int, attributes Attributes) (*SyncPeerModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.SyncPeerRepository.Update(sync_peer_id, attributes)
}

// Destroy deletes the sync peer sync_peer_id.
func (s SyncPeerService) Destroy(sync_peer_id int) error {
	return s.SyncPeerRepository.Destroy(sync_peer_id)
}

// Check asks the sync peer sync_peer_id.
func (s SyncPeerService) Check(sync_peer_id int) (*SyncCheckModel, error) {
	return s.SyncPeerRepository.Check(sync_peer_id)
}

// InboundToken issues a new inbound token for the sync peer sync_peer_id.
func (s SyncPeerService) InboundToken(sync_peer_id int) (*SyncInboundTokenModel, error) {
	return s.SyncPeerRepository.InboundToken(sync_peer_id)
}

// Push queues the pushes to the sync peer sync_peer_id; attributes may be
// empty (the newest release of each software, filtered as configured).
func (s SyncPeerService) Push(sync_peer_id int, attributes Attributes) (*SyncPushModel, error) {
	return s.SyncPeerRepository.Push(sync_peer_id, attributes)
}

// Run queues the pushes to several peers.
func (s SyncPeerService) Run(attributes Attributes) (*SyncPushModel, error) {
	return s.SyncPeerRepository.Run(attributes)
}
