# Changelog

## 1.0.0

For WarpEngine 6.0 and its Admin API v2 (`/api/admin/v2`). It no longer talks to an engine before 6.0.

- `pipeline` is now `ci-repo` (aliases `ci-repos`, `ci-repository`), over the Admin API's `ci_repositories` resource; `list`, `show`, `update`, `sync`, `runs` and `trigger` are unchanged. The runs are still pipeline runs, and the sort key stays `lastPipelineAt`.
- `outlet list|check` is now `feed list|check`, over `/feeds`; a feed's type is `directory` where it was `shelf`. `download list --outlet` is now `--feed`, and a download shows its feed.
- `link` handles only the plain links of a software (`/softwares/{id}/links`, `/links/{id}`): `--kind` and `--ref` are gone, and a link has no kind, ref or public flag any more.
- New `source show|set|delete SOFTWARE`: the source repository of a software, one per software (`GET`/`PUT`/`DELETE /softwares/{id}/source`, body `url` and `ref`), with whether it is public and the last check. It replaces the `link --kind source` of 0.x.
- `fork` works over `/forks`; a fork's fields are `engineUrl`, `name` and `title` (they were `childEngine`, `childName`, `childTitle`).
- `local components` is now `local providers` (`/local_providers`), `local operation ID` is now `local task ID` (`/local_tasks`); the git server's role is `scm`.
- `software create|update --shelf-name` is now `--directory-name` (`directoryName`), and `software show --include` takes `links` instead of `externalLinks`.
- Asset kinds come from the server, so `package` (the former `cartridge`) needs no change here.

## 0.4.0

For the Admin API of warp_engine 5.11.0.

- `outlet list` and `outlet check KEY` (`--strict` exits 1 when the outlet leaves a software out): the outlets and what each offers, the Admin API's `outlets` resource.
- `download list` takes `--outlet` (exact; `none` for the catalog API) and `--client-id`, and a download shows its outlet and Client-Id.

## 0.3.1

- `check makefile` runs make in the project's directory instead of with `-C`, so GNU make 4 no longer puts its "Leaving directory" line where the reason a target fails belongs.

## 0.3.0

Tagged, but its pipeline failed on that test under GNU make 4, so it has no binaries; 0.3.1 is the release.


- `release upload FILE…` and `release publish`: a release made outside the CI (the macOS packages of ebitengine, electron and golang gui programs), through the build endpoints `/build/upload` and `/build/publish` with an application token (upload and update scopes), the way the CI publishes. The upload streams each file and sends its SHA-256, which the server checks; file names must start with `NAME-VERSION`. `publish` takes `--platform`, `--app` and `--commit`. The server and token come from `--server`/`--token`, `WARP_SERVER`/`WARP_TOKEN` (the names the templates' `mac-release` target uses), `WARP_ENGINE_URL`/`WARP_ENGINE_TOKEN`, or the profile.
- `check makefile [DIR]`: checks a project's Makefile against the shared contract of the WarpEngine templates — `build`, `run`, `check`, `clean`, plus `web` when the platform has a web build and `mac-release` when it has kinds the CI does not build (read from the server's `/build/platforms`, for the marker's app variant). Each target is tried with `make -n`; exits 1 when one is missing or fails. The platform comes from `.woodpecker.yaml` or `--platform`/`--app`.

## 0.2.0

For the Admin API of warp_engine 5.10.0.

- The platforms, asset kinds and app variants are no longer built into the command line: `--platform` (software create/update/scaffold, platform links, pipelines), the `--kind` of an asset and the new `--app` are checked against the server's public build description (`GET /build/platforms`), read once per run, so golang, batocera, retroarch, the golang `cli_*` kinds and any platform a host registers itself work. When the description cannot be read, only the shape of a value is checked and the server decides. A pipeline also takes `unknown`.
- `software scaffold --app`: the platform's variant (golang `cli` or `gui`), sent as `app`; checked against the platform's variants.
- Shell completion offers the server's platforms, asset kinds and the `--app` variants of the given `--platform`.
- A value the server does not know exits with 2 (usage) before any request is sent.
- `fork list|show|update|delete` (`update` alias `moderate`): the forks other engines announced, filtered by `--status` and `--software`, moderated to visible or hidden.
- `link create|update` take `--kind link|source` and `--ref`; links show their kind, ref, whether they are public and the last verification.
- `software create|update` take `--forks-allowed` and `--fork-of-engine|name|version|repo`, and show them.
- `software scaffold` takes `--license` and `--fork-of-engine|name|version|repo` to set up a fork.

## 0.1.0

The first release: `warp-engine`, the command line of a WarpEngine server.

- Signing in: device sign-in in the browser, an existing token, or WarpEngine Local's token file; saved profiles; `--profile`, `--server`, `--token` and `WARP_ENGINE_URL`/`WARP_ENGINE_TOKEN` for scripts and CI.
- The whole Admin API as commands: `software`, `release`, `asset`, `link`, `image`, `extra`, `status`, `category`, `platform-link`, `download`, `pipeline`, `file`, `runtime`, `peer`, `pairing`, `sync-event`, `token`, `local`, and `api` for any request.
- Table, JSON and ID output; paging with `--all`; confirmations with `--yes`; exit codes by kind of failure.
- One-line installers for macOS, Linux and Windows; the manual on the wiki as Go doc.
