package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestOutletService_Check_RefusesAnEmptyKey(t *testing.T) {
	called := false
	svc := domain.OutletService{OutletRepository: &mockrepo.MockOutletRepository{CheckFunc: func(key string) (*domain.OutletCheckModel, error) {
		called = true
		return nil, nil
	}}}

	_, err := svc.Check("")

	assert.Error(t, err)
	assert.False(t, called)
}

func TestOutletService_Check_AsksTheRepository(t *testing.T) {
	svc := domain.OutletService{OutletRepository: &mockrepo.MockOutletRepository{}}

	check, err := svc.Check("assembly64")

	assert.NoError(t, err)
	assert.Equal(t, "assembly64", check.Key)
}
