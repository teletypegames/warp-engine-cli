package domain

// FileTypes are the kinds of [FileEntryModel], as the API spells them.
var FileTypes = []string{"file", "directory"}

// FileEntryModel is one entry of a directory of the artifact store: a file
// or a directory.
type FileEntryModel struct {
	// Name is the entry's name in its directory.
	Name string `json:"name"`
	// Path is its path relative to the artifact store.
	Path string `json:"path"`
	// Type is one of [FileTypes].
	Type string `json:"type"`
	// Size is the size of a file in bytes, or nil (a directory).
	Size *int64 `json:"size"`
	// ModifiedAt is when it last changed (RFC 3339), or nil.
	ModifiedAt *string `json:"modifiedAt"`
	// Downloads is how often a file was downloaded, or nil.
	Downloads *int `json:"downloads"`
}

// FilePathModel is the answer of a change in the artifact store: the path
// of the file or directory that was made or renamed.
type FilePathModel struct {
	// Path is relative to the artifact store.
	Path string `json:"path"`
}
