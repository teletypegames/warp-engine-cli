package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

// peerRouter returns a router with only the peer commands over svc.
func peerRouter(svc *mocksvc.MockSyncPeerService) cli.Router {
	return cli.Router{SyncPeerController: cli.SyncPeerController{SyncPeerService: svc}}
}

func TestSyncPeerController_Index_PassesFiltersAndPrintsTable(t *testing.T) {
	var got_filter domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockSyncPeerService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncPeerModel], error) {
			got_filter, got_query = filter, query
			return &domain.PageModel[domain.SyncPeerModel]{Items: []domain.SyncPeerModel{
				{ID: 3, Label: "games", Transport: "http", TransportOptions: map[string]any{"host": "https://games.example.org"}, OutboundEnabled: true},
			}}, nil
		},
	}

	stdout, _, err := execute(peerRouter(svc), "peer", "list", "--outbound", "--transport", "http", "--sort", "-lastSyncedAt")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"outboundEnabled": "true", "transport": "http"}, got_filter)
	assert.Equal(t, domain.ListQuery{Sort: "-lastSyncedAt"}, got_query)
	assert.Contains(t, stdout, "https://games.example.org")
}

func TestSyncPeerController_Create_SendsFieldsAndOptions(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSyncPeerService{
		CreateFunc: func(attributes domain.Attributes) (*domain.SyncPeerModel, error) {
			got = attributes
			return &domain.SyncPeerModel{ID: 4, Label: "mirror"}, nil
		},
	}

	stdout, _, err := execute(peerRouter(svc), "peer", "create", "--label", "mirror", "--transport", "rsync",
		"--option", "host=mirror.example.org", "--option", "user=warp", "--outbound=false", "--channels", "stable,beta")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{
		"label": "mirror", "transport": "rsync", "outboundEnabled": false,
		"allowedChannels":  []string{"stable", "beta"},
		"transportOptions": map[string]any{"host": "mirror.example.org", "user": "warp"},
	}, got)
	assert.Contains(t, stdout, "mirror")
}

func TestSyncPeerController_Update_BadIDIsUsageError(t *testing.T) {
	_, _, err := execute(peerRouter(&mocksvc.MockSyncPeerService{}), "peer", "update", "games", "--label", "x")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSyncPeerController_Destroy_NeedsYes(t *testing.T) {
	var got_id int
	svc := &mocksvc.MockSyncPeerService{DestroyFunc: func(sync_peer_id int) error { got_id = sync_peer_id; return nil }}

	_, _, err := execute(peerRouter(svc), "peer", "delete", "3")
	assert.Error(t, err)
	assert.Equal(t, 0, got_id)

	stdout, _, err := execute(peerRouter(svc), "peer", "delete", "3", "--yes")
	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
	assert.Contains(t, stdout, "deleted")
}

func TestSyncPeerController_InboundToken_PrintsThePlainTokenOnce(t *testing.T) {
	svc := &mocksvc.MockSyncPeerService{
		InboundTokenFunc: func(sync_peer_id int) (*domain.SyncInboundTokenModel, error) {
			return &domain.SyncInboundTokenModel{TokenPrefix: "abcd1234", PlainToken: "abcd1234-the-whole-secret"}, nil
		},
	}

	_, _, err := execute(peerRouter(svc), "peer", "inbound-token", "3")
	assert.Error(t, err)

	stdout, _, err := execute(peerRouter(svc), "peer", "inbound-token", "3", "--yes")
	assert.NoError(t, err)
	assert.Contains(t, stdout, "shown only this once")
	assert.Contains(t, stdout, "    abcd1234-the-whole-secret\n")
}

func TestSyncPeerController_Push_DryRunNeedsNoYes(t *testing.T) {
	var got_id int
	var got domain.Attributes
	svc := &mocksvc.MockSyncPeerService{
		PushFunc: func(sync_peer_id int, attributes domain.Attributes) (*domain.SyncPushModel, error) {
			got_id, got = sync_peer_id, attributes
			return &domain.SyncPushModel{DryRun: true, Releases: []domain.SyncPushReleaseModel{{Software: "pong", Version: "1.0.0"}}}, nil
		},
	}

	stdout, _, err := execute(peerRouter(svc), "peer", "push", "3", "--dry-run", "--software", "pong", "--filtered=false")

	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
	assert.Equal(t, domain.Attributes{"dryRun": true, "software": "pong", "filtered": false}, got)
	assert.Contains(t, stdout, "pong 1.0.0")
}

func TestSyncPeerController_Push_RealPushNeedsYes(t *testing.T) {
	called := false
	svc := &mocksvc.MockSyncPeerService{
		PushFunc: func(sync_peer_id int, attributes domain.Attributes) (*domain.SyncPushModel, error) {
			called = true
			return &domain.SyncPushModel{Queued: 2}, nil
		},
	}

	_, _, err := execute(peerRouter(svc), "peer", "push", "3")
	assert.Error(t, err)
	assert.False(t, called)

	_, _, err = execute(peerRouter(svc), "peer", "push", "3", "--yes")
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSyncPeerController_Run_SendsPeerIDs(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSyncPeerService{
		RunFunc: func(attributes domain.Attributes) (*domain.SyncPushModel, error) {
			got = attributes
			return &domain.SyncPushModel{}, nil
		},
	}

	stdout, _, err := execute(peerRouter(svc), "peer", "run", "--peer", "3", "--peer", "5", "--all", "--yes", "-o", "json")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"peerIds": []int{3, 5}, "all": true}, got)
	var push domain.SyncPushModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &push))
}

func TestSyncPeerController_Check_RefusalExitsNonZero(t *testing.T) {
	svc := &mocksvc.MockSyncPeerService{
		CheckFunc: func(sync_peer_id int) (*domain.SyncCheckModel, error) {
			return &domain.SyncCheckModel{OK: false, Message: "401: token refused"}, nil
		},
	}

	stdout, _, err := execute(peerRouter(svc), "peer", "check", "3")

	assert.ErrorIs(t, err, cli.ErrCheckFailed)
	assert.Contains(t, stdout, "token refused")
}
