package domain

// SyncTransports are the ways a sync peer is reached, as the API spells them.
var SyncTransports = []string{"http", "rsync"}

// SyncPeerModel is another engine this one syncs with: it pushes releases
// to it (outbound), takes its pushes (inbound), or both.
type SyncPeerModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Label is the display name.
	Label string `json:"label"`
	// Transport is one of [SyncTransports].
	Transport string `json:"transport"`
	// TransportOptions are the transport's settings: host for http; host,
	// user, path, port and host_key for rsync.
	TransportOptions map[string]any `json:"transportOptions"`
	// SecretPrefix is the first characters of the stored secret, or nil.
	SecretPrefix *string `json:"secretPrefix"`
	// OutboundEnabled pushes this engine's releases to the peer.
	OutboundEnabled bool `json:"outboundEnabled"`
	// InboundEnabled accepts the peer's pushes.
	InboundEnabled bool `json:"inboundEnabled"`
	// AllowedChannels narrows the channels pushed to it (see [Channels]);
	// empty pushes every channel the engine allows.
	AllowedChannels []string `json:"allowedChannels"`
	// CreatedBy is who paired the peer through the sync pairings, or nil.
	CreatedBy *SyncPeerCreatorModel `json:"createdBy"`
	// InboundTokenPrefix is the first characters of the token the peer
	// pushes here with, or nil.
	InboundTokenPrefix *string `json:"inboundTokenPrefix"`
	// LastSyncedAt is when a push last succeeded (RFC 3339), or nil.
	LastSyncedAt *string `json:"lastSyncedAt"`
	// LastError is the last failure, or nil.
	LastError *string `json:"lastError"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}

// SyncPeerCreatorModel is the account that paired a sync peer.
type SyncPeerCreatorModel struct {
	// Type is the kind of account.
	Type string `json:"type"`
	// ID is the account's identifier.
	ID int `json:"id"`
}

// SyncCheckModel is the answer of checking a sync peer: whether the other
// engine accepts this one.
type SyncCheckModel struct {
	// OK reports whether the peer accepted the connection and the secret.
	OK bool `json:"ok"`
	// Status is the HTTP status the peer answered, or nil.
	Status *int `json:"status"`
	// Version is the peer's engine version, or nil.
	Version *string `json:"version"`
	// Message tells what happened.
	Message string `json:"message"`
}

// SyncInboundTokenModel is a new inbound token of a sync peer: the token the
// peer pushes here with. PlainToken is in this answer only.
type SyncInboundTokenModel struct {
	// TokenPrefix is the first characters of the token, as lists show it.
	TokenPrefix string `json:"tokenPrefix"`
	// PlainToken is the whole token, never shown again.
	PlainToken string `json:"plainToken"`
}

// SyncPushModel is the answer of a push or a sync run: the releases that
// go, to which peers, and how many pushes were queued.
type SyncPushModel struct {
	// DryRun reports that nothing was queued, only planned.
	DryRun bool `json:"dryRun"`
	// Queued is how many pushes were queued (0 on a dry run).
	Queued int `json:"queued"`
	// Peers are the IDs of the peers pushed to.
	Peers []int `json:"peers"`
	// Releases are the releases that go.
	Releases []SyncPushReleaseModel `json:"releases"`
	// Skipped are the releases left out, with the reason.
	Skipped []string `json:"skipped"`
}

// SyncPushReleaseModel is one release of a [SyncPushModel].
type SyncPushReleaseModel struct {
	// Software is the software's name.
	Software string `json:"software"`
	// Version is the release's version.
	Version string `json:"version"`
}
