package domain

import (
	"errors"
	"fmt"
	"time"
)

// ErrOperationTimeout is returned by [LocalOperationService.Wait] when the
// operation did not finish in time; it goes on running on the server.
var ErrOperationTimeout = errors.New("the operation did not finish in time; it goes on on the server")

// LocalOperationPollInterval is the wait between two polls of a local
// operation.
const LocalOperationPollInterval = time.Second

// LocalOperationServiceInterface follows the background jobs of WarpEngine
// Local.
type LocalOperationServiceInterface interface {
	// Show returns the operation operation_id.
	Show(operation_id int) (*LocalOperationModel, error)
	// Wait polls operation until it has succeeded or failed, at most for
	// timeout, calling progress with every state read (the first one too).
	Wait(operation LocalOperationModel, timeout time.Duration, progress func(LocalOperationModel)) (*LocalOperationModel, error)
}

// LocalOperationService is the [LocalOperationServiceInterface] over a
// [LocalOperationRepositoryInterface].
type LocalOperationService struct {
	// LocalOperationRepository talks to the local_operations resource.
	LocalOperationRepository LocalOperationRepositoryInterface
	// Sleep waits between two polls; nil uses time.Sleep (tests replace it).
	Sleep func(time.Duration)
}

// Show gets the operation.
func (s LocalOperationService) Show(operation_id int) (*LocalOperationModel, error) {
	return s.LocalOperationRepository.FindByID(operation_id)
}

// Wait polls every [LocalOperationPollInterval] until operation is
// finished and returns its last state. A failed operation returns that
// state and an error with the server's reason; running out of timeout
// (counted in the intervals slept) returns the last state and
// [ErrOperationTimeout]. progress may be nil.
func (s LocalOperationService) Wait(operation LocalOperationModel, timeout time.Duration, progress func(LocalOperationModel)) (*LocalOperationModel, error) {
	sleep := s.Sleep
	if sleep == nil {
		sleep = time.Sleep
	}
	current := &operation
	waited := time.Duration(0)
	for {
		if progress != nil {
			progress(*current)
		}
		if current.Finished() {
			break
		}
		if waited >= timeout {
			return current, ErrOperationTimeout
		}
		sleep(LocalOperationPollInterval)
		waited += LocalOperationPollInterval
		next, err := s.LocalOperationRepository.FindByID(operation.ID)
		if err != nil {
			return current, err
		}
		current = next
	}
	if current.State == LocalOperationFailed {
		reason := "no reason given"
		if current.Error != nil && *current.Error != "" {
			reason = *current.Error
		}
		return current, fmt.Errorf("operation %d (%s %s) failed: %s", current.ID, current.Kind, current.Target, reason)
	}
	return current, nil
}
