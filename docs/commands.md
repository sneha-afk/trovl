---
layout: default
---

# commands

Complete reference for all trovl commands and options.

## Global Flags

These flags can be used with any command:

| Flag | Description |
|------|-------------|
| `--debug` | Show debug information for troubleshooting |
| `--dry-run` | Walk through an operation without making changes |
| `-h, --help` | Display help information |
| `-v, --verbose` | Show verbose output for actions taken |
| `--version` | Display trovl version |

## Commands

| **Command**  | **Documentation** |
|--------------|-------------------|
| `add`        | [cli/add](./cli/trovl_add.md) |
| `apply`      | [cli/apply](./cli/trovl_apply.md) |
| `plan`       | [cli/plan](./cli/trovl_plan.md) |
| `remove`     | [cli/remove](./cli/trovl_remove.md) |
| `completion` | `trovl completion --help` |
| `help`       | `trovl [command] --help` |

### completion

Generate shell completion scripts for trovl commands.

```bash
trovl completion [shell]
```

**Supported shells:**
- `bash`
- `zsh`
- `fish`
- `powershell`

**Examples:**

```bash
# Generate bash completion
trovl completion bash > /etc/bash_completion.d/trovl

# Generate zsh completion
trovl completion zsh > "${fpath[1]}/_trovl"

# Generate fish completion
trovl completion fish > ~/.config/fish/completions/trovl.fish

# Generate PowerShell completion
trovl completion powershell | Out-String | Invoke-Expression
```

See [detailed documentation](cli/trovl_completion.md).

---

### help

Display help information for any command.

```bash
trovl help [command]
```

**Examples:**

```bash
# General help
trovl help

# Help for specific command
trovl help add
trovl help apply
```

---

## Platform Support

trovl supports the following platforms for `platforms` field in manifests:

- `linux` - Linux systems
- `darwin` - macOS systems
- `windows` - Windows systems
- `all` - All platforms

## Environment Variables <a name="environment-variables"></a>

trovl respects the following environment variables:

- `XDG_CACHE_DIR` - Cache directory for backups: trovl will always respect what is set for `$XDG_CACHE_DIR` first before following back to [Go's defaults](https://pkg.go.dev/os#UserCacheDir):
    - Unix: `$HOME/.cache`
    - Darwin (macOS): `$HOME/Library/Caches`
    - Windows: `%LocalAppData%`

## Configuration Files

trovl does not currently use configuration files. All configuration is done through command-line flags and manifest files.
