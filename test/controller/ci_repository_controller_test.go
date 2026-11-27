package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestCiRepositoryController_Index_PassesFiltersAndPrints(t *testing.T) {
	var got_filter domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockCiRepositoryService{AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.CiRepositoryModel], error) {
		got_filter, got_query = filter, query
		return &domain.PageModel[domain.CiRepositoryModel]{Items: []domain.CiRepositoryModel{{ID: 1, FullName: "games/pong", Platform: "love", Active: true}}}, nil
	}}

	stdout, _, err := execute(cli.Router{CiRepositoryController: cli.CiRepositoryController{CiRepositoryService: svc}},
		"ci-repo", "list", "--platform", "love", "--active", "--sort", "-lastPipelineAt")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"platform": "love", "active": "true"}, got_filter)
	assert.Equal(t, domain.ListQuery{Sort: "-lastPipelineAt"}, got_query)
	assert.Contains(t, stdout, "games/pong")
}

func TestCiRepositoryController_Show_BadIDIsUsageError(t *testing.T) {
	svc := &mocksvc.MockCiRepositoryService{ShowFunc: func(ci_repository_id int) (*domain.CiRepositoryModel, error) {
		t.Fatal("must not reach the service")
		return nil, nil
	}}

	_, _, err := execute(cli.Router{CiRepositoryController: cli.CiRepositoryController{CiRepositoryService: svc}}, "ci-repos", "show", "pong")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestCiRepositoryController_Update_SendsSoftwareReferenceAndPlatform(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockCiRepositoryService{UpdateFunc: func(ci_repository_id int, attributes domain.Attributes) (*domain.CiRepositoryModel, error) {
		assert.Equal(t, 4, ci_repository_id)
		got = attributes
		return &domain.CiRepositoryModel{ID: 4}, nil
	}}
	router := cli.Router{CiRepositoryController: cli.CiRepositoryController{CiRepositoryService: svc}}

	_, _, err := execute(router, "ci-repo", "update", "4", "--software", "pong", "--platform", "unknown")
	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": "pong", "platform": "unknown"}, got)

	_, _, err = execute(router, "ci-repo", "update", "4", "--software", "")
	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": nil}, got)
}

func TestCiRepositoryController_Sync_PrintsChanges(t *testing.T) {
	svc := &mocksvc.MockCiRepositoryService{SyncFunc: func() (*domain.CiRepositorySyncModel, error) {
		return &domain.CiRepositorySyncModel{
			Created:     []domain.CiRepositoryModel{{ID: 2, FullName: "games/new"}},
			Deactivated: []domain.CiRepositoryModel{{ID: 3, FullName: "games/old"}},
		}, nil
	}}

	stdout, _, err := execute(cli.Router{CiRepositoryController: cli.CiRepositoryController{CiRepositoryService: svc}}, "ci-repo", "sync")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "created")
	assert.Contains(t, stdout, "games/new")
	assert.Contains(t, stdout, "deactivated")
}

func TestCiRepositoryController_Runs_PassesPage(t *testing.T) {
	svc := &mocksvc.MockCiRepositoryService{RunsFunc: func(ci_repository_id int, page int) (*domain.PageModel[domain.PipelineRunModel], error) {
		assert.Equal(t, 4, ci_repository_id)
		assert.Equal(t, 2, page)
		return &domain.PageModel[domain.PipelineRunModel]{Items: []domain.PipelineRunModel{{Number: 17, Status: text("success"), Message: text("Release 1.0\n\nbody")}}}, nil
	}}

	stdout, _, err := execute(cli.Router{CiRepositoryController: cli.CiRepositoryController{CiRepositoryService: svc}}, "ci-repo", "runs", "4", "--page", "2")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "17")
	assert.Contains(t, stdout, "Release 1.0")
	assert.NotContains(t, stdout, "body")
}

func TestCiRepositoryController_Trigger_PassesBranch(t *testing.T) {
	svc := &mocksvc.MockCiRepositoryService{TriggerFunc: func(ci_repository_id int, branch string) (*domain.PipelineRunModel, error) {
		assert.Equal(t, "develop", branch)
		return &domain.PipelineRunModel{Number: 18, Status: text("pending")}, nil
	}}

	stdout, _, err := execute(cli.Router{CiRepositoryController: cli.CiRepositoryController{CiRepositoryService: svc}}, "ci-repo", "trigger", "4", "--branch", "develop")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "18")
	assert.Contains(t, stdout, "pending")
}
