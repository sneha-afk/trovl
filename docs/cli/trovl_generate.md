## trovl generate

Generate a blank manifest file with the current schema.

### Synopsis

Generate a blank manifest file with trovl's current schema. By default, this will be generated at the default location of
$XDG_CONFIG_HOME/trovl/manifest.json (see [environment variable usage](../configuration/#environment-variables))

```
trovl generate [flags]
```

### Options

```
  -h, --help          help for generate
  -p, --path string   path to output manifest file
```

### Options inherited from parent commands

```
      --debug     show debug info
      --dry-run   walk through an operation without making changes
  -v, --verbose   have verbose outputs for actions taken
```

### SEE ALSO

* [trovl](trovl.md)	 - A cross-platform symlink manager.

