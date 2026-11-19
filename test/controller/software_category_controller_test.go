package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareCategoryController_Index_PassesNeedsReview(t *testing.T) {
	var got_filter domain.Filter
	svc := &mocksvc.MockSoftwareCategoryService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareCategoryModel], error) {
			got_filter = filter
			return &domain.PageModel[domain.SoftwareCategoryModel]{
				Items: []domain.SoftwareCategoryModel{{ID: 5, Key: "puzzle", Label: "Puzzle", NeedsReview: true}},
			}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareCategoryController: cli.SoftwareCategoryController{SoftwareCategoryService: svc}},
		"categories", "list", "--needs-review")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"needsReview": "true"}, got_filter)
	assert.Contains(t, stdout, "puzzle")
}

func TestSoftwareCategoryController_Update_ApprovesWithNeedsReviewFalse(t *testing.T) {
	var got_id int
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareCategoryService{UpdateFunc: func(category_id int, attributes domain.Attributes) (*domain.SoftwareCategoryModel, error) {
		got_id, got = category_id, attributes
		return &domain.SoftwareCategoryModel{ID: category_id, Key: "puzzle"}, nil
	}}

	stdout, _, err := execute(cli.Router{SoftwareCategoryController: cli.SoftwareCategoryController{SoftwareCategoryService: svc}},
		"category", "update", "5", "--needs-review=false", "--label", "Puzzles")

	assert.NoError(t, err)
	assert.Equal(t, 5, got_id)
	assert.Equal(t, domain.Attributes{"needsReview": false, "label": "Puzzles"}, got)
	assert.Contains(t, stdout, "puzzle")
}

func TestSoftwareCategoryController_Create_BadPositionIsAUsageError(t *testing.T) {
	svc := &mocksvc.MockSoftwareCategoryService{CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareCategoryModel, error) {
		t.Fatal("must not reach the service")
		return nil, nil
	}}

	_, _, err := execute(cli.Router{SoftwareCategoryController: cli.SoftwareCategoryController{SoftwareCategoryService: svc}},
		"category", "create", "--key", "tool", "--label", "Tool", "--position", "first")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSoftwareCategoryController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	called := false
	svc := &mocksvc.MockSoftwareCategoryService{DestroyFunc: func(category_id int) error { called = true; return nil }}
	router := cli.Router{SoftwareCategoryController: cli.SoftwareCategoryController{SoftwareCategoryService: svc}}

	_, _, err := execute(router, "category", "delete", "5")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(router, "category", "delete", "5", "-y")
	assert.NoError(t, err)
	assert.True(t, called)
}
