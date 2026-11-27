package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AdminPrefix is the path the Admin API is mounted at on every WarpEngine
// server; [Client.Admin] puts it in front of the paths it is given.
const AdminPrefix = "/api/admin/v2"

// DefaultTimeout is how long one request may take, the transfer included,
// before [New] gives up on it. An upload of a large file may need more: set
// [Client.HTTP] to a client with a longer timeout.
const DefaultTimeout = 5 * time.Minute

// Client sends requests to one WarpEngine server. The zero value is not
// usable: build one with [New]. The fields may be changed between requests
// (the command line sets BaseURL and Token once it knows the profile), but not
// while a request runs.
type Client struct {
	// BaseURL is the server's address without a trailing slash, for example
	// https://games.example.org.
	BaseURL string
	// Token is the bearer token sent with every request; empty sends none.
	Token string
	// UserAgent is sent as the User-Agent header.
	UserAgent string
	// HTTP is the underlying HTTP client.
	HTTP *http.Client
}

// New returns a Client for the server at base_url, signed in with token (may
// be empty), with [DefaultTimeout] and the User-Agent user_agent.
func New(base_url string, token string, user_agent string) *Client {
	return &Client{
		BaseURL:   strings.TrimRight(base_url, "/"),
		Token:     token,
		UserAgent: user_agent,
		HTTP:      &http.Client{Timeout: DefaultTimeout},
	}
}

// Configure points the client at base_url and signs it in with token. The
// trailing slash of base_url is dropped.
func (c *Client) Configure(base_url string, token string) {
	c.BaseURL = strings.TrimRight(base_url, "/")
	c.Token = token
}

// Do sends method to the engine path (absolute, "/api/auth/device"), with the
// query parameters query (may be nil) and body encoded as JSON (nil sends no
// body). A 2xx answer with a body is decoded into out, unless out is nil. Any
// other status returns an [*Error]; a missing BaseURL returns [ErrNoServer].
func (c *Client) Do(method string, path string, query url.Values, body any, out any) error {
	var reader io.Reader
	content_type := ""
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode the request body: %w", err)
		}
		reader = bytes.NewReader(encoded)
		content_type = "application/json"
	}
	return c.send(method, path, query, reader, content_type, out)
}

// Admin is [Client.Do] for an Admin API path: path is relative to
// [AdminPrefix] ("/softwares/12").
func (c *Client) Admin(method string, path string, query url.Values, body any, out any) error {
	return c.Do(method, AdminPrefix+path, query, body, out)
}

// Upload posts a multipart/form-data request to the Admin API path: the file
// at file_path as the part file_field, and every entry of fields as a plain
// part. The answer is decoded into out like [Client.Do] does.
func (c *Client) Upload(path string, fields map[string]string, file_field string, file_path string, out any) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	for key, value := range fields {
		if err = writer.WriteField(key, value); err != nil {
			return err
		}
	}
	part, err := writer.CreateFormFile(file_field, filepath.Base(file_path))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}
	if err = writer.Close(); err != nil {
		return err
	}
	return c.send(http.MethodPost, AdminPrefix+path, nil, &buffer, writer.FormDataContentType(), out)
}

// send builds and runs one request with the bearer token and reads its
// answer; see [Client.Do].
func (c *Client) send(method string, path string, query url.Values, body io.Reader, content_type string, out any) error {
	request, err := c.prepare(method, path, query, body, content_type)
	if err != nil {
		return err
	}
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return c.exchange(c.HTTP, request, method, path, out)
}

// prepare builds a request to the engine path with the headers every
// request has; [ErrNoServer] without a BaseURL.
func (c *Client) prepare(method string, path string, query url.Values, body io.Reader, content_type string) (*http.Request, error) {
	if c.BaseURL == "" {
		return nil, ErrNoServer
	}
	address := c.BaseURL + path
	if len(query) > 0 {
		address += "?" + query.Encode()
	}
	request, err := http.NewRequest(method, address, body)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", method, path, err)
	}
	request.Header.Set("Accept", "application/json")
	if content_type != "" {
		request.Header.Set("Content-Type", content_type)
	}
	if c.UserAgent != "" {
		request.Header.Set("User-Agent", c.UserAgent)
	}
	return request, nil
}

// exchange runs request with client and reads its answer into out; see
// [Client.Do].
func (c *Client) exchange(client *http.Client, request *http.Request, method string, path string, out any) error {
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("%s %s: read the answer: %w", method, path, err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return newError(method, path, response, payload)
	}
	if out == nil || len(bytes.TrimSpace(payload)) == 0 {
		return nil
	}
	if err = json.Unmarshal(payload, out); err != nil {
		return fmt.Errorf("%s %s: the answer is no JSON this program understands: %w", method, path, err)
	}
	return nil
}
