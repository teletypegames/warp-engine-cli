package service_test

import (
	"errors"
	"testing"
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

// noSleep counts the waits of a [domain.LocalOperationService] without
// waiting.
func noSleep(count *int) func(time.Duration) {
	return func(time.Duration) { *count++ }
}

func TestLocalOperationService_Wait_PollsUntilSucceeded(t *testing.T) {
	states := []string{"running", "running", "succeeded"}
	polls := 0
	repo := &mockrepo.MockLocalOperationRepository{FindByIDFunc: func(operation_id int) (*domain.LocalOperationModel, error) {
		assert.Equal(t, 7, operation_id)
		state := states[polls]
		polls++
		return &domain.LocalOperationModel{ID: 7, State: state, Steps: make([]domain.LocalOperationStepModel, polls)}, nil
	}}
	sleeps := 0
	svc := domain.LocalOperationService{LocalOperationRepository: repo, Sleep: noSleep(&sleeps)}
	var seen []string

	final, err := svc.Wait(domain.LocalOperationModel{ID: 7, State: "queued"}, time.Minute, func(operation domain.LocalOperationModel) {
		seen = append(seen, operation.State)
	})

	assert.NoError(t, err)
	assert.Equal(t, "succeeded", final.State)
	assert.Equal(t, 3, polls)
	assert.Equal(t, 3, sleeps)
	assert.Equal(t, []string{"queued", "running", "running", "succeeded"}, seen)
}

func TestLocalOperationService_Wait_FinishedOperationIsNotPolled(t *testing.T) {
	repo := &mockrepo.MockLocalOperationRepository{FindByIDFunc: func(operation_id int) (*domain.LocalOperationModel, error) {
		t.Fatal("a finished operation must not be polled")
		return nil, nil
	}}
	sleeps := 0
	svc := domain.LocalOperationService{LocalOperationRepository: repo, Sleep: noSleep(&sleeps)}

	final, err := svc.Wait(domain.LocalOperationModel{ID: 7, State: "succeeded"}, time.Minute, nil)

	assert.NoError(t, err)
	assert.Equal(t, 7, final.ID)
	assert.Zero(t, sleeps)
}

func TestLocalOperationService_Wait_FailedIsAnErrorWithTheReason(t *testing.T) {
	reason := "docker is not running"
	repo := &mockrepo.MockLocalOperationRepository{FindByIDFunc: func(operation_id int) (*domain.LocalOperationModel, error) {
		return &domain.LocalOperationModel{ID: 7, Kind: "enable", Target: "gitea", State: "failed", Error: &reason}, nil
	}}
	sleeps := 0
	svc := domain.LocalOperationService{LocalOperationRepository: repo, Sleep: noSleep(&sleeps)}

	final, err := svc.Wait(domain.LocalOperationModel{ID: 7, State: "running"}, time.Minute, nil)

	assert.ErrorContains(t, err, "docker is not running")
	assert.ErrorContains(t, err, "enable gitea")
	assert.Equal(t, "failed", final.State)
}

func TestLocalOperationService_Wait_TimesOut(t *testing.T) {
	polls := 0
	repo := &mockrepo.MockLocalOperationRepository{FindByIDFunc: func(operation_id int) (*domain.LocalOperationModel, error) {
		polls++
		return &domain.LocalOperationModel{ID: 7, State: "running"}, nil
	}}
	sleeps := 0
	svc := domain.LocalOperationService{LocalOperationRepository: repo, Sleep: noSleep(&sleeps)}

	final, err := svc.Wait(domain.LocalOperationModel{ID: 7, State: "queued"}, 3*domain.LocalOperationPollInterval, nil)

	assert.ErrorIs(t, err, domain.ErrOperationTimeout)
	assert.Equal(t, "running", final.State)
	assert.Equal(t, 3, polls)
	assert.Equal(t, 3, sleeps)
}

func TestLocalOperationService_Wait_PollErrorStops(t *testing.T) {
	repo := &mockrepo.MockLocalOperationRepository{FindByIDFunc: func(operation_id int) (*domain.LocalOperationModel, error) {
		return nil, errors.New("connection refused")
	}}
	sleeps := 0
	svc := domain.LocalOperationService{LocalOperationRepository: repo, Sleep: noSleep(&sleeps)}

	final, err := svc.Wait(domain.LocalOperationModel{ID: 7, State: "queued"}, time.Minute, nil)

	assert.ErrorContains(t, err, "connection refused")
	assert.Equal(t, "queued", final.State)
}
