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

func TestRuntimeRepository_All_SendsName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v1/runtimes", r.URL.Path)
		assert.Equal(t, "emulatorjs", r.URL.Query().Get("name"))
		io.WriteString(w, `{"data":[{"id":1,"name":"emulatorjs","sizeBytes":10380000,"components":[{"kind":"core","key":"vice_x64sc"}]}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`)
	}))
	defer server.Close()
	repo := domain.RuntimeRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"name": "emulatorjs"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, "vice_x64sc", page.Items[0].Components[0].Key)
}

func TestRuntimeRepository_FindByID_GetsOne(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v1/runtimes/3", r.URL.Path)
		io.WriteString(w, `{"data":{"id":3,"version":"4.2.3"}}`)
	}))
	defer server.Close()
	repo := domain.RuntimeRepository{Client: api.New(server.URL, "", "test")}

	runtime, err := repo.FindByID(3)

	assert.NoError(t, err)
	assert.Equal(t, "4.2.3", runtime.Version)
}

func TestRuntimeRepository_Install_PostsBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v1/runtimes/install", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"name": "tic80", "dryRun": true}, body)
		io.WriteString(w, `{"data":{"name":"tic80","version":"1.2","fetched":["a.js"],"kept":[],"components":{},"dryRun":true}}`)
	}))
	defer server.Close()
	repo := domain.RuntimeRepository{Client: api.New(server.URL, "", "test")}

	install, err := repo.Install(domain.Attributes{"name": "tic80", "dryRun": true})

	assert.NoError(t, err)
	assert.True(t, install.DryRun)
	assert.Equal(t, []string{"a.js"}, install.Fetched)
}
