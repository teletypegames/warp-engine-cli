package domain

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// AuthRepositoryInterface is the sign-in side of a server: who it is, the
// device sign-in, the current principal and revoking the token. It also
// points the shared client at a server, since every other repository talks
// to the server it is pointed at.
type AuthRepositoryInterface interface {
	// Use points the client at url and signs it in with token.
	Use(url string, token string)
	// Discover reads /.well-known/warp-engine. An engine that does not
	// answer there (older than 5.2) returns a zero ServerModel and no error.
	Discover() (*ServerModel, error)
	// StartDeviceGrant starts a device sign-in for client_name with scopes.
	StartDeviceGrant(client_name string, scopes []string) (*DeviceGrantModel, error)
	// PollDeviceGrant asks how the device sign-in device_code stands.
	PollDeviceGrant(device_code string) (*TokenPollModel, error)
	// Whoami returns who the server takes the current token for.
	Whoami() (*PrincipalModel, error)
	// Revoke revokes the current token on the server.
	Revoke() error
}

// AuthRepository is the [AuthRepositoryInterface] over the engine's HTTP API.
type AuthRepository struct {
	// Client is the shared HTTP client.
	Client *api.Client
}

// Use points Client at url with token.
func (r AuthRepository) Use(url string, token string) {
	r.Client.Configure(url, token)
}

// Discover gets /.well-known/warp-engine; a 404 is an older engine and
// returns an empty ServerModel.
func (r AuthRepository) Discover() (*ServerModel, error) {
	server := ServerModel{}
	err := r.Client.Do(http.MethodGet, "/.well-known/warp-engine", nil, nil, &server)
	if api.StatusOf(err) == http.StatusNotFound {
		return &ServerModel{}, nil
	}
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// StartDeviceGrant posts /api/auth/device.
func (r AuthRepository) StartDeviceGrant(client_name string, scopes []string) (*DeviceGrantModel, error) {
	grant := DeviceGrantModel{}
	body := map[string]any{"client_name": client_name, "scopes": scopes}
	if err := r.Client.Do(http.MethodPost, "/api/auth/device", nil, body, &grant); err != nil {
		return nil, err
	}
	return &grant, nil
}

// PollDeviceGrant posts /api/auth/device/token.
func (r AuthRepository) PollDeviceGrant(device_code string) (*TokenPollModel, error) {
	poll := TokenPollModel{}
	body := map[string]any{"device_code": device_code}
	if err := r.Client.Do(http.MethodPost, "/api/auth/device/token", nil, body, &poll); err != nil {
		return nil, err
	}
	return &poll, nil
}

// Whoami gets the Admin API's /me.
func (r AuthRepository) Whoami() (*PrincipalModel, error) {
	return api.Get[PrincipalModel](r.Client, "/me", nil)
}

// Revoke sends DELETE /api/auth/token.
func (r AuthRepository) Revoke() error {
	return r.Client.Do(http.MethodDelete, "/api/auth/token", nil, nil, nil)
}

// LocalTokenRepositoryInterface reads the token of WarpEngine Local.
type LocalTokenRepositoryInterface interface {
	// Read returns the token, or [ErrNoLocalToken] when there is none.
	Read() (string, error)
}

// ErrNoLocalToken is returned when WarpEngine Local's token file is missing:
// the local engine is not installed for this user.
var ErrNoLocalToken = errors.New("no WarpEngine Local token file: is WarpEngine Local installed?")

// LocalTokenRepository reads the token file WarpEngine Local writes into its
// data directory.
type LocalTokenRepository struct {
	// Directory is the local engine's data directory; see [DefaultLocalDirectory].
	Directory string
}

// DefaultLocalDirectory returns WarpEngine Local's data directory: the
// WARP_ENGINE_LOCAL_DIR environment variable when set, otherwise
// ~/.warpengine.
func DefaultLocalDirectory() string {
	if directory := os.Getenv("WARP_ENGINE_LOCAL_DIR"); directory != "" {
		return directory
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".warpengine"
	}
	return filepath.Join(home, ".warpengine")
}

// Read returns the trimmed content of Directory/local-token.
func (r LocalTokenRepository) Read() (string, error) {
	data, err := os.ReadFile(filepath.Join(r.Directory, "local-token"))
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrNoLocalToken
	}
	if err != nil {
		return "", err
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", ErrNoLocalToken
	}
	return token, nil
}
