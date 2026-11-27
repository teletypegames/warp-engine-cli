package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareForkService_All_ResolvesTheSoftwareName(t *testing.T) {
	var got domain.Filter
	links := &mockrepo.MockSoftwareForkRepository{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareForkModel], error) {
			got = filter
			return &domain.PageModel[domain.SoftwareForkModel]{Items: []domain.SoftwareForkModel{{ID: 1}}}, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) { return &domain.SoftwareModel{ID: 8}, nil },
	}
	svc := domain.SoftwareForkService{SoftwareForkRepository: links, SoftwareRepository: softwares}

	page, err := svc.All(domain.Filter{"software": "pong", "status": "pending"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"softwareId": "8", "status": "pending"}, got)
	assert.Len(t, page.Items, 1)
}

func TestSoftwareForkService_All_WithoutSoftwarePassesTheFilter(t *testing.T) {
	var got domain.Filter
	links := &mockrepo.MockSoftwareForkRepository{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.SoftwareForkModel], error) {
			got = filter
			return &domain.PageModel[domain.SoftwareForkModel]{}, nil
		},
	}
	svc := domain.SoftwareForkService{SoftwareForkRepository: links}

	_, err := svc.All(domain.Filter{"status": "visible"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"status": "visible"}, got)
}

func TestSoftwareForkService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.SoftwareForkService{SoftwareForkRepository: &mockrepo.MockSoftwareForkRepository{}}

	_, err := svc.Update(4, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestSoftwareForkService_Update_DelegatesTheStatus(t *testing.T) {
	var got_id int
	var got domain.Attributes
	svc := domain.SoftwareForkService{SoftwareForkRepository: &mockrepo.MockSoftwareForkRepository{
		UpdateFunc: func(software_fork_id int, attributes domain.Attributes) (*domain.SoftwareForkModel, error) {
			got_id, got = software_fork_id, attributes
			return &domain.SoftwareForkModel{ID: software_fork_id, Status: "visible"}, nil
		},
	}}

	link, err := svc.Update(4, domain.Attributes{"status": "visible"})

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Equal(t, domain.Attributes{"status": "visible"}, got)
	assert.Equal(t, "visible", link.Status)
}

func TestSoftwareForkService_Destroy_DelegatesTheID(t *testing.T) {
	var got_id int
	svc := domain.SoftwareForkService{SoftwareForkRepository: &mockrepo.MockSoftwareForkRepository{
		DestroyFunc: func(software_fork_id int) error { got_id = software_fork_id; return nil },
	}}

	err := svc.Destroy(4)

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
}

func TestSoftwareService_Scaffold_ForkNeedsTheParentsName(t *testing.T) {
	svc := domain.SoftwareService{SoftwareRepository: &mockrepo.MockSoftwareRepository{}}

	_, err := svc.Scaffold(domain.Attributes{"name": "pong-plus", "title": "Pong+", "platform": "love",
		"forkOf": map[string]any{"engine": "https://games.example.org"}})

	assert.ErrorContains(t, err, "name")
}
