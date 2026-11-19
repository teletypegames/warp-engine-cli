package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestPipelineController_Index_PassesFiltersAndPrints(t *testing.T) {
	var got_filter domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockPipelineService{AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PipelineModel], error) {
		got_filter, got_query = filter, query
		return &domain.PageModel[domain.PipelineModel]{Items: []domain.PipelineModel{{ID: 1, FullName: "games/pong", Platform: "love", Active: true}}}, nil
	}}

	stdout, _, err := execute(cli.Router{PipelineController: cli.PipelineController{PipelineService: svc}},
		"pipeline", "list", "--platform", "love", "--active", "--sort", "-lastPipelineAt")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"platform": "love", "active": "true"}, got_filter)
	assert.Equal(t, domain.ListQuery{Sort: "-lastPipelineAt"}, got_query)
	assert.Contains(t, stdout, "games/pong")
}

func TestPipelineController_Show_BadIDIsUsageError(t *testing.T) {
	svc := &mocksvc.MockPipelineService{ShowFunc: func(pipeline_id int) (*domain.PipelineModel, error) {
		t.Fatal("must not reach the service")
		return nil, nil
	}}

	_, _, err := execute(cli.Router{PipelineController: cli.PipelineController{PipelineService: svc}}, "pipelines", "show", "pong")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestPipelineController_Update_SendsSoftwareReferenceAndPlatform(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockPipelineService{UpdateFunc: func(pipeline_id int, attributes domain.Attributes) (*domain.PipelineModel, error) {
		assert.Equal(t, 4, pipeline_id)
		got = attributes
		return &domain.PipelineModel{ID: 4}, nil
	}}
	router := cli.Router{PipelineController: cli.PipelineController{PipelineService: svc}}

	_, _, err := execute(router, "pipeline", "update", "4", "--software", "pong", "--platform", "unknown")
	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": "pong", "platform": "unknown"}, got)

	_, _, err = execute(router, "pipeline", "update", "4", "--software", "")
	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": nil}, got)
}

func TestPipelineController_Sync_PrintsChanges(t *testing.T) {
	svc := &mocksvc.MockPipelineService{SyncFunc: func() (*domain.PipelineSyncModel, error) {
		return &domain.PipelineSyncModel{
			Created:     []domain.PipelineModel{{ID: 2, FullName: "games/new"}},
			Deactivated: []domain.PipelineModel{{ID: 3, FullName: "games/old"}},
		}, nil
	}}

	stdout, _, err := execute(cli.Router{PipelineController: cli.PipelineController{PipelineService: svc}}, "pipeline", "sync")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "created")
	assert.Contains(t, stdout, "games/new")
	assert.Contains(t, stdout, "deactivated")
}

func TestPipelineController_Runs_PassesPage(t *testing.T) {
	svc := &mocksvc.MockPipelineService{RunsFunc: func(pipeline_id int, page int) (*domain.PageModel[domain.PipelineRunModel], error) {
		assert.Equal(t, 4, pipeline_id)
		assert.Equal(t, 2, page)
		return &domain.PageModel[domain.PipelineRunModel]{Items: []domain.PipelineRunModel{{Number: 17, Status: text("success"), Message: text("Release 1.0\n\nbody")}}}, nil
	}}

	stdout, _, err := execute(cli.Router{PipelineController: cli.PipelineController{PipelineService: svc}}, "pipeline", "runs", "4", "--page", "2")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "17")
	assert.Contains(t, stdout, "Release 1.0")
	assert.NotContains(t, stdout, "body")
}

func TestPipelineController_Trigger_PassesBranch(t *testing.T) {
	svc := &mocksvc.MockPipelineService{TriggerFunc: func(pipeline_id int, branch string) (*domain.PipelineRunModel, error) {
		assert.Equal(t, "develop", branch)
		return &domain.PipelineRunModel{Number: 18, Status: text("pending")}, nil
	}}

	stdout, _, err := execute(cli.Router{PipelineController: cli.PipelineController{PipelineService: svc}}, "pipeline", "trigger", "4", "--branch", "develop")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "18")
	assert.Contains(t, stdout, "pending")
}
