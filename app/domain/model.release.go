package domain

// Channels are the release channels, as the API spells them.
var Channels = []string{"stable", "beta", "dev"}

// ReleaseModel is one version of a software, with its files.
type ReleaseModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// SoftwareID is the software it is a version of.
	SoftwareID int `json:"softwareId"`
	// Version is the version, "1.2.0".
	Version string `json:"version"`
	// Channel is one of [Channels].
	Channel string `json:"channel"`
	// CommitSha is the commit it was built from, or nil.
	CommitSha *string `json:"commitSha"`
	// Assets are its files.
	Assets []ReleaseAssetModel `json:"assets,omitempty"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
