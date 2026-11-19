package mocksvc

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockDownloadService mocks [domain.DownloadServiceInterface].
type MockDownloadService struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.DownloadModel], error)
	// ShowFunc replaces Show.
	ShowFunc func(download_id int) (*domain.DownloadModel, error)
	// SummaryFunc replaces Summary.
	SummaryFunc func() (*domain.DownloadSummaryModel, error)
	// TopFilesFunc replaces TopFiles.
	TopFilesFunc func(limit int) ([]domain.FileDownloadsModel, error)
	// TopSoftwaresFunc replaces TopSoftwares.
	TopSoftwaresFunc func(limit int) ([]domain.SoftwareDownloadsModel, error)
	// StatsFunc replaces Stats.
	StatsFunc func(reference string) (*domain.SoftwareDownloadStatsModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockDownloadService) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.DownloadModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.DownloadModel]{}, nil
}

// Show calls ShowFunc, or returns an empty download.
func (m *MockDownloadService) Show(download_id int) (*domain.DownloadModel, error) {
	if m.ShowFunc != nil {
		return m.ShowFunc(download_id)
	}
	return &domain.DownloadModel{}, nil
}

// Summary calls SummaryFunc, or returns an empty summary.
func (m *MockDownloadService) Summary() (*domain.DownloadSummaryModel, error) {
	if m.SummaryFunc != nil {
		return m.SummaryFunc()
	}
	return &domain.DownloadSummaryModel{}, nil
}

// TopFiles calls TopFilesFunc, or returns no files.
func (m *MockDownloadService) TopFiles(limit int) ([]domain.FileDownloadsModel, error) {
	if m.TopFilesFunc != nil {
		return m.TopFilesFunc(limit)
	}
	return nil, nil
}

// TopSoftwares calls TopSoftwaresFunc, or returns no softwares.
func (m *MockDownloadService) TopSoftwares(limit int) ([]domain.SoftwareDownloadsModel, error) {
	if m.TopSoftwaresFunc != nil {
		return m.TopSoftwaresFunc(limit)
	}
	return nil, nil
}

// Stats calls StatsFunc, or returns empty statistics.
func (m *MockDownloadService) Stats(reference string) (*domain.SoftwareDownloadStatsModel, error) {
	if m.StatsFunc != nil {
		return m.StatsFunc(reference)
	}
	return &domain.SoftwareDownloadStatsModel{}, nil
}
