package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockAuthService mocks [domain.AuthServiceInterface]. With ConnectFunc nil,
// Connect succeeds with a profile for https://games.example.org, so the
// controller tests of other commands need not set it.
type MockAuthService struct {
	// ConnectFunc replaces Connect.
	ConnectFunc func(selection domain.Selection) (*domain.ProfileModel, error)
	// DiscoverFunc replaces Discover.
	DiscoverFunc func(server_url string) (*domain.ServerModel, error)
	// StartDeviceLoginFunc replaces StartDeviceLogin.
	StartDeviceLoginFunc func(server_url string, client_name string, scopes []string) (*domain.DeviceGrantModel, error)
	// WaitForApprovalFunc replaces WaitForApproval.
	WaitForApprovalFunc func(grant domain.DeviceGrantModel) (string, error)
	// FinishLoginFunc replaces FinishLogin.
	FinishLoginFunc func(server_url string, token string, name string) (*domain.ProfileModel, error)
	// LoginLocalFunc replaces LoginLocal.
	LoginLocalFunc func(server_url string, name string) (*domain.ProfileModel, error)
	// LogoutFunc replaces Logout.
	LogoutFunc func(name string, revoke bool) (*domain.ProfileModel, error)
	// WhoamiFunc replaces Whoami.
	WhoamiFunc func() (*domain.PrincipalModel, error)
}

// Connect calls ConnectFunc, or returns a test profile.
func (m *MockAuthService) Connect(selection domain.Selection) (*domain.ProfileModel, error) {
	if m.ConnectFunc != nil {
		return m.ConnectFunc(selection)
	}
	return &domain.ProfileModel{Name: "test", URL: "https://games.example.org", Token: "test-token"}, nil
}

// Discover calls DiscoverFunc, or returns an empty server.
func (m *MockAuthService) Discover(server_url string) (*domain.ServerModel, error) {
	if m.DiscoverFunc != nil {
		return m.DiscoverFunc(server_url)
	}
	return &domain.ServerModel{}, nil
}

// StartDeviceLogin calls StartDeviceLoginFunc, or returns an empty grant.
func (m *MockAuthService) StartDeviceLogin(server_url string, client_name string, scopes []string) (*domain.DeviceGrantModel, error) {
	if m.StartDeviceLoginFunc != nil {
		return m.StartDeviceLoginFunc(server_url, client_name, scopes)
	}
	return &domain.DeviceGrantModel{}, nil
}

// WaitForApproval calls WaitForApprovalFunc, or returns "".
func (m *MockAuthService) WaitForApproval(grant domain.DeviceGrantModel) (string, error) {
	if m.WaitForApprovalFunc != nil {
		return m.WaitForApprovalFunc(grant)
	}
	return "", nil
}

// FinishLogin calls FinishLoginFunc, or returns a profile for server_url.
func (m *MockAuthService) FinishLogin(server_url string, token string, name string) (*domain.ProfileModel, error) {
	if m.FinishLoginFunc != nil {
		return m.FinishLoginFunc(server_url, token, name)
	}
	return &domain.ProfileModel{Name: name, URL: server_url, Token: token}, nil
}

// LoginLocal calls LoginLocalFunc, or returns a Local profile for server_url.
func (m *MockAuthService) LoginLocal(server_url string, name string) (*domain.ProfileModel, error) {
	if m.LoginLocalFunc != nil {
		return m.LoginLocalFunc(server_url, name)
	}
	return &domain.ProfileModel{Name: "local", URL: server_url, TokenFile: true}, nil
}

// Logout calls LogoutFunc, or returns a profile named name.
func (m *MockAuthService) Logout(name string, revoke bool) (*domain.ProfileModel, error) {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(name, revoke)
	}
	return &domain.ProfileModel{Name: name}, nil
}

// Whoami calls WhoamiFunc, or returns an empty principal.
func (m *MockAuthService) Whoami() (*domain.PrincipalModel, error) {
	if m.WhoamiFunc != nil {
		return m.WhoamiFunc()
	}
	return &domain.PrincipalModel{}, nil
}
