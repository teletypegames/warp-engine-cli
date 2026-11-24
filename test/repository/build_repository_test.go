package repository_test

import (
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

func TestBuildRepository_Upload_StreamsTheFileWithTheUpdateSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/build/upload", r.URL.Path)
		assert.Equal(t, "pong", r.URL.Query().Get("name"))
		assert.Equal(t, "1.2.0", r.URL.Query().Get("version"))
		assert.Equal(t, "abc", r.URL.Query().Get("sha256"))
		assert.Equal(t, "we_token", r.Header.Get(api.UpdateHeader))
		assert.Empty(t, r.Header.Get("Authorization"))
		assert.Greater(t, r.ContentLength, int64(0))
		file, header, err := r.FormFile("file")
		assert.NoError(t, err)
		assert.Equal(t, "pong-1.2.0-mac-x64.zip", header.Filename)
		content, _ := io.ReadAll(file)
		assert.Equal(t, "zip bytes", string(content))
		io.WriteString(w, `{"file":"pong-1.2.0-mac-x64.zip","size":9,"sha256":"abc"}`)
	}))
	defer server.Close()
	file_path := filepath.Join(t.TempDir(), "pong-1.2.0-mac-x64.zip")
	assert.NoError(t, os.WriteFile(file_path, []byte("zip bytes"), 0o600))
	repo := domain.BuildRepository{Client: api.New(server.URL, "we_token", "test")}

	upload, err := repo.Upload(domain.BuildTargetModel{Name: "pong", Version: "1.2.0"}, file_path, "abc")

	assert.NoError(t, err)
	assert.Equal(t, &domain.BuildUploadModel{File: "pong-1.2.0-mac-x64.zip", Size: 9, SHA256: "abc"}, upload)
}

func TestBuildRepository_Publish_SendsTheTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/build/publish", r.URL.Path)
		assert.Equal(t, "name=gio&platform=golang&version=0.3.0&commit=3fb3ed8&app=gui", sortedQuery(r))
		assert.Equal(t, "we_token", r.Header.Get(api.UpdateHeader))
		io.WriteString(w, `{"published":true,"name":"gio","platform":"golang","version":"0.3.0"}`)
	}))
	defer server.Close()
	repo := domain.BuildRepository{Client: api.New(server.URL, "we_token", "test")}

	published, err := repo.Publish(domain.BuildTargetModel{Name: "gio", Version: "0.3.0", Platform: "golang", App: "gui", Commit: "3fb3ed8"})

	assert.NoError(t, err)
	assert.Equal(t, &domain.BuildPublishModel{Published: true, Name: "gio", Platform: "golang", Version: "0.3.0"}, published)
}

func TestBuildRepository_Publish_ReturnsTheRefusal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error":"Unauthorized"}`)
	}))
	defer server.Close()
	repo := domain.BuildRepository{Client: api.New(server.URL, "bad", "test")}

	_, err := repo.Publish(domain.BuildTargetModel{Name: "pong", Version: "1", Platform: "godot"})

	assert.Equal(t, http.StatusUnauthorized, api.StatusOf(err))
}

// sortedQuery returns the query of r in a fixed order: name, platform,
// version, commit, app.
func sortedQuery(r *http.Request) string {
	query := r.URL.Query()
	text := ""
	for _, key := range []string{"name", "platform", "version", "commit", "app"} {
		if value := query.Get(key); value != "" {
			if text != "" {
				text += "&"
			}
			text += key + "=" + value
		}
	}
	return text
}
