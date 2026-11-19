package repository_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepository_Discover_OldEngineIsEmpty(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	repo := domain.AuthRepository{Client: api.New(server.URL, "", "test")}

	answer, err := repo.Discover()

	assert.NoError(t, err)
	assert.Equal(t, domain.ServerModel{}, *answer)
}

func TestAuthRepository_StartDeviceGrant_SendsClientNameAndScopes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/auth/device", r.URL.Path)
		body := map[string]any{}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "cli", body["client_name"])
		assert.Equal(t, []any{"admin:read"}, body["scopes"])
		io.WriteString(w, `{"deviceCode":"d","userCode":"ABCD-1234","verificationUrl":"v","verificationUrlComplete":"v?code=ABCD-1234","scopes":["admin:read"],"interval":5,"expiresIn":600}`)
	}))
	defer server.Close()
	repo := domain.AuthRepository{Client: api.New(server.URL, "", "test")}

	grant, err := repo.StartDeviceGrant("cli", []string{"admin:read"})

	assert.NoError(t, err)
	assert.Equal(t, "ABCD-1234", grant.UserCode)
	assert.Equal(t, 600, grant.ExpiresIn)
}

func TestLocalTokenRepository_Read_TrimsAndMissing(t *testing.T) {
	directory := t.TempDir()
	repo := domain.LocalTokenRepository{Directory: directory}

	_, err := repo.Read()
	assert.ErrorIs(t, err, domain.ErrNoLocalToken)

	assert.NoError(t, os.WriteFile(filepath.Join(directory, "local-token"), []byte("abc\n"), 0o600))
	token, err := repo.Read()
	assert.NoError(t, err)
	assert.Equal(t, "abc", token)
}

func TestProfileRepository_SaveLoad_RoundTripWithMode0600(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "config.json")
	repo := domain.ProfileRepository{Path: path}
	config := &domain.ConfigModel{Current: "a", Profiles: map[string]domain.ProfileModel{"a": {Name: "a", URL: "u", Token: "t"}}}

	assert.NoError(t, repo.Save(config))
	loaded, err := repo.Load()

	assert.NoError(t, err)
	assert.Equal(t, config, loaded)
	info, _ := os.Stat(path)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestProfileRepository_Load_MissingIsEmptyBrokenIsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	repo := domain.ProfileRepository{Path: path}

	config, err := repo.Load()
	assert.NoError(t, err)
	assert.Empty(t, config.Profiles)

	assert.NoError(t, os.WriteFile(path, []byte("{nope"), 0o600))
	_, err = repo.Load()
	var config_error *domain.ConfigError
	assert.ErrorAs(t, err, &config_error)
}

func TestRawRepository_Send_PathRouting(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		io.WriteString(w, `{}`)
	}))
	defer server.Close()
	repo := domain.RawRepository{Client: api.New(server.URL, "", "test")}

	for _, path := range []string{"/softwares", "me", "/api/auth/device", "/.well-known/warp-engine"} {
		_, err := repo.Send(http.MethodGet, path, nil, nil)
		assert.NoError(t, err)
	}

	assert.Equal(t, []string{"/api/admin/v1/softwares", "/api/admin/v1/me", "/api/auth/device", "/.well-known/warp-engine"}, paths)
}
