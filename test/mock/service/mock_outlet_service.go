package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockOutletService mocks [domain.OutletServiceInterface].
type MockOutletService struct {
	// AllFunc replaces All.
	AllFunc func() ([]domain.OutletModel, error)
	// CheckFunc replaces Check.
	CheckFunc func(key string) (*domain.OutletCheckModel, error)
}

// All calls AllFunc, or returns no outlets.
func (m *MockOutletService) All() ([]domain.OutletModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return nil, nil
}

// Check calls CheckFunc, or returns an empty check.
func (m *MockOutletService) Check(key string) (*domain.OutletCheckModel, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(key)
	}
	return &domain.OutletCheckModel{Key: key}, nil
}
