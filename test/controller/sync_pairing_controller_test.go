package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

// pairingRouter returns a router with only the pairing commands over svc.
func pairingRouter(svc *mocksvc.MockSyncPairingService) cli.Router {
	return cli.Router{SyncPairingController: cli.SyncPairingController{SyncPairingService: svc}}
}

func TestSyncPairingController_Create_PrintsTheTokenOnce(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSyncPairingService{
		CreateFunc: func(attributes domain.Attributes) (*domain.SyncPairingModel, error) {
			got = attributes
			return &domain.SyncPairingModel{ID: 2, Label: "laptop local", InboundTokenPrefix: text("9f8e"), PlainToken: text("9f8e-whole")}, nil
		},
	}

	stdout, _, err := execute(pairingRouter(svc), "pairing", "create", "--label", "laptop local")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"label": "laptop local"}, got)
	assert.Contains(t, stdout, "laptop local")
	assert.Contains(t, stdout, "    9f8e-whole\n")
}

func TestSyncPairingController_Create_JSONCarriesThePlainToken(t *testing.T) {
	svc := &mocksvc.MockSyncPairingService{
		CreateFunc: func(attributes domain.Attributes) (*domain.SyncPairingModel, error) {
			return &domain.SyncPairingModel{ID: 2, PlainToken: text("9f8e-whole")}, nil
		},
	}

	stdout, _, err := execute(pairingRouter(svc), "pairing", "create", "--label", "x", "-o", "json")

	assert.NoError(t, err)
	var pairing domain.SyncPairingModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &pairing))
	assert.Equal(t, "9f8e-whole", *pairing.PlainToken)
	assert.NotContains(t, stdout, "shown only this once")
}

func TestSyncPairingController_Show_PrintsNoTokenBlock(t *testing.T) {
	svc := &mocksvc.MockSyncPairingService{
		ShowFunc: func(sync_pairing_id int) (*domain.SyncPairingModel, error) {
			return &domain.SyncPairingModel{ID: sync_pairing_id, Label: "laptop"}, nil
		},
	}

	stdout, _, err := execute(pairingRouter(svc), "pairing", "show", "2")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "laptop")
	assert.NotContains(t, stdout, "shown only this once")
}

func TestSyncPairingController_TokenAndDelete_NeedYes(t *testing.T) {
	calls := []string{}
	svc := &mocksvc.MockSyncPairingService{
		TokenFunc: func(sync_pairing_id int) (*domain.SyncPairingModel, error) {
			calls = append(calls, "token")
			return &domain.SyncPairingModel{ID: sync_pairing_id, PlainToken: text("fresh")}, nil
		},
		DestroyFunc: func(sync_pairing_id int) error { calls = append(calls, "destroy"); return nil },
	}

	_, _, err := execute(pairingRouter(svc), "pairing", "token", "2")
	assert.Error(t, err)
	_, _, err = execute(pairingRouter(svc), "pairing", "delete", "2")
	assert.Error(t, err)
	assert.Empty(t, calls)

	stdout, _, err := execute(pairingRouter(svc), "pairing", "token", "2", "--yes")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "    fresh\n")
	_, _, err = execute(pairingRouter(svc), "pairing", "rm", "2", "-y")
	assert.NoError(t, err)
	assert.Equal(t, []string{"token", "destroy"}, calls)
}
