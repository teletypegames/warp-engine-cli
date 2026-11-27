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

func TestSoftwareRepository_All_FollowsPagesWithAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/softwares", r.URL.Path)
		assert.Equal(t, "Bearer secret", r.Header.Get("Authorization"))
		assert.Equal(t, "c64", r.URL.Query().Get("platform"))
		if r.URL.Query().Get("page") == "2" {
			io.WriteString(w, `{"data":[{"id":2,"name":"b"}],"meta":{"page":2,"perPage":1,"total":2,"totalPages":2,"nextPage":null}}`)
			return
		}
		io.WriteString(w, `{"data":[{"id":1,"name":"a"}],"meta":{"page":1,"perPage":1,"total":2,"totalPages":2,"nextPage":2}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareRepository{Client: api.New(server.URL, "secret", "test")}

	page, err := repo.All(domain.Filter{"platform": "c64"}, domain.ListQuery{All: true})

	assert.NoError(t, err)
	assert.Len(t, page.Items, 2)
	assert.Equal(t, 2, page.Total)
}

func TestSoftwareRepository_Update_SendsPatchWithAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/admin/v2/softwares/5", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"title": "New", "author": nil}, body)
		io.WriteString(w, `{"data":{"id":5,"title":"New"}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareRepository{Client: api.New(server.URL, "", "test")}

	software, err := repo.Update(5, domain.Attributes{"title": "New", "author": nil})

	assert.NoError(t, err)
	assert.Equal(t, "New", software.Title)
}

func TestSoftwareRepository_FindByName_MissingIs404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":[{"id":1,"name":"pong-deluxe"}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareRepository{Client: api.New(server.URL, "", "test")}

	_, err := repo.FindByName("pong")

	assert.Equal(t, http.StatusNotFound, api.StatusOf(err))
}

func TestSoftwareRepository_Destroy_ErrorCarriesStatusAndMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, `{"error":"The token lacks admin:write"}`)
	}))
	defer server.Close()
	repo := domain.SoftwareRepository{Client: api.New(server.URL, "", "test")}

	err := repo.Destroy(5)

	assert.Equal(t, http.StatusForbidden, api.StatusOf(err))
	assert.ErrorContains(t, err, "admin:write")
}
