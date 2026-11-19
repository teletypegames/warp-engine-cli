package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSyncPeerService_Create_NeedsLabel(t *testing.T) {
	repo := &mockrepo.MockSyncPeerRepository{
		CreateFunc: func(attributes domain.Attributes) (*domain.SyncPeerModel, error) {
			t.Fatal("must not reach the repository")
			return nil, nil
		},
	}
	svc := domain.SyncPeerService{SyncPeerRepository: repo}

	_, err := svc.Create(domain.Attributes{"transport": "http"})

	assert.ErrorContains(t, err, "label")
}

func TestSyncPeerService_Create_Delegates(t *testing.T) {
	var got domain.Attributes
	repo := &mockrepo.MockSyncPeerRepository{
		CreateFunc: func(attributes domain.Attributes) (*domain.SyncPeerModel, error) {
			got = attributes
			return &domain.SyncPeerModel{ID: 1, Label: "games"}, nil
		},
	}
	svc := domain.SyncPeerService{SyncPeerRepository: repo}

	peer, err := svc.Create(domain.Attributes{"label": "games"})

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"label": "games"}, got)
	assert.Equal(t, 1, peer.ID)
}

func TestSyncPeerService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.SyncPeerService{SyncPeerRepository: &mockrepo.MockSyncPeerRepository{}}

	_, err := svc.Update(3, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestSyncPeerService_Push_DelegatesIDAndCriteria(t *testing.T) {
	var got_id int
	var got domain.Attributes
	repo := &mockrepo.MockSyncPeerRepository{
		PushFunc: func(sync_peer_id int, attributes domain.Attributes) (*domain.SyncPushModel, error) {
			got_id, got = sync_peer_id, attributes
			return &domain.SyncPushModel{DryRun: true}, nil
		},
	}
	svc := domain.SyncPeerService{SyncPeerRepository: repo}

	push, err := svc.Push(3, domain.Attributes{"dryRun": true})

	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
	assert.Equal(t, domain.Attributes{"dryRun": true}, got)
	assert.True(t, push.DryRun)
}

func TestSyncPeerService_ActionsDelegate(t *testing.T) {
	calls := []string{}
	repo := &mockrepo.MockSyncPeerRepository{
		CheckFunc: func(sync_peer_id int) (*domain.SyncCheckModel, error) {
			calls = append(calls, "check")
			return &domain.SyncCheckModel{OK: true}, nil
		},
		InboundTokenFunc: func(sync_peer_id int) (*domain.SyncInboundTokenModel, error) {
			calls = append(calls, "token")
			return &domain.SyncInboundTokenModel{PlainToken: "secret"}, nil
		},
		RunFunc: func(attributes domain.Attributes) (*domain.SyncPushModel, error) {
			calls = append(calls, "run")
			return &domain.SyncPushModel{}, nil
		},
		DestroyFunc: func(sync_peer_id int) error { calls = append(calls, "destroy"); return nil },
	}
	svc := domain.SyncPeerService{SyncPeerRepository: repo}

	_, _ = svc.Check(1)
	_, _ = svc.InboundToken(1)
	_, _ = svc.Run(domain.Attributes{})
	_ = svc.Destroy(1)

	assert.Equal(t, []string{"check", "token", "run", "destroy"}, calls)
}
