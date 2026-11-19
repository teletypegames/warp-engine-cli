package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareService_Show_NumericReferenceIsTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareRepository{
		FindByIDFunc: func(software_id int, include []string) (*domain.SoftwareModel, error) {
			got_id = software_id
			return &domain.SoftwareModel{ID: software_id}, nil
		},
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
			t.Fatal("a numeric reference must not be looked up by name")
			return nil, nil
		},
	}
	svc := domain.SoftwareService{SoftwareRepository: repo}

	software, err := svc.Show("12", nil)

	assert.NoError(t, err)
	assert.Equal(t, 12, got_id)
	assert.Equal(t, 12, software.ID)
}

func TestSoftwareService_Show_NameReferenceIsLookedUp(t *testing.T) {
	repo := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
			assert.Equal(t, "pong", name)
			return &domain.SoftwareModel{ID: 3, Name: name}, nil
		},
		FindByIDFunc: func(software_id int, include []string) (*domain.SoftwareModel, error) {
			return &domain.SoftwareModel{ID: software_id, Name: "pong"}, nil
		},
	}
	svc := domain.SoftwareService{SoftwareRepository: repo}

	software, err := svc.Show("pong", []string{"releases"})

	assert.NoError(t, err)
	assert.Equal(t, 3, software.ID)
}

func TestSoftwareService_Create_NeedsNameTitlePlatform(t *testing.T) {
	repo := &mockrepo.MockSoftwareRepository{
		CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareModel, error) {
			t.Fatal("must not reach the repository")
			return nil, nil
		},
	}
	svc := domain.SoftwareService{SoftwareRepository: repo}

	_, err := svc.Create(domain.Attributes{"name": "pong", "title": "Pong"})

	assert.ErrorContains(t, err, "platform")
}

func TestSoftwareService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.SoftwareService{SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Update("3", domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestSoftwareService_Destroy_DelegatesTheResolvedID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareRepository{DestroyFunc: func(software_id int) error { got_id = software_id; return nil }}
	svc := domain.SoftwareService{SoftwareRepository: repo}

	err := svc.Destroy("9")

	assert.NoError(t, err)
	assert.Equal(t, 9, got_id)
}
