package domain

import (
	"net/url"
	"strconv"
)

// PageModel is one answer of a list operation: the records and the paging
// of the page they came from.
type PageModel[T any] struct {
	// Items are the records, in the order the server sent them.
	Items []T
	// Page is the number of the (last) page read, from 1; 0 for a list
	// without paging.
	Page int
	// PerPage is how many records a page holds; 0 for a list without paging.
	PerPage int
	// Total is how many records match on all pages; len(Items) for a list
	// without paging.
	Total int
	// TotalPages is how many pages there are; 1 for a list without paging.
	TotalPages int
}

// ListQuery is the paging and order every paged list takes. The zero value
// asks for the server's first page in its default order.
type ListQuery struct {
	// Page is the page to read, from 1; 0 leaves it to the server.
	Page int
	// PerPage is how many records a page should hold; 0 leaves it to the
	// server, which caps it anyway.
	PerPage int
	// Sort is the order, a field name with an optional leading "-" for
	// descending ("-createdAt"); empty keeps the default.
	Sort string
	// All reads every page from Page on and returns the records of all.
	All bool
}

// Values returns the query parameters of q: page, perPage and sort, each
// only when set.
func (q ListQuery) Values() url.Values {
	values := url.Values{}
	if q.Page > 0 {
		values.Set("page", strconv.Itoa(q.Page))
	}
	if q.PerPage > 0 {
		values.Set("perPage", strconv.Itoa(q.PerPage))
	}
	if q.Sort != "" {
		values.Set("sort", q.Sort)
	}
	return values
}

// Filter is a set of query parameters on top of a [ListQuery]: the
// "name contains", "platform is" … conditions of a list. The keys are the
// API's parameter names (camelCase); an empty value is not sent.
type Filter map[string]string

// Values returns the query of query with every non-empty entry of f added.
func (f Filter) Values(query ListQuery) url.Values {
	values := query.Values()
	for key, value := range f {
		if value != "" {
			values.Set(key, value)
		}
	}
	return values
}
