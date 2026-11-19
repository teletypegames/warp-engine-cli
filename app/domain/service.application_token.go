package domain

import "errors"

// ApplicationTokenServiceInterface manages the application tokens. Creating
// and rotating answer the plain token, which is never shown again.
type ApplicationTokenServiceInterface interface {
	// All lists the tokens matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[ApplicationTokenModel], error)
	// Show returns the token application_token_id.
	Show(application_token_id int) (*ApplicationTokenModel, error)
	// Create issues a token.
	Create(attributes Attributes) (*ApplicationTokenModel, error)
	// Update changes the token application_token_id.
	Update(application_token_id int, attributes Attributes) (*ApplicationTokenModel, error)
	// Destroy revokes the token application_token_id.
	Destroy(application_token_id int) error
	// Rotate replaces the token application_token_id with a new one.
	Rotate(application_token_id int) (*ApplicationTokenModel, error)
}

// ApplicationTokenService is the [ApplicationTokenServiceInterface] over an
// [ApplicationTokenRepositoryInterface].
type ApplicationTokenService struct {
	// ApplicationTokenRepository talks to the application_tokens resource.
	ApplicationTokenRepository ApplicationTokenRepositoryInterface
}

// All lists the tokens.
func (s ApplicationTokenService) All(filter Filter, query ListQuery) (*PageModel[ApplicationTokenModel], error) {
	return s.ApplicationTokenRepository.All(filter, query)
}

// Show gets the token application_token_id.
func (s ApplicationTokenService) Show(application_token_id int) (*ApplicationTokenModel, error) {
	return s.ApplicationTokenRepository.FindByID(application_token_id)
}

// Create checks that attributes hold a name and at least one scope — the
// server would refuse without them — and issues the token.
func (s ApplicationTokenService) Create(attributes Attributes) (*ApplicationTokenModel, error) {
	if name, _ := attributes["name"].(string); name == "" {
		return nil, errors.New("a token needs a name")
	}
	if !tokenScopesGiven(attributes["scopes"]) {
		return nil, errors.New("a token needs at least one scope")
	}
	return s.ApplicationTokenRepository.Create(attributes)
}

// Update changes the token application_token_id; empty attributes are
// [ErrNothingToChange].
func (s ApplicationTokenService) Update(application_token_id int, attributes Attributes) (*ApplicationTokenModel, error) {
	if len(attributes) == 0 {
		return nil, ErrNothingToChange
	}
	return s.ApplicationTokenRepository.Update(application_token_id, attributes)
}

// Destroy revokes the token application_token_id.
func (s ApplicationTokenService) Destroy(application_token_id int) error {
	return s.ApplicationTokenRepository.Destroy(application_token_id)
}

// Rotate replaces the token application_token_id.
func (s ApplicationTokenService) Rotate(application_token_id int) (*ApplicationTokenModel, error) {
	return s.ApplicationTokenRepository.Rotate(application_token_id)
}

// tokenScopesGiven reports whether value is a non-empty list, as the flags
// ([]string) or --data ([]any) give one.
func tokenScopesGiven(value any) bool {
	switch items := value.(type) {
	case []string:
		return len(items) > 0
	case []any:
		return len(items) > 0
	}
	return false
}
