// Command warp-engine manages a WarpEngine server from the command line. It
// speaks the server's Admin API: the catalog (softwares, releases, release
// files, images, links, announced forks), its vocabulary (statuses, categories, platform
// links), the download statistics, the outlets, the pipelines, the artifact store, the
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
// The platforms (--platform), the kinds of a release file (--kind) and the
// app variants of a platform (--app) are the server's own: the command line
// reads them once per run from the server's public build description
// (GET /build/platforms) and refuses any other value before sending
// anything. When that description cannot be read, it only checks that the
// value looks like a name (a-z, 0-9 and _) and leaves the rest to the
// server. The shell completion offers the same values.
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
// --forks-allowed lets the developer tools offer to fork a software (it needs
// a public source link); a fork names its parent with --fork-of-engine,
// --fork-of-name, --fork-of-version and --fork-of-repo:
//
//	warp-engine software update pong --forks-allowed
//	warp-engine software update pong-plus --fork-of-engine https://games.example.org --fork-of-name pong
//
// "software scaffold" sets a new software up from nothing on a server that
// has a forge and a CI: the catalog entry, its git repository and its
// pipeline, and prints the repository's addresses:
//
//	warp-engine software scaffold --name my-game --title "My Game" --platform tic80
//
// --app picks the platform's variant, for a platform that has them (golang:
// cli or gui); without it the server takes the platform's default:
//
//	warp-engine software scaffold --name my-tool --title "My Tool" --platform golang --app cli
//
// It also takes --license, and the --fork-of-… flags to set up a fork of a
// software on another engine (--fork-of-engine and --fork-of-name are
// needed, --fork-of-version and --fork-of-repo are optional):
//
//	warp-engine software scaffold --name pong-plus --title "Pong+" --platform love --license MIT \
//	    --fork-of-engine https://games.example.org --fork-of-name pong --fork-of-version 1.2.0
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
// # Releasing by hand
//
// Some packages the CI cannot build: the macOS ones of ebitengine, electron
// and golang gui programs are built and signed on a Mac. "release upload"
// sends them into the release's folder, as the CI's upload step does, and
// "release publish" registers the release again, which keeps the files the
// CI uploaded and adds the new ones. Every file name starts with
// NAME-VERSION; the command sends each file's SHA-256, which the server
// checks. The mac-release target of the WarpEngine templates runs both.
//
// These commands talk to the build endpoints (/build/upload, /build/publish),
// not to the Admin API: they take an application token with the upload
// (upload) and the update (publish) scope, the one the CI publishes with. The
// server is --server, else WARP_SERVER, else WARP_ENGINE_URL, else the
// profile; the token --token, else WARP_TOKEN, else WARP_ENGINE_TOKEN, else
// the profile's.
//
//	export WARP_SERVER=https://teletypegames.org WARP_TOKEN=we_xxxxxxxx
//	warp-engine release upload --name pong --version 1.2.0 pong-1.2.0-mac-x64.zip pong-1.2.0-mac-arm64.zip
//	warp-engine release publish --name pong --version 1.2.0 --platform ebitengine
//	warp-engine release publish --name gio-demo --version 0.3.0 --platform golang --app gui
//
// # Release files
//
// The files of a release, under asset, point into the artifact store: a kind
// and an absolute path. The kinds come from the server: every kind of every
// platform and app variant it builds for (cartridge, source, html, docs,
// win_x64, linux_x64, mac_arm64, cli_linux_x64 …). A download-only platform
// cannot have an html file.
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
// page, a video. --kind is link (the default) or source: the software's git
// repository, one per software, with --ref naming a branch or tag (empty for
// the default branch). The server checks a source link and shows whether it
// is public, when it checked and why the check failed.
//
//	warp-engine link list pong
//	warp-engine link create pong --label Source --url https://git.example.org/games/pong --kind source
//	warp-engine link update 5 --label "Source code" --ref v1.2.0
//	warp-engine link delete 5 --yes
//
// # Forks
//
// When another engine forks one of the softwares here, it announces the fork,
// and the announcement waits as a pending fork link until it is moderated:
// visible lists the fork with its parent, hidden keeps it out; gone is a fork
// its engine no longer has. list filters on --status and --software (an ID or
// a name) and sorts by createdAt or checkedAt; update (alias moderate) sets
// the status. delete removes the link, but a new announcement brings it back
// as pending.
//
//	warp-engine fork list --status pending
//	warp-engine fork list --software pong --sort -checkedAt
//	warp-engine fork show 4
//	warp-engine fork moderate 4 --status visible
//	warp-engine fork update 5 --status hidden
//	warp-engine fork delete 5 --yes
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
// software's title, --ip-address, --user-agent, --referer, --client-id), by
// the outlet it came through (--outlet, exact; none for the catalog API)
// and by time (--from, inclusive, and --to, exclusive). Times are ISO 8601; a
// bare date means its midnight UTC.
//
//	warp-engine download list --software ra --from 2026-11-01
//	warp-engine download list --outlet assembly64 --client-id Ultimate
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
// # Outlets
//
// The formats the catalog leaves the engine in: the catalog API, the shelves,
// an Assembly64 server and the host's own types. "outlet list" lists them,
// the catalog API first; "outlet check" shows what one offers (a release and
// kind per title) and which softwares it leaves out, and why; with --strict it
// exits 1 when it leaves any out. A server whose outlet configuration is
// invalid answers 409.
//
//	warp-engine outlet list
//	warp-engine outlet check assembly64
//	warp-engine outlet check c64 --strict
//
// # Pipelines
//
// A pipeline is a repository on the CI server that the engine knows: the
// platform it builds for and the software its builds are published as.
// "pipeline sync" reads the repositories from the CI server and prints what was
// created, updated or deactivated. "pipeline update" sets the platform (one of the
// server's, or unknown) or the software; --software takes an ID or a name, and "" detaches it. A software
// has one pipeline, so giving it to one pipeline takes it from the other.
// "pipeline runs" lists the runs newest first, one page of the CI server's
// listing at a time, and "pipeline trigger" starts a run, on main unless
// --branch says otherwise.
//
//	warp-engine pipeline list --platform c64 --active
//	warp-engine pipeline sync
//	warp-engine pipeline update 4 --software blessing-of-ra --platform c64
//	warp-engine pipeline runs 4 --page 2
//	warp-engine pipeline trigger 4 --branch develop
//
// # Files
//
// The artifact store is the directory tree of files the engine serves: release
// files, runtimes, anything uploaded. Paths are relative to its root. "file
// list" shows one directory (the root without DIR). "file upload" sends a local
// file into --dir. "file rename" renames in place: NAME is a new name, not a
// path. "file delete" removes a file or an empty directory and asks first.
// "file mkdir" needs the parent to exist.
//
//	warp-engine file list blessing-of-ra/1.2.0
//	warp-engine file upload ./manual.pdf --dir blessing-of-ra/docs
//	warp-engine file rename blessing-of-ra/docs/manual.pdf handbook.pdf
//	warp-engine file mkdir blessing-of-ra/docs
//	warp-engine file delete blessing-of-ra/docs/old.pdf --yes
//
// # Runtimes
//
// A runtime is the player a platform's games run in on the web. "runtime list"
// and "runtime show" list the installed versions with their components and how
// many releases use them. "runtime install" installs a runtime the engine
// downloads, with the components its platforms need; without VERSION it takes
// the version the engine has chosen, and with --dry-run it only prints what it
// would fetch and keep.
//
//	warp-engine runtime list --name emulatorjs
//	warp-engine runtime install emulatorjs --dry-run
//	warp-engine runtime install tic80 1.2.2715-dev
//
// # Sync peers
//
// A sync peer is another WarpEngine this one pushes releases to (outbound),
// takes pushes from (inbound), or both. The http transport needs the other
// engine's host and the token it issued (--secret); rsync needs host, user and
// path, and takes port and host_key. Transport options are given with --option
// key=value, repeatable; they replace the stored set, so give all of them
// (--data '{"transportOptions":{…}}' does the same in JSON). check asks whether
// the other side accepts this engine and exits 1 when it does not. push queues
// the releases one peer may be missing, and run does it for every outbound peer
// or the --peer ones; both ask first unless --yes, and --dry-run only shows
// what would go. inbound-token issues the token the peer pushes here with: it
// is shown once, and the old one stops working.
//
//	warp-engine peer list --outbound
//	warp-engine peer create --label games --option host=https://games.example.org --secret we_xxx --outbound
//	warp-engine peer update 3 --channels stable,beta
//	warp-engine peer check 3
//	warp-engine peer push 3 --software blessing-of-ra --dry-run
//	warp-engine peer run --peer 3 --peer 5 --all --yes
//	warp-engine peer inbound-token 3 --yes
//	warp-engine peer delete 3 --yes
//
// # Pairings
//
// A pairing is an inbound-only peer you add yourself, so that another engine
// (your laptop's WarpEngine Local, say) may push here. It needs the sync:pair
// scope while the server lets anyone pair, or admin:write. create and token
// print the new token once: set it as the secret of the other engine's peer
// for this one. The server answers 409 while inbound sync is off.
//
//	warp-engine pairing create --label "laptop local"
//	warp-engine pairing list
//	warp-engine pairing token 2 --yes
//	warp-engine pairing delete 2 --yes
//
// # Sync events
//
// Every push of a release between this engine and a peer, in either
// direction, is a sync event, newest first; --software takes an ID or a name.
// "sync-event software" lists the events of one software, and "sync-event
// status" shows where it stands with every peer: the latest event, and why the
// newest release would not go.
//
//	warp-engine sync-event list --status failed
//	warp-engine sync-event list --peer 3 --software blessing-of-ra
//	warp-engine sync-event show 41
//	warp-engine sync-event software blessing-of-ra --direction outbound
//	warp-engine sync-event status blessing-of-ra
//
// # Application tokens
//
// Application tokens are what programs sign in with: pipelines publishing
// releases, store clients, admin tools, other engines. Managing them needs the
// admin:tokens scope (sign in with --scopes admin:read,admin:write,admin:tokens).
// You cannot grant scopes you do not hold, nor an expiry later than your own;
// without --expires-at a new token expires with yours, and --expires-at ""
// means never. create and rotate print the token once, under the record; with
// -o json it is the plainToken field. rotate replaces a token, and the new one
// is pushed to the pipelines that use it. delete (alias revoke) ends it; a
// token cannot revoke itself.
//
//	warp-engine token list --state active
//	warp-engine token create --name ci --scopes admin:read,admin:write --expires-at 2027-06-30T00:00:00Z
//	warp-engine token create --name reader --scopes admin:read -o json | jq -r .plainToken
//	warp-engine token update 12 --scopes admin:read
//	warp-engine token rotate 12 --yes
//	warp-engine token revoke 12 --yes
//
// # WarpEngine Local
//
// WarpEngine Local can run a git server and a CI next to the engine. Only
// WarpEngine Local has the "local" commands; any other server answers 404
// (exit code 4). "local components" shows the components and their state.
// "local enable" starts a component and what it needs; "local disable" stops
// it and what needs it, and keeps its data; "local remove" removes its
// containers, with --purge its data and generated sign-ins too, and asks first.
// "local credentials" prints a component's generated sign-in, password
// included. "local ci-setup" sets a software (ID or name) up on the local git
// server and CI, starting them when they are off; its result holds the push
// URL. "local ci-teardown" stops a software's local CI and keeps its
// repository.
//
// Enabling, disabling, removing and the CI setup run in the background on the
// engine. The commands wait for them: they print each step on standard error
// as it happens, then the finished operation. A failed operation is an error
// (exit code 1), and so is one that has not finished within --timeout (default
// 10m); it goes on running on the engine. --no-wait prints the queued operation
// at once, and "local operation ID --wait" follows it later.
//
//	warp-engine local components
//	warp-engine local enable woodpecker
//	warp-engine local ci-setup my-game
//	warp-engine local enable gitea --no-wait -o id
//	warp-engine local operation 12 --wait
//	warp-engine local remove woodpecker --purge --yes
//
// # Checking a Makefile
//
// The Makefiles of the WarpEngine templates share a contract: build, run,
// check (the CI's check step runs it) and clean in every project; web when
// the platform's releases have a web build; mac-release when the platform
// has kinds the CI does not build. "check makefile" tries each target the
// project in DIR (the current directory without it) needs with make -n,
// which builds nothing, prints them with their state (ok, missing, fails)
// and exits with 1 when one is missing or fails.
//
// The platform and its variant come from the CI marker (.woodpecker.yaml
// "platform:" and "app:"), or --platform and --app; which kinds the platform
// has, from the server found as for "release upload". Without a server only
// the four targets every project has are checked, and a note says so.
//
//	warp-engine check makefile
//	warp-engine check makefile ../pong --server https://teletypegames.org
//	warp-engine check makefile . --platform golang --app gui -o json
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
//	WARP_SERVER            the server of "release upload|publish" and "check makefile", before WARP_ENGINE_URL
//	WARP_TOKEN             the token of "release upload|publish", before WARP_ENGINE_TOKEN
//
// WARP_ENGINE_URL and WARP_ENGINE_TOKEN mean the same to warp-engine-mcp, the
// MCP server of the same family.
package main
