package repository_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareForkRepository_All_SendsFiltersAndSort(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":[{"id":4,"softwareId":8,"engineUrl":"https://other.example.org","name":"pong-plus","status":"pending"}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`, &seen)
	defer server.Close()
	repo := domain.SoftwareForkRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"status": "pending", "softwareId": "8"}, domain.ListQuery{Sort: "-checkedAt"})

	assert.NoError(t, err)
	assert.Equal(t, "GET", seen.method)
	assert.Equal(t, "/api/admin/v2/forks", seen.path)
	assert.Equal(t, "softwareId=8&sort=-checkedAt&status=pending", seen.query)
	assert.Equal(t, "pong-plus", page.Items[0].Name)
	assert.Equal(t, 1, page.Total)
}

func TestSoftwareForkRepository_FindByID_GetsTheLink(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"id":4,"status":"visible","lastError":"timeout"}}`, &seen)
	defer server.Close()
	repo := domain.SoftwareForkRepository{Client: api.New(server.URL, "", "test")}

	link, err := repo.FindByID(4)

	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/forks/4", seen.path)
	assert.Equal(t, "timeout", *link.LastError)
}

func TestSoftwareForkRepository_Update_PatchesTheStatus(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"id":4,"status":"hidden"}}`, &seen)
	defer server.Close()
	repo := domain.SoftwareForkRepository{Client: api.New(server.URL, "", "test")}

	link, err := repo.Update(4, domain.Attributes{"status": "hidden"})

	assert.NoError(t, err)
	assert.Equal(t, "PATCH", seen.method)
	assert.Equal(t, "/api/admin/v2/forks/4", seen.path)
	assert.Equal(t, map[string]any{"status": "hidden"}, seen.body)
	assert.Equal(t, "hidden", link.Status)
}

func TestSoftwareForkRepository_Destroy_SendsDelete(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 204, ``, &seen)
	defer server.Close()
	repo := domain.SoftwareForkRepository{Client: api.New(server.URL, "", "test")}

	err := repo.Destroy(4)

	assert.NoError(t, err)
	assert.Equal(t, "DELETE", seen.method)
	assert.Equal(t, "/api/admin/v2/forks/4", seen.path)
}

func TestSoftwareSourceRepository_Save_PutsURLAndRef(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"id":6,"softwareId":8,"url":"https://git.example.org/pong","ref":"main","public":true,"verifiedAt":"2026-11-20T10:00:00Z","verifyError":null}}`, &seen)
	defer server.Close()
	repo := domain.SoftwareSourceRepository{Client: api.New(server.URL, "", "test")}

	source, err := repo.Save(8, domain.Attributes{"url": "https://git.example.org/pong", "ref": "main"})

	assert.NoError(t, err)
	assert.Equal(t, "PUT", seen.method)
	assert.Equal(t, "/api/admin/v2/softwares/8/source", seen.path)
	assert.Equal(t, "main", seen.body["ref"])
	assert.Equal(t, "https://git.example.org/pong", source.URL)
	assert.True(t, *source.Public)
	assert.Nil(t, source.VerifyError)
}

func TestSoftwareSourceRepository_FindAndDestroy_UseTheSoftwarePath(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"id":6,"softwareId":8,"url":"https://git.example.org/pong","ref":null}}`, &seen)
	defer server.Close()
	repo := domain.SoftwareSourceRepository{Client: api.New(server.URL, "", "test")}

	source, err := repo.Find(8)
	assert.NoError(t, err)
	assert.Equal(t, "GET", seen.method)
	assert.Equal(t, "/api/admin/v2/softwares/8/source", seen.path)
	assert.Nil(t, source.Ref)

	assert.NoError(t, repo.Destroy(8))
	assert.Equal(t, "DELETE", seen.method)
	assert.Equal(t, "/api/admin/v2/softwares/8/source", seen.path)
}
