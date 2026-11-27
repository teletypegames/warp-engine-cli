package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockFeedRepository mocks [domain.FeedRepositoryInterface].
type MockFeedRepository struct {
	// AllFunc replaces All.
	AllFunc func() ([]domain.FeedModel, error)
	// CheckFunc replaces Check.
	CheckFunc func(key string) (*domain.FeedCheckModel, error)
}

// All calls AllFunc, or returns no feeds.
func (m *MockFeedRepository) All() ([]domain.FeedModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return nil, nil
}

// Check calls CheckFunc, or returns an empty check.
func (m *MockFeedRepository) Check(key string) (*domain.FeedCheckModel, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(key)
	}
	return &domain.FeedCheckModel{Key: key}, nil
}
