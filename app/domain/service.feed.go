package domain

import "errors"

// FeedServiceInterface reads the engine's feeds.
type FeedServiceInterface interface {
	// All lists the configured feeds.
	All() ([]FeedModel, error)
	// Check returns what the feed key offers and leaves out.
	Check(key string) (*FeedCheckModel, error)
}

// FeedService is the [FeedServiceInterface] over an
// [FeedRepositoryInterface].
type FeedService struct {
	// FeedRepository talks to the feeds resource.
	FeedRepository FeedRepositoryInterface
}

// All lists the feeds.
func (s FeedService) All() ([]FeedModel, error) {
	return s.FeedRepository.All()
}

// Check refuses an empty key and gets the feed's check.
func (s FeedService) Check(key string) (*FeedCheckModel, error) {
	if key == "" {
		return nil, errors.New("the feed key is empty")
	}
	return s.FeedRepository.Check(key)
}
