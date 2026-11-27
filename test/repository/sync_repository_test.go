package repository_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

// syncRequest is one request a test server saw.
type syncRequest struct {
	method string
	path   string
	query  string
	body   map[string]any
}

// serveSync starts a server that records each request and answers answer with
// status; the caller closes it.
func serveSync(t *testing.T, status int, answer string, seen *syncRequest) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen.method, seen.path, seen.query = r.Method, r.URL.Path, r.URL.RawQuery
		seen.body = nil
		if r.ContentLength > 0 {
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&seen.body))
		}
		w.WriteHeader(status)
		io.WriteString(w, answer)
	}))
}

func TestSyncPeerRepository_All_SendsFilters(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":[{"id":3,"label":"games","transportOptions":{"host":"h"}}],"meta":{"page":1,"perPage":25,"total":1,"totalPages":1,"nextPage":null}}`, &seen)
	defer server.Close()
	repo := domain.SyncPeerRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"outboundEnabled": "true", "transport": "http"}, domain.ListQuery{Sort: "label"})

	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/sync_peers", seen.path)
	assert.Equal(t, "outboundEnabled=true&sort=label&transport=http", seen.query)
	assert.Equal(t, "h", page.Items[0].TransportOptions["host"])
}

func TestSyncPeerRepository_Create_PostsAttributes(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 201, `{"data":{"id":4,"label":"mirror","allowedChannels":["stable"]}}`, &seen)
	defer server.Close()
	repo := domain.SyncPeerRepository{Client: api.New(server.URL, "", "test")}

	peer, err := repo.Create(domain.Attributes{"label": "mirror", "transportOptions": map[string]any{"host": "m"}})

	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, seen.method)
	assert.Equal(t, "/api/admin/v2/sync_peers", seen.path)
	assert.Equal(t, map[string]any{"label": "mirror", "transportOptions": map[string]any{"host": "m"}}, seen.body)
	assert.Equal(t, []string{"stable"}, peer.AllowedChannels)
}

func TestSyncPeerRepository_Actions_HitTheirPaths(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"ok":true,"status":200,"version":"5.5.0","message":"ok","tokenPrefix":"ab","plainToken":"ab-whole","dryRun":true,"queued":0,"peers":[3],"releases":[{"software":"pong","version":"1.0"}],"skipped":[]}}`, &seen)
	defer server.Close()
	repo := domain.SyncPeerRepository{Client: api.New(server.URL, "", "test")}

	check, err := repo.Check(3)
	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/sync_peers/3/check", seen.path)
	assert.True(t, check.OK)

	token, err := repo.InboundToken(3)
	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/sync_peers/3/inbound_token", seen.path)
	assert.Equal(t, "ab-whole", token.PlainToken)

	push, err := repo.Push(3, domain.Attributes{"dryRun": true})
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, seen.method)
	assert.Equal(t, "/api/admin/v2/sync_peers/3/push", seen.path)
	assert.Equal(t, map[string]any{"dryRun": true}, seen.body)
	assert.Equal(t, "pong", push.Releases[0].Software)

	_, err = repo.Run(domain.Attributes{"peerIds": []int{3, 5}})
	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/sync_peers/run", seen.path)
	assert.Equal(t, map[string]any{"peerIds": []any{float64(3), float64(5)}}, seen.body)
}

func TestSyncPeerRepository_Destroy_SendsDelete(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 204, ``, &seen)
	defer server.Close()
	repo := domain.SyncPeerRepository{Client: api.New(server.URL, "", "test")}

	err := repo.Destroy(3)

	assert.NoError(t, err)
	assert.Equal(t, http.MethodDelete, seen.method)
	assert.Equal(t, "/api/admin/v2/sync_peers/3", seen.path)
}

func TestSyncPairingRepository_TokenAndUpdate(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 201, `{"data":{"id":2,"label":"laptop","inboundTokenPrefix":"9f","plainToken":"9f-whole","lastInboundAt":null}}`, &seen)
	defer server.Close()
	repo := domain.SyncPairingRepository{Client: api.New(server.URL, "", "test")}

	pairing, err := repo.Token(2)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, seen.method)
	assert.Equal(t, "/api/admin/v2/sync_pairings/2/token", seen.path)
	assert.Equal(t, "9f-whole", *pairing.PlainToken)

	_, err = repo.Update(2, domain.Attributes{"label": "desk"})
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPatch, seen.method)
	assert.Equal(t, "/api/admin/v2/sync_pairings/2", seen.path)
	assert.Equal(t, map[string]any{"label": "desk"}, seen.body)
}

func TestSyncPairingRepository_All_InboundOffIs409(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 409, `{"error":"inbound sync is off here"}`, &seen)
	defer server.Close()
	repo := domain.SyncPairingRepository{Client: api.New(server.URL, "", "test")}

	_, err := repo.All(domain.Filter{}, domain.ListQuery{})

	assert.Equal(t, http.StatusConflict, api.StatusOf(err))
	assert.ErrorContains(t, err, "inbound sync is off")
}

func TestSyncEventRepository_SoftwareViews(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":[{"peerId":3,"peer":"games","direction":"outbound","reason":null,"event":{"id":9,"status":"succeeded"}}]}`, &seen)
	defer server.Close()
	repo := domain.SyncEventRepository{Client: api.New(server.URL, "", "test")}

	states, err := repo.Status(12)
	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/softwares/12/sync_status", seen.path)
	assert.Equal(t, "succeeded", states[0].Event.Status)

	_, err = repo.AllOfSoftware(12, "inbound")
	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/softwares/12/sync_events", seen.path)
	assert.Equal(t, "direction=inbound", seen.query)

	_, err = repo.AllOfSoftware(12, "")
	assert.NoError(t, err)
	assert.Equal(t, "", seen.query)
}

