package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestReleaseAssetService_All_ParsesTheReleaseID(t *testing.T) {
	var got_id int
	svc := domain.ReleaseAssetService{ReleaseAssetRepository: &mockrepo.MockReleaseAssetRepository{
		AllFunc: func(release_id int) ([]domain.ReleaseAssetModel, error) {
			got_id = release_id
			return []domain.ReleaseAssetModel{{ID: 1}}, nil
		},
	}}

	assets, err := svc.All("42")

	assert.NoError(t, err)
	assert.Equal(t, 42, got_id)
	assert.Len(t, assets, 1)
}

func TestReleaseAssetService_All_RejectsANonNumericRelease(t *testing.T) {
	svc := domain.ReleaseAssetService{ReleaseAssetRepository: &mockrepo.MockReleaseAssetRepository{}}

	_, err := svc.All("latest")

	assert.ErrorContains(t, err, "release ID")
}

func TestReleaseAssetService_Create_NeedsKindAndPath(t *testing.T) {
	svc := domain.ReleaseAssetService{ReleaseAssetRepository: &mockrepo.MockReleaseAssetRepository{
		CreateFunc: func(release_id int, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
			t.Fatal("must not reach the repository")
			return nil, nil
		},
	}}

	_, err := svc.Create("42", domain.Attributes{"kind": "html"})

	assert.ErrorContains(t, err, "path")
}

func TestReleaseAssetService_Create_DelegatesToTheRelease(t *testing.T) {
	var got_id int
	svc := domain.ReleaseAssetService{ReleaseAssetRepository: &mockrepo.MockReleaseAssetRepository{
		CreateFunc: func(release_id int, attributes domain.Attributes) (*domain.ReleaseAssetModel, error) {
			got_id = release_id
			return &domain.ReleaseAssetModel{ID: 5}, nil
		},
	}}

	asset, err := svc.Create("42", domain.Attributes{"kind": "html", "path": "/a/index.html"})

	assert.NoError(t, err)
	assert.Equal(t, 42, got_id)
	assert.Equal(t, 5, asset.ID)
}

func TestReleaseAssetService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.ReleaseAssetService{ReleaseAssetRepository: &mockrepo.MockReleaseAssetRepository{}}

	_, err := svc.Update("5", domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestReleaseAssetService_Destroy_DelegatesTheParsedID(t *testing.T) {
	var got_id int
	svc := domain.ReleaseAssetService{ReleaseAssetRepository: &mockrepo.MockReleaseAssetRepository{
		DestroyFunc: func(asset_id int) error { got_id = asset_id; return nil },
	}}

	err := svc.Destroy("5")

	assert.NoError(t, err)
	assert.Equal(t, 5, got_id)
}
