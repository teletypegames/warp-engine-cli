package domain

// ExternalLinkKinds are the kinds of an external link, as the API spells
// them: "source" is the software's git repository (one per software), "link"
// anything else.
var ExternalLinkKinds = []string{"link", "source"}

// ExternalLinkModel is a link of a software to somewhere else: its source
// repository, a store page, a video.
type ExternalLinkModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// SoftwareID is the software it belongs to.
	SoftwareID int `json:"softwareId"`
	// Label is the link's text.
	Label string `json:"label"`
	// URL is where it points.
	URL string `json:"url"`
	// Kind is one of [ExternalLinkKinds].
	Kind string `json:"kind"`
	// Ref is the branch or tag of a source link, or nil for the default
	// branch.
	Ref *string `json:"ref"`
	// Public tells whether the server found a source link publicly
	// reachable, or nil when it has not checked.
	Public *bool `json:"public"`
	// VerifiedAt is when the server last checked the link (RFC 3339), or nil.
	VerifiedAt *string `json:"verifiedAt"`
	// VerifyError is why the last check failed, or nil.
	VerifyError *string `json:"verifyError"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
