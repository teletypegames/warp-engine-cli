// Package api is the HTTP client of a WarpEngine server: the one place that
// speaks HTTP, the way the database helper library is the one place that
// speaks SQL in a server application. The repositories of package domain
// hold a [Client] and call it; nothing else in the program opens a connection.
//
// # Servers and paths
//
// A [Client] talks to one server, given by its base URL
// (https://games.example.org, or http://127.0.0.1:7470 for WarpEngine Local).
// Two families of paths hang off it:
//
//   - the Admin API under [AdminPrefix] (/api/admin/v1), reached with
//     [Client.Admin], [List], [Get], [Send] and [Delete];
//   - the engine's own endpoints (/.well-known/warp-engine, /api/auth/…),
//     reached with [Client.Do].
//
// Both answer JSON. A request with a body sends it as JSON, except an upload
// ([Client.Upload]), which is multipart/form-data.
//
// # Authentication
//
// When [Client.Token] is set, every request carries it as a bearer token:
//
//	Authorization: Bearer <token>
//
// The token is an application token of the engine, issued by a device sign-in
// or, on WarpEngine Local, read from the token file. The client does not care
// where it came from.
//
// # The Admin API envelope
//
// Every Admin API answer wraps its payload in an envelope, [Envelope]: a
// single record is
//
//	{"data": {…}}
//
// and a list is
//
//	{"data": [{…}, …], "meta": {"page": 1, "perPage": 25, "total": 60, "totalPages": 3, "nextPage": 2}}
//
// [List] follows the nextPage links when asked for every page, so the caller
// gets one slice.
//
// # Errors
//
// An answer outside 2xx becomes an [*Error] with the HTTP status, the server's
// message and, for a failed validation, the field details. A network failure
// is returned as it is, wrapped with the request it belongs to.
package api
