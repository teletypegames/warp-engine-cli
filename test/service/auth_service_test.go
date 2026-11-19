package service_test

import (
	"testing"
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

// authService returns an AuthService over the given mocks that never sleeps.
func authService(auth *mockrepo.MockAuthRepository, profiles *mockrepo.MockProfileRepository, local *mockrepo.MockLocalTokenRepository) domain.AuthService {
	return domain.AuthService{AuthRepository: auth, ProfileRepository: profiles, LocalTokenRepository: local, Sleep: func(time.Duration) {}}
}

// config returns a configuration with profiles, current being the current one.
func config(current string, profiles ...domain.ProfileModel) *domain.ConfigModel {
	result := &domain.ConfigModel{Current: current, Profiles: map[string]domain.ProfileModel{}}
	for _, profile := range profiles {
		result.Profiles[profile.Name] = profile
	}
	return result
}

func TestAuthService_Connect_UsesTheCurrentProfile(t *testing.T) {
	auth := &mockrepo.MockAuthRepository{}
	profiles := &mockrepo.MockProfileRepository{Config: config("prod", domain.ProfileModel{Name: "prod", URL: "https://games.example.org", Token: "t1"})}
	svc := authService(auth, profiles, &mockrepo.MockLocalTokenRepository{})

	profile, err := svc.Connect(domain.Selection{})

	assert.NoError(t, err)
	assert.Equal(t, "prod", profile.Name)
	assert.Equal(t, [][2]string{{"https://games.example.org", "t1"}}, auth.Used)
}

func TestAuthService_Connect_TokenOverridesTheProfile(t *testing.T) {
	auth := &mockrepo.MockAuthRepository{}
	profiles := &mockrepo.MockProfileRepository{Config: config("prod", domain.ProfileModel{Name: "prod", URL: "https://games.example.org", Token: "t1"})}
	svc := authService(auth, profiles, &mockrepo.MockLocalTokenRepository{})

	_, err := svc.Connect(domain.Selection{Token: "t2"})

	assert.NoError(t, err)
	assert.Equal(t, "t2", auth.Used[0][1])
}

func TestAuthService_Connect_ServerTakesTheTokenOfAMatchingProfile(t *testing.T) {
	auth := &mockrepo.MockAuthRepository{}
	profiles := &mockrepo.MockProfileRepository{Config: config("", domain.ProfileModel{Name: "x", URL: "https://games.example.org", Token: "t1"})}
	svc := authService(auth, profiles, &mockrepo.MockLocalTokenRepository{})

	_, err := svc.Connect(domain.Selection{Server: "https://games.example.org/"})

	assert.NoError(t, err)
	assert.Equal(t, [][2]string{{"https://games.example.org", "t1"}}, auth.Used)
}

func TestAuthService_Connect_LoopbackServerReadsTheLocalToken(t *testing.T) {
	auth := &mockrepo.MockAuthRepository{}
	svc := authService(auth, &mockrepo.MockProfileRepository{}, &mockrepo.MockLocalTokenRepository{Token: "local"})

	profile, err := svc.Connect(domain.Selection{Server: "http://127.0.0.1:7470"})

	assert.NoError(t, err)
	assert.Equal(t, "local", profile.Name)
	assert.Equal(t, "local", auth.Used[0][1])
}

func TestAuthService_Connect_TokenFileProfileReadsTheFile(t *testing.T) {
	auth := &mockrepo.MockAuthRepository{}
	profiles := &mockrepo.MockProfileRepository{Config: config("local", domain.ProfileModel{Name: "local", URL: domain.LocalURL, TokenFile: true})}
	svc := authService(auth, profiles, &mockrepo.MockLocalTokenRepository{Token: "from-file"})

	_, err := svc.Connect(domain.Selection{})

	assert.NoError(t, err)
	assert.Equal(t, "from-file", auth.Used[0][1])
}

func TestAuthService_Connect_NothingIsErrNoProfile(t *testing.T) {
	svc := authService(&mockrepo.MockAuthRepository{}, &mockrepo.MockProfileRepository{}, &mockrepo.MockLocalTokenRepository{})

	_, err := svc.Connect(domain.Selection{})

	assert.ErrorIs(t, err, domain.ErrNoProfile)
}

func TestAuthService_WaitForApproval_ReturnsTheTokenWhenApproved(t *testing.T) {
	polls := 0
	auth := &mockrepo.MockAuthRepository{PollDeviceGrantFunc: func(device_code string) (*domain.TokenPollModel, error) {
		polls++
		assert.Equal(t, "dev-1", device_code)
		if polls < 3 {
			return &domain.TokenPollModel{State: domain.TokenPending}, nil
		}
		return &domain.TokenPollModel{State: domain.TokenApproved, Token: "new-token"}, nil
	}}
	svc := authService(auth, &mockrepo.MockProfileRepository{}, &mockrepo.MockLocalTokenRepository{})

	token, err := svc.WaitForApproval(domain.DeviceGrantModel{DeviceCode: "dev-1", Interval: 5, ExpiresIn: 600})

	assert.NoError(t, err)
	assert.Equal(t, "new-token", token)
	assert.Equal(t, 3, polls)
}

func TestAuthService_WaitForApproval_DeniedAndExpired(t *testing.T) {
	denied := &mockrepo.MockAuthRepository{PollDeviceGrantFunc: func(string) (*domain.TokenPollModel, error) {
		return &domain.TokenPollModel{State: domain.TokenDenied}, nil
	}}
	_, err := authService(denied, &mockrepo.MockProfileRepository{}, &mockrepo.MockLocalTokenRepository{}).
		WaitForApproval(domain.DeviceGrantModel{Interval: 1, ExpiresIn: 60})
	assert.ErrorIs(t, err, domain.ErrDenied)

	pending := &mockrepo.MockAuthRepository{}
	_, err = authService(pending, &mockrepo.MockProfileRepository{}, &mockrepo.MockLocalTokenRepository{}).
		WaitForApproval(domain.DeviceGrantModel{Interval: 1, ExpiresIn: 10})
	assert.ErrorIs(t, err, domain.ErrExpired)
}

func TestAuthService_FinishLogin_SavesTheCurrentProfile(t *testing.T) {
	label := "CI bot"
	auth := &mockrepo.MockAuthRepository{WhoamiFunc: func() (*domain.PrincipalModel, error) {
		return &domain.PrincipalModel{Label: &label, Scopes: []string{"admin:read"}}, nil
	}}
	profiles := &mockrepo.MockProfileRepository{}
	svc := authService(auth, profiles, &mockrepo.MockLocalTokenRepository{})

	profile, err := svc.FinishLogin("https://games.example.org/", "tok", "")

	assert.NoError(t, err)
	assert.Equal(t, "games.example.org", profile.Name)
	assert.Equal(t, "games.example.org", profiles.Config.Current)
	assert.Equal(t, "tok", profiles.Config.Profiles["games.example.org"].Token)
	assert.Equal(t, "CI bot", profiles.Config.Profiles["games.example.org"].Label)
}

func TestAuthService_LoginLocal_SavesNoToken(t *testing.T) {
	profiles := &mockrepo.MockProfileRepository{}
	svc := authService(&mockrepo.MockAuthRepository{}, profiles, &mockrepo.MockLocalTokenRepository{Token: "secret"})

	profile, err := svc.LoginLocal(domain.LocalURL, "")

	assert.NoError(t, err)
	assert.Equal(t, "local", profile.Name)
	assert.True(t, profiles.Config.Profiles["local"].TokenFile)
	assert.Empty(t, profiles.Config.Profiles["local"].Token)
}

func TestAuthService_LoginLocal_NoTokenFile(t *testing.T) {
	svc := authService(&mockrepo.MockAuthRepository{}, &mockrepo.MockProfileRepository{}, &mockrepo.MockLocalTokenRepository{})

	_, err := svc.LoginLocal(domain.LocalURL, "")

	assert.ErrorIs(t, err, domain.ErrNoLocalToken)
}

func TestAuthService_Logout_RevokesThenForgets(t *testing.T) {
	revoked := false
	auth := &mockrepo.MockAuthRepository{RevokeFunc: func() error { revoked = true; return nil }}
	profiles := &mockrepo.MockProfileRepository{Config: config("prod", domain.ProfileModel{Name: "prod", URL: "https://games.example.org", Token: "t1"})}
	svc := authService(auth, profiles, &mockrepo.MockLocalTokenRepository{})

	_, err := svc.Logout("", true)

	assert.NoError(t, err)
	assert.True(t, revoked)
	assert.Empty(t, profiles.Config.Profiles)
	assert.Empty(t, profiles.Config.Current)
}

func TestAuthService_Logout_LocalProfileIsNotRevoked(t *testing.T) {
	auth := &mockrepo.MockAuthRepository{RevokeFunc: func() error { t.Fatal("a Local profile must not revoke"); return nil }}
	profiles := &mockrepo.MockProfileRepository{Config: config("local", domain.ProfileModel{Name: "local", URL: domain.LocalURL, TokenFile: true})}
	svc := authService(auth, profiles, &mockrepo.MockLocalTokenRepository{Token: "x"})

	_, err := svc.Logout("local", true)

	assert.NoError(t, err)
}

func TestProfileName_FromURL(t *testing.T) {
	assert.Equal(t, "local", domain.ProfileName("http://localhost:7470"))
	assert.Equal(t, "games.example.org", domain.ProfileName("https://games.example.org"))
	assert.Equal(t, "games.example.org:8443", domain.ProfileName("https://games.example.org:8443/"))
}
