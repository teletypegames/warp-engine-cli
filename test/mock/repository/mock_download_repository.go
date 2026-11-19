package mockrepo

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// MockDownloadRepository mocks [domain.DownloadRepositoryInterface].
type MockDownloadRepository struct {
	// AllFunc replaces All.
	AllFunc func(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.DownloadModel], error)
	// FindByIDFunc replaces FindByID.
	FindByIDFunc func(download_id int) (*domain.DownloadModel, error)
	// SummaryFunc replaces Summary.
	SummaryFunc func() (*domain.DownloadSummaryModel, error)
	// TopFilesFunc replaces TopFiles.
	TopFilesFunc func(limit int) ([]domain.FileDownloadsModel, error)
	// TopSoftwaresFunc replaces TopSoftwares.
	TopSoftwaresFunc func(limit int) ([]domain.SoftwareDownloadsModel, error)
	// StatsFunc replaces Stats.
	StatsFunc func(software_id int) (*domain.SoftwareDownloadStatsModel, error)
}

// All calls AllFunc, or returns an empty page.
func (m *MockDownloadRepository) All(filter domain.Filter, query domain.ListQuery) (*domain.PageModel[domain.DownloadModel], error) {
	if m.AllFunc != nil {
		return m.AllFunc(filter, query)
	}
	return &domain.PageModel[domain.DownloadModel]{}, nil
}

// FindByID calls FindByIDFunc, or returns an empty download.
func (m *MockDownloadRepository) FindByID(download_id int) (*domain.DownloadModel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(download_id)
	}
	return &domain.DownloadModel{}, nil
}

// Summary calls SummaryFunc, or returns an empty summary.
func (m *MockDownloadRepository) Summary() (*domain.DownloadSummaryModel, error) {
	if m.SummaryFunc != nil {
		return m.SummaryFunc()
	}
	return &domain.DownloadSummaryModel{}, nil
}

// TopFiles calls TopFilesFunc, or returns no files.
func (m *MockDownloadRepository) TopFiles(limit int) ([]domain.FileDownloadsModel, error) {
	if m.TopFilesFunc != nil {
		return m.TopFilesFunc(limit)
	}
	return nil, nil
}

// TopSoftwares calls TopSoftwaresFunc, or returns no softwares.
func (m *MockDownloadRepository) TopSoftwares(limit int) ([]domain.SoftwareDownloadsModel, error) {
	if m.TopSoftwaresFunc != nil {
		return m.TopSoftwaresFunc(limit)
	}
	return nil, nil
}

// Stats calls StatsFunc, or returns empty statistics.
func (m *MockDownloadRepository) Stats(software_id int) (*domain.SoftwareDownloadStatsModel, error) {
	if m.StatsFunc != nil {
		return m.StatsFunc(software_id)
	}
	return &domain.SoftwareDownloadStatsModel{}, nil
}
