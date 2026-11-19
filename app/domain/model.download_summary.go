package domain

// DownloadSummaryModel is the number of downloads in total and in the
// current day, week and month.
type DownloadSummaryModel struct {
	// Total is every download ever logged.
	Total int `json:"total"`
	// Today is the downloads of the current day.
	Today int `json:"today"`
	// Week is the downloads of the current week.
	Week int `json:"week"`
	// Month is the downloads of the current month.
	Month int `json:"month"`
}
