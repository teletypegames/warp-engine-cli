package domain

// SoftwareDownloadsModel is how often the files of one software were
// downloaded, all its releases together.
type SoftwareDownloadsModel struct {
	// SoftwareID is the software.
	SoftwareID int `json:"softwareId"`
	// Name is the software's internal name.
	Name string `json:"name"`
	// Title is the software's display name.
	Title string `json:"title"`
	// Downloads is how many of its files were served.
	Downloads int `json:"downloads"`
}
