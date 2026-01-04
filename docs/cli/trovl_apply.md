## trovl apply

Applies a link list specified by schema.

### Synopsis

Applies a link list specified by schema to bulk add links or fix as needed.

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

