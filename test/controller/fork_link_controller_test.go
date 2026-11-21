package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestForkLinkController_Index_PassesFiltersAndPrintsJSON(t *testing.T) {
	var got domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockForkLinkService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ForkLinkModel], error) {
			got, got_query = filter, query
			return &domain.PageModel[domain.ForkLinkModel]{Items: []domain.ForkLinkModel{{ID: 4, ChildName: "pong-plus", Status: "pending"}}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ForkLinkController: cli.ForkLinkController{ForkLinkService: svc}},
		"fork", "list", "--status", "pending", "--software", "pong", "--sort", "-checkedAt", "-o", "json")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"status": "pending", "software": "pong"}, got)
	assert.Equal(t, "-checkedAt", got_query.Sort)
	var links []domain.ForkLinkModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &links))
	assert.Equal(t, "pong-plus", links[0].ChildName)
}

func TestForkLinkController_Show_RejectsANonNumericID(t *testing.T) {
	_, _, err := execute(cli.Router{ForkLinkController: cli.ForkLinkController{ForkLinkService: &mocksvc.MockForkLinkService{}}},
		"fork", "show", "pong")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestForkLinkController_Moderate_SendsTheStatus(t *testing.T) {
	var got_id int
	var got domain.Attributes
	svc := &mocksvc.MockForkLinkService{
		UpdateFunc: func(fork_link_id int, attributes domain.Attributes) (*domain.ForkLinkModel, error) {
			got_id, got = fork_link_id, attributes
			return &domain.ForkLinkModel{ID: 4, Status: "visible"}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ForkLinkController: cli.ForkLinkController{ForkLinkService: svc}},
		"fork", "moderate", "4", "--status", "visible")

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Equal(t, domain.Attributes{"status": "visible"}, got)
	assert.Contains(t, stdout, "visible")
}

func TestForkLinkController_Update_UnknownStatusIsUsageError(t *testing.T) {
	called := false
	svc := &mocksvc.MockForkLinkService{
		UpdateFunc: func(fork_link_id int, attributes domain.Attributes) (*domain.ForkLinkModel, error) {
			called = true
			return &domain.ForkLinkModel{}, nil
		},
	}

	_, _, err := execute(cli.Router{ForkLinkController: cli.ForkLinkController{ForkLinkService: svc}},
		"fork", "update", "4", "--status", "approved")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
	assert.False(t, called)
}

func TestForkLinkController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockForkLinkService{DestroyFunc: func(fork_link_id int) error { called = true; return nil }}
	router := cli.Router{ForkLinkController: cli.ForkLinkController{ForkLinkService: svc}}

	_, _, err := execute(router, "fork", "delete", "4")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "fork", "rm", "4", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestExternalLinkController_Create_SendsKindAndRef(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockExternalLinkService{
		CreateFunc: func(reference string, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
			got = attributes
			return &domain.ExternalLinkModel{ID: 3, Kind: "source"}, nil
		},
	}

	_, _, err := execute(cli.Router{ExternalLinkController: cli.ExternalLinkController{ExternalLinkService: svc}},
		"link", "create", "pong", "--label", "Source", "--url", "https://git.example.org/pong", "--kind", "source", "--ref", "main")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"label": "Source", "url": "https://git.example.org/pong", "kind": "source", "ref": "main"}, got)
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
