package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareLinkController_Index_JSONOutputIsTheLinks(t *testing.T) {
	svc := &mocksvc.MockSoftwareLinkService{
		AllFunc: func(reference string) ([]domain.SoftwareLinkModel, error) {
			assert.Equal(t, "pong", reference)
			return []domain.SoftwareLinkModel{{ID: 2, Label: "Source", URL: "https://git.example.org/pong"}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareLinkController: cli.SoftwareLinkController{SoftwareLinkService: svc}}, "links", "list", "pong", "-o", "json")

	assert.NoError(t, err)
	var links []domain.SoftwareLinkModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &links))
	assert.Equal(t, "Source", links[0].Label)
}

func TestSoftwareLinkController_Update_SendsOnlyGivenFields(t *testing.T) {
	var got_id string
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareLinkService{
		UpdateFunc: func(id string, attributes domain.Attributes) (*domain.SoftwareLinkModel, error) {
			got_id, got = id, attributes
			return &domain.SoftwareLinkModel{ID: 2, Label: "Code"}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareLinkController: cli.SoftwareLinkController{SoftwareLinkService: svc}},
		"link", "update", "2", "--label", "Code")

	assert.NoError(t, err)
	assert.Equal(t, "2", got_id)
	assert.Equal(t, domain.Attributes{"label": "Code"}, got)
	assert.Contains(t, stdout, "Code")
}

func TestSoftwareLinkController_Create_SendsLabelAndURL(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareLinkService{
		CreateFunc: func(reference string, attributes domain.Attributes) (*domain.SoftwareLinkModel, error) {
			got = attributes
			return &domain.SoftwareLinkModel{ID: 3}, nil
		},
	}

	_, _, err := execute(cli.Router{SoftwareLinkController: cli.SoftwareLinkController{SoftwareLinkService: svc}},
		"link", "create", "pong", "--label", "Video", "--url", "https://v.example.org/1")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"label": "Video", "url": "https://v.example.org/1"}, got)
}

func TestSoftwareLinkController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockSoftwareLinkService{DestroyFunc: func(id string) error { called = true; return nil }}
	router := cli.Router{SoftwareLinkController: cli.SoftwareLinkController{SoftwareLinkService: svc}}

	_, _, err := execute(router, "link", "delete", "2")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "link", "delete", "2", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
}
