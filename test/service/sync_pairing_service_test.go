package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSyncPairingService_Create_NeedsLabel(t *testing.T) {
	repo := &mockrepo.MockSyncPairingRepository{
		CreateFunc: func(attributes domain.Attributes) (*domain.SyncPairingModel, error) {
			t.Fatal("must not reach the repository")
			return nil, nil
		},
	}
	svc := domain.SyncPairingService{SyncPairingRepository: repo}

	_, err := svc.Create(domain.Attributes{"label": ""})

	assert.ErrorContains(t, err, "label")
}

func TestSyncPairingService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.SyncPairingService{SyncPairingRepository: &mockrepo.MockSyncPairingRepository{}}

	_, err := svc.Update(2, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestSyncPairingService_Token_DelegatesTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSyncPairingRepository{
		TokenFunc: func(sync_pairing_id int) (*domain.SyncPairingModel, error) {
			got_id = sync_pairing_id
			return &domain.SyncPairingModel{ID: sync_pairing_id}, nil
		},
	}
	svc := domain.SyncPairingService{SyncPairingRepository: repo}

	pairing, err := svc.Token(7)

	assert.NoError(t, err)
	assert.Equal(t, 7, got_id)
	assert.Equal(t, 7, pairing.ID)
}
