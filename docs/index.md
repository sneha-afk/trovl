---
layout: default
---

# trovl

trovl is a cross-platform CLI for managing symbolic links in a predictable, declarative way.

## What it does

- Creates and manages symlinks across Linux, macOS, Windows, and WSL
- Uses a JSON manifest to define links in one place
- Supports platform-specific paths and overrides
- Provides dry-run and backup options for safer changes

trovl is *symlink-first*: the existing filesystem is treated as the source of truth, and actions are taken explicitly via commands or manifests.

It’s inspired by [GNU stow](https://www.gnu.org/software/stow/), but currently uses direct paths instead of directory structure to define links.

## Example

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

* [Installation](./install.md)
* [Quickstart](./quickstart.md)
* [Configuration](./configuration.md)
* [Commands](./commands.md)

## Links

* [GitHub](https://github.com/sneha-afk/trovl)
* [Issues](https://github.com/sneha-afk/trovl/issues)
* [Releases](https://github.com/sneha-afk/trovl/releases)

