package mocksvc

import (
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// MockLocalOperationService mocks [domain.LocalOperationServiceInterface].
type MockLocalOperationService struct {
	// ShowFunc replaces Show.
	ShowFunc func(operation_id int) (*domain.LocalOperationModel, error)
	// WaitFunc replaces Wait.
	WaitFunc func(operation domain.LocalOperationModel, timeout time.Duration, progress func(domain.LocalOperationModel)) (*domain.LocalOperationModel, error)
}

// Show calls ShowFunc, or returns an empty operation.
func (m *MockLocalOperationService) Show(operation_id int) (*domain.LocalOperationModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(operation_id)
	}
	return &domain.LocalOperationModel{}, nil
}

// Wait calls WaitFunc, or returns operation as it is.
func (m *MockLocalOperationService) Wait(operation domain.LocalOperationModel, timeout time.Duration, progress func(domain.LocalOperationModel)) (*domain.LocalOperationModel, error) {
	if m.WaitFunc != nil {
		return m.WaitFunc(operation, timeout, progress)
	}
	return &operation, nil
}
