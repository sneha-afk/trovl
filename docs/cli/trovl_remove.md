## trovl remove

Removes a specified symlink while keeping the target file as-is.

### Synopsis

Removes symlinks while keeping the target file untouched. Validates any argument passed
	in as truly being a symlink to prevent data loss.

```
trovl remove <symlink> [more_symlinks] [flags]
```

### Examples

```
trovl remove ~/.vimrc (where it is a symlink)
```

### Options

```
  -h, --help   help for remove
```

### Options inherited from parent commands

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl](trovl.md)	 - A cross-platform symlink manager.

