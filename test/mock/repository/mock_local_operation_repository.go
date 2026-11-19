package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockLocalOperationRepository mocks [domain.LocalOperationRepositoryInterface].
type MockLocalOperationRepository struct {
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(operation_id int) (*domain.LocalOperationModel, error)
}

// FindByID calls FindByIDFunc, or returns an empty operation.
func (m *MockLocalOperationRepository) FindByID(operation_id int) (*domain.LocalOperationModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(operation_id)
	}
	return &domain.LocalOperationModel{}, nil
}
