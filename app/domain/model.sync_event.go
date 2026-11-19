package domain

// SyncDirections are the directions of a sync event, as the API spells them.
var SyncDirections = []string{"outbound", "inbound"}

// SyncEventStatuses are the states of a sync event, as the API spells them.
var SyncEventStatuses = []string{"queued", "running", "succeeded", "failed", "skipped"}

// SyncEventModel is one push of one release between this engine and a sync
// peer, in either direction, with how it went.
type SyncEventModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// SyncPeerID is the peer, or nil when it was deleted.
	SyncPeerID *int `json:"syncPeerId"`
	// SoftwareID is the software here, or nil.
	SoftwareID *int `json:"softwareId"`
	// Direction is one of [SyncDirections].
	Direction string `json:"direction"`
	// Status is one of [SyncEventStatuses].
	Status string `json:"status"`
	// SoftwareName is the software's name, or nil.
	SoftwareName *string `json:"softwareName"`
	// Platform is the software's platform, or nil.
	Platform *string `json:"platform"`
	// Version is the release's version, or nil.
	Version *string `json:"version"`
	// ReleaseID is the release here, or nil.
	ReleaseID *int `json:"releaseId"`
	// Step is the step it reached or failed at, or nil.
	Step *string `json:"step"`
	// HTTPStatus is the status the other side answered, or nil.
	HTTPStatus *int `json:"httpStatus"`
	// Message tells what happened, or nil.
	Message *string `json:"message"`
	// FilesCount is how many files went, or nil.
	FilesCount *int `json:"filesCount"`
	// Bytes is how many bytes went, or nil.
	Bytes *int64 `json:"bytes"`
	// DurationMs is how long it took in milliseconds, or nil.
	DurationMs *int `json:"durationMs"`
	// Attempt is the try it is on, from 1, or nil.
	Attempt *int `json:"attempt"`
	// TokenPrefix is the first characters of the token an inbound push came
	// with, or nil.
	TokenPrefix *string `json:"tokenPrefix"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}

// SyncStatusModel is where one software stands with one peer in one
// direction: the latest event, and for outbound why the newest release
// would not go.
type SyncStatusModel struct {
	// PeerID is the sync peer.
	PeerID int `json:"peerId"`
	// Peer is the peer's label.
	Peer string `json:"peer"`
	// Direction is one of [SyncDirections].
	Direction string `json:"direction"`
	// Reason is why the newest release would not go (outbound), or nil.
	Reason *string `json:"reason"`
	// Event is the latest event, or nil when there was none.
	Event *SyncEventModel `json:"event"`
}
