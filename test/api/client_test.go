package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"github.com/stretchr/testify/assert"
)

func TestClient_Do_NoServerIsErrNoServer(t *testing.T) {
	err := api.New("", "", "test").Do(http.MethodGet, "/x", nil, nil, nil)

	assert.ErrorIs(t, err, api.ErrNoServer)
}

func TestClient_Do_ErrorCarriesDetailsAndRetryAfter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, `{"error":"Validation failed","details":{"name":["is taken"],"title":["is blank"]}}`)
	}))
	defer server.Close()

	err := api.New(server.URL, "", "test").Admin(http.MethodPost, "/softwares", nil, map[string]any{}, nil)

	var api_error *api.Error
	assert.ErrorAs(t, err, &api_error)
	assert.Equal(t, 422, api_error.Status)
	assert.Equal(t, 30*time.Second, api_error.RetryAfter)
	assert.Equal(t, "422 Unprocessable Entity: Validation failed (name is taken; title is blank)", err.Error())
}

func TestClient_Do_NonJSONErrorUsesTheStatusText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		io.WriteString(w, "<html>bad gateway</html>")
	}))
	defer server.Close()

	err := api.New(server.URL, "", "test").Do(http.MethodGet, "/x", nil, nil, nil)

	assert.Equal(t, 502, api.StatusOf(err))
	assert.ErrorContains(t, err, "Bad Gateway")
}

func TestClient_Do_SendsHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer tok", r.Header.Get("Authorization"))
		assert.Equal(t, "agent/1", r.Header.Get("User-Agent"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := api.New(server.URL+"/", "tok", "agent/1").Do(http.MethodPost, "/x", nil, map[string]any{"a": 1}, nil)

	assert.NoError(t, err)
}

func TestClient_Upload_SendsMultipart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shot.png")
	assert.NoError(t, os.WriteFile(path, []byte("PNG"), 0o600))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/files", r.URL.Path)
		assert.NoError(t, r.ParseMultipartForm(1<<20))
		assert.Equal(t, "games", r.FormValue("dir"))
		file, header, err := r.FormFile("file")
		assert.NoError(t, err)
		data, _ := io.ReadAll(file)
		assert.Equal(t, "PNG", string(data))
		assert.Equal(t, "shot.png", header.Filename)
		io.WriteString(w, `{"data":{"path":"games/shot.png"}}`)
	}))
	defer server.Close()

	var answer api.Envelope[map[string]string]
	err := api.New(server.URL, "", "test").Upload("/files", map[string]string{"dir": "games"}, "file", path, &answer)

	assert.NoError(t, err)
	assert.Equal(t, "games/shot.png", answer.Data["path"])
}

func TestList_WithoutAllReadsOnePage(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		io.WriteString(w, `{"data":[1,2],"meta":{"page":1,"perPage":2,"total":4,"totalPages":2,"nextPage":2}}`)
	}))
	defer server.Close()

	items, meta, err := api.List[int](api.New(server.URL, "", "test"), "/x", nil, false)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, items)
	assert.Equal(t, 2, *meta.NextPage)
	assert.Equal(t, 1, calls)
}
