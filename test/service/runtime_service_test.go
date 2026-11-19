package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestRuntimeService_Install_NeedsName(t *testing.T) {
	repo := &mockrepo.MockRuntimeRepository{InstallFunc: func(attributes domain.Attributes) (*domain.RuntimeInstallModel, error) {
		t.Fatal("must not reach the repository")
		return nil, nil
	}}

	_, err := domain.RuntimeService{RuntimeRepository: repo}.Install("", "", false)

	assert.ErrorContains(t, err, "no runtime")
}

func TestRuntimeService_Install_SendsVersionOnlyWhenGiven(t *testing.T) {
	var got []domain.Attributes
	repo := &mockrepo.MockRuntimeRepository{InstallFunc: func(attributes domain.Attributes) (*domain.RuntimeInstallModel, error) {
		got = append(got, attributes)
		return &domain.RuntimeInstallModel{}, nil
	}}
	svc := domain.RuntimeService{RuntimeRepository: repo}

	_, err := svc.Install("tic80", "", true)
	assert.NoError(t, err)
	_, err = svc.Install("tic80", "1.2", false)
	assert.NoError(t, err)

	assert.Equal(t, []domain.Attributes{
		{"name": "tic80", "dryRun": true},
		{"name": "tic80", "version": "1.2", "dryRun": false},
	}, got)
}

func TestRuntimeService_Show_DelegatesID(t *testing.T) {
	repo := &mockrepo.MockRuntimeRepository{FindByIDFunc: func(runtime_id int) (*domain.RuntimeVersionModel, error) {
		return &domain.RuntimeVersionModel{ID: runtime_id}, nil
	}}

	runtime, err := domain.RuntimeService{RuntimeRepository: repo}.Show(3)

	assert.NoError(t, err)
	assert.Equal(t, 3, runtime.ID)
}
