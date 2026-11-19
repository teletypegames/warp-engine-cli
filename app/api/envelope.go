package api

import (
	"net/http"
	"net/url"
	"strconv"
)

// Envelope is the wrapper of every Admin API answer: Data holds the record
// or the list, Meta the paging of a list (nil for a single record).
type Envelope[T any] struct {
	// Data is the payload: one record, or a slice of them.
	Data T `json:"data"`
	// Meta is the paging of a list answer.
	Meta *PageMeta `json:"meta,omitempty"`
}

// PageMeta is the paging block of a list answer.
type PageMeta struct {
	// Page is the number of this page, from 1.
	Page int `json:"page"`
	// PerPage is how many records a page holds.
	PerPage int `json:"perPage"`
	// Total is how many records match, on all pages together.
	Total int `json:"total"`
	// TotalPages is how many pages there are.
	TotalPages int `json:"totalPages"`
	// NextPage is the number of the next page, or nil on the last one.
	NextPage *int `json:"nextPage"`
}

// List gets the Admin API list at path with query and returns its records
// and the paging of the last page read. With all set it reads every page,
// following NextPage from the page in query (or the first), and returns the
// records of all of them in order. Some lists (the downloads summary, the
// local components) have no paging: their meta is nil.
func List[T any](c *Client, path string, query url.Values, all bool) ([]T, *PageMeta, error) {
	values := url.Values{}
	for key, value := range query {
		values[key] = value
	}
	items := []T{}
	for {
		var envelope Envelope[[]T]
		if err := c.Admin(http.MethodGet, path, values, nil, &envelope); err != nil {
			return nil, nil, err
		}
		items = append(items, envelope.Data...)
		if !all || envelope.Meta == nil || envelope.Meta.NextPage == nil {
			return items, envelope.Meta, nil
		}
		values.Set("page", strconv.Itoa(*envelope.Meta.NextPage))
	}
}

// Get gets the Admin API record at path, with query (may be nil), and
// returns its data.
func Get[T any](c *Client, path string, query url.Values) (*T, error) {
	var envelope Envelope[T]
	if err := c.Admin(http.MethodGet, path, query, nil, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

// Send sends method (POST, PATCH, PUT) to the Admin API path with body as
// JSON and returns the data of the answer. A 204 answer returns a zero T.
func Send[T any](c *Client, method string, path string, body any) (*T, error) {
	var envelope Envelope[T]
	if err := c.Admin(method, path, nil, body, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}

// Delete sends DELETE to the Admin API path with query (may be nil) and
// expects no body back.
func Delete(c *Client, path string, query url.Values) error {
	return c.Admin(http.MethodDelete, path, query, nil, nil)
}
