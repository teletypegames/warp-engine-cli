package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestLocalProviderService_Enable_EmptyKeyIsAnError(t *testing.T) {
	svc := domain.LocalProviderService{LocalProviderRepository: &mockrepo.MockLocalProviderRepository{}}

	_, err := svc.Enable("")

	assert.Error(t, err)
}

func TestLocalProviderService_Remove_PassesPurge(t *testing.T) {
	repo := &mockrepo.MockLocalProviderRepository{DestroyFunc: func(key string, purge bool) (*domain.LocalTaskModel, error) {
		assert.Equal(t, "gitea", key)
		assert.True(t, purge)
		return &domain.LocalTaskModel{ID: 5, Kind: "purge"}, nil
	}}

	task, err := domain.LocalProviderService{LocalProviderRepository: repo}.Remove("gitea", true)

	assert.NoError(t, err)
	assert.Equal(t, 5, task.ID)
}

func TestLocalProviderService_SetupCI_ResolvesSoftwareName(t *testing.T) {
	repo := &mockrepo.MockLocalProviderRepository{SetupCIFunc: func(software_id int) (*domain.LocalTaskModel, error) {
		assert.Equal(t, 9, software_id)
		return &domain.LocalTaskModel{ID: 1}, nil
	}}
	softwares := &mockrepo.MockSoftwareRepository{FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
		return &domain.SoftwareModel{ID: 9, Name: name}, nil
	}}

	_, err := domain.LocalProviderService{LocalProviderRepository: repo, SoftwareRepository: softwares}.SetupCI("pong")

	assert.NoError(t, err)
}

func TestLocalProviderService_TeardownCI_NumericReferenceIsTheID(t *testing.T) {
	var got int
	repo := &mockrepo.MockLocalProviderRepository{TeardownCIFunc: func(software_id int) error { got = software_id; return nil }}

	err := domain.LocalProviderService{LocalProviderRepository: repo, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}.TeardownCI("12")

	assert.NoError(t, err)
	assert.Equal(t, 12, got)
}
