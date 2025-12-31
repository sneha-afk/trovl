# trovl

A simple, cross-platform symlink manager that eliminates the pain of managing symlinks across different operating systems.

## Why trovl?

Tired of dealing with:
- Different symlink syntax across Windows, Linux, and macOS?
- Platform-specific symlink locations or names (like `~/.vimrc` vs `~/_vimrc`)?
- Manually managing dozens of symlinks for various programs and configurations?

**trovl** provides:
- **Cross-platform support** - Windows, Linux, macOS (amd64 & arm64)
- **Platform-specific overrides** - Different link paths per OS
- **True symlinks** - Uses native symlink APIs when possible
- **Schema-based management** - Define all your symlinks in one JSON file to perform bulk operations, perfect for dotfiles!

## Installation

See [INSTALL.md](./INSTALL.md) for detailed installation instructions including pre-built binaries.

**Quick install with Go:**
```bash
go install github.com/sneha-afk/trovl@latest
```

## Quick Start

### `add` a symlink
Can link to files or directories, default to *absolute* path resolution:
```bash
trovl add /path/to/target /path/to/symlink
```

### `remove` a symlink
Safely removes symlinks while keeping the target file:
```bash
trovl remove /path/to/symlink
```

### `apply` bulk operations
```bash
trovl apply .trovl.json
```

#### Defining schema

Define all your symlinks in a `json` file:

```json
{
  "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
  "links": [
    {
      "target": "~/dotfiles/.vimrc",
      "link": "~/.vimrc",
      "platforms": ["all"],
      "platform_overrides": {
        "windows": { "link": "~/_vimrc" }
      }
    },
    {
      "target": "~/dotfiles/.gitconfig",
      "link": "~/.gitconfig",
      "platforms": ["all"]
    }
  ]
}
```

Then apply all links at once:
```bash
trovl apply .trovl.json # can be named anything!
```

See the full [schema documentation](https://github.com/sneha-afk/trovl/blob/main/docs/trovl_schema.json) for all available options.

## Commands

| Command | Description |
|---------|-------------|
| `add` | Create a new symlink pointing to a target |
| `apply` | Apply multiple links from a schema file |
| `remove` | Delete a symlink (preserves target) |
| `completion` | Generate shell completion scripts |
| `help` | Display help for any command |
| `--version` | Show version information |

## Global Flags

- `--debug` - Show debug information (file paths, line numbers)
- `--verbose` - Display verbose output
- `--help` - Show help information

## Shell Completion

Generate completion scripts for your shell:

```bash
# Bash
trovl completion bash > /etc/bash_completion.d/trovl

# Zsh
trovl completion zsh > "${fpath[1]}/_trovl"

# Fish
trovl completion fish > ~/.config/fish/completions/trovl.fish

# PowerShell
trovl completion powershell > trovl.ps1
```

## Contributing

Contributions are welcome! Bug reports and pull requests can be submitted via [issues](https://github.com/sneha-afk/trovl/issues).

**Bug reports:** Please include the output of `go test ./... -v` and your platform details (OS, architecture).

## Development

<details>


**Prerequisites:**
- Go 1.21+
- [Task](https://taskfile.dev) (optional, for task runner)

<summary>Working on trovl</summary>

**Clone and build:**
```bash
git clone https://github.com/sneha-afk/trovl.git
cd trovl
go build
```

**Development tasks** (requires Task):
```bash
task --list      # List all available tasks
task fmt         # Format code
task test        # Run tests
task check       # Format + test
task build       # Build binary
task release     # Build for all platforms
```

**Recommended tooling:**
- [gopls](https://go.dev/gopls/) for IDE support and formatting

**Testing:**
```bash
go test ./... -v
```

</details>

## License

MIT

