package links

import (
	"fmt"
	"io/fs"
	"os"
	"unicode"

	"github.com/sneha-afk/trovl/internal/models"
)

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

// Construct a Link type and validate the target file exists.
func Construct(targetPath, symlinkPath string, linkType models.LinkType) (models.Link, error) {
	if valid, err := ValidatePath(targetPath); !valid || err != nil {
		return models.Link{}, fmt.Errorf("invalid path '%v': %v", targetPath, err)
	}

	if valid, err := ValidatePath(symlinkPath); valid || err == nil {
		fmt.Printf("[WARN] Construct: file %v already exists, should it be overwritten? [y/N]: ", symlinkPath)
		var input = 'n'
		_, err := fmt.Scanf("%c", &input)
		if err != nil {
			return models.Link{}, fmt.Errorf("could not read input, no action taken: %v", err)
		}

		if unicode.ToLower(input) == 'y' {
			fmt.Printf("[INFO] Construct: user accepted overwriting existing file, continuing\n")
			// TODO: double check if Linux allows direct overwriting of files
			if err := os.Remove(symlinkPath); err != nil {
				return models.Link{}, fmt.Errorf("could not deleting existing file: %v", err)
			}
		} else {
			return models.Link{}, fmt.Errorf("user declined overwriting existing file, no action taken")
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
	if valid, err := IsSymlink(path); !valid || err != nil {
		return fmt.Errorf("invalid symlink: %v", err)
	}
	return os.Remove(path)
}
