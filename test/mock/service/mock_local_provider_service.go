package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockLocalProviderService mocks [domain.LocalProviderServiceInterface].
type MockLocalProviderService struct {
	// AllFunc replaces All.
	AllFunc func() (*domain.PageModel[domain.LocalProviderModel], error)
	// EnableFunc replaces Enable.
	EnableFunc func(key string) (*domain.LocalTaskModel, error)
	// DisableFunc replaces Disable.
	DisableFunc func(key string) (*domain.LocalTaskModel, error)
	// RemoveFunc replaces Remove.
	RemoveFunc func(key string, purge bool) (*domain.LocalTaskModel, error)
	// CredentialsFunc replaces Credentials.
	CredentialsFunc func(key string) (*domain.LocalCredentialsModel, error)
	// SetupCIFunc replaces SetupCI.
	SetupCIFunc func(reference string) (*domain.LocalTaskModel, error)
	// TeardownCIFunc replaces TeardownCI.
	TeardownCIFunc func(reference string) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockLocalProviderService) All() (*domain.PageModel[domain.LocalProviderModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return &domain.PageModel[domain.LocalProviderModel]{}, nil
}

// Enable calls EnableFunc, or returns an empty task.
func (m *MockLocalProviderService) Enable(key string) (*domain.LocalTaskModel, error) {
	if m.EnableFunc != nil {
		return m.EnableFunc(key)
	}
	return &domain.LocalTaskModel{}, nil
}

// Disable calls DisableFunc, or returns an empty task.
func (m *MockLocalProviderService) Disable(key string) (*domain.LocalTaskModel, error) {
	if m.DisableFunc != nil {
		return m.DisableFunc(key)
	}
	return &domain.LocalTaskModel{}, nil
}

// Remove calls RemoveFunc, or returns an empty task.
func (m *MockLocalProviderService) Remove(key string, purge bool) (*domain.LocalTaskModel, error) {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(key, purge)
	}
	return &domain.LocalTaskModel{}, nil
}

// Credentials calls CredentialsFunc, or returns empty credentials.
func (m *MockLocalProviderService) Credentials(key string) (*domain.LocalCredentialsModel, error) {
	if m.CredentialsFunc != nil {
		return m.CredentialsFunc(key)
	}
	return &domain.LocalCredentialsModel{}, nil
}

// SetupCI calls SetupCIFunc, or returns an empty task.
func (m *MockLocalProviderService) SetupCI(reference string) (*domain.LocalTaskModel, error) {
	if m.SetupCIFunc != nil {
		return m.SetupCIFunc(reference)
	}
	return &domain.LocalTaskModel{}, nil
}

// TeardownCI calls TeardownCIFunc, or returns nil.
func (m *MockLocalProviderService) TeardownCI(reference string) error {
	if m.TeardownCIFunc != nil {
		return m.TeardownCIFunc(reference)
	}
	return nil
}
