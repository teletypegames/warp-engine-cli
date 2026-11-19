#!/bin/sh
# Installs the warp-engine command line from the release page.
#
#   curl -fsSL https://git.teletypegames.org/engines/warp-engine-cli/raw/branch/master/install.sh | sh
#
# VERSION      release tag, e.g. v0.1.0 (default: the latest)
# INSTALL_DIR  default ~/.local/bin
set -eu

FORGE=https://git.teletypegames.org
REPO=engines/warp-engine-cli
BIN=warp-engine
INSTALL_DIR=${INSTALL_DIR:-$HOME/.local/bin}

say() { printf '%s\n' "$*" >&2; }
die() { say "warp-engine install: $*"; exit 1; }

command -v curl >/dev/null 2>&1 || die "curl is needed"

case "$(uname -s)" in
  Linux) os=linux ;;
  Darwin) os=darwin ;;
  *) die "unsupported system $(uname -s); on Windows use install.ps1" ;;
esac
case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  arm64 | aarch64) arch=arm64 ;;
  *) die "unsupported CPU $(uname -m)" ;;
esac

VERSION=${VERSION:-}
if [ -z "$VERSION" ]; then
  VERSION=$(curl -fsSL "$FORGE/api/v1/repos/$REPO/releases/latest" |
    sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)
  [ -n "$VERSION" ] || die "no release found at $FORGE/$REPO/releases"
fi

asset="$BIN-$os-$arch"
base="$FORGE/$REPO/releases/download/$VERSION"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT INT TERM

say "Downloading $BIN $VERSION ($os/$arch)"
curl -fsSL -o "$tmp/$asset" "$base/$asset" || die "download failed: $base/$asset"
curl -fsSL -o "$tmp/SHA256SUMS" "$base/SHA256SUMS" || die "download failed: $base/SHA256SUMS"

want=$(awk -v f="$asset" '$2 == f || $2 == "*" f { print $1 }' "$tmp/SHA256SUMS")
[ -n "$want" ] || die "$asset is not in SHA256SUMS"
if command -v sha256sum >/dev/null 2>&1; then
  got=$(sha256sum "$tmp/$asset" | awk '{ print $1 }')
else
  got=$(shasum -a 256 "$tmp/$asset" | awk '{ print $1 }')
fi
[ "$want" = "$got" ] || die "checksum mismatch for $asset"

mkdir -p "$INSTALL_DIR"
chmod 755 "$tmp/$asset"
mv -f "$tmp/$asset" "$INSTALL_DIR/$BIN"
say "Installed $INSTALL_DIR/$BIN ($("$INSTALL_DIR/$BIN" version 2>/dev/null || echo "$VERSION"))"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) say "Note: $INSTALL_DIR is not on your PATH; add it, e.g. export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac

say ""
say "Sign in to a server (or, without a URL, to WarpEngine Local on this machine):"
say "  $BIN login https://games.example.org"
say "Manual: https://wiki.teletypegames.org/godoc/warp-engine-cli/latest/cmd/warp-engine/"
