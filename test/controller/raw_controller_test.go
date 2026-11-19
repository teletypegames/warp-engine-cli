package controller_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/service"
	"github.com/stretchr/testify/assert"
)

func TestRawController_Send_PassesQueryAndBodyAndIndents(t *testing.T) {
	var got_query url.Values
	var got_body any
	svc := &mocksvc.MockRawService{SendFunc: func(method string, path string, query url.Values, body any) (json.RawMessage, error) {
		got_query, got_body = query, body
		return json.RawMessage(`{"data":{"id":1}}`), nil
	}}

	stdout, _, err := execute(cli.Router{RawController: cli.RawController{RawService: svc}},
		"api", "PATCH", "/softwares/1", "--query", "a=1", "--query", "a=2", "--data", `{"title":"X"}`)

	assert.NoError(t, err)
	assert.Equal(t, url.Values{"a": {"1", "2"}}, got_query)
	assert.Equal(t, map[string]any{"title": "X"}, got_body)
	assert.Contains(t, stdout, "\n  \"data\"")
}

func TestRawController_Send_BadQueryIsAUsageError(t *testing.T) {
	_, _, err := execute(cli.Router{RawController: cli.RawController{RawService: &mocksvc.MockRawService{}}}, "api", "GET", "/x", "--query", "novalue")

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}
