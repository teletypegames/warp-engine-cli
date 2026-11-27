package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockLocalProviderRepository mocks [domain.LocalProviderRepositoryInterface].
type MockLocalProviderRepository struct {
	// AllFunc replaces All.
	AllFunc func() (*domain.PageModel[domain.LocalProviderModel], error)
	// EnableFunc replaces Enable.
	EnableFunc func(key string) (*domain.LocalTaskModel, error)
	// DisableFunc replaces Disable.
	DisableFunc func(key string) (*domain.LocalTaskModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(key string, purge bool) (*domain.LocalTaskModel, error)
	// CredentialsFunc replaces Credentials.
	CredentialsFunc func(key string) (*domain.LocalCredentialsModel, error)
	// SetupCIFunc replaces SetupCI.
	SetupCIFunc func(software_id int) (*domain.LocalTaskModel, error)
	// TeardownCIFunc replaces TeardownCI.
	TeardownCIFunc func(software_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockLocalProviderRepository) All() (*domain.PageModel[domain.LocalProviderModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return &domain.PageModel[domain.LocalProviderModel]{}, nil
}

// Enable calls EnableFunc, or returns an empty task.
func (m *MockLocalProviderRepository) Enable(key string) (*domain.LocalTaskModel, error) {
	if m.EnableFunc != nil {
		return m.EnableFunc(key)
	}
	return &domain.LocalTaskModel{}, nil
}

// Disable calls DisableFunc, or returns an empty task.
func (m *MockLocalProviderRepository) Disable(key string) (*domain.LocalTaskModel, error) {
	if m.DisableFunc != nil {
		return m.DisableFunc(key)
	}
	return &domain.LocalTaskModel{}, nil
}

// Destroy calls DestroyFunc, or returns an empty task.
func (m *MockLocalProviderRepository) Destroy(key string, purge bool) (*domain.LocalTaskModel, error) {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(key, purge)
	}
	return &domain.LocalTaskModel{}, nil
}

// Credentials calls CredentialsFunc, or returns empty credentials.
func (m *MockLocalProviderRepository) Credentials(key string) (*domain.LocalCredentialsModel, error) {
	if m.CredentialsFunc != nil {
		return m.CredentialsFunc(key)
	}
	return &domain.LocalCredentialsModel{}, nil
}

// SetupCI calls SetupCIFunc, or returns an empty task.
func (m *MockLocalProviderRepository) SetupCI(software_id int) (*domain.LocalTaskModel, error) {
	if m.SetupCIFunc != nil {
		return m.SetupCIFunc(software_id)
	}
	return &domain.LocalTaskModel{}, nil
}

// TeardownCI calls TeardownCIFunc, or returns nil.
func (m *MockLocalProviderRepository) TeardownCI(software_id int) error {
	if m.TeardownCIFunc != nil {
		return m.TeardownCIFunc(software_id)
	}
	return nil
}
