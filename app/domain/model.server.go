package domain

// LocalURL is the address of WarpEngine Local on this machine with its
// default port; "warp-engine login" signs in there when it is given no URL.
const LocalURL = "http://127.0.0.1:7470"

// SignInTokenFile is the sign-in method of an engine whose token is in a
// file on this machine (WarpEngine Local): [ServerModel.SignIn] lists it.
const SignInTokenFile = "token-file"

// SignInDevice is the sign-in method of an engine that approves devices in
// the browser: [ServerModel.SignIn] lists it.
const SignInDevice = "device"

// ServerModel is what an engine says about itself at /.well-known/warp-engine.
// Older engines do not answer there; their fields are then zero.
type ServerModel struct {
	// Engine is always "warp-engine".
	Engine string `json:"engine"`
	// Version is the engine's version, "5.5.3".
	Version string `json:"version"`
	// API is the version of the catalog API (WECP) it serves.
	API int `json:"api"`
	// Name is the host's display name.
	Name string `json:"name"`
	// Profile is "server" for a hosted engine, "local" for WarpEngine Local.
	Profile string `json:"profile"`
	// SignIn lists the sign-in methods: [SignInDevice], [SignInTokenFile].
	SignIn []string `json:"signIn"`
	// AdminAPI is the path of the Admin API, or empty when it is off.
	AdminAPI string `json:"adminApi"`
	// Admin reports whether the admin web UI is on.
	Admin bool `json:"admin"`
}

// Supports reports whether s lists the sign-in method method.
func (s ServerModel) Supports(method string) bool {
	for _, known := range s.SignIn {
		if known == method {
			return true
		}
	}
	return false
}

// PrincipalModel is who the server takes the token for (GET /me).
type PrincipalModel struct {
	// Label is the token's label, or nil.
	Label *string `json:"label"`
	// OwnerType is the kind of account the token belongs to, or nil.
	OwnerType *string `json:"ownerType"`
	// OwnerID is the account's identifier, or nil.
	OwnerID *int `json:"ownerId"`
	// Scopes are the scopes the token carries.
	Scopes []string `json:"scopes"`
	// Unrestricted reports a token that may use every resource.
	Unrestricted bool `json:"unrestricted"`
	// TokenPrefix is the first characters of the token, to recognise it.
	TokenPrefix *string `json:"tokenPrefix"`
	// ExpiresAt is when the token stops working (RFC 3339), or nil for never.
	ExpiresAt *string `json:"expiresAt"`
	// Resources are the Admin API resources the token may use.
	Resources []string `json:"resources"`
	// ReadOnly reports an engine whose Admin API only reads.
	ReadOnly bool `json:"readOnly"`
	// MCP reports whether the engine serves the Admin API over MCP too.
	MCP bool `json:"mcp"`
}

// DeviceGrantModel is a started device sign-in: the code the person types
// in the browser and the code this program polls with.
type DeviceGrantModel struct {
	// DeviceCode is the secret this program polls with.
	DeviceCode string `json:"deviceCode"`
	// UserCode is the code the person checks or types, "ABCD-1234".
	UserCode string `json:"userCode"`
	// VerificationURL is the page that approves the sign-in.
	VerificationURL string `json:"verificationUrl"`
	// VerificationURLComplete is VerificationURL with the code filled in.
	VerificationURLComplete string `json:"verificationUrlComplete"`
	// Scopes are the scopes asked for.
	Scopes []string `json:"scopes"`
	// Interval is how many seconds to wait between two polls.
	Interval int `json:"interval"`
	// ExpiresIn is how many seconds the codes live.
	ExpiresIn int `json:"expiresIn"`
}

// Token states of a device sign-in poll ([TokenPollModel.State]).
const (
	// TokenPending: the person has not answered yet.
	TokenPending = "pending"
	// TokenApproved: approved; the poll carries the token.
	TokenApproved = "approved"
	// TokenDenied: the person refused.
	TokenDenied = "denied"
	// TokenExpired: the codes expired.
	TokenExpired = "expired"
)

// TokenPollModel is one answer of polling a device sign-in.
type TokenPollModel struct {
	// State is one of [TokenPending], [TokenApproved], [TokenDenied], [TokenExpired].
	State string `json:"state"`
	// Token is the new token, only with [TokenApproved].
	Token string `json:"token"`
}

// DefaultScopes are the scopes a device sign-in asks for when none are named:
// reading and changing everything except the tokens.
var DefaultScopes = []string{"admin:read", "admin:write"}

// Scopes the Admin API knows.
const (
	// ScopeRead reads every resource.
	ScopeRead = "admin:read"
	// ScopeWrite creates, changes and deletes them.
	ScopeWrite = "admin:write"
	// ScopeTokens manages application tokens and reads generated sign-ins.
	ScopeTokens = "admin:tokens"
	// ScopePair pairs another engine for sync.
	ScopePair = "sync:pair"
)
