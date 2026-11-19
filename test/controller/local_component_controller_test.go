package controller_test

import (
	"errors"
	"testing"
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

// localRouter returns a router with the local controllers over components
// and operations.
func localRouter(components *mocksvc.MockLocalComponentService, operations *mocksvc.MockLocalOperationService) cli.Router {
	return cli.Router{
		LocalComponentController: cli.LocalComponentController{LocalComponentService: components, LocalOperationService: operations},
		LocalOperationController: cli.LocalOperationController{LocalOperationService: operations},
	}
}

func TestLocalComponentController_Index_PrintsComponents(t *testing.T) {
	components := &mocksvc.MockLocalComponentService{AllFunc: func() (*domain.PageModel[domain.LocalComponentModel], error) {
		return &domain.PageModel[domain.LocalComponentModel]{Items: []domain.LocalComponentModel{{Key: "woodpecker", State: "ready", Requires: []string{"forge"}}}}, nil
	}}

	stdout, _, err := execute(localRouter(components, &mocksvc.MockLocalOperationService{}), "local", "components")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "woodpecker")
	assert.Contains(t, stdout, "forge")
}

func TestLocalComponentController_Enable_WaitsAndPrintsSteps(t *testing.T) {
	components := &mocksvc.MockLocalComponentService{EnableFunc: func(key string) (*domain.LocalOperationModel, error) {
		assert.Equal(t, "gitea", key)
		return &domain.LocalOperationModel{ID: 3, State: "queued"}, nil
	}}
	var got_timeout time.Duration
	operations := &mocksvc.MockLocalOperationService{WaitFunc: func(operation domain.LocalOperationModel, timeout time.Duration, progress func(domain.LocalOperationModel)) (*domain.LocalOperationModel, error) {
		got_timeout = timeout
		progress(operation)
		progress(domain.LocalOperationModel{ID: 3, State: "running", Steps: []domain.LocalOperationStepModel{{At: "x", Text: "pull gitea"}}})
		final := domain.LocalOperationModel{ID: 3, State: "succeeded", Steps: []domain.LocalOperationStepModel{{At: "x", Text: "pull gitea"}, {At: "y", Text: "gitea is ready"}}}
		progress(final)
		return &final, nil
	}}

	stdout, stderr, err := execute(localRouter(components, operations), "local", "enable", "gitea", "--timeout", "2m")

	assert.NoError(t, err)
	assert.Equal(t, 2*time.Minute, got_timeout)
	assert.Equal(t, "x  pull gitea\ny  gitea is ready\n", stderr)
	assert.Contains(t, stdout, "succeeded")
}

func TestLocalComponentController_Enable_NoWaitReturnsQueued(t *testing.T) {
	components := &mocksvc.MockLocalComponentService{EnableFunc: func(key string) (*domain.LocalOperationModel, error) {
		return &domain.LocalOperationModel{ID: 3, State: "queued"}, nil
	}}
	operations := &mocksvc.MockLocalOperationService{WaitFunc: func(operation domain.LocalOperationModel, timeout time.Duration, progress func(domain.LocalOperationModel)) (*domain.LocalOperationModel, error) {
		t.Fatal("--no-wait must not wait")
		return nil, nil
	}}

	stdout, _, err := execute(localRouter(components, operations), "local", "enable", "gitea", "--no-wait", "-o", "id")

	assert.NoError(t, err)
	assert.Equal(t, "3\n", stdout)
}

func TestLocalComponentController_Disable_FailedOperationIsAnError(t *testing.T) {
	components := &mocksvc.MockLocalComponentService{DisableFunc: func(key string) (*domain.LocalOperationModel, error) {
		return &domain.LocalOperationModel{ID: 3, State: "queued"}, nil
	}}
	operations := &mocksvc.MockLocalOperationService{WaitFunc: func(operation domain.LocalOperationModel, timeout time.Duration, progress func(domain.LocalOperationModel)) (*domain.LocalOperationModel, error) {
		return &domain.LocalOperationModel{ID: 3, State: "failed", Error: text("boom")}, errors.New("operation 3 failed: boom")
	}}

	stdout, _, err := execute(localRouter(components, operations), "local", "disable", "gitea")

	assert.ErrorContains(t, err, "boom")
	assert.Equal(t, cli.ExitError, cli.ExitCode(err))
	assert.Contains(t, stdout, "failed")
}

