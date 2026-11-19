// Command warp-engine manages a WarpEngine server from the command line. It
// speaks the server's Admin API: the catalog (softwares, releases, release
// files, images, links), its vocabulary (statuses, categories, platform
// links), the download statistics, the pipelines, the artifact store, the
// runtimes, sync with other engines, the application tokens and, on
// WarpEngine Local, the git server and CI that run beside it. Anything
// without a command of its own is one "warp-engine api" call away.
//
// This page is the manual. The packages of the program are documented
// under app/: [git.teletypegames.org/engines/warp-engine-cli/app/cli] (the
// commands), [git.teletypegames.org/engines/warp-engine-cli/app/domain] (the
// models and services) and [git.teletypegames.org/engines/warp-engine-cli/app/api]
// (the HTTP client).
//
// # Install
//
// On macOS and Linux one line downloads the binary for the machine, checks
// it against the release's SHA256SUMS and puts it in ~/.local/bin:
//
//	curl -fsSL https://git.teletypegames.org/engines/warp-engine-cli/raw/branch/master/install.sh | sh
//
// On Windows, in PowerShell (into %LOCALAPPDATA%\Programs\warp-engine, which
// is added to the user's PATH):
//
//	irm https://git.teletypegames.org/engines/warp-engine-cli/raw/branch/master/install.ps1 | iex
//
// Both take two environment variables:
//
//   - VERSION: a release tag such as v0.1.0; the latest release by default.
//   - INSTALL_DIR: where the binary goes.
//
// Run the line again to update. The binaries are on the release page of the
// repository, static, for Linux, macOS and Windows on amd64 and arm64. With
// Go at hand:
//
//	go install git.teletypegames.org/engines/warp-engine-cli/cmd/warp-engine@latest
//
// # Signing in
//
// Every command but a few talks to one server as one token. "login" signs in
// and saves the server as a profile, which becomes the current one:
//
//	warp-engine login https://games.example.org
//
// The server decides how:
//
//   - A hosted engine signs in a device: the command prints a code and a
//     page, opens the page in the browser (not with --no-browser) and waits
//     until someone with admin rights approves it there. The token asks for
//     admin:read and admin:write; --scopes asks for others (admin:tokens to
//     manage tokens, sync:pair to pair engines).
//   - WarpEngine Local, the engine on this machine, signs in with its token
//     file (~/.warpengine/local-token). Without a URL, login goes there:
//
//     warp-engine login
//
//     The profile keeps no copy of the token: it reads the file on every
//     run, so a reinstall does not break it. --local forces this way.
//   - With --token, a token you already have (made in the admin, or with
//     "warp-engine token create") is checked and saved:
//
//     warp-engine login https://games.example.org --token we_xxxxxxxx --name ci
//
// "whoami" shows who the server takes the token for — label, owner, scopes,
// expiry and the resources it may use —, "server" what the engine says about
// itself (name, version, sign-in methods). "logout" revokes the token on the
// server and forgets the profile; --keep-token only forgets it.
//
// # Profiles
//
// The profiles live in a JSON file readable by you only:
// ~/.config/warp-engine/config.json on Linux, ~/Library/Application
// Support/warp-engine/config.json on macOS, %AppData%\warp-engine\config.json
// on Windows (WARP_ENGINE_CONFIG names another).
//
//	warp-engine profile list          # the current one is marked
//	warp-engine profile use local
//	warp-engine profile show
//	warp-engine profile remove ci     # forgets it, the token stays valid
//
// Which server and token a command uses, first match wins:
//
//  1. --server URL (or WARP_ENGINE_URL): that server, with --token (or
//     WARP_ENGINE_TOKEN) if given, else the token of a profile with the same
//     URL, else, on this machine, the WarpEngine Local token file.
//  2. --profile NAME (or WARP_ENGINE_PROFILE), or the current profile;
//     --token, if given, replaces its token.
//
// So a CI job needs no configuration file at all:
//
//	WARP_ENGINE_URL=https://games.example.org WARP_ENGINE_TOKEN=$TOKEN warp-engine software list
//
// # Commands
//
// Commands are a noun and a verb: "warp-engine software list". The verbs of
// a record are list (ls), show, create, update and delete (rm); a resource's
// own actions keep their names (pipeline trigger, token rotate). "--help"
// on any command lists its flags, and "warp-engine completion bash|zsh|fish|
// powershell" prints a completion script.
//
// A software is named by its ID or its exact name, everything else by ID:
//
//	warp-engine software show blessing-of-ra
//	warp-engine software show 12 --include releases,images
//
// create and update take one flag per field. update sends only the flags
// given and leaves the other fields alone; an empty value clears a field
// that may be empty:
//
//	warp-engine software update pong --status released --highlighted
//	warp-engine software update pong --author ""
//
// Fields without a flag of their own (nested objects) go in --data, a JSON
// object or @file.json (@- reads standard input); the flags win over it:
//
//	warp-engine software update pong --data @pong.json --title "Pong II"
//
// delete, and every command that cannot be undone, asks first; --yes (-y)
// answers for you, and without a terminal the command refuses to run
// without it.
//
// SECTIONS
//
// # The softwares
//
// The entries of the catalog: games, demos and tools.
//
//	warp-engine software list --platform c64 --status released
//	warp-engine software list --query ra --sort -updatedAt --all
//	warp-engine software show pong --include releases,images,externalLinks,extras
//	warp-engine software create --name pong --title Pong --platform love --author "Teletype Games"
//	warp-engine software update pong --category arcade --shelf-name PONG
//	warp-engine software delete pong --yes
//
// "software scaffold" sets a new software up from nothing on a server that
// has a forge and a CI: the catalog entry, its git repository and its
// pipeline, and prints the repository's addresses:
//
//	warp-engine software scaffold --name my-game --title "My Game" --platform tic80
//
// # Any request
//
// "api" sends any request, signed in like every other command, and prints
// the JSON answer as it came. PATH is relative to the Admin API; a path
// under /api/ or /.well-known/ is taken from the server's root.
//
//	warp-engine api GET /softwares --query platform=c64 --query perPage=5
//	warp-engine api PATCH /softwares/12 --data '{"highlighted":true}'
//	warp-engine api GET /openapi.json > openapi.json
//
// The OpenAPI document lists every operation the server offers, with the
// scope each needs.
//
// # Output
//
// --output (-o) picks the format of every command:
//
//   - table, the default: aligned columns for a list, "Field  value" lines
//     for one record. Notes for people (the paging hint, the sign-in code) go
//     to standard error.
//   - json: the records as the API sent them, indented — a list as an array,
//     a record as an object — for jq and scripts.
//   - id: the identifiers only, one per line, for shell loops.
//
// For example, every C64 software's ID into a loop:
//
//	for id in $(warp-engine software list --platform c64 --all -o id); do
//	  warp-engine release list "$id" -o json
//	done
//
// # Paging
//
// A list shows one page of the server's size (and says so on standard error
// when there are more). --page N picks a page, --per-page N its size (the
// server caps it), --all reads every page into one list, and --sort orders
// by a field, a leading "-" reversing it.
//
// # Exit codes
//
//	0  done
//	1  any other error
//	2  wrong command line: an unknown flag, a missing argument, a bad value
//	3  not signed in, or the token was refused (401) or lacks the scope (403)
//	4  no such record, or the resource is not available on this server (404)
//	5  the server refused the change: a conflict (409) or a record that does not validate (422)
//	6  too many requests (429); wait and try again
//	7  the server failed (5xx) or could not be reached
//
// The error itself goes to standard error as "warp-engine: <message>", with
// the server's reasons per field after a failed validation.
//
// # Environment
//
//	WARP_ENGINE_URL        a server instead of a profile (--server)
//	WARP_ENGINE_TOKEN      a token instead of the profile's (--token)
//	WARP_ENGINE_PROFILE    the profile to use (--profile)
//	WARP_ENGINE_CONFIG     the configuration file
//	WARP_ENGINE_LOCAL_DIR  WarpEngine Local's data directory (~/.warpengine)
//
// WARP_ENGINE_URL and WARP_ENGINE_TOKEN mean the same to warp-engine-mcp, the
// MCP server of the same family.
package main
