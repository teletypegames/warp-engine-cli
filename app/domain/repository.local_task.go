package domain

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// LocalTaskRepositoryInterface is the Admin API's local_tasks
// resource (WarpEngine Local only).
type LocalTaskRepositoryInterface interface {
	// FindByID returns the task task_id with its steps so far.
	FindByID(task_id int) (*LocalTaskModel, error)
}

// LocalTaskRepository is the [LocalTaskRepositoryInterface] over
// the Admin API.
type LocalTaskRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// FindByID gets GET /local_tasks/{id}.
func (r LocalTaskRepository) FindByID(task_id int) (*LocalTaskModel, error) {
	return api.Get[LocalTaskModel](r.Client, "/local_tasks/"+strconv.Itoa(task_id), nil)
}
