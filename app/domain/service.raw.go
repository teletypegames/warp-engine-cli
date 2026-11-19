package domain

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// RawMethods are the HTTP methods [RawService.Send] accepts.
var RawMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// RawServiceInterface sends requests the commands do not cover.
type RawServiceInterface interface {
	// Send sends method (any case) to path with query and body and returns
	// the answer.
	Send(method string, path string, query url.Values, body any) (json.RawMessage, error)
}

// RawService is the [RawServiceInterface] over a [RawRepositoryInterface].
type RawService struct {
	// RawRepository sends the request.
	RawRepository RawRepositoryInterface
}

// Send checks that method is one of [RawMethods] and that path is not empty,
// then sends the request.
func (s RawService) Send(method string, path string, query url.Values, body any) (json.RawMessage, error) {
	method = strings.ToUpper(method)
	known := false
	for _, allowed := range RawMethods {
		known = known || allowed == method
	}
	if !known {
		return nil, fmt.Errorf("%q is no method this program sends (%s)", method, strings.Join(RawMethods, ", "))
	}
	if strings.Trim(path, "/") == "" {
		return nil, fmt.Errorf("no path given")
	}
	return s.RawRepository.Send(method, path, query, body)
}
