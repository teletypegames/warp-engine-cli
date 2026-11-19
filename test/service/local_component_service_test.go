package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestLocalComponentService_Enable_EmptyKeyIsAnError(t *testing.T) {
	svc := domain.LocalComponentService{LocalComponentRepository: &mockrepo.MockLocalComponentRepository{}}

	_, err := svc.Enable("")

	assert.Error(t, err)
}

func TestLocalComponentService_Remove_PassesPurge(t *testing.T) {
	repo := &mockrepo.MockLocalComponentRepository{DestroyFunc: func(key string, purge bool) (*domain.LocalOperationModel, error) {
		assert.Equal(t, "gitea", key)
		assert.True(t, purge)
		return &domain.LocalOperationModel{ID: 5, Kind: "purge"}, nil
	}}

	operation, err := domain.LocalComponentService{LocalComponentRepository: repo}.Remove("gitea", true)

	assert.NoError(t, err)
	assert.Equal(t, 5, operation.ID)
}

func TestLocalComponentService_SetupCI_ResolvesSoftwareName(t *testing.T) {
	repo := &mockrepo.MockLocalComponentRepository{SetupCIFunc: func(software_id int) (*domain.LocalOperationModel, error) {
		assert.Equal(t, 9, software_id)
		return &domain.LocalOperationModel{ID: 1}, nil
	}}
	softwares := &mockrepo.MockSoftwareRepository{FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
		return &domain.SoftwareModel{ID: 9, Name: name}, nil
	}}

	_, err := domain.LocalComponentService{LocalComponentRepository: repo, SoftwareRepository: softwares}.SetupCI("pong")

	assert.NoError(t, err)
}

func TestLocalComponentService_TeardownCI_NumericReferenceIsTheID(t *testing.T) {
	var got int
	repo := &mockrepo.MockLocalComponentRepository{TeardownCIFunc: func(software_id int) error { got = software_id; return nil }}

	err := domain.LocalComponentService{LocalComponentRepository: repo, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}.TeardownCI("12")

	assert.NoError(t, err)
	assert.Equal(t, 12, got)
}
