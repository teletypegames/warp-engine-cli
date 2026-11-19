package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestPipelineService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.PipelineService{PipelineRepository: &mockrepo.MockPipelineRepository{}, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Update(4, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestPipelineService_Update_ResolvesSoftwareName(t *testing.T) {
	var got domain.Attributes
	repo := &mockrepo.MockPipelineRepository{UpdateFunc: func(pipeline_id int, attributes domain.Attributes) (*domain.PipelineModel, error) {
		assert.Equal(t, 4, pipeline_id)
		got = attributes
		return &domain.PipelineModel{ID: pipeline_id}, nil
	}}
	softwares := &mockrepo.MockSoftwareRepository{FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
		assert.Equal(t, "pong", name)
		return &domain.SoftwareModel{ID: 9}, nil
	}}
	svc := domain.PipelineService{PipelineRepository: repo, SoftwareRepository: softwares}
	attributes := domain.Attributes{"softwareId": "pong", "platform": "love"}

	_, err := svc.Update(4, attributes)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": 9, "platform": "love"}, got)
	assert.Equal(t, "pong", attributes["softwareId"], "the caller's attributes stay as they were")
}

func TestPipelineService_Update_NilSoftwareDetaches(t *testing.T) {
	var got domain.Attributes
	repo := &mockrepo.MockPipelineRepository{UpdateFunc: func(pipeline_id int, attributes domain.Attributes) (*domain.PipelineModel, error) {
		got = attributes
		return &domain.PipelineModel{}, nil
	}}
	svc := domain.PipelineService{PipelineRepository: repo, SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Update(4, domain.Attributes{"softwareId": nil})

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": nil}, got)
}

func TestPipelineService_Trigger_SendsBranchOnlyWhenGiven(t *testing.T) {
	var got []domain.Attributes
	repo := &mockrepo.MockPipelineRepository{TriggerFunc: func(pipeline_id int, attributes domain.Attributes) (*domain.PipelineRunModel, error) {
		got = append(got, attributes)
		return &domain.PipelineRunModel{Number: 1}, nil
	}}
	svc := domain.PipelineService{PipelineRepository: repo}

	_, err := svc.Trigger(4, "")
	assert.NoError(t, err)
	_, err = svc.Trigger(4, "develop")
	assert.NoError(t, err)

	assert.Equal(t, []domain.Attributes{{}, {"branch": "develop"}}, got)
}

func TestPipelineService_Runs_DelegatesPage(t *testing.T) {
	repo := &mockrepo.MockPipelineRepository{RunsFunc: func(pipeline_id int, page int) (*domain.PageModel[domain.PipelineRunModel], error) {
		assert.Equal(t, 4, pipeline_id)
		assert.Equal(t, 2, page)
		return &domain.PageModel[domain.PipelineRunModel]{Items: []domain.PipelineRunModel{{Number: 3}}}, nil
	}}
	svc := domain.PipelineService{PipelineRepository: repo}

	page, err := svc.Runs(4, 2)

	assert.NoError(t, err)
	assert.Equal(t, 3, page.Items[0].Number)
}
