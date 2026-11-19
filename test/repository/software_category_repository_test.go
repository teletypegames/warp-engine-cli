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

func TestSoftwareCategoryRepository_All_SendsNeedsReview(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v1/software_categories", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("needsReview"))
		io.WriteString(w, `{"data":[{"id":5,"key":"puzzle","label":"Puzzle","position":3,"needsReview":true,"softwareCount":1}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareCategoryRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"needsReview": "true"}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.True(t, page.Items[0].NeedsReview)
}

func TestSoftwareCategoryRepository_Update_SendsPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/admin/v1/software_categories/5", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"needsReview": false}, body)
		io.WriteString(w, `{"data":{"id":5,"key":"puzzle","needsReview":false}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareCategoryRepository{Client: api.New(server.URL, "", "test")}

	category, err := repo.Update(5, domain.Attributes{"needsReview": false})

	assert.NoError(t, err)
	assert.False(t, category.NeedsReview)
}

func TestSoftwareCategoryRepository_FindByID_GetsThePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/admin/v1/software_categories/5", r.URL.Path)
		io.WriteString(w, `{"data":{"id":5,"key":"puzzle"}}`)
	}))
	defer server.Close()
	repo := domain.SoftwareCategoryRepository{Client: api.New(server.URL, "", "test")}

	category, err := repo.FindByID(5)

	assert.NoError(t, err)
	assert.Equal(t, "puzzle", category.Key)
}
