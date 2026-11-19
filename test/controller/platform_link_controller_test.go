package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestPlatformLinkController_Index_PassesPlatformAndPrintsJSON(t *testing.T) {
	var got_filter domain.Filter
	svc := &mocksvc.MockPlatformLinkService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.PlatformLinkModel], error) {
			got_filter = filter
			return &domain.PageModel[domain.PlatformLinkModel]{
				Items: []domain.PlatformLinkModel{{ID: 3, Platform: "c64", Name: "VICE", URL: "https://vice.example.org"}},
			}, nil
		},
	}

	stdout, _, err := execute(cli.Router{PlatformLinkController: cli.PlatformLinkController{PlatformLinkService: svc}},
		"platform-links", "list", "--platform", "c64", "-o", "json")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"platform": "c64"}, got_filter)
	var items []domain.PlatformLinkModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &items))
	assert.Equal(t, "VICE", items[0].Name)
}

func TestPlatformLinkController_Create_SendsTheFields(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockPlatformLinkService{CreateFunc: func(attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
		got = attributes
		return &domain.PlatformLinkModel{ID: 3, Platform: "c64", Name: "VICE"}, nil
	}}

	stdout, _, err := execute(cli.Router{PlatformLinkController: cli.PlatformLinkController{PlatformLinkService: svc}},
		"platform-link", "create", "--platform", "c64", "--name", "VICE", "--url", "https://vice.example.org", "--position", "0")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"platform": "c64", "name": "VICE", "url": "https://vice.example.org", "position": 0}, got)
	assert.Contains(t, stdout, "VICE")
}

func TestPlatformLinkController_Create_RejectsUnknownPlatform(t *testing.T) {
	svc := &mocksvc.MockPlatformLinkService{CreateFunc: func(attributes domain.Attributes) (*domain.PlatformLinkModel, error) {
		t.Fatal("must not reach the service")
		return nil, nil
	}}

	_, _, err := execute(cli.Router{PlatformLinkController: cli.PlatformLinkController{PlatformLinkService: svc}},
		"platform-link", "create", "--platform", "atari", "--name", "X", "--url", "https://x.example.org")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestPlatformLinkController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	var got_id int
	svc := &mocksvc.MockPlatformLinkService{DestroyFunc: func(platform_link_id int) error { got_id = platform_link_id; return nil }}
	router := cli.Router{PlatformLinkController: cli.PlatformLinkController{PlatformLinkService: svc}}

	_, _, err := execute(router, "platform-link", "delete", "3")
	assert.Error(t, err)
	assert.Zero(t, got_id)

	_, _, err = execute(router, "platform-link", "delete", "3", "--yes")
	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
}
