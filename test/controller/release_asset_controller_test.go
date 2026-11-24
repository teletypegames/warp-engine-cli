package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestReleaseAssetController_Index_PrintsTheAssets(t *testing.T) {
	svc := &mocksvc.MockReleaseAssetService{
		AllFunc: func(release string) ([]domain.ReleaseAssetModel, error) {
			assert.Equal(t, "42", release)
			return []domain.ReleaseAssetModel{{ID: 7, ReleaseID: 42, Kind: "cartridge", Path: "/pong/pong.tic"}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ReleaseAssetController: cli.ReleaseAssetController{ReleaseAssetService: svc}}, "assets", "ls", "42")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "/pong/pong.tic")
}

func TestReleaseAssetController_Create_SendsKindAndPath(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockReleaseAssetService{
		CreateFunc: func(release string, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
			got = attributes
			return &domain.ReleaseAssetModel{ID: 8}, nil
		},
	}

	_, _, err := execute(cli.Router{ReleaseAssetController: cli.ReleaseAssetController{ReleaseAssetService: svc}},
		"asset", "create", "42", "--kind", "html", "--path", "/pong/html/index.html")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"kind": "html", "path": "/pong/html/index.html"}, got)
}

func TestReleaseAssetController_Create_UnknownKindIsAUsageError(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockReleaseAssetService{CreateFunc: func(release string, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
		got = attributes
		return nil, &domain.ValueError{Field: "kind", Value: "amiga"}
	}}

	_, _, err := execute(cli.Router{ReleaseAssetController: cli.ReleaseAssetController{ReleaseAssetService: svc}},
		"asset", "create", "42", "--kind", "amiga", "--path", "/x")

	assert.Equal(t, "amiga", got["kind"])
	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestReleaseAssetController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockReleaseAssetService{DestroyFunc: func(id string) error { called = true; return nil }}
	router := cli.Router{ReleaseAssetController: cli.ReleaseAssetController{ReleaseAssetService: svc}}

	_, _, err := execute(router, "asset", "delete", "8")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "asset", "delete", "8", "-y")
	assert.NoError(t, err)
	assert.True(t, called)
}
