package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

// authRouter returns a router whose AuthController is over auth.
func authRouter(auth *mocksvc.MockAuthService) cli.Router {
	return cli.Router{AuthController: cli.AuthController{AuthService: auth, Session: &cli.Session{}}}
}

func TestAuthController_Login_WithoutURLSignsInToLocalWithTheTokenFile(t *testing.T) {
	var got_url string
	auth := &mocksvc.MockAuthService{
		DiscoverFunc: func(server_url string) (*domain.ServerModel, error) {
			return &domain.ServerModel{SignIn: []string{domain.SignInTokenFile}}, nil
		},
		LoginLocalFunc: func(server_url string, name string) (*domain.ProfileModel, error) {
			got_url = server_url
			return &domain.ProfileModel{Name: "local", URL: server_url, TokenFile: true}, nil
		},
	}

	stdout, _, err := execute(authRouter(auth), "login")

	assert.NoError(t, err)
	assert.Equal(t, domain.LocalURL, got_url)
	assert.Contains(t, stdout, `profile "local"`)
}

func TestAuthController_Login_DeviceFlowPrintsTheCodeAndWaits(t *testing.T) {
	var finished_with string
	auth := &mocksvc.MockAuthService{
		DiscoverFunc: func(server_url string) (*domain.ServerModel, error) {
			return &domain.ServerModel{SignIn: []string{domain.SignInDevice}}, nil
		},
		StartDeviceLoginFunc: func(server_url string, client_name string, scopes []string) (*domain.DeviceGrantModel, error) {
			assert.Equal(t, []string{"admin:read"}, scopes)
			return &domain.DeviceGrantModel{UserCode: "ABCD-1234", VerificationURLComplete: "https://games.example.org/device?code=ABCD-1234"}, nil
		},
		WaitForApprovalFunc: func(grant domain.DeviceGrantModel) (string, error) { return "tok", nil },
		FinishLoginFunc: func(server_url string, token string, name string) (*domain.ProfileModel, error) {
			finished_with = token
			return &domain.ProfileModel{Name: "games.example.org", URL: server_url}, nil
		},
	}

	_, stderr, err := execute(authRouter(auth), "login", "https://games.example.org", "--scopes", "admin:read", "--no-browser")

	assert.NoError(t, err)
	assert.Contains(t, stderr, "ABCD-1234")
	assert.Equal(t, "tok", finished_with)
}

func TestAuthController_Login_TokenFlagSkipsTheDeviceFlow(t *testing.T) {
	auth := &mocksvc.MockAuthService{
		DiscoverFunc: func(server_url string) (*domain.ServerModel, error) {
			t.Fatal("--token must not discover")
			return nil, nil
		},
	}

	_, _, err := execute(authRouter(auth), "login", "https://games.example.org", "--token", "we_123", "--name", "ci")

	assert.NoError(t, err)
}

func TestAuthController_Logout_KeepTokenDoesNotRevoke(t *testing.T) {
	var got_revoke bool
	auth := &mocksvc.MockAuthService{LogoutFunc: func(name string, revoke bool) (*domain.ProfileModel, error) {
		got_revoke = revoke
		return &domain.ProfileModel{Name: name}, nil
	}}

	_, _, err := execute(authRouter(auth), "logout", "prod", "--keep-token")

	assert.NoError(t, err)
	assert.False(t, got_revoke)
}

func TestAuthController_Whoami_PrintsTheLabel(t *testing.T) {
	label := "CI bot"
	auth := &mocksvc.MockAuthService{WhoamiFunc: func() (*domain.PrincipalModel, error) {
		return &domain.PrincipalModel{Label: &label, Scopes: []string{"admin:read"}}, nil
	}}

	stdout, _, err := execute(authRouter(auth), "whoami")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "CI bot")
	assert.Contains(t, stdout, "https://games.example.org")
}

func TestAuthController_Connect_FlagsWinOverEnvironment(t *testing.T) {
	t.Setenv("WARP_ENGINE_URL", "https://env.example.org")
	t.Setenv("WARP_ENGINE_TOKEN", "env-token")
	var got domain.Selection
	auth := &mocksvc.MockAuthService{ConnectFunc: func(selection domain.Selection) (*domain.ProfileModel, error) {
		got = selection
		return &domain.ProfileModel{URL: selection.Server}, nil
	}}

	_, _, err := execute(authRouter(auth), "whoami", "--server", "https://flag.example.org")

	assert.NoError(t, err)
	assert.Equal(t, "https://flag.example.org", got.Server)
	assert.Equal(t, "env-token", got.Token)
}
