package repository_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestForkLinkRepository_All_SendsFiltersAndSort(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":[{"id":4,"softwareId":8,"childEngine":"https://other.example.org","childName":"pong-plus","status":"pending"}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`, &seen)
	defer server.Close()
	repo := domain.ForkLinkRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"status": "pending", "softwareId": "8"}, domain.ListQuery{Sort: "-checkedAt"})

	assert.NoError(t, err)
	assert.Equal(t, "GET", seen.method)
	assert.Equal(t, "/api/admin/v1/fork_links", seen.path)
	assert.Equal(t, "softwareId=8&sort=-checkedAt&status=pending", seen.query)
	assert.Equal(t, "pong-plus", page.Items[0].ChildName)
	assert.Equal(t, 1, page.Total)
}

func TestForkLinkRepository_FindByID_GetsTheLink(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"id":4,"status":"visible","lastError":"timeout"}}`, &seen)
	defer server.Close()
	repo := domain.ForkLinkRepository{Client: api.New(server.URL, "", "test")}

	link, err := repo.FindByID(4)

	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v1/fork_links/4", seen.path)
	assert.Equal(t, "timeout", *link.LastError)
}

func TestForkLinkRepository_Update_PatchesTheStatus(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"id":4,"status":"hidden"}}`, &seen)
	defer server.Close()
	repo := domain.ForkLinkRepository{Client: api.New(server.URL, "", "test")}

	link, err := repo.Update(4, domain.Attributes{"status": "hidden"})

	assert.NoError(t, err)
	assert.Equal(t, "PATCH", seen.method)
	assert.Equal(t, "/api/admin/v1/fork_links/4", seen.path)
	assert.Equal(t, map[string]any{"status": "hidden"}, seen.body)
	assert.Equal(t, "hidden", link.Status)
}

func TestForkLinkRepository_Destroy_SendsDelete(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 204, ``, &seen)
	defer server.Close()
	repo := domain.ForkLinkRepository{Client: api.New(server.URL, "", "test")}

	err := repo.Destroy(4)

	assert.NoError(t, err)
	assert.Equal(t, "DELETE", seen.method)
	assert.Equal(t, "/api/admin/v1/fork_links/4", seen.path)
}

func TestExternalLinkRepository_Create_SendsKindAndRef(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 201, `{"data":{"id":6,"softwareId":8,"label":"Source","url":"https://git.example.org/pong","kind":"source","ref":"main","public":true,"verifiedAt":"2026-11-20T10:00:00Z","verifyError":null}}`, &seen)
	defer server.Close()
	repo := domain.ExternalLinkRepository{Client: api.New(server.URL, "", "test")}

	link, err := repo.Create(8, domain.Attributes{"label": "Source", "url": "https://git.example.org/pong", "kind": "source", "ref": "main"})

	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v1/softwares/8/external_links", seen.path)
	assert.Equal(t, "source", seen.body["kind"])
	assert.Equal(t, "main", seen.body["ref"])
	assert.Equal(t, "source", link.Kind)
	assert.True(t, *link.Public)
	assert.Nil(t, link.VerifyError)
}
