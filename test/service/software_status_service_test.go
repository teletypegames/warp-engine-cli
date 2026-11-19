package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareStatusService_All_DelegatesFilterAndQuery(t *testing.T) {
	var got_filter domain.Filter
	var got_query domain.ListQuery
	repo := &mockrepo.MockSoftwareStatusRepository{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareStatusModel], error) {
			got_filter, got_query = filter, query
			return &domain.PageModel[domain.SoftwareStatusModel]{Items: []domain.SoftwareStatusModel{{ID: 1, Key: "released"}}}, nil
		},
	}
	svc := domain.SoftwareStatusService{SoftwareStatusRepository: repo}

	page, err := svc.All(domain.Filter{"defaultVisible": "true"}, domain.ListQuery{Sort: "key"})

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"defaultVisible": "true"}, got_filter)
	assert.Equal(t, domain.ListQuery{Sort: "key"}, got_query)
	assert.Equal(t, "released", page.Items[0].Key)
}

func TestSoftwareStatusService_Show_DelegatesTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareStatusRepository{FindByIDFunc: func(status_id int) (*domain.SoftwareStatusModel, error) {
		got_id = status_id
		return &domain.SoftwareStatusModel{ID: status_id}, nil
	}}
	svc := domain.SoftwareStatusService{SoftwareStatusRepository: repo}

	status, err := svc.Show(4)

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Equal(t, 4, status.ID)
}

func TestSoftwareStatusService_Create_NeedsKeyAndLabel(t *testing.T) {
	repo := &mockrepo.MockSoftwareStatusRepository{CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
		t.Fatal("must not reach the repository")
		return nil, nil
	}}
	svc := domain.SoftwareStatusService{SoftwareStatusRepository: repo}

	_, err := svc.Create(domain.Attributes{"key": "beta"})

	assert.ErrorContains(t, err, "label")
}

func TestSoftwareStatusService_Create_DelegatesAttributes(t *testing.T) {
	var got domain.Attributes
	repo := &mockrepo.MockSoftwareStatusRepository{CreateFunc: func(attributes domain.Attributes) (*domain.SoftwareStatusModel, error) {
		got = attributes
		return &domain.SoftwareStatusModel{ID: 9, Key: "beta"}, nil
	}}
	svc := domain.SoftwareStatusService{SoftwareStatusRepository: repo}

	status, err := svc.Create(domain.Attributes{"key": "beta", "label": "Beta"})

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"key": "beta", "label": "Beta"}, got)
	assert.Equal(t, 9, status.ID)
}

func TestSoftwareStatusService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.SoftwareStatusService{SoftwareStatusRepository: &mockrepo.MockSoftwareStatusRepository{}}

	_, err := svc.Update(3, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestSoftwareStatusService_Destroy_DelegatesTheID(t *testing.T) {
	var got_id int
	repo := &mockrepo.MockSoftwareStatusRepository{DestroyFunc: func(status_id int) error { got_id = status_id; return nil }}
	svc := domain.SoftwareStatusService{SoftwareStatusRepository: repo}

	err := svc.Destroy(7)

	assert.NoError(t, err)
	assert.Equal(t, 7, got_id)
}
