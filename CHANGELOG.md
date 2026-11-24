# Changelog

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
