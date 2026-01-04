---
layout: default
---

# quickstart

This guide will get you up and running with trovl in minutes.

## Prerequisites

- trovl installed on your system (see [Installation](installation.md))
- Basic familiarity with symbolic links
- Editor for manifest files

## Your First Symlink

### Step 1: Create a target file

First, create a file that you want to symlink:

```bash
echo "Hello trovl" > ~/myfile.txt
```

### Step 2: Add a symlink

Use trovl to create a symlink to this file:

```bash
trovl add ~/myfile.txt ~/Desktop/mylink.txt
```

This creates a symlink at `~/Desktop/mylink.txt` that points to `~/myfile.txt`.

### Step 3: Verify the symlink

Check that the symlink works:

```bash
cat ~/Desktop/mylink.txt
# Output: Hello trovl
```

## Using Manifests

For managing multiple symlinks, manifests are more powerful than individual `add` commands since these can be run altogether.

### Step 1: Create a manifest file

Create a file named `manifest.json` (note: the name *does not matter*):

```json
{
  "$schema": "https://github.com/sneha-afk/trovl/raw/main/docs/trovl_schema.json",
  "links": [
    {
      "target": "~/Documents/notes.txt",
      "link": "~/Desktop/notes.txt"
    },
    {
      "target": "~/.config/myapp",
      "link": "~/myapp-config"
    }
  ]
}
```

> Note: trovl never creates dangling links, so if these two targes do not exist, you will probably get an error.

### Step 2: Preview the changes

Before applying, check what will happen:

```bash
trovl plan manifest.json
```

This shows you what symlinks will be created without modifying anything.

### Step 3: Apply the manifest

Create the symlinks:

```bash
trovl apply manifest.json
```

If files already exist at the link locations, trovl will prompt you to choose whether to backup or skip them.

If a file already exists at the link location, two things happen depending on what type of file it is:
- A directory: nope, trovl will exit with an error.
- An ordinary file: trovl will ask if you would like to *backup* the file (see more details in [`add`](./cli/trovl_add.md))
- A symlink: trovl will ask if you would like to *overwrite* the symlink with a new one

## Platform-Specific Configuration

One of trovl's most powerful features is platform-specific configuration.

Create a manifest that works across different operating systems:

```json
{
  "links": [
    {
      "target": "~/dotfiles/bashrc",
      "link": "~/.bashrc",
      "platforms": ["linux", "darwin"]
    },
    {
      "target": "~/dotfiles/config",
      "link": "~/.config/myapp",
      "platforms": ["all"],
      "platform_overrides": {
        "windows": {
          "link": "~/AppData/Local/myapp/config"
        }
      }
    }
  ]
}
```

When you run `trovl apply manifest.json`, it will:
- Create the appropriate symlinks for your current platform
- Skip symlinks that don't apply to your OS
- Use overrides when specified

## Common Options

### Dry run

Test operations without making changes:

```bash
trovl apply manifest.json --dry-run
```

### Verbose output

See detailed information about what trovl is doing:

```bash
trovl apply manifest.json --verbose
```

## Next Steps

- Explore the [Commands](commands.md) reference for all available options
- Learn about manifest schema details in the CLI documentation
