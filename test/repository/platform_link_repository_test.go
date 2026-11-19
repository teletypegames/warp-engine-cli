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

func TestPlatformLinkRepository_All_SendsPlatform(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v1/platform_links", r.URL.Path)
		assert.Equal(t, "c64", r.URL.Query().Get("platform"))
		io.WriteString(w, `{"data":[{"id":3,"platform":"c64","name":"VICE","url":"https://vice.example.org","position":0}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`)
	}))
	defer server.Close()
	repo := domain.PlatformLinkRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"platform": "c64"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, "https://vice.example.org", page.Items[0].URL)
}

func TestPlatformLinkRepository_Create_PostsAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v1/platform_links", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"platform": "c64", "name": "VICE", "url": "https://vice.example.org"}, body)
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"id":3,"platform":"c64","name":"VICE","url":"https://vice.example.org","position":0}}`)
	}))
	defer server.Close()
	repo := domain.PlatformLinkRepository{Client: api.New(server.URL, "", "test")}

	link, err := repo.Create(domain.Attributes{"platform": "c64", "name": "VICE", "url": "https://vice.example.org"})

	assert.NoError(t, err)
	assert.Equal(t, 3, link.ID)
}

func TestPlatformLinkRepository_Destroy_SendsDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/admin/v1/platform_links/3", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	repo := domain.PlatformLinkRepository{Client: api.New(server.URL, "", "test")}

	assert.NoError(t, repo.Destroy(3))
}
