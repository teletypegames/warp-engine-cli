package domain

import (
	"net/url"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// DownloadRepositoryInterface is the Admin API's download log and its
// statistics.
type DownloadRepositoryInterface interface {
	// All lists the downloads matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[DownloadModel], error)
	// FindByID returns the download download_id.
	FindByID(download_id int) (*DownloadModel, error)
	// Summary returns the downloads in total, today, this week and this month.
	Summary() (*DownloadSummaryModel, error)
	// TopFiles returns the limit most downloaded files; 0 leaves the number
	// to the server.
	TopFiles(limit int) ([]FileDownloadsModel, error)
	// TopSoftwares returns the limit most downloaded softwares; 0 leaves the
	// number to the server.
	TopSoftwares(limit int) ([]SoftwareDownloadsModel, error)
	// Stats returns the downloads of the software software_id by release
	// and by file.
	Stats(software_id int) (*SoftwareDownloadStatsModel, error)
}

// DownloadRepository is the [DownloadRepositoryInterface] over the Admin API.
type DownloadRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /downloads.
func (r DownloadRepository) All(filter Filter, query ListQuery) (*PageModel[DownloadModel], error) {
	return list[DownloadModel](r.Client, "/downloads", filter, query)
}

// FindByID gets GET /downloads/{id}.
func (r DownloadRepository) FindByID(download_id int) (*DownloadModel, error) {
	return api.Get[DownloadModel](r.Client, "/downloads/"+strconv.Itoa(download_id), nil)
}

// Summary gets GET /downloads/summary.
func (r DownloadRepository) Summary() (*DownloadSummaryModel, error) {
	return api.Get[DownloadSummaryModel](r.Client, "/downloads/summary", nil)
}

// TopFiles gets GET /downloads/top_files?limit=….
func (r DownloadRepository) TopFiles(limit int) ([]FileDownloadsModel, error) {
	items, _, err := api.List[FileDownloadsModel](r.Client, "/downloads/top_files", topLimitQuery(limit), false)
	return items, err
}

// TopSoftwares gets GET /downloads/top_softwares?limit=….
func (r DownloadRepository) TopSoftwares(limit int) ([]SoftwareDownloadsModel, error) {
	items, _, err := api.List[SoftwareDownloadsModel](r.Client, "/downloads/top_softwares", topLimitQuery(limit), false)
	return items, err
}

// Stats gets GET /softwares/{id}/download_stats.
func (r DownloadRepository) Stats(software_id int) (*SoftwareDownloadStatsModel, error) {
	return api.Get[SoftwareDownloadStatsModel](r.Client, "/softwares/"+strconv.Itoa(software_id)+"/download_stats", nil)
}

// topLimitQuery returns the query of a top list: limit when it is set, nothing
// otherwise.
func topLimitQuery(limit int) url.Values {
	if limit < 1 {
		return nil
	}
	return url.Values{"limit": {strconv.Itoa(limit)}}
}
