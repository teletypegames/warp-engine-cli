package repository_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestBuildPlatformRepository_All_ReadsThePublicDescription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/build/platforms", r.URL.Path)
		io.WriteString(w, `{"platforms":[{"name":"golang","label":"Go","kinds":[{"name":"cli_linux_x64","required":true}],`+
			`"apps":[{"name":"cli","default":true,"kinds":[{"name":"cli_linux_x64"}]},{"name":"gui","default":false,"kinds":[{"name":"linux_x64"}]}]}]}`)
	}))
	defer server.Close()
	repo := domain.BuildPlatformRepository{Client: api.New(server.URL, "", "test")}

	platforms, err := repo.All()

	assert.NoError(t, err)
	assert.Equal(t, []domain.BuildPlatformModel{{
		Name: "golang", Label: "Go", Kinds: []domain.BuildKindModel{{Name: "cli_linux_x64"}},
		Apps: []domain.BuildAppModel{
			{Name: "cli", Default: true, Kinds: []domain.BuildKindModel{{Name: "cli_linux_x64"}}},
			{Name: "gui", Kinds: []domain.BuildKindModel{{Name: "linux_x64"}}},
		},
	}}, platforms)
}

func TestBuildPlatformRepository_All_ReturnsTheError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	repo := domain.BuildPlatformRepository{Client: api.New(server.URL, "", "test")}

	_, err := repo.All()

	assert.Error(t, err)
}
