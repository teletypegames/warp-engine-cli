package service_test

import (
	"net/http"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestDownloadService_Stats_NameReferenceIsLookedUp(t *testing.T) {
	var got_id int
	downloads := &mockrepo.MockDownloadRepository{StatsFunc: func(software_id int) (*domain.SoftwareDownloadStatsModel, error) {
		got_id = software_id
		return &domain.SoftwareDownloadStatsModel{Total: 12}, nil
	}}
	softwares := &mockrepo.MockSoftwareRepository{FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
		assert.Equal(t, "pong", name)
		return &domain.SoftwareModel{ID: 3, Name: name}, nil
	}}
	svc := domain.DownloadService{DownloadRepository: downloads, SoftwareRepository: softwares}

	stats, err := svc.Stats("pong")

	assert.NoError(t, err)
	assert.Equal(t, 3, got_id)
	assert.Equal(t, 12, stats.Total)
}

func TestDownloadService_Stats_NumericReferenceIsTheID(t *testing.T) {
	var got_id int
	downloads := &mockrepo.MockDownloadRepository{StatsFunc: func(software_id int) (*domain.SoftwareDownloadStatsModel, error) {
		got_id = software_id
		return &domain.SoftwareDownloadStatsModel{}, nil
	}}
	softwares := &mockrepo.MockSoftwareRepository{FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
		t.Fatal("a numeric reference must not be looked up by name")
		return nil, nil
	}}
	svc := domain.DownloadService{DownloadRepository: downloads, SoftwareRepository: softwares}

	_, err := svc.Stats("42")

	assert.NoError(t, err)
	assert.Equal(t, 42, got_id)
}

func TestDownloadService_Stats_UnknownNameStopsBeforeTheRequest(t *testing.T) {
	downloads := &mockrepo.MockDownloadRepository{StatsFunc: func(software_id int) (*domain.SoftwareDownloadStatsModel, error) {
		t.Fatal("must not reach the downloads")
		return nil, nil
	}}
	softwares := &mockrepo.MockSoftwareRepository{FindByNameFunc: func(name string) (*domain.SoftwareModel, error) {
		return nil, &api.Error{Status: http.StatusNotFound, Message: "no software named \"nope\""}
	}}
	svc := domain.DownloadService{DownloadRepository: downloads, SoftwareRepository: softwares}

	_, err := svc.Stats("nope")

	assert.Equal(t, http.StatusNotFound, api.StatusOf(err))
}

func TestDownloadService_TopFiles_DelegatesTheLimit(t *testing.T) {
	var got_limit int
	downloads := &mockrepo.MockDownloadRepository{TopFilesFunc: func(limit int) ([]domain.FileDownloadsModel, error) {
		got_limit = limit
		return []domain.FileDownloadsModel{{FilePath: "a.zip", Downloads: 3}}, nil
	}}
	svc := domain.DownloadService{DownloadRepository: downloads}

	files, err := svc.TopFiles(20)

	assert.NoError(t, err)
	assert.Equal(t, 20, got_limit)
	assert.Len(t, files, 1)
}

func TestDownloadService_TopSoftwares_NegativeLimitIsAnError(t *testing.T) {
	downloads := &mockrepo.MockDownloadRepository{TopSoftwaresFunc: func(limit int) ([]domain.SoftwareDownloadsModel, error) {
		t.Fatal("must not reach the repository")
		return nil, nil
	}}
	svc := domain.DownloadService{DownloadRepository: downloads}

	_, err := svc.TopSoftwares(-1)

	assert.ErrorContains(t, err, "positive")
}

func TestDownloadService_Summary_Delegates(t *testing.T) {
	downloads := &mockrepo.MockDownloadRepository{
		SummaryFunc: func() (*domain.DownloadSummaryModel, error) { return &domain.DownloadSummaryModel{Total: 5}, nil },
	}
	svc := domain.DownloadService{DownloadRepository: downloads}

	summary, err := svc.Summary()

	assert.NoError(t, err)
	assert.Equal(t, 5, summary.Total)
}

func TestDownloadService_Show_DelegatesTheID(t *testing.T) {
	var got_id int
	downloads := &mockrepo.MockDownloadRepository{FindByIDFunc: func(download_id int) (*domain.DownloadModel, error) {
		got_id = download_id
		return &domain.DownloadModel{ID: download_id}, nil
	}}
	svc := domain.DownloadService{DownloadRepository: downloads}

	_, err := svc.Show(8)

	assert.NoError(t, err)
	assert.Equal(t, 8, got_id)
}
