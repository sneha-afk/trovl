---
layout: default
---

# configuration

## Environment Variables <a name="environment-variables"></a>

trovl respects the following environment variables:

- `XDG_CACHE_HOME` - Cache directory for backups: trovl will always respect what is set for `$XDG_CACHE_HOME` first before following back to [Go's defaults](https://pkg.go.dev/os#UserCacheDir):
    - Unix: `$HOME/.cache`
    - Darwin (macOS): `$HOME/Library/Caches`
    - Windows: `%LocalAppData%`
- `XDG_CONFIG_HOME` - Config directory will be located at `$XDG_CONFIG_HOME/trovl`
    - On all systems, the value set for `$XDG_CONFIG_HOME` is scouted first before falling back to `~/.config`

## Default Manifest

On commands where trovl expects manifest files, if no such path is passed in, trovl will attempt to read from `$XDG_CONFIG_HOME/trovl/manifest.json`.

> Tip: perhaps your dotfile setup can be stored in the default manifest!
