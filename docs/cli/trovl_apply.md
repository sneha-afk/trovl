---
title: "trovl apply"
parent: Commands
slug: "trovl_apply"
description: "CLI reference for trovl apply"
---

## trovl apply

Applies a manifest specified by schema (default: `$XDG_CONFIG_HOME/trovl/manifest.json`)

### Synopsis

Applies a manifest specified by schema to bulk add or fix links as needed.

By default, trovl looks for a manifest in `$XDG_CONFIG_HOME/trovl/manifest.json` If `$XDG_CONFIG_HOME` is not set, trovl then checks `~/.config/trovl/manifest.json` on all systems. If any manifest is specified into the command, the default manifest file is not applied(i.e, this process happens when invoking `trovl apply` with no arguments).
See [trovl's use of environment variables](/trovl/configuration/#environment-variables) to learn more on how these are determined.

Similar to the add command:
- If a symlink already exists at the specified location, the user will be prompted on if they want to overwrite it with the new link.
- If a directory already exists at the specified location for the symlink, an error will occur.
- If a single, ordinary file already exists at the specified location for the symlink, the user will be prompted on if they want to backup the file.

When backing up a file that would be overwritten by this new symlink, trovl always uses `$XDG_CACHE_HOME` first, before
falling back to OS defaults. The backup directory is `$XDG_CACHE_HOME/trovl/backups`.

```
trovl apply <manifest_file> [more_manifests] [flags]
```

### Examples

```
trovl apply .trovl
```

### Options

```
      --backup              backup existing single files if a symlink would overwrite it
      --backup-dir string   specify where to backup files (default: $XDG_CACHE_HOME/trovl/backups)
  -h, --help                help for apply
      --no-backup           do not backup existing files and abandon symlink creation
      --no-overwrite        do not overwrite any existing symlinks
      --overwrite           overwrite any existing symlinks
```

### Options inherited from parent commands

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl](trovl.md)	 - A cross-platform symlink manager.

