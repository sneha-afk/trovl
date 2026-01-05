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

## Manifests

A manifest contains a list of links to apply. These take in fields for which platforms the link should be applied on, and
overrides for certain platforms. A manifest file is in JSON format for ease of use, see [`trovl apply`](/trovl/cli/trovl_apply/)
to see how to "run" a manifest.

### Writing a manifest

Defining the `$schema` property allows an IDE and/or LSPs to validate your JSON file according to a defined schema:

```
{
  "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
}
```

After this, the manifest is simply a list of links to define:
```json
{
  "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
  "links": [
    {
      "target": "example_target",
      "link": "example_symlink",
      "platforms": [
        "all"
      ],
      "relative": false,
      "platform_overrides": {
        "linux": {
          "link": "example_override"
        }
      }
    }
  ]
}
```

The *only* required attributes of a link are the *target* (what file/directory does the symlink point to?) and the *link* (position and name of the symlink, its path).

The default attributes for the other keys are:
- `relative = false`: links should be constructed using their absolute paths from the root of the user's filesystem
- `platforms = ["all"]`: this link applies to all platforms
- `platform_overrides = {}`: no overrides for any platform

trovl supports the following platforms for `platforms` field in manifests:

- `linux` - Linux systems
- `darwin` - macOS systems
- `windows` - Windows systems
- `all` - All platforms


### Default manifest

On commands where trovl expects manifest files, if no such path is passed in, trovl will attempt to read from `$XDG_CONFIG_HOME/trovl/manifest.json`.

> Tip: perhaps your dotfile setup can be stored in the default manifest!

You can generate a starting manifest by using `trovl generate` which will be in the default location, or specify where with `trovl generate <path>`.


### Example

The following manifest sets up `.vimrc` with an override for Windows, and a relative link for trovl's configuration directory to home (by running this command *at home*, hence the `./`):
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

This example is hosted at [example_manifest.json](https://github.com/sneha-afk/trovl/blob/main/docs/example_manifest.json)

