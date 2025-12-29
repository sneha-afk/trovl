#!/usr/bin/env bash
set -euo pipefail

BINARY_NAME=trovl
LD_FLAGS="-s -w"
VERSION=$(git describe --tags --dirty --always || echo dev)
DIST_DIR=dist

mkdir -p $DIST_DIR

TARGETS=(
  "windows amd64"
  "windows arm64"
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
)

for TARGET in "${TARGETS[@]}"; do
  read -r OS ARCH <<<"$TARGET"
  EXT=""
  [[ "$OS" == "windows" ]] && EXT=".exe"
  echo ": Building $OS/$ARCH"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build -ldflags "$LD_FLAGS -X main.version=$VERSION" -o $DIST_DIR/${BINARY_NAME}_${OS}_${ARCH}${EXT} .
done

# Checksums
if command -v sha256sum &>/dev/null; then
  sha256sum $DIST_DIR/* > $DIST_DIR/checksums.txt
else
  shasum -a 256 $DIST_DIR/* > $DIST_DIR/checksums.txt
fi

echo "Release $VERSION builds complete"
