<#
.SYNOPSIS
    Install script for trovl

.DESCRIPTION
    Downloads and installs the latest version of trovl for Windows.
    Automatically detects architecture and installs to ~/.local/bin by default.

.PARAMETER InstallDir
    Custom installation directory (default: $env:USERPROFILE\.local\bin)

.EXAMPLE
    irm https://raw.githubusercontent.com/sneha-afk/trovl/main/install.ps1 | iex

.EXAMPLE
    .\install.ps1 -InstallDir "C:\bin"
#>

param(
    [string]$InstallDir = "$env:USERPROFILE\.local\bin"
)

$ErrorActionPreference = "Stop"

$Repo = "sneha-afk/trovl"
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { "$env:USERPROFILE\.local\bin" }
$BinaryName = "trovl.exe"

function Write-Info { Write-Host "[INFO] $args" -ForegroundColor Green }
function Write-Warn { Write-Host "[WARN] $args" -ForegroundColor Yellow }
function Write-Err { Write-Host "[ERROR] $args" -ForegroundColor Red; exit 1 }

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { Write-Err "Unsupported architecture: $arch. See https://github.com/$Repo for manual installation." }
    }
}

Write-Info "Installing $BinaryName..."

$Arch = Get-Architecture
Write-Info "Detected: Windows ($Arch)"

$Filename = "trovl_windows_$Arch.zip"
$Url = "https://github.com/$Repo/releases/latest/download/$Filename"

$TmpDir = New-Item -ItemType Directory -Path ([System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), [System.Guid]::NewGuid()))

try {
    Write-Info "Downloading from $Url..."
    Invoke-WebRequest -Uri $Url -OutFile "$TmpDir\$Filename"

    Write-Info "Extracting..."
    Expand-Archive -Path "$TmpDir\$Filename" -DestinationPath $TmpDir -Force

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Move-Item -Path "$TmpDir\trovl.exe" -Destination "$InstallDir\$BinaryName" -Force

    Write-Info "Installed to $InstallDir\$BinaryName"

    if ($env:PATH -notlike "*$InstallDir*") {
        Write-Warn "$InstallDir is not in your PATH"
        Write-Warn "Run this command to add it:"
        Write-Host "  setx PATH `"`$env:PATH;$InstallDir`"" -ForegroundColor Cyan
    }

    try {
        $version = & "$InstallDir\$BinaryName" --version 2>&1
        Write-Info "Installation successful!"
        Write-Host $version
    } catch {
        Write-Warn "Binary installed but verification failed"
    }
} finally {
    Remove-Item -Path $TmpDir -Recurse -Force
}
