## trovl add

Adds a symlink that points to the target file

### Synopsis

When possible, add a true symlink (as in, not a junction or hard link) to a target file.

When backing up a file that would be overwritten by this new symlink, trovl always uses $XDG_CACHE_HOME first, before
falling back to OS defaults. See [trovl's use of environment variables](../commands.md/#environment-variables) to learn more.


```
trovl add <target> <symlink> [target2, symlink2, ...] [flags]
```

### Examples

```
trovl add ~/dotfiles/.vimrc ~/.vimrc
```

### Options

```
      --backup         backup existing single files if a symlink would overwrite it
  -h, --help           help for add
      --no-backup      do not backup existing files and abandon symlink creation
      --no-overwrite   do not overwrite any existing symlinks
      --overwrite      overwrite any existing symlinks
      --relative       retain relative paths to target
```

### Options inherited from parent commands

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl](trovl.md)	 - A cross-platform symlink manager.

