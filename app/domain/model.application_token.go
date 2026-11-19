package domain

// ApplicationTokenStates are the values of the "state" filter of the
// application tokens.
var ApplicationTokenStates = []string{"active", "expired"}

// ApplicationTokenKinds are the values of the "kind" filter of the
// application tokens: publishing (CI), clients (catalog readers), admin
// (Admin API) and sync (other engines).
var ApplicationTokenKinds = []string{"publishing", "clients", "admin", "sync"}

// ApplicationTokenModel is a token a program signs in to the engine with:
// a pipeline publishing releases, a store client, an admin tool, another
// engine. PlainToken is only in the answer that created or rotated it.
type ApplicationTokenModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Name says what uses it.
	Name string `json:"name"`
	// TokenPrefix is the first characters of the token, as lists show it.
	TokenPrefix string `json:"tokenPrefix"`
	// OwnerType is the kind of account that owns it.
	OwnerType string `json:"ownerType"`
	// OwnerID is the owning account.
	OwnerID int `json:"ownerId"`
	// Scopes are what it may do ("admin:read", "admin:write", "admin:tokens" …).
	Scopes []string `json:"scopes"`
	// Unrestricted lifts the owner restriction: it reaches every owner's
	// records.
	Unrestricted bool `json:"unrestricted"`
	// ExpiresAt is when it stops working (RFC 3339), or nil for never.
	ExpiresAt *string `json:"expiresAt"`
	// LastUsedAt is when it was last used (RFC 3339), or nil.
	LastUsedAt *string `json:"lastUsedAt"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
	// PlainToken is the whole token, only in the answer of a create or a
	// rotate; nil otherwise.
	PlainToken *string `json:"plainToken,omitempty"`
	// Sync counts the CI secrets the new token was written to, only in the
	// answer of a create or a rotate; nil otherwise.
	Sync *ApplicationTokenSyncModel `json:"sync,omitempty"`
}

// ApplicationTokenSyncModel counts the pipelines a new token was pushed to
// as a CI secret.
type ApplicationTokenSyncModel struct {
	// Synced is how many took it.
	Synced int `json:"synced"`
	// Failed is how many did not.
	Failed int `json:"failed"`
}
