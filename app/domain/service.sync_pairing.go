package domain

import "errors"

// SyncPairingServiceInterface manages the caller's own sync pairings: the
// inbound peers it added so that another engine may push here.
type SyncPairingServiceInterface interface {
	// All lists the pairings, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SyncPairingModel], error)
	// Show returns the pairing sync_pairing_id.
	Show(sync_pairing_id int) (*SyncPairingModel, error)
	// Create pairs and returns the pairing with its new token.
	Create(attributes Attributes) (*SyncPairingModel, error)
	// Update renames the pairing sync_pairing_id.
	Update(sync_pairing_id int, attributes Attributes) (*SyncPairingModel, error)
	// Destroy unpairs the pairing sync_pairing_id.
	Destroy(sync_pairing_id int) error
	// Token issues a new token for the pairing sync_pairing_id.
	Token(sync_pairing_id int) (*SyncPairingModel, error)
}

// SyncPairingService is the [SyncPairingServiceInterface] over a
// [SyncPairingRepositoryInterface].
type SyncPairingService struct {
	// SyncPairingRepository talks to the sync_pairings resource.
	SyncPairingRepository SyncPairingRepositoryInterface
}

// All lists the pairings.
func (s SyncPairingService) All(filter Filter, query ListQuery) (*PageModel[SyncPairingModel], error) {
	return s.SyncPairingRepository.All(filter, query)
}

// Show gets the pairing sync_pairing_id.
func (s SyncPairingService) Show(sync_pairing_id int) (*SyncPairingModel, error) {
	return s.SyncPairingRepository.FindByID(sync_pairing_id)
}

// Create checks that attributes hold a label — the server would refuse
// without one — and pairs.
func (s SyncPairingService) Create(attributes Attributes) (*SyncPairingModel, error) {
	if label, _ := attributes["label"].(string); label == "" {
		return nil, errors.New("a pairing needs a label")
	}
	return s.SyncPairingRepository.Create(attributes)
}

// Update renames the pairing sync_pairing_id; empty attributes are
// [ErrNothingToChange].
func (s SyncPairingService) Update(sync_pairing_id int, attributes Attributes) (*SyncPairingModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.SyncPairingRepository.Update(sync_pairing_id, attributes)
}

// Destroy unpairs the pairing sync_pairing_id.
func (s SyncPairingService) Destroy(sync_pairing_id int) error {
	return s.SyncPairingRepository.Destroy(sync_pairing_id)
}

// Token issues a new token for the pairing sync_pairing_id.
func (s SyncPairingService) Token(sync_pairing_id int) (*SyncPairingModel, error) {
	return s.SyncPairingRepository.Token(sync_pairing_id)
}
