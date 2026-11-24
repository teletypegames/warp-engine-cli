package domain

// PipelineUnknownPlatform is the platform of a pipeline whose repository
// the engine could not place; a pipeline takes it besides the server's
// platforms.
const PipelineUnknownPlatform = "unknown"

// PipelineModel is one CI repository the engine knows: a repository on the
// CI server, with the platform it builds for and the software it publishes.
type PipelineModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// RepoID is the repository's identifier on the CI server.
	RepoID int `json:"repoId"`
	// RepoOwner is the owner (user or organisation) of the repository.
	RepoOwner string `json:"repoOwner"`
	// RepoName is the repository's name.
	RepoName string `json:"repoName"`
	// FullName is "owner/name".
	FullName string `json:"fullName"`
	// Platform is one of the server's platforms or [PipelineUnknownPlatform].
	Platform string `json:"platform"`
	// Active tells whether the repository is still active on the CI server.
	Active bool `json:"active"`
	// SoftwareID is the software its builds are published as, or nil.
	SoftwareID *int `json:"softwareId"`
	// LastPipelineStatus is the status of the newest run ("success",
	// "failure" …), or nil.
	LastPipelineStatus *string `json:"lastPipelineStatus"`
	// LastPipelineAt is when the newest run started (RFC 3339), or nil.
	LastPipelineAt *string `json:"lastPipelineAt"`
	// CreatedAt is when the engine first saw it (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}

// PipelineSyncModel is the answer of syncing the pipelines with the CI
// server: the pipelines it created, updated and deactivated.
type PipelineSyncModel struct {
	// Created are the repositories new to the engine.
	Created []PipelineModel `json:"created"`
	// Updated are the known repositories that changed.
	Updated []PipelineModel `json:"updated"`
	// Deactivated are the repositories gone from the CI server.
	Deactivated []PipelineModel `json:"deactivated"`
}

// PipelineRunModel is one run of a pipeline on the CI server. Runs have no
// identifier of their own in the engine: Number is the CI server's.
type PipelineRunModel struct {
	// Number is the run's number on the CI server.
	Number int `json:"number"`
	// Status is the state of the run ("pending", "running", "success" …), or nil.
	Status *string `json:"status"`
	// Branch is the branch it ran on, or nil.
	Branch *string `json:"branch"`
	// Message is the commit message, or nil.
	Message *string `json:"message"`
	// CreatedAt is when it was started (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// URL is the run's page on the CI server, or nil.
	URL *string `json:"url"`
	// Commit is the commit it built, or nil.
	Commit *string `json:"commit"`
	// Event is what started it ("push", "tag", "manual" …), or nil.
	Event *string `json:"event"`
	// Ref is the git ref it built, or nil.
	Ref *string `json:"ref"`
}
