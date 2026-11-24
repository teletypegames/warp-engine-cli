package cli

import (
	"strconv"
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// DownloadSerializer turns the flags of the download commands into the
// API's terms and the download log and its statistics into tables and
// records.
type DownloadSerializer struct{}

// DownloadSorts are the orders "download list --sort" takes.
var DownloadSorts = []string{"createdAt"}

// Filters are the filter flags of "download list".
func (s DownloadSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "file-path", Key: "filePath", Usage: "file path contains"},
		{Flag: "software", Key: "software", Usage: "software title contains"},
		{Flag: "ip-address", Key: "ipAddress", Usage: "IP address contains"},
		{Flag: "user-agent", Key: "userAgent", Usage: "User-Agent contains"},
		{Flag: "referer", Key: "referer", Usage: "Referer contains"},
		{Flag: "outlet", Key: "outlet", Usage: "the outlet it came through, exactly; none for the catalog API"},
		{Flag: "client-id", Key: "clientId", Usage: "Client-Id contains"},
		{Flag: "from", Key: "from", Usage: "downloaded at or after: 2026-11-01T12:00:00Z, or 2026-11-01 for its midnight UTC"},
		{Flag: "to", Key: "to", Usage: "downloaded before: an ISO 8601 time, or a date for its midnight UTC"},
	}
}

// FlagsToFilter reads the filter flags of "download list". A bare date
// given to --from or --to (2026-11-01) becomes its midnight UTC
// (2026-11-01T00:00:00Z), since the server takes only full times.
func (s DownloadSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	filter := FlagsToFilter(flags, s.Filters())
	for _, key := range []string{"from", "to"} {
		if day, err := time.Parse(time.DateOnly, filter[key]); err == nil {
			filter[key] = day.Format(time.RFC3339)
		}
	}
	return filter
}

// ModelToTable renders downloads as a table: ID, time, file, software,
// release, IP address.
func (s DownloadSerializer) ModelToTable(downloads []domain.DownloadModel) Table {
	table := Table{Headers: []string{"ID", "Time", "File", "Software", "Release", "IP"}}
	for _, download := range downloads {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(download.ID), Text(download.CreatedAt), Truncate(download.FilePath, 60),
			Number(download.SoftwareID), Number(download.ReleaseID), Text(download.IPAddress),
		})
	}
	return table
}

// ModelToRecord renders one download with every field.
func (s DownloadSerializer) ModelToRecord(download domain.DownloadModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(download.ID)},
		{Name: "Time", Value: Text(download.CreatedAt)},
		{Name: "File", Value: download.FilePath},
		{Name: "Software", Value: Number(download.SoftwareID)},
		{Name: "Release", Value: Number(download.ReleaseID)},
		{Name: "IP address", Value: Text(download.IPAddress)},
		{Name: "User-Agent", Value: Text(download.UserAgent)},
		{Name: "Referer", Value: Text(download.Referer)},
		{Name: "Outlet", Value: Text(download.Outlet)},
		{Name: "Client-Id", Value: Text(download.ClientID)},
	}
}

// SummaryToRecord renders the download counts: total, today, this week,
// this month.
func (s DownloadSerializer) SummaryToRecord(summary domain.DownloadSummaryModel) Record {
	return Record{
		{Name: "Total", Value: strconv.Itoa(summary.Total)},
		{Name: "Today", Value: strconv.Itoa(summary.Today)},
		{Name: "This week", Value: strconv.Itoa(summary.Week)},
		{Name: "This month", Value: strconv.Itoa(summary.Month)},
	}
}

// FilesToTable renders the top files as a table: file, downloads, last
// download. The file path is the first column, the one -o id prints.
func (s DownloadSerializer) FilesToTable(files []domain.FileDownloadsModel) Table {
	table := Table{Headers: []string{"File", "Downloads", "Last"}}
	for _, file := range files {
		table.Rows = append(table.Rows, []string{file.FilePath, strconv.Itoa(file.Downloads), Text(file.LastDownloadAt)})
	}
	return table
}

// SoftwaresToTable renders the top softwares as a table: ID, name, title,
// downloads.
func (s DownloadSerializer) SoftwaresToTable(softwares []domain.SoftwareDownloadsModel) Table {
	table := Table{Headers: []string{"ID", "Name", "Title", "Downloads"}}
	for _, software := range softwares {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(software.SoftwareID), software.Name, Truncate(software.Title, 40), strconv.Itoa(software.Downloads),
		})
	}
	return table
}

// StatsToRecord renders the downloads of one software: the total, then a
// line per release ("Release 1.2.0 (#5)") followed by an indented line per
// file with its count and last download.
func (s DownloadSerializer) StatsToRecord(stats domain.SoftwareDownloadStatsModel) Record {
	record := Record{{Name: "Total", Value: strconv.Itoa(stats.Total)}}
	for _, release := range stats.Releases {
		record = append(record, Field{
			Name:  "Release " + release.Version + " (#" + strconv.Itoa(release.ReleaseID) + ")",
			Value: strconv.Itoa(release.Downloads),
		})
		for _, file := range release.Files {
			record = append(record, Field{
				Name:  "  " + file.FilePath,
				Value: strconv.Itoa(file.Downloads) + " (last " + Text(file.LastDownloadAt) + ")",
			})
		}
	}
	return record
}
