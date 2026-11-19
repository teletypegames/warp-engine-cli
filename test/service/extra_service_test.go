package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestExtraService_Show_ResolvesTheSoftwareName(t *testing.T) {
	var got_id int
	var got_association string
	extras := &mockrepo.MockExtraRepository{
		FindFunc: func(software_id int, association string) (*domain.ExtraModel, error) {
			got_id, got_association = software_id, association
			return &domain.ExtraModel{Association: association}, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) { return &domain.SoftwareModel{ID: 4}, nil },
	}
	svc := domain.ExtraService{ExtraRepository: extras, SoftwareRepository: softwares}

	extra, err := svc.Show("pong", "seo_metadata")

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Equal(t, "seo_metadata", got_association)
	assert.Equal(t, "seo_metadata", extra.Association)
}

func TestExtraService_Update_NeedsAValue(t *testing.T) {
	svc := domain.ExtraService{ExtraRepository: &mockrepo.MockExtraRepository{
		UpdateFunc: func(software_id int, association string, attributes domain.Attributes) (*domain.ExtraModel, error) {
			t.Fatal("must not reach the repository")
			return nil, nil
		},
	}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Update("3", "seo_metadata", domain.Attributes{})

	assert.ErrorContains(t, err, "value")
}

func TestExtraService_Update_DelegatesTheBody(t *testing.T) {
	var got domain.Attributes
	svc := domain.ExtraService{ExtraRepository: &mockrepo.MockExtraRepository{
		UpdateFunc: func(software_id int, association string, attributes domain.Attributes) (*domain.ExtraModel, error) {
			got = attributes
			return &domain.ExtraModel{}, nil
		},
	}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Update("3", "seo_metadata", domain.Attributes{"value": map[string]any{"title": "Pong"}})

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"value": map[string]any{"title": "Pong"}}, got)
}

func TestExtraService_Show_NeedsAnAssociation(t *testing.T) {
	svc := domain.ExtraService{ExtraRepository: &mockrepo.MockExtraRepository{}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Show("3", "")

	assert.ErrorContains(t, err, "association")
}
