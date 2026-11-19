package domain

// SoftwareCategoryModel is one word of the catalog's category vocabulary
// ("game", "demo", "tool" …). A category key a software brings that the
// vocabulary does not know yet is registered with NeedsReview set.
type SoftwareCategoryModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Key is the word a software's Category holds.
	Key string `json:"key"`
	// Label is the display name.
	Label string `json:"label"`
	// Position orders the categories, lowest first.
	Position int `json:"position"`
	// NeedsReview marks a category registered on the fly that nobody has
	// looked at yet.
	NeedsReview bool `json:"needsReview"`
	// SoftwareCount is how many softwares are in this category.
	SoftwareCount int `json:"softwareCount"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
