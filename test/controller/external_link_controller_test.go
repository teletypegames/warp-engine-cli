package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestExternalLinkController_Index_JSONOutputIsTheLinks(t *testing.T) {
	svc := &mocksvc.MockExternalLinkService{
		AllFunc: func(reference string) ([]domain.ExternalLinkModel, error) {
			assert.Equal(t, "pong", reference)
			return []domain.ExternalLinkModel{{ID: 2, Label: "Source", URL: "https://git.example.org/pong"}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ExternalLinkController: cli.ExternalLinkController{ExternalLinkService: svc}}, "links", "list", "pong", "-o", "json")

	assert.NoError(t, err)
	var links []domain.ExternalLinkModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &links))
	assert.Equal(t, "Source", links[0].Label)
}

func TestExternalLinkController_Update_SendsOnlyGivenFields(t *testing.T) {
	var got_id string
	var got domain.Attributes
	svc := &mocksvc.MockExternalLinkService{
		UpdateFunc: func(id string, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
			got_id, got = id, attributes
			return &domain.ExternalLinkModel{ID: 2, Label: "Code"}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ExternalLinkController: cli.ExternalLinkController{ExternalLinkService: svc}},
		"link", "update", "2", "--label", "Code")

	assert.NoError(t, err)
	assert.Equal(t, "2", got_id)
	assert.Equal(t, domain.Attributes{"label": "Code"}, got)
	assert.Contains(t, stdout, "Code")
}

func TestExternalLinkController_Create_SendsLabelAndURL(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockExternalLinkService{
		CreateFunc: func(reference string, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
			got = attributes
			return &domain.ExternalLinkModel{ID: 3}, nil
		},
	}

	_, _, err := execute(cli.Router{ExternalLinkController: cli.ExternalLinkController{ExternalLinkService: svc}},
		"link", "create", "pong", "--label", "Video", "--url", "https://v.example.org/1")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"label": "Video", "url": "https://v.example.org/1"}, got)
}

func TestExternalLinkController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockExternalLinkService{DestroyFunc: func(id string) error { called = true; return nil }}
	router := cli.Router{ExternalLinkController: cli.ExternalLinkController{ExternalLinkService: svc}}

	_, _, err := execute(router, "link", "delete", "2")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "link", "delete", "2", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
}
