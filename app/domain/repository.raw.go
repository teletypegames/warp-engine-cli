package domain

import (
	"encoding/json"
	"net/url"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// RawRepositoryInterface sends any request to the server and returns its
// answer untouched: the escape hatch for what has no command of its own.
type RawRepositoryInterface interface {
	// Send sends method to path with query and body (nil for none) and
	// returns the JSON answer as it came (nil for an empty answer).
	Send(method string, path string, query url.Values, body any) (json.RawMessage, error)
}

// RawRepository is the [RawRepositoryInterface] over the HTTP client.
type RawRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// Send sends the request. A path under /api/ or /.well-known/ is taken from
// the server's root; any other path is relative to the Admin API
// ([api.AdminPrefix]), so "/softwares" means /api/admin/v2/softwares.
func (r RawRepository) Send(method string, path string, query url.Values, body any) (json.RawMessage, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/.well-known/") {
		path = api.AdminPrefix + path
	}
	var answer json.RawMessage
	if err := r.Client.Do(method, path, query, body, &answer); err != nil {
		return nil, err
	}
	return answer, nil
}
