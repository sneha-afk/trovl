## trovl plan

Describes what will happen during an `apply` without modifying the filesystem

### Synopsis

Describes the actions that will happen when a manifest file is applied. This is essentially
an alias for running
  trovl apply --dry-run

```
trovl plan <manifest_file> [more_manifests] [flags]
```

### Examples

```
trovl plan .trovl
```

### Options

```
  -h, --help   help for plan
```

### Options inherited from parent commands

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl](trovl.md)	 - A cross-platform symlink manager.

