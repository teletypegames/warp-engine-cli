package domain

// RuntimeSources are where a runtime version came from, as the API spells
// them.
var RuntimeSources = []string{"download", "upload"}

// RuntimeVersionModel is one installed version of a runtime: the player a
// platform's games run in on the web (an emulator, an engine build).
type RuntimeVersionModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Name is the runtime's name.
	Name string `json:"name"`
	// Version is the version.
	Version string `json:"version"`
	// Source is one of [RuntimeSources].
	Source string `json:"source"`
	// Path is where it lies, relative to the artifact store.
	Path string `json:"path"`
	// SizeBytes is its size in bytes.
	SizeBytes int64 `json:"sizeBytes"`
	// Components are the parts installed with it.
	Components []RuntimeComponentModel `json:"components"`
	// Releases is how many releases are bound to it.
	Releases int `json:"releases"`
	// CreatedAt is when it was installed (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}

// RuntimeComponentModel is one part of a [RuntimeVersionModel].
type RuntimeComponentModel struct {
	// Kind is what sort of part it is.
	Kind string `json:"kind"`
	// Key names the part within its kind.
	Key string `json:"key"`
}

// RuntimeInstallModel is the answer of installing a runtime: what was
// fetched and what was already there.
type RuntimeInstallModel struct {
	// Name is the runtime installed.
	Name string `json:"name"`
	// Version is the version installed.
	Version string `json:"version"`
	// Fetched are the files downloaded (or, on a dry run, to be downloaded).
	Fetched []string `json:"fetched"`
	// Kept are the files that were there already.
	Kept []string `json:"kept"`
	// Components are the components its platforms need, by kind.
	Components map[string]any `json:"components"`
	// DryRun tells that nothing was installed, only reported.
	DryRun bool `json:"dryRun"`
}
