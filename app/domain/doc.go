// Package domain is the business layer of the command line: what a
// WarpEngine server holds (the models), how to reach it (the repositories)
// and what the commands do with it (the services). It knows nothing of flags
// or terminals; package cli drives it, and package api carries its requests.
//
// # Layers
//
// Every resource of the Admin API is one entity with three parts, each in its
// own file named {type}.{entity}.go:
//
//   - the model (model.software.go): a struct named after the entity with
//     the suffix Model, whose JSON tags are the API's field names; a nullable
//     field is a pointer.
//   - the repository (repository.software.go): an interface, then a struct
//     over the shared [api.Client]. Its methods map one to one onto the
//     Admin API's operations: All lists, FindByID gets one, Create, Update and
//     Destroy change, and the resource's own actions keep their names
//     (Trigger, Rotate, Check).
//   - the service (service.software.go): an interface, then a struct over
//     the repository interfaces it needs. All, Show, Create, Update and
//     Destroy plus the domain actions. A service checks what would only cost
//     a round trip (a required field, an empty update), resolves references
//     and runs the flows that take more than one request: a device sign-in,
//     waiting for a WarpEngine Local operation.
//
// A controller of package cli calls services, a service calls repositories,
// a repository calls the client; nothing skips a layer, and every dependency
// is an interface, so the tests replace any layer with a mock.
//
// # Wiring
//
// [NewDomain] is the composition root: it builds every repository around one
// [api.Client], the configuration file and the WarpEngine Local data
// directory, builds the services around the repositories, and returns them
// all in a [Domain]. The client is shared and mutable: [AuthService.Connect]
// points it at the server and token a command asked for, and every
// repository then talks to that server.
//
// # Lists, filters and attributes
//
// A list takes a [Filter] (the query parameters that narrow it, keyed by the
// API's names) and a [ListQuery] (page, page size, order, or every page),
// and answers a [PageModel]. A create or update takes [Attributes]: only the
// fields present are sent, so an update leaves the rest alone, and a nil
// value clears a nullable field.
//
// # References
//
// A software can be named by its numeric ID or its exact name; the services
// resolve it with [ResolveSoftwareID] before they call the repository. Other
// records are named by ID only ([ParseID]).
//
// # Profiles and signing in
//
// The servers the command line knows are profiles ([ProfileModel]) in a JSON
// configuration file ([DefaultConfigPath]), written with mode 0600 because
// it holds tokens. [AuthService] signs in three ways — a device sign-in
// approved in the browser, an existing token, or WarpEngine Local's token
// file ([LocalTokenRepository]) — and decides, for every command, which
// server and token to use ([AuthService.Connect]).
//
// # Errors
//
// The repositories return the [*api.Error] of a refused request as it is,
// so the command line can tell a 404 from a 422. The services add their own
// sentinel errors: [ErrNoProfile], [ErrNoLocalToken], [ErrDenied],
// [ErrExpired], [ErrNothingToChange].
package domain
