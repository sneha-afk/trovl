---
title: "trovl"
parent: Commands
slug: "trovl"
description: "CLI reference for trovl"
---

## trovl

A cross-platform symlink manager.

### Synopsis

trovl is a cross-platform symlink manager that aims to make file management easier and more efficient.
It features configurable paths for files and directories that vary in location depending on the system,
and true-symlinking when possible.
	

### Options

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -h, --help      help for trovl
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl add](trovl_add.md)	 - Adds a symlink that points to the target file
* [trovl apply](trovl_apply.md)	 - Applies a manifest specified by schema (default: `$XDG_CONFIG_HOME/trovl/manifest.json`)
* [trovl generate](trovl_generate.md)	 - Generate a blank manifest file with the current schema (default: `$XDG_CONFIG_HOME/trovl/manifest.json`).
* [trovl plan](trovl_plan.md)	 - Describes what will happen during an `apply` without modifying the filesystem
* [trovl remove](trovl_remove.md)	 - Removes a specified symlink while keeping the target file as-is.

