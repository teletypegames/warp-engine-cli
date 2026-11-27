package mocksvc

import (
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// MockLocalTaskService mocks [domain.LocalTaskServiceInterface].
type MockLocalTaskService struct {
	// ShowFunc replaces Show.
	ShowFunc func(task_id int) (*domain.LocalTaskModel, error)
	// WaitFunc replaces Wait.
	WaitFunc func(task domain.LocalTaskModel, timeout time.Duration, progress func(domain.LocalTaskModel)) (*domain.LocalTaskModel, error)
}

// Show calls ShowFunc, or returns an empty task.
func (m *MockLocalTaskService) Show(task_id int) (*domain.LocalTaskModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(task_id)
	}
	return &domain.LocalTaskModel{}, nil
}

// Wait calls WaitFunc, or returns task as it is.
func (m *MockLocalTaskService) Wait(task domain.LocalTaskModel, timeout time.Duration, progress func(domain.LocalTaskModel)) (*domain.LocalTaskModel, error) {
	if m.WaitFunc != nil {
		return m.WaitFunc(task, timeout, progress)
	}
	return &task, nil
}
