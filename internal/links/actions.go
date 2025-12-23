package links

import (
	"os"

	"github.com/sneha-afk/trovl/internal/models"
)

// Add a symlink specified by the Link class.
// Wrapper around os.Symlink which is already OS-agnostic
func Add(link models.Link) error {
	err := os.Symlink(link.Target, link.LinkMount)
	return err
}
