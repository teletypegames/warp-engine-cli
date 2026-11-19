package domain

// SyncPairingModel is an inbound-only sync peer the caller paired with this
// engine itself: another engine that pushes here with the pairing's token.
// The server offers the sync pairings only while inbound sync is on.
type SyncPairingModel struct {
	// ID is the identifier (the sync peer's).
	ID int `json:"id"`
	// Label says what the other engine is ("laptop local").
	Label string `json:"label"`
	// InboundTokenPrefix is the first characters of the token the other
	// engine pushes with, or nil.
	InboundTokenPrefix *string `json:"inboundTokenPrefix"`
	// PlainToken is the whole new token, only in the answer that issued it
	// (create and token); nil otherwise.
	PlainToken *string `json:"plainToken"`
	// LastInboundAt is when the other engine last pushed (RFC 3339), or nil.
	LastInboundAt *string `json:"lastInboundAt"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
