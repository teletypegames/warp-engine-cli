package mocksvc

import (
	"encoding/json"
	"net/url"
)

// MockRawService mocks [domain.RawServiceInterface].
type MockRawService struct {
	// SendFunc replaces Send.
	SendFunc func(method string, path string, query url.Values, body any) (json.RawMessage, error)
}

// Send calls SendFunc, or returns nil.
func (m *MockRawService) Send(method string, path string, query url.Values, body any) (json.RawMessage, error) {
	if m.SendFunc != nil {
		return m.SendFunc(method, path, query, body)
	}
	return nil, nil
}
