package domain

import (
	"net/http"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// ApplicationTokenRepositoryInterface is the Admin API's application_tokens
// resource.
type ApplicationTokenRepositoryInterface interface {
	// All lists the tokens matching filter, paged by query.
	All(filter Filter, query ListQuery) (*PageModel[ApplicationTokenModel], error)
	// FindByID returns the token application_token_id.
	FindByID(application_token_id int) (*ApplicationTokenModel, error)
	// Create issues a token from attributes; the answer holds the plain token.
	Create(attributes Attributes) (*ApplicationTokenModel, error)
	// Update changes the token application_token_id.
	Update(application_token_id int, attributes Attributes) (*ApplicationTokenModel, error)
	// Destroy revokes the token application_token_id.
	Destroy(application_token_id int) error
	// Rotate replaces the token application_token_id with a new one; the
	// answer is the new token with its plain token.
	Rotate(application_token_id int) (*ApplicationTokenModel, error)
}

// ApplicationTokenRepository is the [ApplicationTokenRepositoryInterface]
// over the Admin API.
type ApplicationTokenRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /application_tokens.
func (r ApplicationTokenRepository) All(filter Filter, query ListQuery) (*PageModel[ApplicationTokenModel], error) {
	return list[ApplicationTokenModel](r.Client, "/application_tokens", filter, query)
}

// FindByID gets GET /application_tokens/{id}.
func (r ApplicationTokenRepository) FindByID(application_token_id int) (*ApplicationTokenModel, error) {
	return api.Get[ApplicationTokenModel](r.Client, applicationTokenPath(application_token_id), nil)
}

// Create posts POST /application_tokens.
func (r ApplicationTokenRepository) Create(attributes Attributes) (*ApplicationTokenModel, error) {
	return api.Send[ApplicationTokenModel](r.Client, http.MethodPost, "/application_tokens", attributes)
}

// Update sends PATCH /application_tokens/{id}.
func (r ApplicationTokenRepository) Update(application_token_id int, attributes Attributes) (*ApplicationTokenModel, error) {
	return api.Send[ApplicationTokenModel](r.Client, http.MethodPatch, applicationTokenPath(application_token_id), attributes)
}

// Destroy sends DELETE /application_tokens/{id}.
func (r ApplicationTokenRepository) Destroy(application_token_id int) error {
	return api.Delete(r.Client, applicationTokenPath(application_token_id), nil)
}

// Rotate posts POST /application_tokens/{id}/rotate.
func (r ApplicationTokenRepository) Rotate(application_token_id int) (*ApplicationTokenModel, error) {
	return api.Send[ApplicationTokenModel](r.Client, http.MethodPost, applicationTokenPath(application_token_id)+"/rotate", nil)
}

// applicationTokenPath is the path of the token application_token_id.
func applicationTokenPath(application_token_id int) string {
	return "/application_tokens/" + strconv.Itoa(application_token_id)
}
