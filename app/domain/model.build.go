package domain

// BuildTargetModel names the release the build endpoints (/build/upload,
// /build/publish) work on: what a CI pipeline sends, and what "release
// upload" and "release publish" send from a developer's machine.
type BuildTargetModel struct {
	// Name is the software's name in the catalog ("pong").
	Name string
	// Version is the release's version ("1.2.0").
	Version string
	// Platform is the software's platform ("ebitengine"); publish only.
	Platform string
	// App is the platform variant ("gui"); empty for the default one.
	App string
	// Commit is the commit the files were built from; may be empty.
	Commit string
}

// BuildUploadModel is the answer of POST /build/upload: the file as the
// server stored it in the release's folder.
type BuildUploadModel struct {
	// File is the stored file name.
	File string `json:"file"`
	// Size is its size in bytes.
	Size int64 `json:"size"`
	// SHA256 is the hex SHA-256 digest the server computed.
	SHA256 string `json:"sha256"`
}

// BuildPublishModel is the answer of POST /build/publish.
type BuildPublishModel struct {
	// Published is true when the release was registered.
	Published bool `json:"published"`
	// Name is the software's name.
	Name string `json:"name"`
	// Platform is the software's platform.
	Platform string `json:"platform"`
	// Version is the published version.
	Version string `json:"version"`
}