func TestLocalComponentController_Remove_NeedsYesAndPassesPurge(t *testing.T) {
	var got_purge *bool
	components := &mocksvc.MockLocalComponentService{RemoveFunc: func(key string, purge bool) (*domain.LocalOperationModel, error) {
		got_purge = &purge
		return &domain.LocalOperationModel{ID: 4, State: "succeeded"}, nil
	}}
	router := localRouter(components, &mocksvc.MockLocalOperationService{})

	_, _, err := execute(router, "local", "remove", "woodpecker", "--purge")
	assert.Error(t, err, "no --yes and no answer: nothing is done")
	assert.Nil(t, got_purge)

	_, _, err = execute(router, "local", "rm", "woodpecker", "--purge", "--yes")
	assert.NoError(t, err)
	assert.True(t, *got_purge)
}

func TestLocalComponentController_Credentials_PrintsSignIn(t *testing.T) {
	components := &mocksvc.MockLocalComponentService{CredentialsFunc: func(key string) (*domain.LocalCredentialsModel, error) {
		return &domain.LocalCredentialsModel{URL: "http://127.0.0.1:7471", User: "warp-admin", Password: "secret"}, nil
	}}

	stdout, _, err := execute(localRouter(components, &mocksvc.MockLocalOperationService{}), "local", "credentials", "gitea")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "warp-admin")
	assert.Contains(t, stdout, "secret")
}

func TestLocalComponentController_SetupCI_PassesReferenceAndPrintsResult(t *testing.T) {
	components := &mocksvc.MockLocalComponentService{SetupCIFunc: func(reference string) (*domain.LocalOperationModel, error) {
		assert.Equal(t, "pong", reference)
		return &domain.LocalOperationModel{ID: 5, State: "queued"}, nil
	}}
	operations := &mocksvc.MockLocalOperationService{WaitFunc: func(operation domain.LocalOperationModel, timeout time.Duration, progress func(domain.LocalOperationModel)) (*domain.LocalOperationModel, error) {
		return &domain.LocalOperationModel{ID: 5, State: "succeeded", Result: map[string]any{"pushUrl": "http://warp@127.0.0.1:7471/games/pong.git"}}, nil
	}}

	stdout, _, err := execute(localRouter(components, operations), "local", "ci-setup", "pong")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "Result pushUrl")
	assert.Contains(t, stdout, "http://warp@127.0.0.1:7471/games/pong.git")
}

func TestLocalComponentController_TeardownCI_NeedsYes(t *testing.T) {
	var got string
	components := &mocksvc.MockLocalComponentService{TeardownCIFunc: func(reference string) error { got = reference; return nil }}
	router := localRouter(components, &mocksvc.MockLocalOperationService{})

	_, _, err := execute(router, "local", "ci-teardown", "pong")
	assert.Error(t, err, "no --yes and no answer: nothing is done")
	assert.Empty(t, got)

	stdout, _, err := execute(router, "local", "ci-teardown", "pong", "--yes")
	assert.NoError(t, err)
	assert.Equal(t, "pong", got)
	assert.Contains(t, stdout, "stopped")
}

func TestLocalOperationController_Show_WaitsOnlyWithWait(t *testing.T) {
	waited := false
	operations := &mocksvc.MockLocalOperationService{
		ShowFunc: func(operation_id int) (*domain.LocalOperationModel, error) {
			assert.Equal(t, 12, operation_id)
			return &domain.LocalOperationModel{ID: 12, State: "running"}, nil
		},
		WaitFunc: func(operation domain.LocalOperationModel, timeout time.Duration, progress func(domain.LocalOperationModel)) (*domain.LocalOperationModel, error) {
			waited = true
			return &domain.LocalOperationModel{ID: 12, State: "succeeded"}, nil
		},
	}
	router := localRouter(&mocksvc.MockLocalComponentService{}, operations)

	stdout, _, err := execute(router, "local", "operation", "12")
	assert.NoError(t, err)
	assert.False(t, waited)
	assert.Contains(t, stdout, "running")

	stdout, _, err = execute(router, "local", "operation", "12", "--wait")
	assert.NoError(t, err)
	assert.True(t, waited)
	assert.Contains(t, stdout, "succeeded")
}
