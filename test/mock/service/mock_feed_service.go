package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockFeedService mocks [domain.FeedServiceInterface].
type MockFeedService struct {
	// AllFunc replaces All.
	AllFunc func() ([]domain.FeedModel, error)
	// CheckFunc replaces Check.
	CheckFunc func(key string) (*domain.FeedCheckModel, error)
}

// All calls AllFunc, or returns no feeds.
func (m *MockFeedService) All() ([]domain.FeedModel, error) {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return nil, nil
}

// Check calls CheckFunc, or returns an empty check.
func (m *MockFeedService) Check(key string) (*domain.FeedCheckModel, error) {
	if m.CheckFunc != nil {
		return m.CheckFunc(key)
	}
	return &domain.FeedCheckModel{Key: key}, nil
}
