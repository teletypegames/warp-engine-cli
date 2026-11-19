package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareCategoryService_Create_NeedsKeyAndLabel(t *testing.T) {
	repo := &mockrepo.MockSoftwareCategoryRepository{CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareCategoryModel, error) {
		t.Fatal("must not reach the repository")
		return nil, nil
	}}
	svc := domain.SoftwareCategoryService{SoftwareCategoryRepository: repo}

	_, err := svc.Create(domain.Attributes{"label": "Tool", "key": ""})

	assert.ErrorContains(t, err, "key")
}

func TestSoftwareCategoryService_Update_DelegatesIDAndAttributes(t *testing.T) {
	var got_id int
	var got domain.Attributes
	repo := &mockrepo.MockSoftwareCategoryRepository{UpdateFunc: func(category_id int, attributes domain.Attributes) (*domain.SoftwareCategoryModel, error) {
		got_id, got = category_id, attributes
		return &domain.SoftwareCategoryModel{ID: category_id}, nil
	}}
	svc := domain.SoftwareCategoryService{SoftwareCategoryRepository: repo}

	_, err := svc.Update(5, domain.Attributes{"needsReview": false})

	assert.NoError(t, err)
	assert.Equal(t, 5, got_id)
	assert.Equal(t, domain.Attributes{"needsReview": false}, got)
}

func TestSoftwareCategoryService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.SoftwareCategoryService{SoftwareCategoryRepository: &mockrepo.MockSoftwareCategoryRepository{}}

	_, err := svc.Update(5, nil)

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestSoftwareCategoryService_Show_DelegatesTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareCategoryRepository{FindByIDFunc: func(category_id int) (*domain.SoftwareCategoryModel, error) {
		got_id = category_id
		return &domain.SoftwareCategoryModel{ID: category_id}, nil
	}}
	svc := domain.SoftwareCategoryService{SoftwareCategoryRepository: repo}

	_, err := svc.Show(2)

	assert.NoError(t, err)
	assert.Equal(t, 2, got_id)
}

func TestSoftwareCategoryService_Destroy_DelegatesTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareCategoryRepository{DestroyFunc: func(category_id int) error { got_id = category_id; return nil }}
	svc := domain.SoftwareCategoryService{SoftwareCategoryRepository: repo}

	assert.NoError(t, svc.Destroy(3))
	assert.Equal(t, 3, got_id)
}
