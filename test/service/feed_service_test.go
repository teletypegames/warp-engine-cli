package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestFeedService_Check_RefusesAnEmptyKey(t *testing.T) {
	called := false
	svc := domain.FeedService{FeedRepository: &mockrepo.MockFeedRepository{CheckFunc: func(key string) (*domain.FeedCheckModel, error) {
		called = true
		return nil, nil
	}}}

	_, err := svc.Check("")

	assert.Error(t, err)
	assert.False(t, called)
}

func TestFeedService_Check_AsksTheRepository(t *testing.T) {
	svc := domain.FeedService{FeedRepository: &mockrepo.MockFeedRepository{}}

	check, err := svc.Check("assembly64")

	assert.NoError(t, err)
	assert.Equal(t, "assembly64", check.Key)
}
