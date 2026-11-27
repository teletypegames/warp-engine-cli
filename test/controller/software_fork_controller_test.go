package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareForkController_Index_PassesFiltersAndPrintsJSON(t *testing.T) {
	var got domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockSoftwareForkService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareForkModel], error) {
			got, got_query = filter, query
			return &domain.PageModel[domain.SoftwareForkModel]{Items: []domain.SoftwareForkModel{{ID: 4, Name: "pong-plus", Status: "pending"}}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareForkController: cli.SoftwareForkController{SoftwareForkService: svc}},
		"fork", "list", "--status", "pending", "--software", "pong", "--sort", "-checkedAt", "-o", "json")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"status": "pending", "software": "pong"}, got)
	assert.Equal(t, "-checkedAt", got_query.Sort)
	var links []domain.SoftwareForkModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &links))
	assert.Equal(t, "pong-plus", links[0].Name)
}

func TestSoftwareForkController_Show_RejectsANonNumericID(t *testing.T) {
	_, _, err := execute(cli.Router{SoftwareForkController: cli.SoftwareForkController{SoftwareForkService: &mocksvc.MockSoftwareForkService{}}},
		"fork", "show", "pong")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSoftwareForkController_Moderate_SendsTheStatus(t *testing.T) {
	var got_id int
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareForkService{
		UpdateFunc: func(software_fork_id int, attributes domain.Attributes) (*domain.SoftwareForkModel, error) {
			got_id, got = software_fork_id, attributes
			return &domain.SoftwareForkModel{ID: 4, Status: "visible"}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareForkController: cli.SoftwareForkController{SoftwareForkService: svc}},
		"fork", "moderate", "4", "--status", "visible")

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Equal(t, domain.Attributes{"status": "visible"}, got)
	assert.Contains(t, stdout, "visible")
}

func TestSoftwareForkController_Update_UnknownStatusIsUsageError(t *testing.T) {
	called := false
	svc := &mocksvc.MockSoftwareForkService{
		UpdateFunc: func(software_fork_id int, attributes domain.Attributes) (*domain.SoftwareForkModel, error) {
			called = true
			return &domain.SoftwareForkModel{}, nil
		},
	}

	_, _, err := execute(cli.Router{SoftwareForkController: cli.SoftwareForkController{SoftwareForkService: svc}},
		"fork", "update", "4", "--status", "approved")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
	assert.False(t, called)
}

func TestSoftwareForkController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockSoftwareForkService{DestroyFunc: func(software_fork_id int) error { called = true; return nil }}
	router := cli.Router{SoftwareForkController: cli.SoftwareForkController{SoftwareForkService: svc}}

	_, _, err := execute(router, "fork", "delete", "4")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "fork", "rm", "4", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSoftwareSourceController_Set_SendsURLAndRef(t *testing.T) {
	var got domain.Attributes
	var software string
	svc := &mocksvc.MockSoftwareSourceService{
		SetFunc: func(reference string, attributes domain.Attributes) (*domain.SoftwareSourceModel, error) {
			software, got = reference, attributes
			return &domain.SoftwareSourceModel{ID: 3, URL: "https://git.example.org/pong"}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareSourceController: cli.SoftwareSourceController{SoftwareSourceService: svc}},
		"source", "set", "pong", "--url", "https://git.example.org/pong", "--ref", "main")

	assert.NoError(t, err)
	assert.Equal(t, "pong", software)
	assert.Equal(t, domain.Attributes{"url": "https://git.example.org/pong", "ref": "main"}, got)
	assert.Contains(t, stdout, "https://git.example.org/pong")
}

func TestSoftwareSourceController_Delete_NeedsYes(t *testing.T) {
	called := false
	svc := &mocksvc.MockSoftwareSourceService{DestroyFunc: func(reference string) error { called = true; return nil }}
	router := cli.Router{SoftwareSourceController: cli.SoftwareSourceController{SoftwareSourceService: svc}}

	_, _, err := execute(router, "source", "delete", "pong")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "source", "rm", "pong", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSoftwareController_Scaffold_NestsTheForkOfFlags(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareService{
		ScaffoldFunc: func(attributes domain.Attributes) (*domain.ScaffoldModel, error) {
			got = attributes
			return &domain.ScaffoldModel{}, nil
		},
	}

	_, _, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}},
		"software", "scaffold", "--name", "pong-plus", "--title", "Pong+", "--platform", "love", "--license", "MIT",
		"--fork-of-engine", "https://games.example.org", "--fork-of-name", "pong")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{
		"name": "pong-plus", "title": "Pong+", "platform": "love", "license": "MIT",
		"forkOf": map[string]any{"engine": "https://games.example.org", "name": "pong"},
	}, got)
}

func TestSoftwareController_Update_SendsTheForkFields(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareService{
		UpdateFunc: func(reference string, attributes domain.Attributes) (*domain.SoftwareModel, error) {
			got = attributes
			return &domain.SoftwareModel{ID: 1}, nil
		},
	}

	_, _, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}},
		"software", "update", "pong", "--forks-allowed", "--fork-of-repo", "")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"forksAllowed": true, "forkOfRepo": nil}, got)
}
