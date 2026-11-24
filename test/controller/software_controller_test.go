package controller_test

import (
	"encoding/json"
	"errors"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareController_Index_PrintsTableAndPassesFilters(t *testing.T) {
	var got_filter domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockSoftwareService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareModel], error) {
			got_filter, got_query = filter, query
			return &domain.PageModel[domain.SoftwareModel]{
				Items: []domain.SoftwareModel{{ID: 7, Name: "pong", Title: "Pong", Platform: "love", Status: text("released")}},
				Page:  1, TotalPages: 2, Total: 30,
			}, nil
		},
	}

	stdout, stderr, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}},
		"software", "list", "--platform", "love", "--query", "po", "--per-page", "5", "--sort", "-name")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"platform": "love", "q": "po"}, got_filter)
	assert.Equal(t, domain.ListQuery{PerPage: 5, Sort: "-name"}, got_query)
	assert.Contains(t, stdout, "pong")
	assert.Contains(t, stdout, "released")
	assert.Contains(t, stderr, "Page 1 of 2")
}

func TestSoftwareController_Index_JSONOutputIsTheItems(t *testing.T) {
	svc := &mocksvc.MockSoftwareService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareModel], error) {
			return &domain.PageModel[domain.SoftwareModel]{Items: []domain.SoftwareModel{{ID: 7, Name: "pong"}}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}}, "software", "list", "-o", "json")

	assert.NoError(t, err)
	var items []domain.SoftwareModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &items))
	assert.Equal(t, 7, items[0].ID)
}

func TestSoftwareController_Show_PassesReferenceAndInclude(t *testing.T) {
	var got_reference string
	var got_include []string
	svc := &mocksvc.MockSoftwareService{
		ShowFunc: func(reference string, include []string) (*domain.SoftwareModel, error) {
			got_reference, got_include = reference, include
			return &domain.SoftwareModel{ID: 7, Name: "pong", Title: "Pong"}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}},
		"software", "show", "pong", "--include", "releases,images")

	assert.NoError(t, err)
	assert.Equal(t, "pong", got_reference)
	assert.Equal(t, []string{"releases", "images"}, got_include)
	assert.Contains(t, stdout, "Pong")
}

func TestSoftwareController_Update_SendsOnlyGivenFieldsAndNullForEmpty(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareService{
		UpdateFunc: func(reference string, attributes domain.Attributes) (*domain.SoftwareModel, error) {
			got = attributes
			return &domain.SoftwareModel{ID: 7}, nil
		},
	}

	_, _, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}},
		"software", "update", "7", "--title", "Pong II", "--author", "", "--highlighted")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"title": "Pong II", "author": nil, "highlighted": true}, got)
}

func TestSoftwareController_Create_UnknownPlatformIsAUsageError(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareService{CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareModel, error) {
		got = attributes
		return nil, &domain.ValueError{Field: "platform", Value: "atari", Allowed: []string{"c64"}}
	}}

	_, _, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}},
		"software", "create", "--name", "x", "--title", "X", "--platform", "atari")

	assert.Equal(t, "atari", got["platform"])
	assert.ErrorContains(t, err, "platform takes one of c64")
	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSoftwareController_Scaffold_SendsTheApp(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareService{ScaffoldFunc: func(attributes domain.Attributes) (*domain.ScaffoldModel, error) {
		got = attributes
		return &domain.ScaffoldModel{}, nil
	}}

	_, _, err := execute(cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}},
		"software", "scaffold", "--name", "tool", "--title", "Tool", "--platform", "golang", "--app", "cli")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"name": "tool", "title": "Tool", "platform": "golang", "app": "cli"}, got)
}

func TestSoftwareController_Complete_PlatformsAndApps(t *testing.T) {
	platforms := &mocksvc.MockBuildPlatformService{
		PlatformsFunc: func() []string { return []string{"c64", "golang", "godot"} },
		AppsFunc: func(platform string) ([]string, bool) {
			assert.Equal(t, "golang", platform)
			return []string{"cli", "gui"}, true
		},
	}
	router := cli.Router{
		SoftwareController:      cli.SoftwareController{SoftwareService: &mocksvc.MockSoftwareService{}},
		BuildPlatformController: cli.BuildPlatformController{BuildPlatformService: platforms},
	}

	stdout, _, err := execute(router, "__complete", "software", "create", "--platform", "go")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "golang\ngodot\n")

	stdout, _, err = execute(router, "__complete", "software", "scaffold", "--platform", "golang", "--app", "")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "cli\ngui\n")
}

func TestSoftwareController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockSoftwareService{DestroyFunc: func(reference string) error { called = true; return nil }}
	router := cli.Router{SoftwareController: cli.SoftwareController{SoftwareService: svc}}

	_, _, err := execute(router, "software", "delete", "7")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "software", "delete", "7", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSoftwareController_Show_ConnectErrorStopsTheCommand(t *testing.T) {
	auth := &mocksvc.MockAuthService{ConnectFunc: func(selection domain.Selection) (*domain.ProfileModel, error) {
		return nil, domain.ErrNoProfile
	}}
	router := cli.Router{
		AuthController:     cli.AuthController{AuthService: auth, Session: &cli.Session{}},
		SoftwareController: cli.SoftwareController{SoftwareService: &mocksvc.MockSoftwareService{}},
	}

	_, _, err := execute(router, "software", "show", "7")

	assert.True(t, errors.Is(err, domain.ErrNoProfile))
	assert.Equal(t, cli.ExitAuth, cli.ExitCode(err))
}
