package domain

// Platforms are the platforms a software can be built for, as the API
// spells them.
var Platforms = []string{"tic80", "ebitengine", "love", "c64", "plus4", "godot", "bevy", "phaser", "electron", "dos"}

// SoftwareModel is one entry of the catalog: a game, demo or tool, with
// the releases, images and links that belong to it.
type SoftwareModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Name is the internal name, unique in the catalog: a-z, 0-9 and dashes.
	Name string `json:"name"`
	// Title is the display name.
	Title string `json:"title"`
	// Author is who made it, or nil.
	Author *string `json:"author"`
	// Desc is the short description, or nil.
	Desc *string `json:"desc"`
	// Story is the long description, or nil.
	Story *string `json:"story"`
	// License is the licence, or nil.
	License *string `json:"license"`
	// Platform is one of [Platforms].
	Platform string `json:"platform"`
	// Status is a key of the software statuses ([SoftwareStatusModel]), or nil.
	Status *string `json:"status"`
	// Category is a key of the software categories ([SoftwareCategoryModel]), or nil.
	Category *string `json:"category"`
	// Highlighted marks a software the stores show first, or nil.
	Highlighted *bool `json:"highlighted"`
	// ShelfName is the file name on a shelf (A-Z, 0-9, -), or nil to use Name.
	ShelfName *string `json:"shelfName"`
	// OwnerType is the kind of account that owns it, or nil.
	OwnerType *string `json:"ownerType"`
	// OwnerID is the owning account, or nil.
	OwnerID *int `json:"ownerId"`
	// ForksAllowed tells whether the developer tools may offer to fork it
	// (it needs a public source link).
	ForksAllowed bool `json:"forksAllowed"`
	// ForkOfEngine is, for a fork, the parent engine's base URL, or nil.
	ForkOfEngine *string `json:"forkOfEngine"`
	// ForkOfName is, for a fork, the parent's name on its engine, or nil.
	ForkOfName *string `json:"forkOfName"`
	// ForkOfVersion is, for a fork, the parent's version it was forked
	// from, or nil.
	ForkOfVersion *string `json:"forkOfVersion"`
	// ForkOfRepo is, for a fork, the parent's source repository, or nil.
	ForkOfRepo *string `json:"forkOfRepo"`
	// CreatedAt is when it was created (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
	// Releases are its releases, only when asked for with the include "releases".
	Releases []ReleaseModel `json:"releases,omitempty"`
	// Images are its images, only with the include "images".
	Images []SoftwareImageModel `json:"images,omitempty"`
	// ExternalLinks are its links, only with the include "externalLinks".
	ExternalLinks []ExternalLinkModel `json:"externalLinks,omitempty"`
	// Extras are the host's additional resources, one key each, only with
	// the include "extras".
	Extras map[string]any `json:"extras,omitempty"`
}

// SoftwareIncludes are the related records [SoftwareRepository.FindByID] can
// embed in a software.
var SoftwareIncludes = []string{"releases", "images", "externalLinks", "extras"}

// ScaffoldModel is the answer of setting a new software up from nothing:
// the software, its repository on the forge and its pipeline.
type ScaffoldModel struct {
	// Software is the created software.
	Software SoftwareModel `json:"software"`
	// Repository describes the created forge repository (its URLs and name).
	Repository map[string]any `json:"repository"`
	// PipelineID is the CI pipeline set up for it, or nil.
	PipelineID *int `json:"pipelineId"`
	// Warnings are the steps that did not fully succeed.
	Warnings []string `json:"warnings"`
}
