package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ErrNoServer is returned by every request of a [Client] that has no BaseURL:
// the command line has no profile and no --server yet.
var ErrNoServer = errors.New("no server: sign in with \"warp-engine login <url>\" or pass --server")

// Error is an answer of the server outside 2xx.
type Error struct {
	// Method and Path are the request that failed ("PATCH", "/api/admin/v1/softwares/3").
	Method string
	// Path is the path of the request, without the query.
	Path string
	// Status is the HTTP status code, for example 404.
	Status int
	// Message is the server's "error" text, or the status text when the
	// answer had none.
	Message string
	// Details lists, per field, why a record did not validate (422), or is nil.
	Details map[string][]string
	// RetryAfter is how long to wait after a 429, from the Retry-After header;
	// zero when the server did not say.
	RetryAfter time.Duration
}

// Error renders the status, the message and the details on one line:
// "404 Not Found: Software not found".
func (e *Error) Error() string {
	text := fmt.Sprintf("%d %s: %s", e.Status, http.StatusText(e.Status), e.Message)
	if len(e.Details) == 0 {
		return text
	}
	fields := make([]string, 0, len(e.Details))
	for field := range e.Details {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		parts = append(parts, field+" "+strings.Join(e.Details[field], ", "))
	}
	return text + " (" + strings.Join(parts, "; ") + ")"
}

// StatusOf returns the HTTP status of err when it is, or wraps, an [*Error],
// and 0 otherwise.
func StatusOf(err error) int {
	var api_error *Error
	if errors.As(err, &api_error) {
		return api_error.Status
	}
	return 0
}

// newError builds the [*Error] of a non-2xx response whose body is payload.
func newError(method string, path string, response *http.Response, payload []byte) *Error {
	result := &Error{Method: method, Path: path, Status: response.StatusCode}
	var body struct {
		Error   string              `json:"error"`
		Details map[string][]string `json:"details"`
	}
	if json.Unmarshal(payload, &body) == nil && body.Error != "" {
		result.Message = body.Error
		result.Details = body.Details
	} else {
		result.Message = http.StatusText(response.StatusCode)
	}
	if seconds, err := strconv.Atoi(response.Header.Get("Retry-After")); err == nil {
		result.RetryAfter = time.Duration(seconds) * time.Second
	}
	return result
}
