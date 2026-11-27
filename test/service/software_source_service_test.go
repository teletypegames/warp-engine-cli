package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareSourceService_Set_NeedsURL(t *testing.T) {
	svc := domain.SoftwareSourceService{SoftwareSourceRepository: &mockrepo.MockSoftwareSourceRepository{}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Set("8", domain.Attributes{"ref": "main"})

	assert.ErrorContains(t, err, "url")
}

func TestSoftwareSourceService_Set_ResolvesTheSoftwareName(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareSourceRepository{SaveFunc: func(software_id int, attributes domain.Attributes) (*domain.SoftwareSourceModel, error) {
		got_id = software_id
		return &domain.SoftwareSourceModel{SoftwareID: software_id}, nil
	}}
	software := &mockrepo.MockSoftwareRepository{FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
		return &domain.SoftwareModel{ID: 8, Name: name}, nil
	}}
	svc := domain.SoftwareSourceService{SoftwareSourceRepository: repo, SoftwareRepository: software}

	source, err := svc.Set("pong", domain.Attributes{"url": "https://git.example.org/pong"})

	assert.NoError(t, err)
	assert.Equal(t, 8, got_id)
	assert.Equal(t, 8, source.SoftwareID)
}

func TestSoftwareSourceService_Destroy_NumericReferenceIsTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareSourceRepository{DestroyFunc: func(software_id int) error { got_id = software_id; return nil }}
	svc := domain.SoftwareSourceService{SoftwareSourceRepository: repo, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	assert.NoError(t, svc.Destroy("12"))
	assert.Equal(t, 12, got_id)
}
