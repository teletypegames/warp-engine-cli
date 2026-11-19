package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

// tokenRouter returns a router with only the token commands over svc.
func tokenRouter(svc *mocksvc.MockApplicationTokenService) cli.Router {
	return cli.Router{ApplicationTokenController: cli.ApplicationTokenController{ApplicationTokenService: svc}}
}

func TestApplicationTokenController_Index_PassesFilters(t *testing.T) {
	var got domain.Filter
	svc := &mocksvc.MockApplicationTokenService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ApplicationTokenModel], error) {
			got = filter
			return &domain.PageModel[domain.ApplicationTokenModel]{Items: []domain.ApplicationTokenModel{
				{ID: 2, Name: "hello-c64 CI", TokenPrefix: "304f835d", Scopes: []string{"update", "upload"}},
			}}, nil
		},
	}

	stdout, _, err := execute(tokenRouter(svc), "tokens", "list", "--kind", "publishing", "--state", "active")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"kind": "publishing", "state": "active"}, got)
	assert.Contains(t, stdout, "update, upload")
	assert.NotContains(t, stdout, "shown only this once")
}

func TestApplicationTokenController_Create_SendsFlagsAndPrintsTheTokenOnce(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockApplicationTokenService{
		CreateFunc: func(attributes domain.Attributes) (*domain.ApplicationTokenModel, error) {
			got = attributes
			return &domain.ApplicationTokenModel{ID: 5, Name: "ci", TokenPrefix: "aa11", PlainToken: text("aa11-whole-token")}, nil
		},
	}

	stdout, _, err := execute(tokenRouter(svc), "token", "create", "--name", "ci", "--scopes", "admin:read,admin:write",
		"--expires-at", "2027-01-01T00:00:00Z", "--owner-id", "1")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{
		"name": "ci", "scopes": []string{"admin:read", "admin:write"}, "expiresAt": "2027-01-01T00:00:00Z", "ownerId": 1,
	}, got)
	assert.Contains(t, stdout, "shown only this once")
	assert.Contains(t, stdout, "    aa11-whole-token\n")
}

func TestApplicationTokenController_Create_JSONHasThePlainToken(t *testing.T) {
	svc := &mocksvc.MockApplicationTokenService{
		CreateFunc: func(attributes domain.Attributes) (*domain.ApplicationTokenModel, error) {
			return &domain.ApplicationTokenModel{ID: 5, PlainToken: text("aa11-whole-token")}, nil
		},
	}

	stdout, _, err := execute(tokenRouter(svc), "token", "create", "--name", "ci", "--scopes", "admin:read", "-o", "json")

	assert.NoError(t, err)
	var token map[string]any
	assert.NoError(t, json.Unmarshal([]byte(stdout), &token))
	assert.Equal(t, "aa11-whole-token", token["plainToken"])
}

func TestApplicationTokenController_Update_EmptyExpiryClearsIt(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockApplicationTokenService{
		UpdateFunc: func(application_token_id int, attributes domain.Attributes) (*domain.ApplicationTokenModel, error) {
			got = attributes
			return &domain.ApplicationTokenModel{ID: application_token_id}, nil
		},
	}

	_, _, err := execute(tokenRouter(svc), "token", "update", "5", "--expires-at", "", "--unrestricted=false")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"expiresAt": nil, "unrestricted": false}, got)
}

func TestApplicationTokenController_RotateAndRevoke_NeedYes(t *testing.T) {
	calls := []string{}
	svc := &mocksvc.MockApplicationTokenService{
		RotateFunc: func(application_token_id int) (*domain.ApplicationTokenModel, error) {
			calls = append(calls, "rotate")
			return &domain.ApplicationTokenModel{ID: 6, PlainToken: text("bb22-rotated")}, nil
		},
		DestroyFunc: func(application_token_id int) error { calls = append(calls, "destroy"); return nil },
	}

	_, _, err := execute(tokenRouter(svc), "token", "rotate", "5")
	assert.Error(t, err)
	_, _, err = execute(tokenRouter(svc), "token", "revoke", "5")
	assert.Error(t, err)
	assert.Empty(t, calls)

	stdout, _, err := execute(tokenRouter(svc), "token", "rotate", "5", "--yes")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "    bb22-rotated\n")
	stdout, _, err = execute(tokenRouter(svc), "token", "delete", "5", "--yes")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "revoked")
	assert.Equal(t, []string{"rotate", "destroy"}, calls)
}

func TestApplicationTokenController_GlobalTokenFlagStillWorks(t *testing.T) {
	var got domain.Selection
	auth := &mocksvc.MockAuthService{ConnectFunc: func(selection domain.Selection) (*domain.ProfileModel, error) {
		got = selection
		return &domain.ProfileModel{Name: "x", URL: "https://games.example.org"}, nil
	}}
	router := tokenRouter(&mocksvc.MockApplicationTokenService{})
	router.AuthController = cli.AuthController{AuthService: auth, Session: &cli.Session{}}

	_, _, err := execute(router, "token", "list", "--token", "we_bearer")

	assert.NoError(t, err)
	assert.Equal(t, "we_bearer", got.Token)
}
