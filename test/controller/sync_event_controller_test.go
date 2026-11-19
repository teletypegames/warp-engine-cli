package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

// eventRouter returns a router with only the sync-event commands over svc.
func eventRouter(svc *mocksvc.MockSyncEventService) cli.Router {
	return cli.Router{SyncEventController: cli.SyncEventController{SyncEventService: svc}}
}

func TestSyncEventController_Index_MapsSoftwareNameAndID(t *testing.T) {
	filters := []domain.Filter{}
	svc := &mocksvc.MockSyncEventService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncEventModel], error) {
			filters = append(filters, filter)
			return &domain.PageModel[domain.SyncEventModel]{Items: []domain.SyncEventModel{
				{ID: 9, Direction: "outbound", Status: "failed", SoftwareName: text("pong"), Message: text("502 from peer")},
			}}, nil
		},
	}

	stdout, _, err := execute(eventRouter(svc), "sync-events", "list", "--software", "pong", "--status", "failed", "--peer", "3")
	assert.NoError(t, err)
	_, _, err = execute(eventRouter(svc), "sync-event", "ls", "--software", "12")
	assert.NoError(t, err)

	assert.Equal(t, domain.Filter{"software": "pong", "status": "failed", "syncPeerId": "3"}, filters[0])
	assert.Equal(t, domain.Filter{"softwareId": "12"}, filters[1])
	assert.Contains(t, stdout, "502 from peer")
}

func TestSyncEventController_Software_PassesReferenceAndDirection(t *testing.T) {
	var got_reference, got_direction string
	svc := &mocksvc.MockSyncEventService{
		SoftwareFunc: func(reference string, direction string) ([]domain.SyncEventModel, error) {
			got_reference, got_direction = reference, direction
			return []domain.SyncEventModel{{ID: 1, Direction: "inbound", Status: "succeeded"}}, nil
		},
	}

	stdout, _, err := execute(eventRouter(svc), "sync-event", "software", "pong", "--direction", "inbound")

	assert.NoError(t, err)
	assert.Equal(t, "pong", got_reference)
	assert.Equal(t, "inbound", got_direction)
	assert.Contains(t, stdout, "succeeded")
}

func TestSyncEventController_Software_RejectsUnknownDirection(t *testing.T) {
	_, _, err := execute(eventRouter(&mocksvc.MockSyncEventService{}), "sync-event", "software", "pong", "--direction", "up")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSyncEventController_Status_PrintsEveryPeer(t *testing.T) {
	svc := &mocksvc.MockSyncEventService{
		StatusFunc: func(reference string) ([]domain.SyncStatusModel, error) {
			assert.Equal(t, "pong", reference)
			return []domain.SyncStatusModel{
				{PeerID: 3, Peer: "games", Direction: "outbound", Event: &domain.SyncEventModel{Status: "succeeded", Version: text("1.2.0")}},
				{PeerID: 4, Peer: "mirror", Direction: "outbound", Reason: text("channel beta not allowed")},
			}, nil
		},
	}

	stdout, _, err := execute(eventRouter(svc), "sync-event", "status", "pong")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "1.2.0")
	assert.Contains(t, stdout, "channel beta not allowed")
}

func TestSyncEventController_Show_PassesTheID(t *testing.T) {
	var got_id int
	svc := &mocksvc.MockSyncEventService{
		ShowFunc: func(sync_event_id int) (*domain.SyncEventModel, error) {
			got_id = sync_event_id
			return &domain.SyncEventModel{ID: sync_event_id, Status: "failed"}, nil
		},
	}

	_, _, err := execute(eventRouter(svc), "sync-event", "show", "41")

	assert.NoError(t, err)
	assert.Equal(t, 41, got_id)
}
