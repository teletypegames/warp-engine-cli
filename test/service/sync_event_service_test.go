package service_test

import (
	"net/http"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSyncEventService_Software_ResolvesNameAndPassesDirection(t *testing.T) {
	var got_id int
	var got_direction string
	events := &mockrepo.MockSyncEventRepository{
		AllOfSoftwareFunc: func(software_id int, direction string) ([]domain.SyncEventModel, error) {
			got_id, got_direction = software_id, direction
			return []domain.SyncEventModel{{ID: 1}}, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
			assert.Equal(t, "pong", name)
			return &domain.SoftwareModel{ID: 12, Name: name}, nil
		},
	}
	svc := domain.SyncEventService{SyncEventRepository: events, SoftwareRepository: softwares}

	items, err := svc.Software("pong", "outbound")

	assert.NoError(t, err)
	assert.Equal(t, 12, got_id)
	assert.Equal(t, "outbound", got_direction)
	assert.Len(t, items, 1)
}

func TestSyncEventService_Status_NumericReferenceIsTheID(t *testing.T) {
	var got_id int
	events := &mockrepo.MockSyncEventRepository{
		StatusFunc: func(software_id int) ([]domain.SyncStatusModel, error) {
			got_id = software_id
			return []domain.SyncStatusModel{{PeerID: 3}}, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
			t.Fatal("a numeric reference must not be looked up by name")
			return nil, nil
		},
	}
	svc := domain.SyncEventService{SyncEventRepository: events, SoftwareRepository: softwares}

	states, err := svc.Status("9")

	assert.NoError(t, err)
	assert.Equal(t, 9, got_id)
	assert.Equal(t, 3, states[0].PeerID)
}

func TestSyncEventService_Status_UnknownNameStopsBeforeTheRequest(t *testing.T) {
	events := &mockrepo.MockSyncEventRepository{
		StatusFunc: func(software_id int) ([]domain.SyncStatusModel, error) {
			t.Fatal("must not reach the sync events")
			return nil, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
			return nil, &api.Error{Status: http.StatusNotFound, Message: "no software named " + name}
		},
	}
	svc := domain.SyncEventService{SyncEventRepository: events, SoftwareRepository: softwares}

	_, err := svc.Status("nope")

	assert.Equal(t, http.StatusNotFound, api.StatusOf(err))
}

func TestSyncEventService_All_PassesFilterAndQuery(t *testing.T) {
	var got_filter domain.Filter
	events := &mockrepo.MockSyncEventRepository{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SyncEventModel], error) {
			got_filter = filter
			return &domain.PageModel[domain.SyncEventModel]{}, nil
		},
	}
	svc := domain.SyncEventService{SyncEventRepository: events, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.All(domain.Filter{"status": "failed"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"status": "failed"}, got_filter)
}
