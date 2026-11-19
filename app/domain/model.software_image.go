package domain

// SoftwareImageModel is a picture of a software: a screenshot or a cover,
// bound to the host's image library.
type SoftwareImageModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// SoftwareID is the software it shows.
	SoftwareID int `json:"softwareId"`
	// ImageID is the image in the host's image library.
	ImageID int `json:"imageId"`
	// IsDefault marks the image the stores show first.
	IsDefault bool `json:"isDefault"`
	// Position orders the images, from 0.
	Position int `json:"position"`
	// Origin is who set it: "local", or "sync" for one another engine pushed.
	Origin string `json:"origin"`
	// URL is where the image can be fetched, or nil.
	URL *string `json:"url"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
