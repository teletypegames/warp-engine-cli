package domain

// BuildPlatformModel is one platform the server builds and publishes for,
// as its public build description (GET /build/platforms) gives it: the
// truth about which platforms, asset kinds and app variants exist on that
// server. Only the fields the command line uses are read.
type BuildPlatformModel struct {
	// Name is the platform as the API spells it ("tic80", "golang" …).
	Name string `json:"name"`
	// Label is the display name ("TIC-80").
	Label string `json:"label"`
	// Kinds are the asset kinds a release of the platform has.
	Kinds []BuildKindModel `json:"kinds"`
	// Apps are the platform's variants ("cli", "gui" for golang); empty for
	// a platform without variants.
	Apps []BuildAppModel `json:"apps"`
}

// BuildKindModel is one asset kind of a [BuildPlatformModel] or a
// [BuildAppModel].
type BuildKindModel struct {
	// Name is the kind as the API spells it ("cartridge", "cli_linux_x64" …).
	Name string `json:"name"`
	// Built is false for a kind the CI does not build, which is uploaded by
	// hand (the macOS packages, signed only on a Mac).
	Built bool `json:"built"`
}

// BuildAppModel is one variant of a [BuildPlatformModel].
type BuildAppModel struct {
	// Name is the variant as the API spells it ("cli").
	Name string `json:"name"`
	// Default marks the variant used when none is asked for.
	Default bool `json:"default"`
	// Kinds are the asset kinds a release of the variant has.
	Kinds []BuildKindModel `json:"kinds"`
}
