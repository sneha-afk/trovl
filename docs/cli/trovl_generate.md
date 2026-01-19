---
title: "trovl generate"
parent: Commands
slug: "trovl_generate"
description: "CLI reference for trovl generate"
---

## trovl generate

Generate a blank manifest file with the current schema (default: `$XDG_CONFIG_HOME/trovl/manifest.json`).

### Synopsis

Generate a blank manifest file with trovl's current schema. By default, this will be generated at the default location of `$XDG_CONFIG_HOME/trovl/manifest.json` (see [environment variable usage](/trovl/configuration/#environment-variables))

```
trovl generate [optional path(s)] [flags]
```

### Examples

```
trovl generate     # Default location
trovl generate here.json
```

### Options

```
  -h, --help   help for generate
```

### Options inherited from parent commands

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl](trovl.md)	 - A cross-platform symlink manager.

