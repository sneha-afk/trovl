## trovl apply

Applies a manifest specified by schema (default: $XDG_CONFIG_HOME/trovl/manifest.json)

### Synopsis

Applies a manifest specified by schema to bulk add or fix links as needed.

By default, trovl looks for a manifest in $XDG_CONFIG_HOME/trovl/manifest.json (typically ~/.config/trovl/manifest.json). If $XDG_CONFIG_HOME
is not set, trovl then checks ~/.config/trovl/manifest.json (on all OSes). If any manifest is specified into the command, the default
manifest file is not applied (i.e, this process happens only upon trovl apply)

When backing up a file that would be overwritten by this new symlink, trovl always uses $XDG_CACHE_HOME first, before
falling back to OS defaults. See [trovl's use of environment variables](../configuration/#environment-variables) to learn more.


```
trovl apply <manifest_file> [more_manifests] [flags]
```

### Examples

```
trovl apply .trovl
```

### Options

```
      --backup         backup existing single files if a symlink would overwrite it
  -h, --help           help for apply
      --no-backup      do not backup existing files and abandon symlink creation
      --no-overwrite   do not overwrite any existing symlinks
      --overwrite      overwrite any existing symlinks
```

### Options inherited from parent commands

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl](trovl.md)	 - A cross-platform symlink manager.

