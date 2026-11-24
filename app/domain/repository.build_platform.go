package domain

import (
	"net/http"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// BuildPlatformRepositoryInterface is the server's public build description:
// the platforms it builds for. It needs no sign-in.
type BuildPlatformRepositoryInterface interface {
	// All returns every platform of the server.
	All() ([]BuildPlatformModel, error)
}

// BuildPlatformRepository is the [BuildPlatformRepositoryInterface] over
// GET /build/platforms (outside the Admin API).
type BuildPlatformRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// All gets GET /build/platforms and returns its platforms.
func (r BuildPlatformRepository) All() ([]BuildPlatformModel, error) {
	var answer struct {
		Platforms []BuildPlatformModel `json:"platforms"`
	}
	if err := r.Client.Do(http.MethodGet, "/build/platforms", nil, nil, &answer); err != nil {
		return nil, err
	}
	return answer.Platforms, nil
}
