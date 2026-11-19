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
//
//   - WarpEngine Local, the engine on this machine, signs in with its token
//     file (~/.warpengine/local-token). Without a URL, login goes there:
//
//     warp-engine login
//
//     The profile keeps no copy of the token: it reads the file on every
//     run, so a reinstall does not break it. --local forces this way.
//
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
// # Releases
//
// A release is one version of a software, with its files. The channel follows
// from the version and cannot be set: dev-… is dev, a pre-release such as
// 1.3.0-beta.1 is beta, anything else stable. The commit is set by the CI.
// Adding a version that was deleted brings it back. Deleting a release deletes
// its files too.
//
//	warp-engine release list pong --channel beta --sort -version
//	warp-engine release list pong --commit-sha 3fb3ed8 -o id
//	warp-engine release show 42
//	warp-engine release create pong --version 1.2.0
//	warp-engine release update 42 --version 1.2.1
//	warp-engine release delete 42 --yes
//
// # Release files
//
// The files of a release, under asset, point into the artifact store: a kind
// (cartridge, source, html, docs, win_x86, win_x64, linux_x86, linux_x64,
// linux_arm64, mac_x64, mac_arm64, mac_universal) and an absolute path. A
// download-only platform cannot have an html file.
//
//	warp-engine asset list 42
//	warp-engine asset show 7
//	warp-engine asset create 42 --kind html --path /pong/1.2.0/html/index.html
//	warp-engine asset update 7 --path /pong/1.2.0/pong.love
//	warp-engine asset delete 7 --yes
//
// # Links
//
// The external links of a software, under link: its source repository, a store
// page, a video.
//
//	warp-engine link list pong
//	warp-engine link create pong --label Source --url https://git.example.org/games/pong
//	warp-engine link update 5 --label "Source code"
//	warp-engine link delete 5 --yes
//
// # Images
//
// The screenshots and covers of a software, under image, in order. upload puts
// a file into the host's image library and binds it in one step; create binds
// an image the library already has, by its ID. update makes an image the
// default or moves it. delete only unbinds it: the library keeps the picture.
//
//	warp-engine image list pong
//	warp-engine image upload pong cover.png --default
//	warp-engine image upload pong shot-2.png --position 2
//	warp-engine image create pong --image-id 118 --position 3
//	warp-engine image update 31 --position 0
//	warp-engine image delete 31 --yes
//
// # Extras
//
// A hosted engine can declare additional resources of its own for a software,
// such as its SEO metadata; "software show SOFTWARE --include extras" lists
// them, and an engine without any answers 404 (exit code 4). update replaces
// the value: a JSON object, or an array of objects for a resource that holds
// many; an entry with "id" and "_destroy": true removes one.
//
//	warp-engine extra show pong seo_metadata -o json
//	warp-engine extra update pong seo_metadata --value '{"title":"Pong"}'
//	warp-engine extra update pong seo_metadata --value @seo.json
//
// # Statuses and categories
//
// A software's status and category are keys of two vocabularies on the
// server. The statuses say what state a software is in and whether the stores
// show it by default; the categories group the catalog. A category key that a
// software brings and the vocabulary does not know yet is registered with
// "needs review" set, waiting for someone to approve it. Records are named by
// their numeric ID.
//
//	warp-engine status list
//	warp-engine status create --key beta --label Beta --badge BETA --default-visible
//	warp-engine status update 4 --position 2 --badge ""
//	warp-engine status delete 4 --yes
//	warp-engine category list --needs-review
//	warp-engine category update 9 --label Utility --needs-review=false
//
// Both lists take --sort position or key; statuses filter on
// --default-visible, categories on --needs-review (both also take =false).
//
// # Platform links
//
// A platform link is shown with every software of one platform: the emulator,
// runtime or documentation a player needs.
//
//	warp-engine platform-link list --platform c64
//	warp-engine platform-link create --platform c64 --name VICE --url https://vice-emu.sourceforge.io/
//	warp-engine platform-link update 3 --position 0
//
// # Downloads
//
// The server logs every file it serves. "download list" reads the log, newest
// first, and filters it by substrings (--file-path, --software for the
// software's title, --ip-address, --user-agent, --referer) and by time
// (--from, inclusive, and --to, exclusive). Times are ISO 8601; a bare date
// means its midnight UTC.
//
//	warp-engine download list --software ra --from 2026-11-01
//	warp-engine download list --file-path .d64 --all -o json
//	warp-engine download show 1234
//
// The statistics: the counts of today, this week, this month and in total, the
// most downloaded files and softwares (--limit, at most 100), and one
// software's downloads by release and by file:
//
//	warp-engine download summary
//	warp-engine download top-files --limit 20
//	warp-engine download top-softwares
//	warp-engine download stats blessing-of-ra
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
