# Installs the warp-engine command line from the release page.
#
#   irm https://git.teletypegames.org/engines/warp-engine-cli/raw/branch/master/install.ps1 | iex
#
# VERSION and INSTALL_DIR as in install.sh; INSTALL_DIR defaults to
# %LOCALAPPDATA%\Programs\warp-engine, which is added to the user's PATH.
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$Forge = 'https://git.teletypegames.org'
$Repo = 'engines/warp-engine-cli'
$Bin = 'warp-engine'
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\$Bin" }

$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { 'amd64' }
    'ARM64' { 'arm64' }
    default { throw "warp-engine install: unsupported CPU $($env:PROCESSOR_ARCHITECTURE)" }
}

$version = $env:VERSION
if (-not $version) {
    $version = (Invoke-RestMethod "$Forge/api/v1/repos/$Repo/releases/latest").tag_name
    if (-not $version) { throw "warp-engine install: no release found at $Forge/$Repo/releases" }
}

$asset = "$Bin-windows-$arch.exe"
$base = "$Forge/$Repo/releases/download/$version"
$tmp = Join-Path ([IO.Path]::GetTempPath()) ([IO.Path]::GetRandomFileName())
New-Item -ItemType Directory -Path $tmp | Out-Null
try {
    Write-Host "Downloading $Bin $version (windows/$arch)"
    Invoke-WebRequest -UseBasicParsing -Uri "$base/$asset" -OutFile (Join-Path $tmp $asset)
    $sums = (Invoke-WebRequest -UseBasicParsing -Uri "$base/SHA256SUMS").Content
    if ($sums -is [byte[]]) { $sums = [Text.Encoding]::UTF8.GetString($sums) }
    $line = $sums -split "`n" | Where-Object { ($_ -split '\s+')[1] -in @($asset, "*$asset") } | Select-Object -First 1
    if (-not $line) { throw "warp-engine install: $asset is not in SHA256SUMS" }
    $want = ($line -split '\s+')[0].ToLower()
    $got = (Get-FileHash -Algorithm SHA256 (Join-Path $tmp $asset)).Hash.ToLower()
    if ($want -ne $got) { throw "warp-engine install: checksum mismatch for $asset" }

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $exe = Join-Path $InstallDir "$Bin.exe"
    Move-Item -Force (Join-Path $tmp $asset) $exe
    Write-Host "Installed $exe"
} finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (-not (($userPath -split ';') -contains $InstallDir)) {
    [Environment]::SetEnvironmentVariable('Path', (@($userPath, $InstallDir) | Where-Object { $_ }) -join ';', 'User')
    Write-Host "Added $InstallDir to your PATH (new terminals)."
}

Write-Host ""
Write-Host "Sign in to a server (or, without a URL, to WarpEngine Local on this machine):"
Write-Host "  $Bin login https://games.example.org"
Write-Host "Manual: https://wiki.teletypegames.org/godoc/warp-engine-cli/latest/cmd/warp-engine/"
