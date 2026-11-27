package repository_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

// taskJSON is a queued local task as the server answers it.
const taskJSON = `{"data":{"id":7,"kind":"enable","target":"gitea","state":"queued","steps":[],"result":null,"error":null,"createdAt":null,"updatedAt":null}}`

func TestLocalProviderRepository_All_ReadsListWithoutMeta(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/local_providers", r.URL.Path)
		io.WriteString(w, `{"data":[{"key":"gitea","state":"ready","requires":[]},{"key":"woodpecker","state":"off","requires":["forge"]}]}`)
	}))
	defer server.Close()
	repo := domain.LocalProviderRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All()

	assert.NoError(t, err)
	assert.Len(t, page.Items, 2)
	assert.Equal(t, 2, page.Total)
}

func TestLocalProviderRepository_EnableAndDisable_PostToTheAction(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		paths = append(paths, r.URL.Path)
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, taskJSON)
	}))
	defer server.Close()
	repo := domain.LocalProviderRepository{Client: api.New(server.URL, "", "test")}

	task, err := repo.Enable("gitea")
	assert.NoError(t, err)
	assert.Equal(t, 7, task.ID)
	assert.Equal(t, "queued", task.State)
	_, err = repo.Disable("woodpecker")
	assert.NoError(t, err)

	assert.Equal(t, []string{"/api/admin/v2/local_providers/gitea/enable", "/api/admin/v2/local_providers/woodpecker/disable"}, paths)
}

func TestLocalProviderRepository_Destroy_SendsPurgeOnlyWhenSet(t *testing.T) {
	var queries []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/admin/v2/local_providers/gitea", r.URL.Path)
		queries = append(queries, r.URL.RawQuery)
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, taskJSON)
	}))
	defer server.Close()
	repo := domain.LocalProviderRepository{Client: api.New(server.URL, "", "test")}

	task, err := repo.Destroy("gitea", false)
	assert.NoError(t, err)
	assert.Equal(t, 7, task.ID)
	_, err = repo.Destroy("gitea", true)
	assert.NoError(t, err)

	assert.Equal(t, []string{"", "purge=true"}, queries)
}

func TestLocalProviderRepository_Credentials_GetsSignIn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/local_providers/gitea/credentials", r.URL.Path)
		io.WriteString(w, `{"data":{"url":"http://127.0.0.1:7471","user":"warp-admin","password":"secret"}}`)
	}))
	defer server.Close()
	repo := domain.LocalProviderRepository{Client: api.New(server.URL, "", "test")}

	credentials, err := repo.Credentials("gitea")

	assert.NoError(t, err)
	assert.Equal(t, "secret", credentials.Password)
}

func TestLocalProviderRepository_SetupAndTeardownCI(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/softwares/9/local_ci", r.URL.Path)
		calls = append(calls, r.Method)
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, taskJSON)
	}))
	defer server.Close()
	repo := domain.LocalProviderRepository{Client: api.New(server.URL, "", "test")}

	task, err := repo.SetupCI(9)
	assert.NoError(t, err)
	assert.Equal(t, 7, task.ID)
	assert.NoError(t, repo.TeardownCI(9))

	assert.Equal(t, []string{http.MethodPost, http.MethodDelete}, calls)
}

func TestLocalProviderRepository_OffLocalIs404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"error":"Not found"}`)
	}))
	defer server.Close()
	repo := domain.LocalProviderRepository{Client: api.New(server.URL, "", "test")}

	_, err := repo.All()

	assert.Equal(t, http.StatusNotFound, api.StatusOf(err))
}

func TestLocalTaskRepository_FindByID_ReadsStepsAndResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/local_tasks/7", r.URL.Path)
		io.WriteString(w, `{"data":{"id":7,"kind":"local_ci","target":"pong","state":"succeeded","steps":[{"at":"2026-11-19T12:41:28Z","text":"create the repository"}],"result":{"pushUrl":"http://x/pong.git"},"error":null}}`)
	}))
	defer server.Close()
	repo := domain.LocalTaskRepository{Client: api.New(server.URL, "", "test")}

	task, err := repo.FindByID(7)

	assert.NoError(t, err)
	assert.True(t, task.Finished())
	assert.Equal(t, "create the repository", task.Steps[0].Text)
	assert.Equal(t, "http://x/pong.git", task.Result["pushUrl"])
}
