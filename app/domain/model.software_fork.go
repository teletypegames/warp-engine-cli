package domain

// SoftwareForkStatuses are the states of an announced fork, as the API spells them:
// pending waits for moderation, visible lists the fork with its parent,
// hidden keeps it out, gone is a fork its engine no longer has.
var SoftwareForkStatuses = []string{"pending", "visible", "hidden", "gone"}

// SoftwareForkModel is a fork of one of this engine's softwares that another
// engine announced: the fork's engine and name, and how it was moderated.
type SoftwareForkModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// SoftwareID is the parent software on this engine.
	SoftwareID int `json:"softwareId"`
	// SoftwareName is the parent software's name, or nil.
	SoftwareName *string `json:"softwareName"`
	// EngineURL is the base URL of the engine the fork lives on.
	EngineURL string `json:"engineUrl"`
	// Name is the fork's name on its engine.
	Name string `json:"name"`
	// Title is the fork's display name, or nil.
	Title *string `json:"title"`
	// Status is one of [SoftwareForkStatuses].
	Status string `json:"status"`
	// CheckedAt is when this engine last checked the fork (RFC 3339), or nil.
	CheckedAt *string `json:"checkedAt"`
	// LastError is why the last check failed, or nil.
	LastError *string `json:"lastError"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
