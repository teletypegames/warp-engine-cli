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

func TestOutletRepository_All_ReadsTheOutlets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v1/outlets", r.URL.Path)
		io.WriteString(w, `{"data":[{"key":"wecp","type":"wecp","platforms":[],"description":"wecp (wecp)","materialized":false,"served":true,"details":{"version":4}},`+
			`{"key":"assembly64","type":"assembly64","platforms":["c64"],"description":"assembly64 (assembly64: c64)","materialized":false,"served":true,"details":{"hosts":["a64.example.org"]}}]}`)
	}))
	defer server.Close()
	repo := domain.OutletRepository{Client: api.New(server.URL, "", "test")}

	outlets, err := repo.All()

	assert.NoError(t, err)
	assert.Equal(t, []string{"wecp", "assembly64"}, []string{outlets[0].Key, outlets[1].Key})
	assert.Equal(t, []string{"c64"}, outlets[1].Platforms)
	assert.Equal(t, []any{"a64.example.org"}, outlets[1].Details["hosts"])
}

func TestOutletRepository_Check_EscapesTheKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/v1/outlets/c64/check", r.URL.Path)
		io.WriteString(w, `{"data":{"key":"c64","type":"shelf","offer":true,"items":[{"softwareId":1,"name":"ariadne","title":"Ariadne","version":"1.0","kind":"cartridge"}],"skipped":[{"softwareId":2,"name":"empty","reason":"no cartridge"}]}}`)
	}))
	defer server.Close()
	repo := domain.OutletRepository{Client: api.New(server.URL, "", "test")}

	check, err := repo.Check("c64")

	assert.NoError(t, err)
	assert.True(t, check.Offer)
	assert.Equal(t, "ariadne", check.Items[0].Name)
	assert.Equal(t, "no cartridge", check.Skipped[0].Reason)
}
