package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// UpdateHeader is the header the build endpoints (/build/upload,
// /build/publish) read their credential from: an application token with the
// upload or update scope, sent instead of the bearer token.
const UpdateHeader = "X-Update-Secret"

// UploadTimeout is how long one upload to a build endpoint may take, the
// transfer included: a release file can be far larger than an Admin API
// request.
const UploadTimeout = time.Hour

// BuildUpload posts the file at file_path to the engine path (absolute,
// "/build/upload") as the multipart part "file", with the query parameters
// query, signed with the client's Token in [UpdateHeader]. The file is
// streamed, not read into memory. The answer is decoded into out like
// [Client.Do] does.
func (c *Client) BuildUpload(path string, query url.Values, file_path string, out any) error {
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}

	var head bytes.Buffer
	writer := multipart.NewWriter(&head)
	if _, err = writer.CreateFormFile("file", filepath.Base(file_path)); err != nil {
		return err
	}
	prefix := append([]byte(nil), head.Bytes()...)
	head.Reset()
	if err = writer.Close(); err != nil {
		return err
	}
	suffix := head.Bytes()

	body := io.MultiReader(bytes.NewReader(prefix), file, bytes.NewReader(suffix))
	request, err := c.prepare(http.MethodPost, path, query, body, writer.FormDataContentType())
	if err != nil {
		return err
	}
	request.ContentLength = int64(len(prefix)) + info.Size() + int64(len(suffix))
	c.sign(request)
	client := *c.HTTP
	client.Timeout = UploadTimeout
	return c.exchange(&client, request, http.MethodPost, path, out)
}

// BuildPost posts to the engine path (absolute, "/build/publish") with the
// query parameters query and no body, signed with the client's Token in
// [UpdateHeader]. The answer is decoded into out like [Client.Do] does.
func (c *Client) BuildPost(path string, query url.Values, out any) error {
	request, err := c.prepare(http.MethodPost, path, query, nil, "")
	if err != nil {
		return err
	}
	c.sign(request)
	return c.exchange(c.HTTP, request, http.MethodPost, path, out)
}

// sign puts the client's Token in [UpdateHeader] of request.
func (c *Client) sign(request *http.Request) {
	if c.Token != "" {
		request.Header.Set(UpdateHeader, c.Token)
	}
}