func TestSyncEventRepository_FindByID(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":{"id":9,"direction":"inbound","status":"failed","bytes":1024}}`, &seen)
	defer server.Close()
	repo := domain.SyncEventRepository{Client: api.New(server.URL, "", "test")}

	event, err := repo.FindByID(9)

	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/sync_events/9", seen.path)
	assert.Equal(t, int64(1024), *event.Bytes)
}

func TestApplicationTokenRepository_CreateAndRotate(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 201, `{"data":{"id":5,"name":"ci","tokenPrefix":"aa","ownerType":"Owner","ownerId":1,"scopes":["admin:read"],"unrestricted":false,"plainToken":"aa-whole","sync":{"synced":1,"failed":0}}}`, &seen)
	defer server.Close()
	repo := domain.ApplicationTokenRepository{Client: api.New(server.URL, "", "test")}

	token, err := repo.Create(domain.Attributes{"name": "ci", "scopes": []string{"admin:read"}, "expiresAt": nil})
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, seen.method)
	assert.Equal(t, "/api/admin/v2/application_tokens", seen.path)
	assert.Equal(t, map[string]any{"name": "ci", "scopes": []any{"admin:read"}, "expiresAt": nil}, seen.body)
	assert.Equal(t, "aa-whole", *token.PlainToken)
	assert.Equal(t, 1, token.Sync.Synced)

	_, err = repo.Rotate(5)
	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/application_tokens/5/rotate", seen.path)
}

func TestApplicationTokenRepository_All_SendsFilters(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 200, `{"data":[],"meta":{"page":1,"perPage":25,"total":0,"totalPages":0,"nextPage":null}}`, &seen)
	defer server.Close()
	repo := domain.ApplicationTokenRepository{Client: api.New(server.URL, "", "test")}

	page, err := repo.All(domain.Filter{"kind": "admin", "state": ""}, domain.ListQuery{})

	assert.NoError(t, err)
	assert.Equal(t, "/api/admin/v2/application_tokens", seen.path)
	assert.Equal(t, "kind=admin", seen.query)
	assert.Empty(t, page.Items)
}

func TestApplicationTokenRepository_Destroy_SelfIs409(t *testing.T) {
	var seen syncRequest
	server := serveSync(t, 409, `{"error":"A token cannot revoke itself here"}`, &seen)
	defer server.Close()
	repo := domain.ApplicationTokenRepository{Client: api.New(server.URL, "", "test")}

	err := repo.Destroy(1)

	assert.Equal(t, http.MethodDelete, seen.method)
	assert.Equal(t, "/api/admin/v2/application_tokens/1", seen.path)
	assert.Equal(t, http.StatusConflict, api.StatusOf(err))
}
