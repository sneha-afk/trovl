---
layout: default
---

# installation

`trovl` can be installed in multiple ways: using **install scripts (recommended)**, Go toolchain, `eget`, or manually with pre-built binaries.

> **Tip:** The install scripts handle OS and architecture detection automatically, so you don’t have to worry about downloading the wrong binary.

---

## Compatibility

`trovl` currently supports the following operating systems and architectures:

| OS      | Architectures          | Notes |
|---------|----------------------|-------|
| Linux   | amd64, arm64          | Adjust download/extract paths if needed |
| macOS   | amd64 (Intel), arm64 (Apple Silicon) | Pick the correct binary for your architecture |
| Windows | amd64, arm64          | PowerShell and CMD instructions provided separately |

> ⚠️ Make sure to adjust URLs if you are not using the automated install scripts.

---

## Quick Install (Recommended)

### 1. Using Install Scripts (All Platforms)

This is the easiest method as it automatically detects your OS and architecture.

#### Linux/macOS

```bash
curl -fsSL https://raw.githubusercontent.com/sneha-afk/trovl/main/install.sh | sh
````

Custom installation directory:

```bash
curl -fsSL https://raw.githubusercontent.com/sneha-afk/trovl/main/install.sh | INSTALL_DIR=/usr/local/bin sh
```

#### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/sneha-afk/trovl/main/install.ps1 | iex
```

Custom installation directory:

```powershell
.\install.ps1 -InstallDir "C:\bin"
```

**Verify installation:**

```bash
trovl --version
```

---

### 2. Using Go Toolchain

**Requirements:** Go 1.21+

```bash
go install github.com/sneha-afk/trovl@latest
```

**Verify installation:**

```bash
trovl --version
```

---

### 3. Using `eget` (all platforms)

```bash
eget sneha-afk/trovl
```

**Verify installation:**

```bash
trovl --version
```

[Install `eget`](https://github.com/zyedidia/eget) if needed.

---

## Manual Installation (Pre-built Binaries)

Pre-built binaries are available for Linux, macOS, and Windows.

**Download from:** [GitHub Releases](https://github.com/sneha-afk/trovl/releases)

**Recommended installation location:** to simplify your ACLs

| OS          | Default Location           |
| ----------- | -------------------------- |
| Linux/macOS | `~/.local/bin`             |
| Windows     | `%USERPROFILE%\.local\bin` |

> Make sure the directory is in your `PATH`.

```bash
# Linux/macOS
export PATH="$HOME/.local/bin:$PATH"
```

```powershell
# Windows PowerShell
setx PATH "$env:PATH;$env:USERPROFILE\.local\bin"
```

---

### Linux

```bash
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_linux_amd64.tar.gz
tar -xzf trovl_linux_amd64.tar.gz

mkdir -p "$HOME/.local/bin"
mv trovl "$HOME/.local/bin/"

trovl --version
```

### macOS

For Apple Silicon (`arm64`):

```bash
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_macos_arm64.tar.gz
tar -xzf trovl_macos_arm64.tar.gz

mkdir -p "$HOME/.local/bin"
mv trovl "$HOME/.local/bin/"

trovl --version
```

For Intel (`amd64`), replace `arm64` with `amd64`.

### Windows

#### PowerShell

```powershell
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_windows_amd64.zip
Expand-Archive trovl_windows_amd64.zip

mkdir "$env:USERPROFILE\.local\bin" -Force
move trovl.exe "$env:USERPROFILE\.local\bin\"

setx PATH "$env:PATH;$env:USERPROFILE\.local\bin"

trovl --version
```

#### Command Prompt

```cmd
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_windows_amd64.zip
tar -xf trovl_windows_amd64.zip

mkdir "%USERPROFILE%\.local\bin"
move trovl.exe "%USERPROFILE%\.local\bin\"

setx PATH "%PATH%;%USERPROFILE%\.local\bin"

trovl --version
```

---

## Verifying Checksums (Optional)

Ensures downloads haven’t been corrupted or tampered with.

```bash
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/checksums.txt
```

### Linux/macOS

```bash
sha256sum trovl_linux_amd64.tar.gz
sha256sum -c checksums.txt
```

### Windows

#### PowerShell

```powershell
(Get-FileHash trovl_windows_amd64.zip -Algorithm SHA256).Hash -eq (Select-String trovl_windows_amd64.zip checksums.txt).Line.Split()[0]
```

#### Command Prompt

```cmd
certutil -hashfile trovl_windows_amd64.zip SHA256
```

Compare output with `checksums.txt`.

---

## Updating

Repeat the installation steps or use the same tool:

```bash
# Go
go install github.com/sneha-afk/trovl@latest

# eget
eget sneha-afk/trovl
```

Check your version:

```bash
trovl --version
```

---

## Uninstalling

### Linux/macOS

```bash
rm ~/.local/bin/trovl
```

### Windows

#### PowerShell

```powershell
del "$env:USERPROFILE\.local\bin\trovl.exe"
```

#### Command Prompt

```cmd
del "%USERPROFILE%\.local\bin\trovl.exe"
```

---

## Need Help?

If you encounter issues, open an issue: [github.com/sneha-afk/trovl/issues](https://github.com/sneha-afk/trovl/issues)
