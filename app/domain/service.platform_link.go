package domain

import "errors"

// PlatformLinkServiceInterface manages the platform links: the links shown with every software of a platform.
type PlatformLinkServiceInterface interface {
	// All lists the platform links matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[PlatformLinkModel], error)
	// Show returns the platform link platform_link_id.
	Show(platform_link_id int) (*PlatformLinkModel, error)
	// Create creates a platform link.
	Create(attributes Attributes) (*PlatformLinkModel, error)
	// Update changes the platform link platform_link_id.
	Update(platform_link_id int, attributes Attributes) (*PlatformLinkModel, error)
	// Destroy deletes the platform link platform_link_id.
	Destroy(platform_link_id int) error
}

// PlatformLinkService is the [PlatformLinkServiceInterface] over a [PlatformLinkRepositoryInterface].
type PlatformLinkService struct {
	// PlatformLinkRepository talks to the platform links resource.
	PlatformLinkRepository PlatformLinkRepositoryInterface
}

// All lists the platform links.
func (s PlatformLinkService) All(filter Filter, query ListQuery) (*PageModel[PlatformLinkModel], error) {
	return s.PlatformLinkRepository.All(filter, query)
}

// Show gets the platform link platform_link_id.
func (s PlatformLinkService) Show(platform_link_id int) (*PlatformLinkModel, error) {
	return s.PlatformLinkRepository.FindByID(platform_link_id)
}

// Create checks that attributes hold platform, name and url — the server would
// refuse without them — and creates the platform link.
func (s PlatformLinkService) Create(attributes Attributes) (*PlatformLinkModel, error) {
	for _, key := range []string{"platform", "name", "url"} {
		if !attributes.Has(key) || attributes[key] == nil || attributes[key] == "" {
			return nil, errors.New("a platform link needs a " + key)
		}
	}
	return s.PlatformLinkRepository.Create(attributes)
}

// Update changes the platform link platform_link_id; empty attributes are
// [ErrNothingToChange], since there would be nothing to send.
func (s PlatformLinkService) Update(platform_link_id int, attributes Attributes) (*PlatformLinkModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.PlatformLinkRepository.Update(platform_link_id, attributes)
}

// Destroy deletes the platform link platform_link_id.
func (s PlatformLinkService) Destroy(platform_link_id int) error {
	return s.PlatformLinkRepository.Destroy(platform_link_id)
}
