package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestProfileController_Index_JSONHidesTheTokens(t *testing.T) {
	svc := &mocksvc.MockProfileService{AllFunc: func() (*[]domain.ProfileModel, string, error) {
		return &[]domain.ProfileModel{{Name: "prod", URL: "https://games.example.org", Token: "secret"}}, "prod", nil
	}}
	auth := &mocksvc.MockAuthService{ConnectFunc: func(domain.Selection) (*domain.ProfileModel, error) {
		t.Fatal("profile list must not connect")
		return nil, nil
	}}
	router := cli.Router{AuthController: cli.AuthController{AuthService: auth, Session: &cli.Session{}}, ProfileController: cli.ProfileController{ProfileService: svc}}

	stdout, _, err := execute(router, "profile", "list", "-o", "json")

	assert.NoError(t, err)
	assert.NotContains(t, stdout, "secret")
	var profiles []domain.ProfileModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &profiles))
	assert.Equal(t, "prod", profiles[0].Name)
}

func TestProfileController_Use_PassesTheName(t *testing.T) {
	var got string
	svc := &mocksvc.MockProfileService{UseFunc: func(name string) error { got = name; return nil }}

	_, _, err := execute(cli.Router{ProfileController: cli.ProfileController{ProfileService: svc}}, "profile", "use", "prod")

	assert.NoError(t, err)
	assert.Equal(t, "prod", got)
}
