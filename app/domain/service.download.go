package domain

import "fmt"

// DownloadServiceInterface reads the download log and its statistics.
type DownloadServiceInterface interface {
	// All lists the downloads matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[DownloadModel], error)
	// Show returns the download download_id.
	Show(download_id int) (*DownloadModel, error)
	// Summary returns the downloads in total, today, this week and this month.
	Summary() (*DownloadSummaryModel, error)
	// TopFiles returns the limit most downloaded files (0: the server's
	// default; the server caps it at 100).
	TopFiles(limit int) ([]FileDownloadsModel, error)
	// TopSoftwares returns the limit most downloaded softwares (0: the
	// server's default; the server caps it at 100).
	TopSoftwares(limit int) ([]SoftwareDownloadsModel, error)
	// Stats returns the downloads of the software reference (an ID or an
	// exact name) by release and by file.
	Stats(reference string) (*SoftwareDownloadStatsModel, error)
}

// DownloadService is the [DownloadServiceInterface] over a
// [DownloadRepositoryInterface], with a [SoftwareRepositoryInterface] to
// resolve software references.
type DownloadService struct {
	// DownloadRepository talks to the downloads resource.
	DownloadRepository DownloadRepositoryInterface
	// SoftwareRepository looks softwares up by name.
	SoftwareRepository SoftwareRepositoryInterface
}

// All lists the downloads.
func (s DownloadService) All(filter Filter, query ListQuery) (*PageModel[DownloadModel], error) {
	return s.DownloadRepository.All(filter, query)
}

// Show gets the download download_id.
func (s DownloadService) Show(download_id int) (*DownloadModel, error) {
	return s.DownloadRepository.FindByID(download_id)
}

// Summary gets the download counts.
func (s DownloadService) Summary() (*DownloadSummaryModel, error) {
	return s.DownloadRepository.Summary()
}

// TopFiles checks that limit is not negative and gets the top files.
func (s DownloadService) TopFiles(limit int) ([]FileDownloadsModel, error) {
	if err := checkTopLimit(limit); err != nil {
		return nil, err
	}
	return s.DownloadRepository.TopFiles(limit)
}

// TopSoftwares checks that limit is not negative and gets the top softwares.
func (s DownloadService) TopSoftwares(limit int) ([]SoftwareDownloadsModel, error) {
	if err := checkTopLimit(limit); err != nil {
		return nil, err
	}
	return s.DownloadRepository.TopSoftwares(limit)
}

// Stats resolves reference and gets the software's download statistics.
func (s DownloadService) Stats(reference string) (*SoftwareDownloadStatsModel, error) {
	software_id, err := ResolveSoftwareID(s.SoftwareRepository, reference)
	if err != nil {
		return nil, err
	}
	return s.DownloadRepository.Stats(software_id)
}

// checkTopLimit returns an error for a negative top-list limit, which the
// server would refuse; 0 means none is sent.
func checkTopLimit(limit int) error {
	if limit < 0 {
		return fmt.Errorf("the limit is a positive number, not %d", limit)
	}
	return nil
}
