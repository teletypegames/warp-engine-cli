package domain

// SoftwareDownloadStatsModel is the downloads of one software, broken down
// by release and by file.
type SoftwareDownloadStatsModel struct {
	// Total is every download of the software's files.
	Total int `json:"total"`
	// Releases are the downloads of each release.
	Releases []ReleaseDownloadsModel `json:"releases"`
}

// ReleaseDownloadsModel is the downloads of one release in a
// [SoftwareDownloadStatsModel].
type ReleaseDownloadsModel struct {
	// ReleaseID is the release.
	ReleaseID int `json:"releaseId"`
	// Version is the release's version.
	Version string `json:"version"`
	// Downloads is every download of the release's files.
	Downloads int `json:"downloads"`
	// Files are the downloads of each of its files.
	Files []FileDownloadsModel `json:"files"`
}
