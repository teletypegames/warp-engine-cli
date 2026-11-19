package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestReleaseController_Index_PassesSoftwareFiltersAndSort(t *testing.T) {
	var got_reference string
	var got_filter domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockReleaseService{
		AllFunc: func(reference string, filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ReleaseModel], error) {
			got_reference, got_filter, got_query = reference, filter, query
			return &domain.PageModel[domain.ReleaseModel]{
				Items: []domain.ReleaseModel{{ID: 42, Version: "1.2.0-beta.1", Channel: "beta", CommitSha: text("3fb3ed8aa1b2c3d4")}},
				Page:  1, TotalPages: 1, Total: 1,
			}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ReleaseController: cli.ReleaseController{ReleaseService: svc}},
		"release", "list", "pong", "--channel", "beta", "--commit-sha", "3fb3ed8", "--sort", "-version")

	assert.NoError(t, err)
	assert.Equal(t, "pong", got_reference)
	assert.Equal(t, domain.Filter{"channel": "beta", "commitSha": "3fb3ed8"}, got_filter)
	assert.Equal(t, domain.ListQuery{Sort: "-version"}, got_query)
	assert.Contains(t, stdout, "1.2.0-beta.1")
	assert.Contains(t, stdout, "3fb3ed8aa1")
}

func TestReleaseController_Create_SendsTheVersion(t *testing.T) {
	var got_reference string
	var got domain.Attributes
	svc := &mocksvc.MockReleaseService{
		CreateFunc: func(reference string, attributes domain.Attributes) (*domain.ReleaseModel, error) {
			got_reference, got = reference, attributes
			return &domain.ReleaseModel{ID: 43, Version: "1.3.0", Channel: "stable"}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ReleaseController: cli.ReleaseController{ReleaseService: svc}},
		"release", "create", "pong", "--version", "1.3.0", "-o", "id")

	assert.NoError(t, err)
	assert.Equal(t, "pong", got_reference)
	assert.Equal(t, domain.Attributes{"version": "1.3.0"}, got)
	assert.Equal(t, "43\n", stdout)
}

func TestReleaseController_Show_PrintsTheAssets(t *testing.T) {
	svc := &mocksvc.MockReleaseService{
		ShowFunc: func(id string) (*domain.ReleaseModel, error) {
			assert.Equal(t, "42", id)
			return &domain.ReleaseModel{ID: 42, Version: "1.2.0", Assets: []domain.ReleaseAssetModel{{ID: 7, Kind: "html", Path: "/pong/index.html"}}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ReleaseController: cli.ReleaseController{ReleaseService: svc}}, "release", "show", "42")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "#7 html /pong/index.html")
}

func TestReleaseController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockReleaseService{DestroyFunc: func(id string) error { called = true; return nil }}
	router := cli.Router{ReleaseController: cli.ReleaseController{ReleaseService: svc}}

	_, _, err := execute(router, "release", "delete", "42")
	assert.Error(t, err)
	assert.False(t, called)

	stdout, _, err := execute(router, "release", "rm", "42", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
	assert.Contains(t, stdout, "Release 42 deleted.")
}

func TestReleaseController_Index_NeedsTheSoftware(t *testing.T) {
	_, _, err := execute(cli.Router{ReleaseController: cli.ReleaseController{ReleaseService: &mocksvc.MockReleaseService{}}}, "release", "list")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}
