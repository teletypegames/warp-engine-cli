package domain

// LocalProviderKeys are the providers WarpEngine Local can run, as the
// API spells them.
var LocalProviderKeys = []string{"gitea", "woodpecker"}

// LocalProviderStates are the states of a [LocalProviderModel].
var LocalProviderStates = []string{"off", "starting", "ready", "stopped", "failed", "unavailable"}

// LocalProviderModel is one provider WarpEngine Local can run next to
// the engine: the git server, the CI.
type LocalProviderModel struct {
	// Key is the identifier, one of [LocalProviderKeys].
	Key string `json:"key"`
	// Label is the display name.
	Label string `json:"label"`
	// Role is what it is for ("scm" for the git server, "ci").
	Role string `json:"role"`
	// Requires are the roles of the providers it needs; enabling it
	// enables them too.
	Requires []string `json:"requires"`
	// Description says what it does.
	Description string `json:"description"`
	// State is one of [LocalProviderStates].
	State string `json:"state"`
	// Message explains the state (why it failed or is unavailable), or nil.
	Message *string `json:"message"`
}

// LocalCredentialsModel is the generated sign-in of a provider, to open
// it in a browser.
type LocalCredentialsModel struct {
	// URL is the provider's address.
	URL string `json:"url"`
	// User is the user name.
	User string `json:"user"`
	// Password is the password.
	Password string `json:"password"`
}
