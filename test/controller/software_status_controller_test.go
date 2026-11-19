package controller_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareStatusController_Index_PrintsTableAndPassesFilters(t *testing.T) {
	var got_filter domain.Filter
	var got_query domain.ListQuery
	svc := &mocksvc.MockSoftwareStatusService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareStatusModel], error) {
			got_filter, got_query = filter, query
			return &domain.PageModel[domain.SoftwareStatusModel]{
				Items: []domain.SoftwareStatusModel{{ID: 2, Key: "released", Label: "Released", DefaultVisible: true, SoftwareCount: 14}},
			}, nil
		},
	}

	stdout, _, err := execute(cli.Router{SoftwareStatusController: cli.SoftwareStatusController{SoftwareStatusService: svc}},
		"statuses", "ls", "--default-visible", "--sort", "-key")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"defaultVisible": "true"}, got_filter)
	assert.Equal(t, domain.ListQuery{Sort: "-key"}, got_query)
	assert.Contains(t, stdout, "released")
	assert.Contains(t, stdout, "14")
}

func TestSoftwareStatusController_Show_ParsesTheID(t *testing.T) {
	var got_id int
	svc := &mocksvc.MockSoftwareStatusService{ShowFunc: func(status_id int) (*domain.SoftwareStatusModel, error) {
		got_id = status_id
		return &domain.SoftwareStatusModel{ID: status_id, Key: "beta", Badge: text("BETA")}, nil
	}}

	stdout, _, err := execute(cli.Router{SoftwareStatusController: cli.SoftwareStatusController{SoftwareStatusService: svc}}, "status", "show", "4")

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Contains(t, stdout, "BETA")
}

func TestSoftwareStatusController_Show_NonNumericIDIsAUsageError(t *testing.T) {
	svc := &mocksvc.MockSoftwareStatusService{ShowFunc: func(status_id int) (*domain.SoftwareStatusModel, error) {
		t.Fatal("must not reach the service")
		return nil, nil
	}}

	_, _, err := execute(cli.Router{SoftwareStatusController: cli.SoftwareStatusController{SoftwareStatusService: svc}}, "status", "show", "released")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSoftwareStatusController_Create_SendsTheFields(t *testing.T) {
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareStatusService{CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
		got = attributes
		return &domain.SoftwareStatusModel{ID: 9, Key: "beta"}, nil
	}}

	stdout, _, err := execute(cli.Router{SoftwareStatusController: cli.SoftwareStatusController{SoftwareStatusService: svc}},
		"status", "create", "--key", "beta", "--label", "Beta", "--position", "3", "--default-visible", "-o", "id")

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"key": "beta", "label": "Beta", "position": 3, "defaultVisible": true}, got)
	assert.Equal(t, "9\n", stdout)
}

func TestSoftwareStatusController_Update_EmptyBadgeSendsNull(t *testing.T) {
	var got_id int
	var got domain.Attributes
	svc := &mocksvc.MockSoftwareStatusService{UpdateFunc: func(status_id int, attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
		got_id, got = status_id, attributes
		return &domain.SoftwareStatusModel{ID: status_id}, nil
	}}

	_, _, err := execute(cli.Router{SoftwareStatusController: cli.SoftwareStatusController{SoftwareStatusService: svc}},
		"status", "update", "4", "--badge", "")

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Equal(t, domain.Attributes{"badge": nil}, got)
}

func TestSoftwareStatusController_Destroy_NeedsYesWithoutTerminal(t *testing.T) {
	var got_id int
	svc := &mocksvc.MockSoftwareStatusService{DestroyFunc: func(status_id int) error { got_id = status_id; return nil }}
	router := cli.Router{SoftwareStatusController: cli.SoftwareStatusController{SoftwareStatusService: svc}}

	_, _, err := execute(router, "status", "delete", "4")
	assert.Error(t, err)
	assert.Zero(t, got_id)

	stdout, _, err := execute(router, "status", "rm", "4", "--yes")
	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Contains(t, stdout, "Status 4 deleted.")
}
