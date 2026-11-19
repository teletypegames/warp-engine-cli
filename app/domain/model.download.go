package domain

// DownloadModel is one entry of the download log: a file served once.
type DownloadModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// FilePath is the path of the file that was served.
	FilePath string `json:"filePath"`
	// ReleaseID is the release the file belongs to, or nil.
	ReleaseID *int `json:"releaseId"`
	// SoftwareID is the software the file belongs to, or nil.
	SoftwareID *int `json:"softwareId"`
	// IPAddress is the address of the client, or nil.
	IPAddress *string `json:"ipAddress"`
	// UserAgent is the client's User-Agent header, or nil.
	UserAgent *string `json:"userAgent"`
	// Referer is the client's Referer header, or nil.
	Referer *string `json:"referer"`
	// CreatedAt is when it was downloaded (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when the entry last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
