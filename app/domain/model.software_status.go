package domain

// SoftwareStatusModel is one word of the catalog's status vocabulary
// ("released", "development" …): what a software's Status names, how the
// stores label it and whether they show such softwares by default.
type SoftwareStatusModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Key is the word a software's Status holds: a-z, 0-9, _ and -.
	Key string `json:"key"`
	// Label is the display name.
	Label string `json:"label"`
	// Position orders the statuses, lowest first.
	Position int `json:"position"`
	// DefaultVisible tells whether the stores show softwares of this status
	// without being asked to.
	DefaultVisible bool `json:"defaultVisible"`
	// Badge is the short tag the stores print on such a software, or nil.
	Badge *string `json:"badge"`
	// SoftwareCount is how many softwares have this status.
	SoftwareCount int `json:"softwareCount"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
