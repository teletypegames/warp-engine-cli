package domain

import (
	"errors"
	"fmt"
	"time"
)

// ErrTaskTimeout is returned by [LocalTaskService.Wait] when the
// task did not finish in time; it goes on running on the server.
var ErrTaskTimeout = errors.New("the task did not finish in time; it goes on on the server")

// LocalTaskPollInterval is the wait between two polls of a local
// task.
const LocalTaskPollInterval = time.Second

// LocalTaskServiceInterface follows the background jobs of WarpEngine
// Local.
type LocalTaskServiceInterface interface {
	// Show returns the task task_id.
	Show(task_id int) (*LocalTaskModel, error)
	// Wait polls task until it has succeeded or failed, at most for
	// timeout, calling progress with every state read (the first one too).
	Wait(task LocalTaskModel, timeout time.Duration, progress func(LocalTaskModel)) (*LocalTaskModel, error)
}

// LocalTaskService is the [LocalTaskServiceInterface] over a
// [LocalTaskRepositoryInterface].
type LocalTaskService struct {
	// LocalTaskRepository talks to the local_tasks resource.
	LocalTaskRepository LocalTaskRepositoryInterface
	// Sleep waits between two polls; nil uses time.Sleep (tests replace it).
	Sleep func(time.Duration)
}

// Show gets the task.
func (s LocalTaskService) Show(task_id int) (*LocalTaskModel, error) {
	return s.LocalTaskRepository.FindByID(task_id)
}

// Wait polls every [LocalTaskPollInterval] until task is
// finished and returns its last state. A failed task returns that
// state and an error with the server's reason; running out of timeout
// (counted in the intervals slept) returns the last state and
// [ErrTaskTimeout]. progress may be nil.
func (s LocalTaskService) Wait(task LocalTaskModel, timeout time.Duration, progress func(LocalTaskModel)) (*LocalTaskModel, error) {
	sleep := s.Sleep
	if sleep == nil {
		sleep = time.Sleep
	}
	current := &task
	waited := time.Duration(0)
	for {
		if progress != nil {
			progress(*current)
		}
		if current.Finished() {
			break
		}
		if waited >= timeout {
			return current, ErrTaskTimeout
		}
		sleep(LocalTaskPollInterval)
		waited += LocalTaskPollInterval
		next, err := s.LocalTaskRepository.FindByID(task.ID)
		if err != nil {
			return current, err
		}
		current = next
	}
	if current.State == LocalTaskFailed {
		reason := "no reason given"
		if current.Error != nil && *current.Error != "" {
			reason = *current.Error
		}
		return current, fmt.Errorf("task %d (%s %s) failed: %s", current.ID, current.Kind, current.Target, reason)
	}
	return current, nil
}
