package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockLocalComponentRepository mocks [domain.LocalComponentRepositoryInterface].
type MockLocalComponentRepository struct {
	// AllFunc replaces All.
	AllFunc func() (*domain.PageModel[domain.LocalComponentModel], error)
	// EnableFunc replaces Enable.
	EnableFunc func(key string) (*domain.LocalOperationModel, error)
	// DisableFunc replaces Disable.
	DisableFunc func(key string) (*domain.LocalOperationModel, error)
	// DestroyFunc replaces Destroy.
	DestroyFunc func(key string, purge bool) (*domain.LocalOperationModel, error)
	// CredentialsFunc replaces Credentials.
	CredentialsFunc func(key string) (*domain.LocalCredentialsModel, error)
	// SetupCIFunc replaces SetupCI.
	SetupCIFunc func(software_id int) (*domain.LocalOperationModel, error)
	// TeardownCIFunc replaces TeardownCI.
	TeardownCIFunc func(software_id int) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockLocalComponentRepository) All() (*domain.PageModel[domain.LocalComponentModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return &domain.PageModel[domain.LocalComponentModel]{}, nil
}

// Enable calls EnableFunc, or returns an empty operation.
func (m *MockLocalComponentRepository) Enable(key string) (*domain.LocalOperationModel, error) {
	if m.EnableFunc != nil {
		return m.EnableFunc(key)
	}
	return &domain.LocalOperationModel{}, nil
}

// Disable calls DisableFunc, or returns an empty operation.
func (m *MockLocalComponentRepository) Disable(key string) (*domain.LocalOperationModel, error) {
	if m.DisableFunc != nil {
		return m.DisableFunc(key)
	}
	return &domain.LocalOperationModel{}, nil
}

// Destroy calls DestroyFunc, or returns an empty operation.
func (m *MockLocalComponentRepository) Destroy(key string, purge bool) (*domain.LocalOperationModel, error) {
	if m.DestroyFunc != nil {
		return m.DestroyFunc(key, purge)
	}
	return &domain.LocalOperationModel{}, nil
}

// Credentials calls CredentialsFunc, or returns empty credentials.
func (m *MockLocalComponentRepository) Credentials(key string) (*domain.LocalCredentialsModel, error) {
	if m.CredentialsFunc != nil {
		return m.CredentialsFunc(key)
	}
	return &domain.LocalCredentialsModel{}, nil
}

// SetupCI calls SetupCIFunc, or returns an empty operation.
func (m *MockLocalComponentRepository) SetupCI(software_id int) (*domain.LocalOperationModel, error) {
	if m.SetupCIFunc != nil {
		return m.SetupCIFunc(software_id)
	}
	return &domain.LocalOperationModel{}, nil
}

// TeardownCI calls TeardownCIFunc, or returns nil.
func (m *MockLocalComponentRepository) TeardownCI(software_id int) error {
	if m.TeardownCIFunc != nil {
		return m.TeardownCIFunc(software_id)
	}
	return nil
}
