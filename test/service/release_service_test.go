package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestReleaseService_All_ResolvesTheSoftwareName(t *testing.T) {
	var got_id int
	var got_filter domain.Filter
	releases := &mockrepo.MockReleaseRepository{
		AllFunc: func(software_id int, filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ReleaseModel], error) {
			got_id, got_filter = software_id, filter
			return &domain.PageModel[domain.ReleaseModel]{Items: []domain.ReleaseModel{{ID: 1}}}, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
			assert.Equal(t, "pong", name)
			return &domain.SoftwareModel{ID: 3}, nil
		},
	}
	svc := domain.ReleaseService{ReleaseRepository: releases, SoftwareRepository: softwares}

	page, err := svc.All("pong", domain.Filter{"channel": "beta"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
	assert.Equal(t, domain.Filter{"channel": "beta"}, got_filter)
	assert.Len(t, page.Items, 1)
}

func TestReleaseService_Show_RejectsANonNumericID(t *testing.T) {
	svc := domain.ReleaseService{ReleaseRepository: &mockrepo.MockReleaseRepository{
		FindByIDFunc: func(release_id int) (*domain.ReleaseModel, error) {
			t.Fatal("must not reach the repository")
			return nil, nil
		},
	}}

	_, err := svc.Show("pong")

	assert.ErrorContains(t, err, "release ID")
}

func TestReleaseService_Create_NeedsAVersion(t *testing.T) {
	svc := domain.ReleaseService{ReleaseRepository: &mockrepo.MockReleaseRepository{}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Create("3", domain.Attributes{})

	assert.ErrorContains(t, err, "version")
}

func TestReleaseService_Create_DelegatesToTheSoftware(t *testing.T) {
	var got_id int
	var got domain.Attributes
	releases := &mockrepo.MockReleaseRepository{
		CreateFunc: func(software_id int, attributes domain.Attributes) (*domain.ReleaseModel, error) {
			got_id, got = software_id, attributes
			return &domain.ReleaseModel{ID: 9, Version: "1.0.0"}, nil
		},
	}
	svc := domain.ReleaseService{ReleaseRepository: releases, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	release, err := svc.Create("3", domain.Attributes{"version": "1.0.0"})

	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
	assert.Equal(t, domain.Attributes{"version": "1.0.0"}, got)
	assert.Equal(t, 9, release.ID)
}

func TestReleaseService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.ReleaseService{ReleaseRepository: &mockrepo.MockReleaseRepository{}}

	_, err := svc.Update("4", domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestReleaseService_Destroy_DelegatesTheParsedID(t *testing.T) {
	var got_id int
	svc := domain.ReleaseService{ReleaseRepository: &mockrepo.MockReleaseRepository{
		DestroyFunc: func(release_id int) error { got_id = release_id; return nil },
	}}

	err := svc.Destroy("42")

	assert.NoError(t, err)
	assert.Equal(t, 42, got_id)
}
