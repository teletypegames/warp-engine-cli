package domain

import "errors"

// SoftwareStatusServiceInterface manages the software statuses, the vocabulary of a software's Status.
type SoftwareStatusServiceInterface interface {
	// All lists the software statuses matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[SoftwareStatusModel], error)
	// Show returns the status status_id.
	Show(status_id int) (*SoftwareStatusModel, error)
	// Create creates a status.
	Create(attributes Attributes) (*SoftwareStatusModel, error)
	// Update changes the status status_id.
	Update(status_id int, attributes Attributes) (*SoftwareStatusModel, error)
	// Destroy deletes the status status_id.
	Destroy(status_id int) error
}

// SoftwareStatusService is the [SoftwareStatusServiceInterface] over a [SoftwareStatusRepositoryInterface].
type SoftwareStatusService struct {
	// SoftwareStatusRepository talks to the software statuses resource.
	SoftwareStatusRepository SoftwareStatusRepositoryInterface
}

// All lists the software statuses.
func (s SoftwareStatusService) All(filter Filter, query ListQuery) (*PageModel[SoftwareStatusModel], error) {
	return s.SoftwareStatusRepository.All(filter, query)
}

// Show gets the status status_id.
func (s SoftwareStatusService) Show(status_id int) (*SoftwareStatusModel, error) {
	return s.SoftwareStatusRepository.FindByID(status_id)
}

// Create checks that attributes hold key and label — the server would
// refuse without them — and creates the status.
func (s SoftwareStatusService) Create(attributes Attributes) (*SoftwareStatusModel, error) {
	for _, key := range []string{"key", "label"} {
		if !attributes.Has(key) || attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("a status needs a " + key)
		}
	}
	return s.SoftwareStatusRepository.Create(attributes)
}

// Update changes the status status_id; empty attributes are
// [ErrNothingToChange], since there would be nothing to send.
func (s SoftwareStatusService) Update(status_id int, attributes Attributes) (*SoftwareStatusModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.SoftwareStatusRepository.Update(status_id, attributes)
}

// Destroy deletes the status status_id.
func (s SoftwareStatusService) Destroy(status_id int) error {
	return s.SoftwareStatusRepository.Destroy(status_id)
}
