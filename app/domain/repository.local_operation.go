package domain

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// LocalOperationRepositoryInterface is the Admin API's local_operations
// resource (WarpEngine Local only).
type LocalOperationRepositoryInterface interface {
	// FindByID returns the operation operation_id with its steps so far.
	FindByID(operation_id int) (*LocalOperationModel, error)
}

// LocalOperationRepository is the [LocalOperationRepositoryInterface] over
// the Admin API.
type LocalOperationRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// FindByID gets GET /local_operations/{id}.
func (r LocalOperationRepository) FindByID(operation_id int) (*LocalOperationModel, error) {
	return api.Get[LocalOperationModel](r.Client, "/local_operations/"+strconv.Itoa(operation_id), nil)
}
