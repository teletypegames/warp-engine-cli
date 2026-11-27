package repository_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestFileRepository_All_SendsDirOnlyWhenSet(t *testing.T) {
	var queries []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/files", r.URL.Path)
		queries = append(queries, r.URL.RawQuery)
		io.WriteString(w, `{"data":[{"name":"a.zip","path":"pong/a.zip","type":"file","size":1024,"modifiedAt":null,"downloads":3}]}`)
	}))
	defer server.Close()
	repo := domain.FileRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All("pong")
	assert.NoError(t, err)
	_, err = repo.All("")
	assert.NoError(t, err)

	assert.Equal(t, []string{"dir=pong", ""}, queries)
	assert.Equal(t, int64(1024), *page.Items[0].Size)
	assert.Equal(t, 1, page.TotalPages)
}

func TestFileRepository_Upload_SendsMultipart(t *testing.T) {
	local := filepath.Join(t.TempDir(), "manual.pdf")
	assert.NoError(t, os.WriteFile(local, []byte("PDF"), 0o600))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/files", r.URL.Path)
		assert.NoError(t, r.ParseMultipartForm(1<<20))
		assert.Equal(t, "docs", r.FormValue("dir"))
		file, header, err := r.FormFile("file")
		assert.NoError(t, err)
		content, _ := io.ReadAll(file)
		assert.Equal(t, "manual.pdf", header.Filename)
		assert.Equal(t, "PDF", string(content))
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"path":"docs/manual.pdf"}}`)
	}))
	defer server.Close()
	repo := domain.FileRepository{Client: api.New(server.URL, "", "test")}

	result, err := repo.Upload("docs", local)

	assert.NoError(t, err)
	assert.Equal(t, "docs/manual.pdf", result.Path)
}

func TestFileRepository_Rename_SendsPathAndNewName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/admin/v2/files", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"path": "docs/a.pdf", "newName": "b.pdf"}, body)
		io.WriteString(w, `{"data":{"path":"docs/b.pdf"}}`)
	}))
	defer server.Close()
	repo := domain.FileRepository{Client: api.New(server.URL, "", "test")}

	result, err := repo.Rename("docs/a.pdf", "b.pdf")

	assert.NoError(t, err)
	assert.Equal(t, "docs/b.pdf", result.Path)
}

func TestFileRepository_Destroy_SendsPathQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/admin/v2/files", r.URL.Path)
		assert.Equal(t, "docs/a b.pdf", r.URL.Query().Get("path"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	repo := domain.FileRepository{Client: api.New(server.URL, "", "test")}

	assert.NoError(t, repo.Destroy("docs/a b.pdf"))
}

func TestFileRepository_CreateDirectory_PostsDirAndName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/files/directories", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"dir": "pong", "name": "docs"}, body)
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"path":"pong/docs"}}`)
	}))
	defer server.Close()
	repo := domain.FileRepository{Client: api.New(server.URL, "", "test")}

	result, err := repo.CreateDirectory("pong", "docs")

	assert.NoError(t, err)
	assert.Equal(t, "pong/docs", result.Path)
}
