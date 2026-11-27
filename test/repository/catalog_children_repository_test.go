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

// catalogServer starts a test server answering every request with handler and
// returns a client for it; the server stops with the test.
func catalogServer(t *testing.T, handler http.HandlerFunc) *api.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return api.New(server.URL, "secret", "test")
}

// catalogBody decodes the JSON body of r.
func catalogBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	body := map[string]any{}
	assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
	return body
}

func TestReleaseRepository_All_SendsFiltersAndSort(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/admin/v2/softwares/3/releases", r.URL.Path)
		assert.Equal(t, "beta", r.URL.Query().Get("channel"))
		assert.Equal(t, "-version", r.URL.Query().Get("sort"))
		io.WriteString(w, `{"data":[{"id":42,"softwareId":3,"version":"1.0-beta","channel":"beta","commitSha":null,"assets":[]}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`)
	})

	page, err := domain.ReleaseRepository{Client: client}.All(3, domain.Filter{"channel": "beta"}, domain.ListQuery{Sort: "-version"})

	assert.NoError(t, err)
	assert.Equal(t, 42, page.Items[0].ID)
	assert.Equal(t, 1, page.TotalPages)
}

func TestReleaseRepository_Create_PostsTheVersion(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/softwares/3/releases", r.URL.Path)
		assert.Equal(t, map[string]any{"version": "dev-abc"}, catalogBody(t, r))
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"id":43,"version":"dev-abc","channel":"dev"}}`)
	})

	release, err := domain.ReleaseRepository{Client: client}.Create(3, domain.Attributes{"version": "dev-abc"})

	assert.NoError(t, err)
	assert.Equal(t, "dev", release.Channel)
}

func TestReleaseRepository_Destroy_SendsDelete(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/admin/v2/releases/43", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	assert.NoError(t, domain.ReleaseRepository{Client: client}.Destroy(43))
}

func TestReleaseAssetRepository_All_ReadsTheUnpagedList(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/releases/42/assets", r.URL.Path)
		io.WriteString(w, `{"data":[{"id":7,"releaseId":42,"kind":"html","path":"/p"}]}`)
	})

	assets, err := domain.ReleaseAssetRepository{Client: client}.All(42)

	assert.NoError(t, err)
	assert.Equal(t, "html", assets[0].Kind)
}

func TestReleaseAssetRepository_Update_PatchesTheAsset(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/admin/v2/assets/7", r.URL.Path)
		assert.Equal(t, map[string]any{"path": "/q"}, catalogBody(t, r))
		io.WriteString(w, `{"data":{"id":7,"path":"/q"}}`)
	})

	asset, err := domain.ReleaseAssetRepository{Client: client}.Update(7, domain.Attributes{"path": "/q"})

	assert.NoError(t, err)
	assert.Equal(t, "/q", asset.Path)
}

func TestSoftwareLinkRepository_Create_PostsLabelAndURL(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/softwares/3/links", r.URL.Path)
		assert.Equal(t, map[string]any{"label": "Source", "url": "https://x"}, catalogBody(t, r))
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"id":2,"softwareId":3,"label":"Source","url":"https://x"}}`)
	})

	link, err := domain.SoftwareLinkRepository{Client: client}.Create(3, domain.Attributes{"label": "Source", "url": "https://x"})

	assert.NoError(t, err)
	assert.Equal(t, 2, link.ID)
}

func TestSoftwareLinkRepository_FindByID_GetsTheLink(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/links/2", r.URL.Path)
		io.WriteString(w, `{"data":{"id":2,"label":"Source"}}`)
	})

	link, err := domain.SoftwareLinkRepository{Client: client}.FindByID(2)

	assert.NoError(t, err)
	assert.Equal(t, "Source", link.Label)
}

func TestSoftwareImageRepository_Upload_SendsMultipart(t *testing.T) {
	file_path := filepath.Join(t.TempDir(), "cover.png")
	assert.NoError(t, os.WriteFile(file_path, []byte("PNGDATA"), 0o600))
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/admin/v2/softwares/3/images/upload", r.URL.Path)
		assert.NoError(t, r.ParseMultipartForm(1<<20))
		assert.Equal(t, "true", r.FormValue("isDefault"))
		assert.Equal(t, "2", r.FormValue("position"))
		file, header, err := r.FormFile("file")
		assert.NoError(t, err)
		data, _ := io.ReadAll(file)
		assert.Equal(t, "cover.png", header.Filename)
		assert.Equal(t, "PNGDATA", string(data))
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"id":31,"softwareId":3,"imageId":118,"isDefault":true,"position":2,"origin":"local","url":null}}`)
	})

	image, err := domain.SoftwareImageRepository{Client: client}.Upload(3, file_path, domain.Attributes{"isDefault": true, "position": 2})

	assert.NoError(t, err)
	assert.Equal(t, 118, image.ImageID)
	assert.True(t, image.IsDefault)
}

func TestSoftwareImageRepository_Create_PostsTheImageID(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/softwares/3/images", r.URL.Path)
		assert.Equal(t, map[string]any{"imageId": float64(118)}, catalogBody(t, r))
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"data":{"id":31,"imageId":118}}`)
	})

	image, err := domain.SoftwareImageRepository{Client: client}.Create(3, domain.Attributes{"imageId": 118})

	assert.NoError(t, err)
	assert.Equal(t, 31, image.ID)
}

func TestSoftwareImageRepository_Destroy_SendsDelete(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/admin/v2/software_images/31", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	assert.NoError(t, domain.SoftwareImageRepository{Client: client}.Destroy(31))
}

func TestExtraRepository_Update_PatchesTheValue(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/admin/v2/softwares/3/extras/seo_metadata", r.URL.Path)
		assert.Equal(t, map[string]any{"value": map[string]any{"title": "Pong"}}, catalogBody(t, r))
		io.WriteString(w, `{"data":{"association":"seo_metadata","multiple":false,"value":{"title":"Pong"}}}`)
	})

	extra, err := domain.ExtraRepository{Client: client}.Update(3, "seo_metadata", domain.Attributes{"value": map[string]any{"title": "Pong"}})

	assert.NoError(t, err)
	assert.Equal(t, "seo_metadata", extra.Association)
	assert.False(t, extra.Multiple)
}

func TestExtraRepository_Find_MissingResourceIs404(t *testing.T) {
	client := catalogServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/softwares/3/extras/seo_metadata", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"error":"Not found"}`)
	})

	_, err := domain.ExtraRepository{Client: client}.Find(3, "seo_metadata")

	assert.Equal(t, http.StatusNotFound, api.StatusOf(err))
}
