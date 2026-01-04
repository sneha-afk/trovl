---
layout: default
---

# installation

## Using Go toolchain (recommended!)

**Requirements:** Go 1.21+

```bash
go install github.com/sneha-afk/trovl@latest
```

This downloads the source, builds `trovl`, and installs it to `$GOBIN` (defaults to `$GOPATH/bin`).

**Verify installation:**

```bash
trovl version
```

---

## Platform-Specific Installation

Pre-built binaries are available for **Linux** (amd64, arm64), **macOS** (Intel & Apple Silicon), and **Windows** (amd64, arm64). All binaries are statically linked with **no additional dependencies**.

| OS | Architectures |
|----|---------------|
| Linux | amd64, arm64 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Windows | amd64, arm64 |

**Download from:** [github.com/sneha-afk/trovl/releases](https://github.com/sneha-afk/trovl/releases)

**Installation location:** The instructions below install the binary to `~/.local/bin` (or `%USERPROFILE%\.local\bin` on Windows), a user-local directory that doesn't require root/administrator privileges. This keeps the installation isolated to your user account and avoids potential conflicts with system packages.

> **Note:** Ensure `~/.local/bin` is in your `PATH`. You can verify if it is already there: `echo $PATH`. If needed, add it to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):
> ```bash
> export PATH="$HOME/.local/bin:$PATH"
> ```
> ```powershell
> setx PATH "$env:PATH;$env:USERPROFILE\.local\bin"
> ```

> Note: for safety, any Windows paths have $env:USERPROFILE specified, though modern builds automatically expand `~` when specified.

---

### Linux

#### Pre-built Binary

```bash
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_linux_amd64.tar.gz
tar -xzf trovl_linux_amd64.tar.gz

mkdir -p "$HOME/.local/bin"
mv trovl "$HOME/.local/bin/"

trovl version
```

---

### macOS

#### Pre-built Binary

Assuming an `arm64` installation:

```bash
# Download and extract
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_macos_arm64.tar.gz
tar -xzf trovl_macos_arm64.tar.gz

# Move to PATH
mkdir -p "$HOME/.local/bin"
mv trovl "$HOME/.local/bin/"

# Verify
trovl version
```

---

### Windows

#### Pre-built Binary

`arm64` binaries are also available!

##### **PowerShell:**

```powershell
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_windows_amd64.zip
Expand-Archive trovl_windows_amd64.zip

mkdir "$env:USERPROFILE\.local\bin" -Force
move trovl.exe "$env:USERPROFILE\.local\bin\"

setx PATH "$env:PATH;$env:USERPROFILE\.local\bin"

trovl version
```

##### **Command Prompt:**

```cmd
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_windows_amd64.zip
tar -xf trovl_windows_amd64.zip

mkdir "%USERPROFILE%\.local\bin"
move trovl.exe "%USERPROFILE%\.local\bin\"

setx PATH "%PATH%;%USERPROFILE%\.local\bin"

trovl version
```

---

## Verifying Checksums (Optional)

Verifying checksums ensures the downloaded file hasn't been tampered with or corrupted during transfer. Compare the computed hash against the value shown on the GitHub releases page, or download the `checksums.txt` file for automated verification.

**Download checksums file:**
```bash
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/checksums.txt
```

### Linux/macOS

Compute checksum:
```bash
sha256sum trovl_linux_amd64.tar.gz
```

Automatically verify against checksums file:
```bash
sha256sum -c checksums.txt
```

### Windows

#### PowerShell

Compute checksum:
```powershell
Get-FileHash trovl_windows_amd64.zip -Algorithm SHA256
```

**Automatically verify against checksums file:**
```powershell
(Get-FileHash trovl_windows_amd64.zip).Hash -eq (Select-String trovl_windows_amd64.zip checksums.txt).Line.Split()[0]
```

#### Command Prompt

Compute checksum:
```cmd
certutil -hashfile trovl_windows_amd64.zip SHA256
```

Compare the output hash with the value on the GitHub releases page or in `checksums.txt`.

---

## Updating

To update `trovl`, repeat the installation steps with the latest release. The new binary will overwrite the old one.

**For `go install` users:**

```bash
go install github.com/sneha-afk/trovl@latest
```

**Tip:** You can check your current version and compare it against the latest release on GitHub:

```bash
trovl version
```

---

## Uninstalling

Remove the binary from its installation location:

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

If you encounter any issues, please open an issue: [github.com/sneha-afk/trovl/issues](https://github.com/sneha-afk/trovl/issues)
