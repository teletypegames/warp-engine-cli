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

func TestSoftwareStatusRepository_All_SendsFilterAndSort(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/admin/v2/software_statuses", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("defaultVisible"))
		assert.Equal(t, "-key", r.URL.Query().Get("sort"))
		io.WriteString(w, `{"data":[{"id":1,"key":"released","label":"Released","position":1,"defaultVisible":true,"badge":null,"softwareCount":4,"createdAt":null,"updatedAt":null}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareStatusRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"defaultVisible": "true"}, domain.ListQuery{Sort: "-key"})

	assert.NoError(t, err)
	assert.Equal(t, "released", page.Items[0].Key)
	assert.Equal(t, 4, page.Items[0].SoftwareCount)
	assert.Nil(t, page.Items[0].Badge)
}

func TestSoftwareStatusRepository_Create_PostsAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/software_statuses", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"key": "beta", "label": "Beta", "defaultVisible": true}, body)
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"id":9,"key":"beta","label":"Beta","defaultVisible":true}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareStatusRepository{Client: api.New(server.URL, "", "test")}

	status, err := repo.Create(domain.Attributes{"key": "beta", "label": "Beta", "defaultVisible": true})

	assert.NoError(t, err)
	assert.Equal(t, 9, status.ID)
}

func TestSoftwareStatusRepository_Destroy_ConflictCarriesTheMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/admin/v2/software_statuses/2", r.URL.Path)
		w.WriteHeader(http.StatusConflict)
		io.WriteString(w, `{"error":"still used by 4 softwares"}`)
	}))
	defer server.Close()
	repo := domain.SoftwareStatusRepository{Client: api.New(server.URL, "", "test")}

	err := repo.Destroy(2)

	assert.Equal(t, http.StatusConflict, api.StatusOf(err))
	assert.ErrorContains(t, err, "still used")
}
