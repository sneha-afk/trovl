---
layout: default
title: Configuration
nav_order: 4
---

# configuration

## Environment Variables <a name="environment-variables"></a>

`trovl` adheres to the XDG specification where applicable.

On **all platforms**, XDG environment variables are **always honored first**.

Fallback paths are used **only when the corresponding XDG variable is unset**.

### `XDG_CACHE_HOME`

Defines the base directory for cache and backup data.

* If set, `trovl` uses `XDG_CACHE_HOME` directly.
* If unset, `trovl` falls back to Go’s platform defaults (`os.UserCacheDir`):
    * **Linux / Unix:** `$HOME/.cache`
    * **macOS (Darwin):** `$HOME/Library/Caches`
    * **Windows:** `%LocalAppData%`

Backups are stored at `<cache-dir>/trovl/backups`.

### `XDG_CONFIG_HOME`

Defines the base directory for configuration files.

* If set, the config directory is `XDG_CONFIG_HOME/trovl`.
* If unset, the config directory falls back to `~/.config/trovl` on all platforms.

## Manifests

A **manifest** describes which symlinks `trovl` should create and on which platforms.

Manifests are JSON files. See [`trovl apply`](/trovl/cli/trovl_apply/) for details on applying one.

### Minimal manifest

Only two fields are required per link:

- `target`: what the symlink points to
- `link`: where the symlink is created

```json
{
  "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
  "links": [
    {
      "target": "~/dotfiles/.vimrc",
      "link": "~/.vimrc"
    }
  ]
}
```

### Optional fields

The default values are shown:

* `relative = false`: use absolute paths
* `platforms = ["all"]`: apply everywhere
* `platform_overrides = {}`: no per-platform overrides

Supported platform values:

* `linux`
* `darwin`
* `windows`
* `wsl`
* `all` (implicit if no platforms list is specified)

---

### Default manifest location

If no manifest path is provided, `trovl` reads from **`$XDG_CONFIG_HOME/trovl/manifest.json`**

{: .highlight }
> Tip: This works well as a dotfiles manifest!

Generate one with [`trovl generate`](/trovl/docs/cli/trovl_generate.md)

```bash
# to the default location
trovl generate
```

---

### Example

This manifest (located at [`example_manifest.json`](https://github.com/sneha-afk/trovl/blob/main/docs/example_manifest.json))

* Links `.vimrc`, with a Windows override
* Uses a relative link for trovl’s config directory

```json
{
  "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
  "links": [
    {
      "target": "~/dotfiles/dot-home/.vimrc",
      "link": "~/.vimrc",
      "platform_overrides": {
        "windows": {
          "link": "~/_vimrc"
        }
      }
    },
    {
      "target": "./.config/trovl",
      "link": "~/.trovl_dir",
      "relative": true
    }
  ]
}
```

