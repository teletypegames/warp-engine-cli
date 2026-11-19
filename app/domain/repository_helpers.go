package domain

import (
	"net/http"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
)

// list reads the Admin API list at path with filter and query into a
// [PageModel]; with query.All it reads every page. A list without paging
// answers a page with Total = len(Items) and TotalPages = 1.
func list[T any](client *api.Client, path string, filter Filter, query ListQuery) (*PageModel[T], error) {
	items, meta, err := api.List[T](client, path, filter.Values(query), query.All)
	if err != nil {
		return nil, err
	}
	page := &PageModel[T]{Items: items, Total: len(items), TotalPages: 1}
	if meta != nil {
		page.Page = meta.Page
		page.PerPage = meta.PerPage
		page.Total = meta.Total
		page.TotalPages = meta.TotalPages
	}
	return page, nil
}

// notFound returns the [*api.Error] of a record this program looked for on
// the client side (a name lookup) and did not find: status 404 with message.
func notFound(path string, message string) error {
	return &api.Error{Method: http.MethodGet, Path: api.AdminPrefix + path, Status: http.StatusNotFound, Message: message}
}
