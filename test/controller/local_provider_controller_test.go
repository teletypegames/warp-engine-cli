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

// localRouter returns a router with the local controllers over providers
// and tasks.
func localRouter(providers *mocksvc.MockLocalProviderService, tasks *mocksvc.MockLocalTaskService) cli.Router {
	return cli.Router{
		LocalProviderController: cli.LocalProviderController{LocalProviderService: providers, LocalTaskService: tasks},
		LocalTaskController:     cli.LocalTaskController{LocalTaskService: tasks},
	}
}

func TestLocalProviderController_Index_PrintsProviders(t *testing.T) {
	providers := &mocksvc.MockLocalProviderService{AllFunc: func() (*domain.PageModel[domain.LocalProviderModel], error) {
		return &domain.PageModel[domain.LocalProviderModel]{Items: []domain.LocalProviderModel{{Key: "woodpecker", State: "ready", Requires: []string{"forge"}}}}, nil
	}}

	stdout, _, err := execute(localRouter(providers, &mocksvc.MockLocalTaskService{}), "local", "providers")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "woodpecker")
	assert.Contains(t, stdout, "forge")
}

func TestLocalProviderController_Enable_WaitsAndPrintsSteps(t *testing.T) {
	providers := &mocksvc.MockLocalProviderService{EnableFunc: func(key string) (*domain.LocalTaskModel, error) {
		assert.Equal(t, "gitea", key)
		return &domain.LocalTaskModel{ID: 3, State: "queued"}, nil
	}}
	var got_timeout time.Duration
	tasks := &mocksvc.MockLocalTaskService{WaitFunc: func(task domain.LocalTaskModel, timeout time.Duration, progress func(domain.LocalTaskModel)) (*domain.LocalTaskModel, error) {
		got_timeout = timeout
		progress(task)
		progress(domain.LocalTaskModel{ID: 3, State: "running", Steps: []domain.LocalTaskStepModel{{At: "x", Text: "pull gitea"}}})
		final := domain.LocalTaskModel{ID: 3, State: "succeeded", Steps: []domain.LocalTaskStepModel{{At: "x", Text: "pull gitea"}, {At: "y", Text: "gitea is ready"}}}
		progress(final)
		return &final, nil
	}}

	stdout, stderr, err := execute(localRouter(providers, tasks), "local", "enable", "gitea", "--timeout", "2m")

	assert.NoError(t, err)
	assert.Equal(t, 2*time.Minute, got_timeout)
	assert.Equal(t, "x  pull gitea\ny  gitea is ready\n", stderr)
	assert.Contains(t, stdout, "succeeded")
}

func TestLocalProviderController_Enable_NoWaitReturnsQueued(t *testing.T) {
	providers := &mocksvc.MockLocalProviderService{EnableFunc: func(key string) (*domain.LocalTaskModel, error) {
		return &domain.LocalTaskModel{ID: 3, State: "queued"}, nil
	}}
	tasks := &mocksvc.MockLocalTaskService{WaitFunc: func(task domain.LocalTaskModel, timeout time.Duration, progress func(domain.LocalTaskModel)) (*domain.LocalTaskModel, error) {
		t.Fatal("--no-wait must not wait")
		return nil, nil
	}}

	stdout, _, err := execute(localRouter(providers, tasks), "local", "enable", "gitea", "--no-wait", "-o", "id")

	assert.NoError(t, err)
	assert.Equal(t, "3\n", stdout)
}

func TestLocalProviderController_Disable_FailedTaskIsAnError(t *testing.T) {
	providers := &mocksvc.MockLocalProviderService{DisableFunc: func(key string) (*domain.LocalTaskModel, error) {
		return &domain.LocalTaskModel{ID: 3, State: "queued"}, nil
	}}
	tasks := &mocksvc.MockLocalTaskService{WaitFunc: func(task domain.LocalTaskModel, timeout time.Duration, progress func(domain.LocalTaskModel)) (*domain.LocalTaskModel, error) {
		return &domain.LocalTaskModel{ID: 3, State: "failed", Error: text("boom")}, errors.New("task 3 failed: boom")
	}}

	stdout, _, err := execute(localRouter(providers, tasks), "local", "disable", "gitea")

	assert.ErrorContains(t, err, "boom")
	assert.Equal(t, cli.ExitError, cli.ExitCode(err))
	assert.Contains(t, stdout, "failed")
}

