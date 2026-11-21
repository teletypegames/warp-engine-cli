# Changelog

## Unreleased

For the Admin API of warp_engine 5.6.0.

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
