package mockrepo

import (
	"encoding/json"
	"net/url"
)

// MockRawRepository mocks [domain.RawRepositoryInterface].
type MockRawRepository struct {
	// SendFunc replaces Send.
	SendFunc func(method string, path string, query url.Values, body any) (json.RawMessage, error)
}

// Send calls SendFunc, or returns nil.
func (m *MockRawRepository) Send(method string, path string, query url.Values, body any) (json.RawMessage, error) {
	if m.SendFunc != nil {
		return m.SendFunc(method, path, query, body)
	}
	return nil, nil
}
