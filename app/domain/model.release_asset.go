package domain

// ReleaseAssetModel is one file of a release in the artifact store.
type ReleaseAssetModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// ReleaseID is the release it belongs to.
	ReleaseID int `json:"releaseId"`
	// Kind is what the file is and which machine it runs on: one of the
	// kinds of the server's platforms ([BuildKindModel]).
	Kind string `json:"kind"`
	// Path is the file, relative to the artifact store.
	Path string `json:"path"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