func TestLocalProviderController_Remove_NeedsYesAndPassesPurge(t *testing.T) {
	var got_purge *bool
	providers := &mocksvc.MockLocalProviderService{RemoveFunc: func(key string, purge bool) (*domain.LocalTaskModel, error) {
		got_purge = &purge
		return &domain.LocalTaskModel{ID: 4, State: "succeeded"}, nil
	}}
	router := localRouter(providers, &mocksvc.MockLocalTaskService{})

	_, _, err := execute(router, "local", "remove", "woodpecker", "--purge")
	assert.Error(t, err, "no --yes and no answer: nothing is done")
	assert.Nil(t, got_purge)

	_, _, err = execute(router, "local", "rm", "woodpecker", "--purge", "--yes")
	assert.NoError(t, err)
	assert.True(t, *got_purge)
}

func TestLocalProviderController_Credentials_PrintsSignIn(t *testing.T) {
	providers := &mocksvc.MockLocalProviderService{CredentialsFunc: func(key string) (*domain.LocalCredentialsModel, error) {
		return &domain.LocalCredentialsModel{URL: "http://127.0.0.1:7471", User: "warp-admin", Password: "secret"}, nil
	}}

	stdout, _, err := execute(localRouter(providers, &mocksvc.MockLocalTaskService{}), "local", "credentials", "gitea")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "warp-admin")
	assert.Contains(t, stdout, "secret")
}

func TestLocalProviderController_SetupCI_PassesReferenceAndPrintsResult(t *testing.T) {
	providers := &mocksvc.MockLocalProviderService{SetupCIFunc: func(reference string) (*domain.LocalTaskModel, error) {
		assert.Equal(t, "pong", reference)
		return &domain.LocalTaskModel{ID: 5, State: "queued"}, nil
	}}
	tasks := &mocksvc.MockLocalTaskService{WaitFunc: func(task domain.LocalTaskModel, timeout time.Duration, progress func(domain.LocalTaskModel)) (*domain.LocalTaskModel, error) {
		return &domain.LocalTaskModel{ID: 5, State: "succeeded", Result: map[string]any{"pushUrl": "http://warp@127.0.0.1:7471/games/pong.git"}}, nil
	}}

	stdout, _, err := execute(localRouter(providers, tasks), "local", "ci-setup", "pong")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "Result pushUrl")
	assert.Contains(t, stdout, "http://warp@127.0.0.1:7471/games/pong.git")
}

func TestLocalProviderController_TeardownCI_NeedsYes(t *testing.T) {
	var got string
	providers := &mocksvc.MockLocalProviderService{TeardownCIFunc: func(reference string) error { got = reference; return nil }}
	router := localRouter(providers, &mocksvc.MockLocalTaskService{})

	_, _, err := execute(router, "local", "ci-teardown", "pong")
	assert.Error(t, err, "no --yes and no answer: nothing is done")
	assert.Empty(t, got)

	stdout, _, err := execute(router, "local", "ci-teardown", "pong", "--yes")
	assert.NoError(t, err)
	assert.Equal(t, "pong", got)
	assert.Contains(t, stdout, "stopped")
}

func TestLocalTaskController_Show_WaitsOnlyWithWait(t *testing.T) {
	waited := false
	tasks := &mocksvc.MockLocalTaskService{
		ShowFunc: func(task_id int) (*domain.LocalTaskModel, error) {
			assert.Equal(t, 12, task_id)
			return &domain.LocalTaskModel{ID: 12, State: "running"}, nil
		},
		WaitFunc: func(task domain.LocalTaskModel, timeout time.Duration, progress func(domain.LocalTaskModel)) (*domain.LocalTaskModel, error) {
			waited = true
			return &domain.LocalTaskModel{ID: 12, State: "succeeded"}, nil
		},
	}
	router := localRouter(&mocksvc.MockLocalProviderService{}, tasks)

	stdout, _, err := execute(router, "local", "task", "12")
	assert.NoError(t, err)
	assert.False(t, waited)
	assert.Contains(t, stdout, "running")

	stdout, _, err = execute(router, "local", "task", "12", "--wait")
	assert.NoError(t, err)
	assert.True(t, waited)
	assert.Contains(t, stdout, "succeeded")
}
