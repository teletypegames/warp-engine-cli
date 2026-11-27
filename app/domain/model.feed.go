package domain

// FeedModel is one configured feed of the engine: a format the catalog
// leaves it in (the catalog API, a directory, an Assembly64 server …).
type FeedModel struct {
	// Key is the feed key, unique among the feeds.
	Key string `json:"key"`
	// Type is the feed type: wecp, directory, assembly64 or one of the host's.
	Type string `json:"type"`
	// Platforms are the platforms it carries; empty for every platform.
	Platforms []string `json:"platforms"`
	// Description is a one-line description.
	Description string `json:"description"`
	// Static tells a feed that writes files.
	Static bool `json:"static"`
	// Dynamic tells a feed that answers requests.
	Dynamic bool `json:"dynamic"`
	// Details is the rest of its descriptor (hosts, path, kinds …).
	Details map[string]any `json:"details"`
}

// FeedCheckModel is what one feed selects and leaves out.
type FeedCheckModel struct {
	// Key is the feed key.
	Key string `json:"key"`
	// Type is the feed type.
	Type string `json:"type"`
	// Selection is false for the catalog API, which decides per caller.
	Selection bool `json:"selection"`
	// Items are the titles on offer.
	Items []FeedItemModel `json:"items"`
	// Skipped are the softwares it leaves out, with the reason.
	Skipped []FeedSkipModel `json:"skipped"`
}

// FeedItemModel is one title a feed offers: a release and the asset
// kind it serves.
type FeedItemModel struct {
	// SoftwareID is the software's identifier.
	SoftwareID int `json:"softwareId"`
	// Name is the software name.
	Name string `json:"name"`
	// Title is the software title.
	Title string `json:"title"`
	// Version is the release offered.
	Version string `json:"version"`
	// Kind is the asset kind served.
	Kind string `json:"kind"`
}

// FeedSkipModel is one software a feed leaves out.
type FeedSkipModel struct {
	// SoftwareID is the software's identifier.
	SoftwareID int `json:"softwareId"`
	// Name is the software name.
	Name string `json:"name"`
	// Reason is why it is left out.
	Reason string `json:"reason"`
}
