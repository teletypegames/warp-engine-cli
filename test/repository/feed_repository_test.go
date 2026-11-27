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

func TestFeedRepository_All_ReadsTheFeeds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/feeds", r.URL.Path)
		io.WriteString(w, `{"data":[{"key":"wecp","type":"wecp","platforms":[],"description":"wecp (wecp)","static":false,"dynamic":true,"details":{"version":4}},`+
			`{"key":"assembly64","type":"assembly64","platforms":["c64"],"description":"assembly64 (assembly64: c64)","static":false,"dynamic":true,"details":{"hosts":["a64.example.org"]}}]}`)
	}))
	defer server.Close()
	repo := domain.FeedRepository{Client: api.New(server.URL, "", "test")}

	feeds, err := repo.All()

	assert.NoError(t, err)
	assert.Equal(t, []string{"wecp", "assembly64"}, []string{feeds[0].Key, feeds[1].Key})
	assert.Equal(t, []string{"c64"}, feeds[1].Platforms)
	assert.Equal(t, []any{"a64.example.org"}, feeds[1].Details["hosts"])
}

func TestFeedRepository_Check_EscapesTheKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v2/feeds/c64/check", r.URL.Path)
		io.WriteString(w, `{"data":{"key":"c64","type":"directory","selection":true,"items":[{"softwareId":1,"name":"ariadne","title":"Ariadne","version":"1.0","kind":"package"}],"skipped":[{"softwareId":2,"name":"empty","reason":"no cartridge"}]}}`)
	}))
	defer server.Close()
	repo := domain.FeedRepository{Client: api.New(server.URL, "", "test")}

	check, err := repo.Check("c64")

	assert.NoError(t, err)
	assert.True(t, check.Selection)
	assert.Equal(t, "ariadne", check.Items[0].Name)
	assert.Equal(t, "no cartridge", check.Skipped[0].Reason)
}
