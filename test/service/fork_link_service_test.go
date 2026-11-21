package service_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestForkLinkService_All_ResolvesTheSoftwareName(t *testing.T) {
	var got domain.Filter
	links := &mockrepo.MockForkLinkRepository{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ForkLinkModel], error) {
			got = filter
			return &domain.PageModel[domain.ForkLinkModel]{Items: []domain.ForkLinkModel{{ID: 1}}}, nil
		},
	}
	softwares := &mockrepo.MockSoftwareRepository{
		FindByNameFunc: func(name string) (*domain.SoftwareModel, error) { return &domain.SoftwareModel{ID: 8}, nil },
	}
	svc := domain.ForkLinkService{ForkLinkRepository: links, SoftwareRepository: softwares}

	page, err := svc.All(domain.Filter{"software": "pong", "status": "pending"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"softwareId": "8", "status": "pending"}, got)
	assert.Len(t, page.Items, 1)
}

func TestForkLinkService_All_WithoutSoftwarePassesTheFilter(t *testing.T) {
	var got domain.Filter
	links := &mockrepo.MockForkLinkRepository{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.ForkLinkModel], error) {
			got = filter
			return &domain.PageModel[domain.ForkLinkModel]{}, nil
		},
	}
	svc := domain.ForkLinkService{ForkLinkRepository: links}

	_, err := svc.All(domain.Filter{"status": "visible"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{"status": "visible"}, got)
}

func TestForkLinkService_Update_EmptyAttributesIsAnError(t *testing.T) {
	svc := domain.ForkLinkService{ForkLinkRepository: &mockrepo.MockForkLinkRepository{}}

	_, err := svc.Update(4, domain.Attributes{})

	assert.ErrorIs(t, err, domain.ErrNothingToChange)
}

func TestForkLinkService_Update_DelegatesTheStatus(t *testing.T) {
	var got_id int
	var got domain.Attributes
	svc := domain.ForkLinkService{ForkLinkRepository: &mockrepo.MockForkLinkRepository{
		UpdateFunc: func(fork_link_id int, attributes domain.Attributes) (*domain.ForkLinkModel, error) {
			got_id, got = fork_link_id, attributes
			return &domain.ForkLinkModel{ID: fork_link_id, Status: "visible"}, nil
		},
	}}

	link, err := svc.Update(4, domain.Attributes{"status": "visible"})

	assert.NoError(t, err)
	assert.Equal(t, 4, got_id)
	assert.Equal(t, domain.Attributes{"status": "visible"}, got)
	assert.Equal(t, "visible", link.Status)
}

func TestForkLinkService_Destroy_DelegatesTheID(t *testing.T) {
	var got_id int
	svc := domain.ForkLinkService{ForkLinkRepository: &mockrepo.MockForkLinkRepository{
		DestroyFunc: func(fork_link_id int) error { got_id = fork_link_id; return nil },
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
