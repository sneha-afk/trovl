package links

import (
	"log"
	"os"

	"github.com/sneha-afk/trovl/internal/models"
)

func Add(link models.Link) error {
	if link.Type == models.LinkFile {
		if err := os.Symlink(link.Target, link.LinkMount); err != nil {
			return err
		}
		log.Printf("Add: completed %v -> %v\n", link.LinkMount, link.Target)
	} else {
		// TODO: deal with directories
	}
	return nil
}
