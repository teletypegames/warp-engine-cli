package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestDownloadSerializer_FlagsToFilter_MapsTheKeys(t *testing.T) {
	serializer := cli.DownloadSerializer{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFilters(flags, serializer.Filters())
	assert.NoError(t, flags.Parse([]string{"--file-path", ".prg", "--ip-address", "10."}))

	assert.Equal(t, domain.Filter{"filePath": ".prg", "ipAddress": "10."}, serializer.FlagsToFilter(flags))
}

func TestDownloadSerializer_FlagsToFilter_DateBecomesMidnightUTC(t *testing.T) {
	serializer := cli.DownloadSerializer{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFilters(flags, serializer.Filters())
	assert.NoError(t, flags.Parse([]string{"--from", "2026-11-01", "--to", "2026-11-02T06:00:00+01:00"}))

	assert.Equal(t, domain.Filter{"from": "2026-11-01T00:00:00Z", "to": "2026-11-02T06:00:00+01:00"}, serializer.FlagsToFilter(flags))
}

func TestDownloadSerializer_ModelToTable_MapsFields(t *testing.T) {
	software_id, created := 7, "2026-11-19T10:00:00Z"
	table := cli.DownloadSerializer{}.ModelToTable([]domain.DownloadModel{
		{ID: 99, FilePath: "pong.zip", SoftwareID: &software_id, CreatedAt: &created},
	})

	assert.Equal(t, []string{"99", "2026-11-19T10:00:00Z", "pong.zip", "7", "-", "-"}, table.Rows[0])
}

func TestDownloadSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.DownloadSerializer{}.ModelToTable([]domain.DownloadModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestDownloadSerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.DownloadSerializer{}.ModelToRecord(domain.DownloadModel{ID: 99})

	assert.Equal(t, "ID", record[0].Name)
	assert.Equal(t, "99", record[0].Value)
}

func TestDownloadSerializer_SummaryToRecord_TotalFirst(t *testing.T) {
	record := cli.DownloadSerializer{}.SummaryToRecord(domain.DownloadSummaryModel{Total: 10, Today: 1, Week: 2, Month: 3})

	assert.Equal(t, cli.Field{Name: "Total", Value: "10"}, record[0])
	assert.Equal(t, cli.Field{Name: "This month", Value: "3"}, record[3])
}

func TestDownloadSerializer_FilesToTable_FilledAndEmpty(t *testing.T) {
	table := cli.DownloadSerializer{}.FilesToTable([]domain.FileDownloadsModel{{FilePath: "a.zip", Downloads: 4}})
	assert.Equal(t, []string{"a.zip", "4", "-"}, table.Rows[0])

	empty := cli.DownloadSerializer{}.FilesToTable(nil)
	assert.Empty(t, empty.Rows)
	assert.NotEmpty(t, empty.Headers)
}

func TestDownloadSerializer_SoftwaresToTable_FilledAndEmpty(t *testing.T) {
	table := cli.DownloadSerializer{}.SoftwaresToTable([]domain.SoftwareDownloadsModel{{SoftwareID: 7, Name: "pong", Title: "Pong", Downloads: 30}})
	assert.Equal(t, []string{"7", "pong", "Pong", "30"}, table.Rows[0])

	empty := cli.DownloadSerializer{}.SoftwaresToTable(nil)
	assert.Empty(t, empty.Rows)
	assert.NotEmpty(t, empty.Headers)
}

func TestDownloadSerializer_StatsToRecord_ReleasesThenFiles(t *testing.T) {
	last := "2026-11-18T09:00:00Z"
	record := cli.DownloadSerializer{}.StatsToRecord(domain.SoftwareDownloadStatsModel{Total: 30, Releases: []domain.ReleaseDownloadsModel{
		{ReleaseID: 5, Version: "1.0.0", Downloads: 30, Files: []domain.FileDownloadsModel{{FilePath: "pong.zip", Downloads: 30, LastDownloadAt: &last}}},
	}})

	assert.Equal(t, cli.Record{
		{Name: "Total", Value: "30"},
		{Name: "Release 1.0.0 (#5)", Value: "30"},
		{Name: "  pong.zip", Value: "30 (last 2026-11-18T09:00:00Z)"},
	}, record)
}

func TestDownloadSerializer_FlagsToFilter_TakesTheOutletAndClientID(t *testing.T) {
	serializer := cli.DownloadSerializer{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFilters(flags, serializer.Filters())
	assert.NoError(t, flags.Parse([]string{"--outlet", "assembly64", "--client-id", "Ultimate"}))

	assert.Equal(t, domain.Filter{"outlet": "assembly64", "clientId": "Ultimate"}, serializer.FlagsToFilter(flags))
}
