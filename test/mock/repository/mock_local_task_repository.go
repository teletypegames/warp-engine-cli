package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockLocalTaskRepository mocks [domain.LocalTaskRepositoryInterface].
type MockLocalTaskRepository struct {
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(task_id int) (*domain.LocalTaskModel, error)
}

// FindByID calls FindByIDFunc, or returns an empty task.
func (m *MockLocalTaskRepository) FindByID(task_id int) (*domain.LocalTaskModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(task_id)
	}
	return &domain.LocalTaskModel{}, nil
}
