package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestExtraController_Show_PassesSoftwareAndAssociation(t *testing.T) {
	svc := &mocksvc.MockExtraService{
		ShowFunc: func(reference string, association string) (*domain.ExtraModel, error) {
			assert.Equal(t, "pong", reference)
			assert.Equal(t, "seo_metadata", association)
			return &domain.ExtraModel{Association: association, Value: map[string]any{"title": "Pong"}}, nil
		},
	}

	stdout, _, err := execute(cli.Router{ExtraController: cli.ExtraController{ExtraService: svc}}, "extra", "show", "pong", "seo_metadata")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "seo_metadata")
	assert.Contains(t, stdout, `"title": "Pong"`)
}

func TestExtraController_Update_ValueBecomesTheBodysValue(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockExtraService{
		UpdateFunc: func(reference string, association string, attributes domain.Attributes) (*domain.ExtraModel, error) {
			got = attributes
			return &domain.ExtraModel{Association: association}, nil
		},
	}

	_, _, err := execute(cli.Router{ExtraController: cli.ExtraController{ExtraService: svc}},
		"extras", "update", "pong", "seo_metadata", "--value", `[{"id":3,"_destroy":true}]`)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"value": []any{map[string]any{"id": float64(3), "_destroy": true}}}, got)
}

func TestExtraController_Update_BadValueIsAUsageError(t *testing.T) {
	_, _, err := execute(cli.Router{ExtraController: cli.ExtraController{ExtraService: &mocksvc.MockExtraService{}}},
		"extra", "update", "pong", "seo_metadata", "--value", "{nope")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}
