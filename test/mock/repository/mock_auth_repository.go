package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockAuthRepository mocks [domain.AuthRepositoryInterface]. Use records
// every call, so a test can check where the client was pointed.
type MockAuthRepository struct {
	// Used lists the url and token of every Use call, in order.
	Used [][2]string
	// DiscoverFunc replaces Discover.
	DiscoverFunc func() (*domain.ServerModel, error)
	// StartDeviceGrantFunc replaces StartDeviceGrant.
	StartDeviceGrantFunc func(client_name string, scopes []string) (*domain.DeviceGrantModel, error)
	// PollDeviceGrantFunc replaces PollDeviceGrant.
	PollDeviceGrantFunc func(device_code string) (*domain.TokenPollModel, error)
	// WhoamiFunc replaces Whoami.
	WhoamiFunc func() (*domain.PrincipalModel, error)
	// RevokeFunc replaces Revoke.
	RevokeFunc func() error
}

// Use records url and token in Used.
func (m *MockAuthRepository) Use(url string, token string) {
	m.Used = append(m.Used, [2]string{url, token})
}

// Discover calls DiscoverFunc, or returns an empty server.
func (m *MockAuthRepository) Discover() (*domain.ServerModel, error) {
	if m.DiscoverFunc != nil {
		return m.DiscoverFunc()
	}
	return &domain.ServerModel{}, nil
}

// StartDeviceGrant calls StartDeviceGrantFunc, or returns an empty grant.
func (m *MockAuthRepository) StartDeviceGrant(client_name string, scopes []string) (*domain.DeviceGrantModel, error) {
	if m.StartDeviceGrantFunc != nil {
		return m.StartDeviceGrantFunc(client_name, scopes)
	}
	return &domain.DeviceGrantModel{}, nil
}

// PollDeviceGrant calls PollDeviceGrantFunc, or answers pending.
func (m *MockAuthRepository) PollDeviceGrant(device_code string) (*domain.TokenPollModel, error) {
	if m.PollDeviceGrantFunc != nil {
		return m.PollDeviceGrantFunc(device_code)
	}
	return &domain.TokenPollModel{State: domain.TokenPending}, nil
}

// Whoami calls WhoamiFunc, or returns an empty principal.
func (m *MockAuthRepository) Whoami() (*domain.PrincipalModel, error) {
	if m.WhoamiFunc != nil {
		return m.WhoamiFunc()
	}
	return &domain.PrincipalModel{}, nil
}

// Revoke calls RevokeFunc, or returns nil.
func (m *MockAuthRepository) Revoke() error {
	if m.RevokeFunc != nil {
		return m.RevokeFunc()
	}
	return nil
}
