package domain

// Attributes are the fields a create or an update sends, keyed by the API's
// field name (camelCase). Only the fields present are sent, so an update
// changes those and leaves the rest alone. A nil value sends JSON null, which
// clears a nullable field.
//
//	domain.Attributes{"title": "Blessing of Ra", "author": nil}
type Attributes map[string]any

// Has reports whether a has the field key, nil values included.
func (a Attributes) Has(key string) bool {
	_, ok := a[key]
	return ok
}
