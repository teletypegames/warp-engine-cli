package domain

// OutletModel is one configured outlet of the engine: a format the catalog
// leaves it in (the catalog API, a shelf, an Assembly64 server …).
type OutletModel struct {
	// Key is the outlet key, unique among the outlets.
	Key string `json:"key"`
	// Type is the outlet type: wecp, shelf, assembly64 or one of the host's.
	Type string `json:"type"`
	// Platforms are the platforms it carries; empty for every platform.
	Platforms []string `json:"platforms"`
	// Description is a one-line description.
	Description string `json:"description"`
	// Materialized tells an outlet that writes files.
	Materialized bool `json:"materialized"`
	// Served tells an outlet that answers requests.
	Served bool `json:"served"`
	// Details is the rest of its descriptor (hosts, path, kinds …).
	Details map[string]any `json:"details"`
}

// OutletCheckModel is what one outlet offers and leaves out.
type OutletCheckModel struct {
	// Key is the outlet key.
	Key string `json:"key"`
	// Type is the outlet type.
	Type string `json:"type"`
	// Offer is false for the catalog API, which decides per caller.
	Offer bool `json:"offer"`
	// Items are the titles on offer.
	Items []OutletItemModel `json:"items"`
	// Skipped are the softwares it leaves out, with the reason.
	Skipped []OutletSkipModel `json:"skipped"`
}

// OutletItemModel is one title an outlet offers: a release and the asset
// kind it serves.
type OutletItemModel struct {
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

// OutletSkipModel is one software an outlet leaves out.
type OutletSkipModel struct {
	// SoftwareID is the software's identifier.
	SoftwareID int `json:"softwareId"`
	// Name is the software name.
	Name string `json:"name"`
	// Reason is why it is left out.
	Reason string `json:"reason"`
}
