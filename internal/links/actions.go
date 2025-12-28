package links

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"unicode"

	"github.com/sneha-afk/trovl/internal/models"
)

var ErrDeclinedOverwrite = errors.New("user declined overwriting existing file, no action taken")

type ConstructOptions struct {
	OverwriteForceYes bool
	OverwriteForceNo  bool
}

// ValidatePath checks if a file at the given path exists and is openable.
func ValidatePath(path string) (bool, error) {
	file, err := os.Open(path)
	file.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func IsSymlink(symlinkPath string) (bool, error) {
	if valid, err := ValidatePath(symlinkPath); !valid || err != nil {
		return false, err
	}

	symlinkInfo, err := os.Lstat(symlinkPath)
	if err != nil {
		return false, err
	}

	if symlinkInfo.Mode()&fs.ModeSymlink == 0 {
		return false, fmt.Errorf("%v is not a symlink", symlinkPath)
	}
	return true, nil
}

// ValidateSymlink first ensures the symlink is indeed one at all, and that it is pointing
// to a valid target file that exists.
func ValidateSymlink(symlinkPath string) (bool, error) {
	if valid, err := IsSymlink(symlinkPath); !valid || err != nil {
		return false, err
	}

	targetPath, err := os.Readlink(symlinkPath)
	if err != nil {
		return false, fmt.Errorf("target file is not readable: %v", err)
	}

	if valid, err := ValidatePath(targetPath); !valid || err != nil {
		return false, fmt.Errorf("could not validate target: %v", err)
	}

	return true, nil
}

// CleanLink defaults to using an absolute filepath, only relative if specified
// Guaranteed that filepath.Clean has been called before returning
func CleanLink(raw string, useRelative bool) (string, error) {
	var ret string
	var err error = nil

	// Handle issues with not dealing with "~" correctly
	if len(raw) > 0 && raw[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		raw = filepath.Join(usr.HomeDir, raw[1:])
	}

	if useRelative {
		ret = filepath.Clean(raw)
	} else {
		ret, err = filepath.Abs(raw)
	}
	return ret, err
}

// Construct a Link type and validate the target file exists.
func Construct(targetPath, symlinkPath string, useRelative bool, opts *ConstructOptions) (models.Link, error) {
	targetPath, err := CleanLink(targetPath, useRelative)
	if err != nil {
		return models.Link{}, fmt.Errorf("invalid path (target): %v", err)
	}
	symlinkPath, err = CleanLink(symlinkPath, useRelative)
	if err != nil {
		return models.Link{}, fmt.Errorf("invalid path (symlink): %v", err)
	}

	if valid, err := ValidatePath(targetPath); !valid || err != nil {
		return models.Link{}, fmt.Errorf("invalid path '%v': %v", targetPath, err)
	}

	targetFile, err := os.Open(targetPath)
	if err != nil {
		return models.Link{}, err
	}
	targetFileInfo, err := targetFile.Stat()
	if err != nil {
		return models.Link{}, fmt.Errorf("could not get target file info: %v", err)
	}

	var linkType models.LinkType
	if targetFileInfo.IsDir() {
		linkType = models.LinkDirectory
	} else {
		linkType = models.LinkFile
	}

	targetFile.Close()

	if valid, err := ValidatePath(symlinkPath); valid || err == nil {
		shouldOverwrite := false

		if opts != nil && opts.OverwriteForceYes {
			shouldOverwrite = true
		} else if opts != nil && opts.OverwriteForceNo {
			shouldOverwrite = false
		} else {
			// Ask from stdin
			fmt.Printf("[WARN] Construct: file %v already exists, should it be overwritten? [y/N]: ", symlinkPath)
			var input = 'n'
			if _, err := fmt.Scanf("%c\n", &input); err != nil {
				return models.Link{}, fmt.Errorf("could not read input, no action taken: %v", err)
			}
			shouldOverwrite = unicode.ToLower(input) == 'y'
		}

		if shouldOverwrite {
			fmt.Printf("[INFO] Construct: overwriting existing file\n")
			if err := os.Remove(symlinkPath); err != nil {
				return models.Link{}, fmt.Errorf("could not delete existing file: %v", err)
			}
		} else {
			fmt.Printf("[INFO] Construct: declined to overwrite existing file\n")
			return models.Link{}, ErrDeclinedOverwrite
		}
	}

	return models.Link{
		Target:    targetPath,
		LinkMount: symlinkPath,
		Type:      linkType,
	}, nil
}

// Add a symlink specified by the Link class.
// Wrapper around os.Symlink which is already OS-agnostic
// Precondition: there is no existing file where the symlink was specified
func Add(link models.Link) error {
	err := os.Symlink(link.Target, link.LinkMount)
	return err
}

// RemoveByPath takes in the path to a symlink to remove, while keeping the original
// file intact (note: target file is not checked for existence as the symlink is being removed.)
func RemoveByPath(path string) error {
	path, err := CleanLink(path, true)
	if err != nil {
		return fmt.Errorf("invalid path (symlink): %v", err)
	}
	if valid, err := IsSymlink(path); !valid || err != nil {
		return fmt.Errorf("invalid symlink: %v", err)
	}
	return os.Remove(path)
}
