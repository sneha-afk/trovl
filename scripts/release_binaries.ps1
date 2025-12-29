$BINARY_NAME = "trovl"
$LD_FLAGS = "-s -w"
$DIST_DIR = "dist"

try {
    $VERSION = git describe --tags --dirty --always 2>$null
    if (-not $VERSION) { $VERSION = "dev" }
} catch {
    $VERSION = "dev"
}


if (-not (Test-Path $DIST_DIR)) { New-Item -ItemType Directory -Path $DIST_DIR | Out-Null }

$Targets = @(
    @{ OS="windows"; ARCH="amd64"; EXT=".exe" },
    @{ OS="windows"; ARCH="arm64"; EXT=".exe" },
    @{ OS="linux"; ARCH="amd64"; EXT="" },
    @{ OS="linux"; ARCH="arm64"; EXT="" },
    @{ OS="darwin"; ARCH="amd64"; EXT="" },
    @{ OS="darwin"; ARCH="arm64"; EXT="" }
)

foreach ($t in $Targets) {
    Write-Host ": Building $($t.OS)/$($t.ARCH)"
    $env:CGO_ENABLED="0"
    $env:GOOS=$t.OS
    $env:GOARCH=$t.ARCH
    $outFile = Join-Path $DIST_DIR "$BINARY_NAME`_$($t.OS)_$($t.ARCH)$($t.EXT)"
    go build -ldflags "$LD_FLAGS -X main.version=$VERSION" -o $outFile .
}

# Checksums
try { Get-Command sha256sum -ErrorAction Stop | Out-Null; & sha256sum "$DIST_DIR\*" > "$DIST_DIR\checksums.txt" }
catch { Get-Command shasum -ErrorAction SilentlyContinue | Out-Null; shasum -a 256 "$DIST_DIR\*" > "$DIST_DIR\checksums.txt" }

Write-Host "Release $VERSION builds complete"
