package domain

// PlatformLinkModel is a link shown with every software of a platform: the
// emulator, the runtime or the documentation a player of that platform needs.
type PlatformLinkModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Platform is one of the server's platforms ([BuildPlatformModel]).
	Platform string `json:"platform"`
	// Name is the link's text.
	Name string `json:"name"`
	// URL is where it points.
	URL string `json:"url"`
	// Position orders the links of a platform, lowest first.
	Position int `json:"position"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
