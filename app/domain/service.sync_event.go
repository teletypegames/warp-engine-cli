package domain

// SyncEventServiceInterface reads the sync events: what was pushed to and
// taken from the sync peers, and where each software stands.
type SyncEventServiceInterface interface {
	// All lists the sync events matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SyncEventModel], error)
	// Show returns the sync event sync_event_id.
	Show(sync_event_id int) (*SyncEventModel, error)
	// Software lists the sync events of the software reference (an ID or
	// an exact name), in direction ("" for both).
	Software(reference string, direction string) ([]SyncEventModel, error)
	// Status returns where the software reference stands with every peer.
	Status(reference string) ([]SyncStatusModel, error)
}

// SyncEventService is the [SyncEventServiceInterface] over a
// [SyncEventRepositoryInterface]; software references are resolved through
// a [SoftwareRepositoryInterface].
type SyncEventService struct {
	// SyncEventRepository talks to the sync_events resource.
	SyncEventRepository SyncEventRepositoryInterface
	// SoftwareRepository resolves software names.
	SoftwareRepository SoftwareRepositoryInterface
}

// All lists the sync events.
func (s SyncEventService) All(filter Filter, query ListQuery) (*PageModel[SyncEventModel], error) {
	return s.SyncEventRepository.All(filter, query)
}

// Show gets the sync event sync_event_id.
func (s SyncEventService) Show(sync_event_id int) (*SyncEventModel, error) {
	return s.SyncEventRepository.FindByID(sync_event_id)
}

// Software resolves reference and lists its sync events.
func (s SyncEventService) Software(reference string, direction string) ([]SyncEventModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SyncEventRepository.AllOfSoftware(software_id, direction)
}

// Status resolves reference and returns its sync state per peer.
func (s SyncEventService) Status(reference string) ([]SyncStatusModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.SyncEventRepository.Status(software_id)
}
