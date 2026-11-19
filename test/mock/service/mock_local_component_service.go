package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockLocalComponentService mocks [domain.LocalComponentServiceInterface].
type MockLocalComponentService struct {
	// AllFunc replaces All.
	AllFunc func() (*domain.PageModel[domain.LocalComponentModel], error)
	// EnableFunc replaces Enable.
	EnableFunc func(key string) (*domain.LocalOperationModel, error)
	// DisableFunc replaces Disable.
	DisableFunc func(key string) (*domain.LocalOperationModel, error)
	// RemoveFunc replaces Remove.
	RemoveFunc func(key string, purge bool) (*domain.LocalOperationModel, error)
	// CredentialsFunc replaces Credentials.
	CredentialsFunc func(key string) (*domain.LocalCredentialsModel, error)
	// SetupCIFunc replaces SetupCI.
	SetupCIFunc func(reference string) (*domain.LocalOperationModel, error)
	// TeardownCIFunc replaces TeardownCI.
	TeardownCIFunc func(reference string) error
}

// All calls AllFunc, or returns an empty page.
func (m *MockLocalComponentService) All() (*domain.PageModel[domain.LocalComponentModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return &domain.PageModel[domain.LocalComponentModel]{}, nil
}

// Enable calls EnableFunc, or returns an empty operation.
func (m *MockLocalComponentService) Enable(key string) (*domain.LocalOperationModel, error) {
	if m.EnableFunc != nil {
		return m.EnableFunc(key)
	}
	return &domain.LocalOperationModel{}, nil
}

// Disable calls DisableFunc, or returns an empty operation.
func (m *MockLocalComponentService) Disable(key string) (*domain.LocalOperationModel, error) {
	if m.DisableFunc != nil {
		return m.DisableFunc(key)
	}
	return &domain.LocalOperationModel{}, nil
}

// Remove calls RemoveFunc, or returns an empty operation.
func (m *MockLocalComponentService) Remove(key string, purge bool) (*domain.LocalOperationModel, error) {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(key, purge)
	}
	return &domain.LocalOperationModel{}, nil
}

// Credentials calls CredentialsFunc, or returns empty credentials.
func (m *MockLocalComponentService) Credentials(key string) (*domain.LocalCredentialsModel, error) {
	if m.CredentialsFunc != nil {
		return m.CredentialsFunc(key)
	}
	return &domain.LocalCredentialsModel{}, nil
}

// SetupCI calls SetupCIFunc, or returns an empty operation.
func (m *MockLocalComponentService) SetupCI(reference string) (*domain.LocalOperationModel, error) {
	if m.SetupCIFunc != nil {
		return m.SetupCIFunc(reference)
	}
	return &domain.LocalOperationModel{}, nil
}

// TeardownCI calls TeardownCIFunc, or returns nil.
func (m *MockLocalComponentService) TeardownCI(reference string) error {
	if m.TeardownCIFunc != nil {
		return m.TeardownCIFunc(reference)
	}
	return nil
}
