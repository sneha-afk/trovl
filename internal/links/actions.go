package links

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/sneha-afk/trovl/internal/models"
)

// ValidatePath checks if a file at the given path exists and is openable.
func ValidatePath(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	file.Close()
	return true, nil
}

func ValidateSymlink(symlinkPath string) (bool, error) {
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

	// TODO: reeval this logic, does it matter if the target is invalid?
	targetPath, err := os.Readlink(symlinkPath)
	if err != nil {
		return false, err
	}

	if valid, err := ValidatePath(targetPath); !valid || err != nil {
		return false, err
	}

	return true, nil
}

// Construct a Link type and validate the target file exists.
func Construct(targetPath, symlinkPath string, linkType models.LinkType) models.Link {
	if valid, err := ValidatePath(targetPath); !valid || err != nil {
		log.Fatalf("[ERROR] Construct: invalid path '%v': %v\n", targetPath, err)
	}

	// TODO: if the symlink file already exists, ask for user input to overwrite?
	if valid, err := ValidatePath(symlinkPath); valid || err == nil {
		log.Printf("[WARN] Construct: file %v already exists, overwriting with new symlink", symlinkPath)
	}

	return models.Link{
		Target:    targetPath,
		LinkMount: symlinkPath,
		Type:      linkType,
	}
}

// Add a symlink specified by the Link class.
// Wrapper around os.Symlink which is already OS-agnostic
func Add(link models.Link) error {
	// TODO: use linktype to debug Windows directory caveats
	err := os.Symlink(link.Target, link.LinkMount)
	return err
}

// RemoveByPath takes in the path to a symlink to remove, while keeping the original
// file intact.
func RemoveByPath(path string) error {
	if valid, err := ValidateSymlink(path); !valid || err != nil {
		log.Fatalln("[ERROR]: RemoveByPath: invalid symlink: ", err)
	}
	return os.Remove(path)
}
