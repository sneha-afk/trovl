# trovl

trovl is a command-line tool designed to simplify the management of symbolic links across different operating systems. It helps you maintain consistent file configurations across Linux, macOS, and Windows by handling platform-specific path differences automatically.

## Why trovl?

Tired of dealing with:
- Different symlink syntax across Windows, Linux, and macOS?
- Platform-specific symlink locations or names (like `~/.vimrc` vs `~/_vimrc`)?
- Manually managing dozens of symlinks for various programs and configurations?

### And why not others?

`trovl` is *symlink-first* and keeps the current filesystem as the source of truth for actions to take. Running `trovl` is primarily through explicit commands and manifest-driven sequences which allows for easy migrations and clear intentions.

[GNU `stow`](https://www.gnu.org/software/stow/) is the most direct inspiration for `trovl` and its design principles.


## Key Features

- **Cross-platform support**: Works seamlessly on Linux, macOS, and Windows on both amd64 and arm64 architectures
- **Manifest-based configuration**: Define all your symlinks in a single JSON file
- **Platform-specific overrides**: Different symlink locations per operating system
- **Safe operations**: Built-in backup and dry-run modes
- **True symlinking**: Uses native symlink functionality when available

## Use Cases

- **Dotfile management**: Keep your configuration files synchronized across machines
- **Development environments**: Link project files to standard locations
- **System configuration**: Manage system-wide symlinks declaratively
- **Cross-platform workflows**: Maintain the same setup on different operating systems

## Quick Example

```json
{
  "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
  "links": [
    {
      "target": "~/.config/myapp/config.json",
      "link": "~/myapp.json",
      "platforms": ["linux"]
    }
  ]
}
```

```bash
trovl apply manifest.json
```

## Getting Started

- [Installation](installation.md) - Install trovl on your system
- [Quick Start](quickstart.md) - Get up and running in minutes
- [Commands](commands.md) - Complete command reference
- [Examples](examples.md) - Real-world usage examples

## Project Links

- [GitHub Repository](https://github.com/sneha-afk/trovl)
- [Issue Tracker](https://github.com/sneha-afk/trovl/issues)
- [Releases](https://github.com/sneha-afk/trovl/releases)
