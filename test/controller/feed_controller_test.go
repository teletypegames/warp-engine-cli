package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestFeedController_Index_ListsTheFeeds(t *testing.T) {
	svc := &mocksvc.MockFeedService{AllFunc: func() ([]domain.FeedModel, error) {
		return []domain.FeedModel{
			{Key: "wecp", Type: "wecp", Dynamic: true, Description: "wecp (wecp)"},
			{Key: "assembly64", Type: "assembly64", Platforms: []string{"c64"}, Dynamic: true, Description: "assembly64 (assembly64: c64)"},
		}, nil
	}}

	stdout, _, err := execute(cli.Router{FeedController: cli.FeedController{FeedService: svc}}, "feed", "list")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "assembly64")
	assert.Contains(t, stdout, "c64")
}

func TestFeedController_Check_PrintsTheSelectionAndFailsStrictOnASkip(t *testing.T) {
	var got_key string
	svc := &mocksvc.MockFeedService{CheckFunc: func(key string) (*domain.FeedCheckModel, error) {
		got_key = key
		return &domain.FeedCheckModel{Key: key, Type: "assembly64", Selection: true,
			Items:   []domain.FeedItemModel{{Name: "ariadne", Version: "1.0", Kind: "cartridge"}},
			Skipped: []domain.FeedSkipModel{{Name: "empty", Reason: "no cartridge"}}}, nil
	}}
	router := cli.Router{FeedController: cli.FeedController{FeedService: svc}}

	stdout, _, err := execute(router, "feeds", "check", "assembly64")
	assert.NoError(t, err)
	assert.Equal(t, "assembly64", got_key)
	assert.Contains(t, stdout, "ariadne")
	assert.Contains(t, stdout, "no cartridge")

	_, _, err = execute(router, "feed", "check", "assembly64", "--strict")
	assert.Error(t, err)
}
