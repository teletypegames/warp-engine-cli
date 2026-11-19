package domain

// AssetKinds are the kinds of release asset, as the API spells them: what
// the file is and which machine it runs on.
var AssetKinds = []string{"cartridge", "source", "html", "docs", "win_x86", "win_x64", "linux_x86", "linux_x64", "linux_arm64", "mac_x64", "mac_arm64", "mac_universal"}

// ReleaseAssetModel is one file of a release in the artifact store.
type ReleaseAssetModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// ReleaseID is the release it belongs to.
	ReleaseID int `json:"releaseId"`
	// Kind is one of [AssetKinds].
	Kind string `json:"kind"`
	// Path is the file, relative to the artifact store.
	Path string `json:"path"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
