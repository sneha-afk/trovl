#!/bin/sh
set -e
# Install script for trovl
# Usage: curl -fsSL https://raw.githubusercontent.com/sneha-afk/trovl/main/install.sh | sh
# Variables:
#        INSTALL_DIR=/custom/path sh install.sh

REPO="sneha-afk/trovl"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="trovl"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { printf "${GREEN}[INFO]${NC} %s\n" "$1"; }
warn() { printf "${YELLOW}[WARN]${NC} %s\n" "$1"; }
error() { printf "${RED}[ERROR]${NC} %s\n" "$1" >&2; exit 1; }

detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$OS" in
        linux*) echo "linux" ;;
        darwin*) echo "macos" ;;
        *) error "Unsupported OS: $OS. See https://github.com/$REPO for manual installation." ;;
    esac
}

detect_arch() {
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64|amd64) echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *) error "Unsupported architecture: $ARCH. See https://github.com/$REPO for manual installation." ;;
    esac
}

main() {
    info "Installing $BINARY_NAME..."

    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Detected: $OS ($ARCH)"

    FILENAME="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/${REPO}/releases/latest/download/${FILENAME}"

    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    info "Downloading from $URL..."
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$URL" -o "$TMP_DIR/$FILENAME"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$URL" -O "$TMP_DIR/$FILENAME"
    else
        error "Neither curl nor wget found. Please install one."
    fi

    info "Extracting..."
    tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

    mkdir -p "$INSTALL_DIR"
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    info "Installed to $INSTALL_DIR/$BINARY_NAME"

    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        warn "$INSTALL_DIR is not in your PATH"
        warn "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi

    if "$INSTALL_DIR/$BINARY_NAME" --version >/dev/null 2>&1; then
        info "Installation successful!"
        "$INSTALL_DIR/$BINARY_NAME" --version
    else
        warn "Binary installed but verification failed"
    fi
}

main "$@"
