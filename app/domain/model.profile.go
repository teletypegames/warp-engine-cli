package domain

// ProfileModel is one server the command line has signed in to: where it is
// and the token it signs in with. Profiles are kept in the configuration file
// (see [ProfileRepository]).
type ProfileModel struct {
	// Name is what the profile is called on the command line (--profile).
	Name string `json:"name"`
	// URL is the server's base URL, https://games.example.org.
	URL string `json:"url"`
	// Token is the application token; empty for a WarpEngine Local profile,
	// whose token is read from the token file on every run.
	Token string `json:"token,omitempty"`
	// Label is who the token belongs to, as the server named it at sign-in.
	Label string `json:"label,omitempty"`
	// Scopes are the scopes the token was granted.
	Scopes []string `json:"scopes,omitempty"`
	// TokenFile marks a WarpEngine Local profile: the token is read from the
	// local engine's token file, not from Token.
	TokenFile bool `json:"tokenFile,omitempty"`
}

// ConfigModel is the whole configuration file: the profiles and which one is
// used when none is named.
type ConfigModel struct {
	// Current is the name of the default profile; empty when there is none.
	Current string `json:"current"`
	// Profiles are the saved profiles, keyed by name.
	Profiles map[string]ProfileModel `json:"profiles"`
}

// Selection is what the command line was told about the server to use,
// before a profile is chosen. Empty fields are not given. See
// [AuthService.Connect] for the order they win in.
type Selection struct {
	// Profile is the profile name (--profile, WARP_ENGINE_PROFILE).
	Profile string
	// Server is a server URL (--server, WARP_ENGINE_URL).
	Server string
	// Token is a token (--token, WARP_ENGINE_TOKEN).
	Token string
}
