# warp-engine-cli

`warp-engine`, the command line of a [WarpEngine](https://wiki.teletypegames.org/engines/warp-engine)
server. It drives the whole Admin API from a terminal or a script: the catalog
(softwares, releases, files, images, links), its vocabulary, pipelines, the
artifact store, runtimes, downloads, sync with other engines, application
tokens and the components of WarpEngine Local. Every list pages or reads
everything, every command prints a table or, with `-o json`, JSON.

## Install

macOS and Linux:

```sh
curl -fsSL https://git.teletypegames.org/engines/warp-engine-cli/raw/branch/master/install.sh | sh
```

Windows (PowerShell):

```powershell
irm https://git.teletypegames.org/engines/warp-engine-cli/raw/branch/master/install.ps1 | iex
```

The installer downloads the latest release binary for your system, checks it
against `SHA256SUMS` and puts it in `~/.local/bin` (Windows:
`%LOCALAPPDATA%\Programs\warp-engine`, added to your PATH). `VERSION=v0.1.0`
picks a release, `INSTALL_DIR=…` another directory. With Go installed:
`go install git.teletypegames.org/engines/warp-engine-cli/cmd/warp-engine@latest`.

Then sign in — to a server, or without a URL to WarpEngine Local on this machine:

```sh
warp-engine login https://games.example.org
warp-engine software list
```

**Documentation:** <https://wiki.teletypegames.org/godoc/warp-engine-cli/latest/cmd/warp-engine/>
— the manual (every command, profiles, output, exit codes) and the package
reference, built from the package comments on every release tag.

## License

MIT — see [`LICENSE.md`](LICENSE.md). Copyright © 2026 Teletype Games.
