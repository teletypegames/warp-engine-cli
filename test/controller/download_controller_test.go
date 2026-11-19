package controller_test

import (
	"encoding/json"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestDownloadController_Index_PassesFiltersAndPaging(t *testing.T) {
	var got_filter domain.Filter
	var got_query domain.ListQuery
	software_id := 7
	svc := &mocksvc.MockDownloadService{
		AllFunc: func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.DownloadModel], error) {
			got_filter, got_query = filter, query
			return &domain.PageModel[domain.DownloadModel]{
				Items: []domain.DownloadModel{{ID: 99, FilePath: "pong/1.0.0/pong.zip", SoftwareID: &software_id, IPAddress: text("10.0.0.1")}},
				Page:  1, TotalPages: 3, Total: 60,
			}, nil
		},
	}

	stdout, stderr, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: svc}},
		"downloads", "list", "--software", "pong", "--file-path", ".zip", "--from", "2026-11-01", "--to", "2026-11-02",
		"--ip-address", "10.0", "--user-agent", "curl", "--referer", "example", "--per-page", "20", "--sort", "createdAt")

	assert.NoError(t, err)
	assert.Equal(t, domain.Filter{
		"software": "pong", "filePath": ".zip", "from": "2026-11-01T00:00:00Z", "to": "2026-11-02T00:00:00Z",
		"ipAddress": "10.0", "userAgent": "curl", "referer": "example",
	}, got_filter)
	assert.Equal(t, domain.ListQuery{PerPage: 20, Sort: "createdAt"}, got_query)
	assert.Contains(t, stdout, "pong/1.0.0/pong.zip")
	assert.Contains(t, stdout, "10.0.0.1")
	assert.Contains(t, stderr, "Page 1 of 3")
}

func TestDownloadController_Show_ParsesTheID(t *testing.T) {
	var got_id int
	svc := &mocksvc.MockDownloadService{ShowFunc: func(download_id int) (*domain.DownloadModel, error) {
		got_id = download_id
		return &domain.DownloadModel{ID: download_id, FilePath: "a.zip", UserAgent: text("curl/8")}, nil
	}}

	stdout, _, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: svc}}, "download", "show", "99")

	assert.NoError(t, err)
	assert.Equal(t, 99, got_id)
	assert.Contains(t, stdout, "curl/8")
}

func TestDownloadController_Summary_PrintsTheCounts(t *testing.T) {
	svc := &mocksvc.MockDownloadService{SummaryFunc: func() (*domain.DownloadSummaryModel, error) {
		return &domain.DownloadSummaryModel{Total: 1234, Today: 5, Week: 40, Month: 321}, nil
	}}

	stdout, _, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: svc}}, "download", "summary")

	assert.NoError(t, err)
	assert.Contains(t, stdout, "1234")
	assert.Contains(t, stdout, "This month")
	assert.Contains(t, stdout, "321")
}

func TestDownloadController_TopFiles_PassesTheLimit(t *testing.T) {
	var got_limit int
	svc := &mocksvc.MockDownloadService{TopFilesFunc: func(limit int) ([]domain.FileDownloadsModel, error) {
		got_limit = limit
		return []domain.FileDownloadsModel{{FilePath: "pong/1.0.0/pong.zip", Downloads: 17}}, nil
	}}

	stdout, _, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: svc}}, "download", "top-files", "--limit", "20")

	assert.NoError(t, err)
	assert.Equal(t, 20, got_limit)
	assert.Contains(t, stdout, "pong/1.0.0/pong.zip")
	assert.Contains(t, stdout, "17")
}

func TestDownloadController_TopSoftwares_WithoutLimitLeavesItToTheServer(t *testing.T) {
	got_limit := -1
	svc := &mocksvc.MockDownloadService{TopSoftwaresFunc: func(limit int) ([]domain.SoftwareDownloadsModel, error) {
		got_limit = limit
		return []domain.SoftwareDownloadsModel{{SoftwareID: 7, Name: "pong", Title: "Pong", Downloads: 30}}, nil
	}}

	stdout, _, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: svc}}, "download", "top-softwares", "-o", "id")

	assert.NoError(t, err)
	assert.Equal(t, 0, got_limit)
	assert.Equal(t, "7\n", stdout)
}

func TestDownloadController_Stats_PassesTheReference(t *testing.T) {
	var got_reference string
	svc := &mocksvc.MockDownloadService{StatsFunc: func(reference string) (*domain.SoftwareDownloadStatsModel, error) {
		got_reference = reference
		return &domain.SoftwareDownloadStatsModel{Total: 30, Releases: []domain.ReleaseDownloadsModel{
			{ReleaseID: 5, Version: "1.0.0", Downloads: 30, Files: []domain.FileDownloadsModel{{FilePath: "pong.zip", Downloads: 30}}},
		}}, nil
	}}

	stdout, _, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: svc}}, "download", "stats", "pong")

	assert.NoError(t, err)
	assert.Equal(t, "pong", got_reference)
	assert.Contains(t, stdout, "Release 1.0.0 (#5)")
	assert.Contains(t, stdout, "pong.zip")
}

func TestDownloadController_Stats_JSONIsTheModel(t *testing.T) {
	svc := &mocksvc.MockDownloadService{StatsFunc: func(reference string) (*domain.SoftwareDownloadStatsModel, error) {
		return &domain.SoftwareDownloadStatsModel{Total: 3, Releases: []domain.ReleaseDownloadsModel{}}, nil
	}}

	stdout, _, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: svc}}, "download", "stats", "7", "-o", "json")

	assert.NoError(t, err)
	var stats domain.SoftwareDownloadStatsModel
	assert.NoError(t, json.Unmarshal([]byte(stdout), &stats))
	assert.Equal(t, 3, stats.Total)
}

func TestDownloadController_Stats_NeedsASoftware(t *testing.T) {
	_, _, err := execute(cli.Router{DownloadController: cli.DownloadController{DownloadService: &mocksvc.MockDownloadService{}}}, "download", "stats")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}
