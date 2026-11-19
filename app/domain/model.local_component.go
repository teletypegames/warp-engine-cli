package domain

// LocalComponentKeys are the components WarpEngine Local can run, as the
// API spells them.
var LocalComponentKeys = []string{"gitea", "woodpecker"}

// LocalComponentStates are the states of a [LocalComponentModel].
var LocalComponentStates = []string{"off", "starting", "ready", "stopped", "failed", "unavailable"}

// LocalComponentModel is one component WarpEngine Local can run next to
// the engine: the git server, the CI.
type LocalComponentModel struct {
	// Key is the identifier, one of [LocalComponentKeys].
	Key string `json:"key"`
	// Label is the display name.
	Label string `json:"label"`
	// Role is what it is for ("git", "ci").
	Role string `json:"role"`
	// Requires are the keys of the components it needs; enabling it
	// enables them too.
	Requires []string `json:"requires"`
	// Description says what it does.
	Description string `json:"description"`
	// State is one of [LocalComponentStates].
	State string `json:"state"`
	// Message explains the state (why it failed or is unavailable), or nil.
	Message *string `json:"message"`
}

// LocalCredentialsModel is the generated sign-in of a component, to open
// it in a browser.
type LocalCredentialsModel struct {
	// URL is the component's address.
	URL string `json:"url"`
	// User is the user name.
	User string `json:"user"`
	// Password is the password.
	Password string `json:"password"`
}
