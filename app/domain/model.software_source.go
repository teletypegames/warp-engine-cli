package domain

// SoftwareSourceModel is the source repository of a software: one per
// software, the git repository its forks start from.
type SoftwareSourceModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// SoftwareID is the software it belongs to.
	SoftwareID int `json:"softwareId"`
	// URL is the repository's address.
	URL string `json:"url"`
	// Ref is the branch or tag, or nil for the default branch.
	Ref *string `json:"ref"`
	// Public tells whether the server found the repository publicly
	// reachable, or nil when it has not checked.
	Public *bool `json:"public"`
	// VerifiedAt is when the server last checked it (RFC 3339), or nil.
	VerifiedAt *string `json:"verifiedAt"`
	// VerifyError is why the last check failed, or nil.
	VerifyError *string `json:"verifyError"`
	// CreatedAt is when it was set (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}
