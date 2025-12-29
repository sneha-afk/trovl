# trovl

Do you find yourself wrangling symlinks across different OSes with varying syntax, positions, or otherwise pain?
`trovl` aims to be a simple, cross-platform symlink manager with:
- Cross-platform support: Windows, Linux, macOS in both amd/arm (64-bit systems only)
- Override a symlink's position and name by OS
- True symlinks when possible
- Schema-based bulk creations of symlinks (perfect for your dotfiles!)


## Installation

### From Source

```bash
go install github.com/sneha-afk/trovl@latest
```

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/sneha-afk/trovl/releases).

For example, getting the latest `linux-amd64` build with `curl`: simply change the ending OS and architecture as needed
```bash
curl -LO https://github.com/sneha-afk/trovl/releases/latest/download/trovl_linux_amd64.exe
```

## Commands

- `add` - Create a new symlink pointing to a target file or directory
- `apply` - Apply a list of links defined in a schema file
- `remove` - Delete a symlink while preserving the target
- `completion` - Generate shell completion scripts for bash, zsh, fish, or powershell
- `help` - Display help information for any command
- `--version` - See current version

## Global Flags

- `--debug` - Show debug information including file paths and line numbers
- `--verbose` - Display verbose output for all operations
- `--help` - Show help information

## Usage

### Adding Links

Create a symlink that points to a target file or directory:


```bash
trovl add /path/to/target /path/to/symlink
```

### Removing Links

Remove a symlink without touching the target file:

```bash
trovl remove mylink
```

### Applying Schemas

Apply multiple links from a JSON schema file:

```bash
trovl apply schema.json
```

See the [schema](https://github.com/sneha-afk/trovl/blob/main/docs/trovl_schema.json) to see all possible options.


The following is an example `.trovl.json`:

```json
{
    "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
    "links": [
        {
            "target": "~/dotfiles/dot-home/.vimrc",
            "link": "~/.vimrc",
            "platforms": ["all"], // its fine if this includes the platforms overwritten!
            "platform_overrides": {
                "windows": { "link": "~/_vimrc" }
            }
        },
        {
            "target": "~/dotfiles/dot-home/.gitconfig",
            "link": "~/.gitconfig",
            "platforms": ["all"]
        }
    ]
}

```


## Development

### Prerequisites

- Go 1.25 or later
- [Task](https://taskfile.dev) (optional, for task runner)

Run `task --list` to see all available development tasks.

<details>
<summary>
Building and Workflows
</summary>

### Building from Source

```bash
git clone https://github.com/yourusername/trovl.git
cd trovl

go build
# Or using Task
task build
```

### Development Workflow

[gopls](https://go.dev/gopls/) is recommended for automatic formatting and completions.

```bash
# Run formatter
task fmt

# Run tests
task test

# Format + run tests
task check
```

### Building Release Binaries

```bash
# Build for all platforms
task release

# Binaries will be in dist/ with checksums
```

</details>

## Contributing and Bugs

`trovl` is developed on Windows 11 and verified both on Windows and Ubuntu via WSL on x86-64 architecture.

Any bug reports and contributions are more than welcome, especially on platforms I have not gotten to testing.
A simple dump of `go test ./... -v` is the best!

## License

MIT

