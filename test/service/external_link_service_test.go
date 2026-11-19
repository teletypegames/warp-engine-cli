package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestExternalLinkService_All_ResolvesTheSoftwareName(t *testing.T) {
	var got_id int
	links := &mockrepo.MockExternalLinkRepository{
		AllFunc: func(software_id int) ([]domain.ExternalLinkModel, error) {
			got_id = software_id
			return []domain.ExternalLinkModel{{ID: 1}}, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) { return &domain.SoftwareModel{ID: 8}, nil },
	}
	svc := domain.ExternalLinkService{ExternalLinkRepository: links, SoftwareRepository: softwares}

	result, err := svc.All("pong")

	assert.NoError(t, err)
	assert.Equal(t, 8, got_id)
	assert.Len(t, result, 1)
}

func TestExternalLinkService_Create_NeedsLabelAndURL(t *testing.T) {
	svc := domain.ExternalLinkService{ExternalLinkRepository: &mockrepo.MockExternalLinkRepository{}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Create("3", domain.Attributes{"label": "Source"})

	assert.ErrorContains(t, err, "url")
}

func TestExternalLinkService_Create_DelegatesToTheSoftware(t *testing.T) {
	var got_id int
	links := &mockrepo.MockExternalLinkRepository{
		CreateFunc: func(software_id int, attributes domain.Attributes) (*domain.ExternalLinkModel, error) {
			got_id = software_id
			return &domain.ExternalLinkModel{ID: 2}, nil
		},
	}
	svc := domain.ExternalLinkService{ExternalLinkRepository: links, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Create("3", domain.Attributes{"label": "Source", "url": "https://x"})

	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
}

func TestExternalLinkService_Show_RejectsANonNumericID(t *testing.T) {
	svc := domain.ExternalLinkService{ExternalLinkRepository: &mockrepo.MockExternalLinkRepository{}}

	_, err := svc.Show("0")

	assert.ErrorContains(t, err, "link ID")
}

func TestExternalLinkService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.ExternalLinkService{ExternalLinkRepository: &mockrepo.MockExternalLinkRepository{}}

	_, err := svc.Update("2", domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestExternalLinkService_Destroy_DelegatesTheParsedID(t *testing.T) {
	var got_id int
	svc := domain.ExternalLinkService{ExternalLinkRepository: &mockrepo.MockExternalLinkRepository{
		DestroyFunc: func(link_id int) error { got_id = link_id; return nil },
	}}

	err := svc.Destroy("2")

	assert.NoError(t, err)
	assert.Equal(t, 2, got_id)
}
