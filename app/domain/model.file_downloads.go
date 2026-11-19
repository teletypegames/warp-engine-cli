package domain

// FileDownloadsModel is how often one file was downloaded.
type FileDownloadsModel struct {
	// FilePath is the path of the file.
	FilePath string `json:"filePath"`
	// Downloads is how many times it was served.
	Downloads int `json:"downloads"`
	// LastDownloadAt is when it was last served (RFC 3339), or nil.
	LastDownloadAt *string `json:"lastDownloadAt"`
}
