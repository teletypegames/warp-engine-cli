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

func TestDownloadRepository_All_SendsFiltersAndReadsMeta(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/downloads", r.URL.Path)
		assert.Equal(t, "pong", r.URL.Query().Get("software"))
		assert.Equal(t, "2026-11-01", r.URL.Query().Get("from"))
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		io.WriteString(w, `{"data":[{"id":99,"filePath":"pong.zip","releaseId":5,"softwareId":7,"ipAddress":"10.0.0.1","userAgent":null,"referer":null,"createdAt":"2026-11-19T10:00:00Z","updatedAt":null}],"meta":{"page":2,"perPage":1,"total":3,"totalPages":3,"nextPage":3}}`)
	}))
	defer server.Close()
	repo := domain.DownloadRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"software": "pong", "from": "2026-11-01"}, domain.ListQuery{Page: 2})

	assert.NoError(t, err)
	assert.Equal(t, 99, page.Items[0].ID)
	assert.Equal(t, 7, *page.Items[0].SoftwareID)
	assert.Equal(t, 3, page.TotalPages)
}

func TestDownloadRepository_FindByID_GetsThePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/downloads/99", r.URL.Path)
		io.WriteString(w, `{"data":{"id":99,"filePath":"pong.zip"}}`)
	}))
	defer server.Close()
	repo := domain.DownloadRepository{Client: api.New(server.URL, "", "test")}

	download, err := repo.FindByID(99)

	assert.NoError(t, err)
	assert.Equal(t, "pong.zip", download.FilePath)
}

func TestDownloadRepository_Summary_ReadsTheCounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/downloads/summary", r.URL.Path)
		io.WriteString(w, `{"data":{"total":100,"today":1,"week":7,"month":30}}`)
	}))
	defer server.Close()
	repo := domain.DownloadRepository{Client: api.New(server.URL, "", "test")}

	summary, err := repo.Summary()

	assert.NoError(t, err)
	assert.Equal(t, domain.DownloadSummaryModel{Total: 100, Today: 1, Week: 7, Month: 30}, *summary)
}

func TestDownloadRepository_TopFiles_SendsTheLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/downloads/top_files", r.URL.Path)
		assert.Equal(t, "20", r.URL.Query().Get("limit"))
		io.WriteString(w, `{"data":[{"filePath":"pong.zip","downloads":17,"lastDownloadAt":"2026-11-19T10:00:00Z"}]}`)
	}))
	defer server.Close()
	repo := domain.DownloadRepository{Client: api.New(server.URL, "", "test")}

	files, err := repo.TopFiles(20)

	assert.NoError(t, err)
	assert.Equal(t, 17, files[0].Downloads)
}

func TestDownloadRepository_TopSoftwares_WithoutLimitSendsNone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/downloads/top_softwares", r.URL.Path)
		assert.False(t, r.URL.Query().Has("limit"))
		io.WriteString(w, `{"data":[{"softwareId":7,"name":"pong","title":"Pong","downloads":30}]}`)
	}))
	defer server.Close()
	repo := domain.DownloadRepository{Client: api.New(server.URL, "", "test")}

	softwares, err := repo.TopSoftwares(0)

	assert.NoError(t, err)
	assert.Equal(t, "pong", softwares[0].Name)
}

func TestDownloadRepository_Stats_GetsTheSoftwarePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/softwares/7/download_stats", r.URL.Path)
		io.WriteString(w, `{"data":{"total":30,"releases":[{"releaseId":5,"version":"1.0.0","downloads":30,"files":[{"filePath":"pong.zip","downloads":30,"lastDownloadAt":null}]}]}}`)
	}))
	defer server.Close()
	repo := domain.DownloadRepository{Client: api.New(server.URL, "", "test")}

	stats, err := repo.Stats(7)

	assert.NoError(t, err)
	assert.Equal(t, 30, stats.Total)
	assert.Equal(t, "pong.zip", stats.Releases[0].Files[0].FilePath)
}
