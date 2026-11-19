package service_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"git.teletypegames.org/engines/warp-engine-cli/test/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestRawService_Send_UppercasesTheMethod(t *testing.T) {
	var got string
	repo := &mockrepo.MockRawRepository{SendFunc: func(method string, path string, query url.Values, body any) (json.RawMessage, error) {
		got = method
		return json.RawMessage(`{}`), nil
	}}

	_, err := domain.RawService{RawRepository: repo}.Send("get", "/softwares", nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "GET", got)
}

func TestRawService_Send_RejectsUnknownMethodAndEmptyPath(t *testing.T) {
	svc := domain.RawService{RawRepository: &mockrepo.MockRawRepository{}}

	_, err := svc.Send("TRACE", "/x", nil, nil)
	assert.Error(t, err)

	_, err = svc.Send("GET", "/", nil, nil)
	assert.Error(t, err)
}
