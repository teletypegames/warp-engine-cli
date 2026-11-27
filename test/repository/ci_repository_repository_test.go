package repository_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestCiRepositoryRepository_All_SendsFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/admin/v2/ci_repositories", r.URL.Path)
		assert.Equal(t, "c64", r.URL.Query().Get("platform"))
		assert.Equal(t, "true", r.URL.Query().Get("active"))
		assert.Equal(t, "-repoName", r.URL.Query().Get("sort"))
		io.WriteString(w, `{"data":[{"id":1,"fullName":"games/hello-c64","softwareId":null}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`)
	}))
	defer server.Close()
	repo := domain.CiRepositoryRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"platform": "c64", "active": "true"}, domain.ListQuery{Sort: "-repoName"})

	assert.NoError(t, err)
	assert.Equal(t, "games/hello-c64", page.Items[0].FullName)
	assert.Nil(t, page.Items[0].SoftwareID)
}

func TestCiRepositoryRepository_Update_SendsPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/admin/v2/ci_repositories/4", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"softwareId": 9.0}, body)
		io.WriteString(w, `{"data":{"id":4,"softwareId":9}}`)
	}))
	defer server.Close()
	repo := domain.CiRepositoryRepository{Client: api.New(server.URL, "", "test")}

	ci_repository, err := repo.Update(4, domain.Attributes{"softwareId": 9})

	assert.NoError(t, err)
	assert.Equal(t, 9, *ci_repository.SoftwareID)
}

func TestCiRepositoryRepository_Sync_PostsAndReadsChanges(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/ci_repositories/sync", r.URL.Path)
		io.WriteString(w, `{"data":{"created":[{"id":2}],"updated":[],"deactivated":[{"id":3}]}}`)
	}))
	defer server.Close()
	repo := domain.CiRepositoryRepository{Client: api.New(server.URL, "", "test")}

	sync, err := repo.Sync()

	assert.NoError(t, err)
	assert.Equal(t, 2, sync.Created[0].ID)
	assert.Equal(t, 3, sync.Deactivated[0].ID)
}

func TestCiRepositoryRepository_Runs_SendsPageWithoutMeta(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/ci_repositories/4/runs", r.URL.Path)
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		io.WriteString(w, `{"data":[{"number":17,"status":"success"}]}`)
	}))
	defer server.Close()
	repo := domain.CiRepositoryRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.Runs(4, 2)

	assert.NoError(t, err)
	assert.Equal(t, 17, page.Items[0].Number)
	assert.Equal(t, 2, page.Page)
}

func TestCiRepositoryRepository_Trigger_PostsBranch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/ci_repositories/4/trigger", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"branch": "develop"}, body)
		io.WriteString(w, `{"data":{"number":18,"status":"pending"}}`)
	}))
	defer server.Close()
	repo := domain.CiRepositoryRepository{Client: api.New(server.URL, "", "test")}

	run, err := repo.Trigger(4, domain.Attributes{"branch": "develop"})

	assert.NoError(t, err)
	assert.Equal(t, 18, run.Number)
}
