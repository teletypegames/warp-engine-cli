package domain

// ExtraModel is what one of the host's additional resources holds for a
// software (the "extras" of a hosted engine, such as its SEO metadata).
type ExtraModel struct {
	// Association is the resource's name, as the host declares it ("seo_metadata").
	Association string `json:"association"`
	// Multiple tells that the resource holds a list of objects rather than one.
	Multiple bool `json:"multiple"`
	// Value is an object, or an array of objects when Multiple.
	Value any `json:"value"`
}
