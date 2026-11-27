package domain

// SoftwareLinkModel is a plain link of a software to somewhere else: a store
// page, a video, a manual. The source repository is not a link: it is the
// software's [SoftwareSourceModel].
type SoftwareLinkModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// SoftwareID is the software it belongs to.
	SoftwareID int `json:"softwareId"`
	// Label is the link's text.
	Label string `json:"label"`
	// URL is where it points.
	URL string `json:"url"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
