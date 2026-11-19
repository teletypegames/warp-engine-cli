// Package cli is the command line: the command tree, the controllers that run
// the commands, and the serializers that turn flags into API attributes and
// records into tables. It plays the part an HTTP layer plays in a server:
// commands instead of routes, flags instead of request bodies, a terminal
// instead of a response. The user manual is the documentation of the
// warp-engine command; this is how the layer is built.
//
// # Start
//
// [Start] is called by main with the [domain.Domain]: it builds one
// controller per entity around its service, hands them to the [Router],
// runs the command line and turns the error into an exit code ([ExitCode]).
//
// # The router
//
// [Router] has one field per controller. [Router.Init] builds the root
// command with the global flags (--profile, --server, --token, --output),
// the help groups, and one command group per controller that is set (a nil
// controller leaves its commands out, which keeps the tests small). The
// commands of an entity are added by a method in router.{entity}.go, each as
// a [Route]: its use line, help, flags and the controller action it runs.
//
// Before a command runs, the root connects: [AuthController.Connect] picks
// the server and token and points the shared client at them. Commands marked
// Offline (login, logout, profile, version, help, completion) skip that.
//
// # Controllers
//
// A controller is an interface plus a struct holding the service interface
// it uses. Its actions have the cobra signature
//
//	func(cmd *cobra.Command, args []string) error
//
// and are named like REST actions: Index lists, Show shows, Create, Update
// and Destroy change, and the resource's own actions keep their names. An
// action reads its flags through the serializer, calls the service once and
// prints the answer. A destructive action asks first ([Confirm]) unless
// --yes is given, and refuses to guess without a terminal.
//
// # Serializers
//
// A serializer is a stateless struct with:
//
//   - Filters, the [FilterFlag] list of the list command, and Fields, the
//     [FieldFlag] list of create and update, which [DeclareFilters] and
//     [DeclareFields] turn into flags;
//   - FlagsToFilter and RequestToAttributes, which read the flags that were
//     set into a [domain.Filter] and [domain.Attributes] — an unset flag
//     sends nothing, an empty nullable flag sends null, a value outside an
//     enum is a [*UsageError];
//   - ModelToTable and ModelToRecord, which render the models as a [Table]
//     (the identifier in the first column) or a [Record] (the identifier
//     first).
//
// # Output
//
// [PrintList] and [PrintRecord] print in the format of --output: aligned
// columns (table, the default), the records as the API sent them (json), or
// just the identifiers (id). Notes for people — the paging hint
// ([PrintPaging]), the sign-in code — go to standard error, so standard
// output stays clean for pipes. [PrintDone] confirms an action in table
// output only.
//
// # Exit codes
//
// [ExitCode] maps an error to the code the program exits with: [ExitUsage]
// for a bad command line, [ExitAuth] when not signed in or refused,
// [ExitNotFound], [ExitRejected] for a conflict or a failed validation,
// [ExitRateLimited], [ExitServer] for a failing or unreachable server, and
// [ExitError] for anything else.
package cli
