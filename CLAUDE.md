# CLAUDE.md

warp-engine-cli: the command line of a WarpEngine server, over its Admin API
(`/api/admin/v2`, WarpEngine 6) and its sign-in endpoints. Go module
`git.teletypegames.org/engines/warp-engine-cli`, binary `warp-engine`
(`cmd/warp-engine`). The documentation lives in the source: the package
comments are the manual and the doc comments the reference. Every release tag
publishes it to https://wiki.teletypegames.org/godoc/warp-engine-cli/latest/cmd/warp-engine/ ;
`README.md` only installs and links there.

The Admin API's contract is its OpenAPI document,
`GET <server>/api/admin/v2/openapi.json` (needs a token; WarpEngine Local's is
in `~/.warpengine/local-token`). The engine's source is
`services/teletypegames/libs/ruby/warp_engine/lib/warp_engine/admin_api/definitions/`
in the TTG workspace.

## Layout and layers

The structure follows the Go rules of LibrePM's `STRUCTURE.md`, with the HTTP
server replaced by a command line:

| Package | Path | Role |
|---|---|---|
| `api` | `app/api/` | the HTTP client (the "database driver"): `Client`, the envelope, `List`/`Get`/`Send`/`Delete`, `Error` |
| `domain` | `app/domain/` | models, repositories (over `*api.Client`), services, `domain.go` wiring |
| `cli` | `app/cli/` | controllers (cobra actions), serializers (flags → attributes, models → tables), the router, output |
| `main` | `cmd/warp-engine/` | builds the client, `domain.NewDomain`, `cli.Start` |

Nothing skips a layer: a controller calls a service interface, a service calls
repository interfaces, a repository calls the `api` client. The domain never
imports `cli`; nothing but `api` imports `net/http` for requests.

File names are `{type}.{entity}.go`, multi-word entities with underscores:
`model.software_image.go`, `repository.software_image.go`,
`service.software_image.go`, `controller.software_image.go`,
`serializer.software_image.go`, `router.software_image.go`.

### Per entity

- **Model** (`domain/model.X.go`): `XModel` struct with JSON tags exactly as
  the OpenAPI schema (camelCase); nullable fields are pointers; enums as
  exported `[]string` vars.
- **Repository** (`domain/repository.X.go`): `XRepositoryInterface` first,
  then `XRepository{Client *api.Client}`. Methods: `All(filter Filter, query
  ListQuery) (*PageModel[XModel], error)` (via the `list` helper),
  `FindByID(id int)`, `Create(attributes Attributes)`, `Update(id int,
  attributes Attributes)`, `Destroy(id int)`, plus the resource's actions
  named after them (`Trigger`, `Rotate` …). Nested lists take the parent ID
  first: `All(software_id int, …)`.
- **Service** (`domain/service.X.go`): `XServiceInterface` +
  `XService{XRepository …}`. Methods `All`, `Show`, `Create`, `Update`,
  `Destroy` + domain actions. A software is always given as a reference
  string (ID or exact name) and resolved with `ResolveSoftwareID`; other IDs
  are parsed by the controller or with `ParseID`. Validate what the server
  would refuse anyway only when it saves a round trip (required fields,
  empty updates → `ErrNothingToChange`).
- **Controller** (`cli/controller.X.go`): `XControllerInterface` +
  `XController{XService domain.XServiceInterface}`. Actions have the cobra
  signature `func(cmd *cobra.Command, args []string) error`: `Index`, `Show`,
  `Create`, `Update`, `Destroy`, plus named actions. They read flags through
  the serializer and print with `PrintList`/`PrintRecord`/`PrintPaging`/
  `PrintDone`. Deletes and other destructive actions ask with `Confirm` and
  take `--yes` (`DeclareYes`).
- **Serializer** (`cli/serializer.X.go`): stateless `XSerializer{}` with
  `Filters() []FilterFlag`, `Fields() []FieldFlag`, `FlagsToFilter`,
  `RequestToAttributes`, `ModelToTable([]XModel) Table` (ID first column),
  `ModelToRecord(XModel) Record` (ID first field).
- **Routes** (`cli/router.X.go`): `func (r Router) x(root *cobra.Command)`
  adding a `group` with `add(… Route{…})` per command; the `Router` struct in
  `router.go` has one `XController XControllerInterface` field and `Init`
  calls `r.x(root)` when it is not nil. Commands are `<noun> <verb>`:
  `list` (alias `ls`), `show`, `create`, `update`, `delete` (alias `rm`), and
  kebab-case actions. Every route has a Short, the non-obvious ones an
  Example.
- **Wiring**: `domain/domain.go` (repository + service fields and
  construction) and `cli/cli.go` (controller construction).

### Naming

As STRUCTURE.md: structs PascalCase with suffix (`ReleaseService`),
interfaces with `Interface`, methods PascalCase and never repeating the
entity (`Create`, not `CreateRelease`), **local variables and parameters
snake_case** (`software_id`), value receivers with one letter, JSON tags
camelCase. Errors are assigned, never shadowed, before being returned.

## Tests

`go test ./...`; everything under `test/`:

- `test/mock/repository` (`package mockrepo`) and `test/mock/service`
  (`package mocksvc`): hand-written mocks with one `XFunc` field per method,
  nil = zero value and no error. File `mock_X_repository.go` /
  `mock_X_service.go`, struct `MockXRepository` / `MockXService`. No
  testify/mock.
- `test/service` — services over mock repositories:
  `TestXService_Method_Scenario`.
- `test/controller` — the real command tree (`execute(cli.Router{XController:
  …}, "x", "list", …)` from `helpers_test.go`) over mock services: flags reach
  the service, output is printed, exit codes.
- `test/serializer` — pure functions; `ModelToTable` with a filled and an
  empty slice.
- `test/repository` — repositories against `httptest.Server`: method, path,
  query, body, envelope.

Assertions with `testify/assert` only.

## Documenting code: write for Go doc

CI fails on any exported identifier without a doc comment (`revive` with
`revive.toml`). Every change keeps the code fully documented:

- Every exported type, function, method, constant, variable and struct field
  gets a doc comment starting with its name: what it is or does, what it
  takes and returns, its errors. Not its history: the why goes in the commit.
- Every package has its comment in `doc.go`. The user manual is the package
  comment of `cmd/warp-engine` (installing, signing in, profiles, every
  command group, output, exit codes, environment); the layer docs are in
  `app/*/doc.go`. A new command or flag goes into the manual too.
- Go doc syntax, not Markdown: `[Name]`/`[pkg.Name]` links, `# Heading` (only
  in package comments), lists with `-`, code as an indented block. No
  backticks, no `**`.
- Before committing:

      gofmt -l . && go vet ./... && go test ./...
      revive -config revive.toml -set_exit_status ./...     # must print nothing
      doc2go -out /tmp/wec-doc -home git.teletypegames.org/engines/warp-engine-cli ./...

## Releases

A release is a `vX.Y.Z` tag on master, with a `## X.Y.Z` section in
`CHANGELOG.md` (it becomes the release notes). Its pipeline runs the checks,
builds the static binaries `warp-engine-<os>-<arch>[.exe]` for linux, darwin
and windows on amd64 and arm64 with `-X …/app/cli.Version=X.Y.Z`, attaches
them with `SHA256SUMS` to the Gitea release (`install.sh` and `install.ps1`
download from there, the same scheme as engines/warp-engine-mcp), and
publishes the Go doc as the generic package
`engines/warp-engine-cli-godoc/<version>`, which the wiki-fetcher unpacks
under `/godoc/`.
