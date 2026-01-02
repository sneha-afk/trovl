/*
Package links deals with core actions of handling symlinks.
*/
package links

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"unicode"

	"github.com/sneha-afk/trovl/internal/state"
	"github.com/sneha-afk/trovl/internal/utils"
)

type LinkType int

const (
	LinkFile LinkType = iota
	LinkDirectory
)

type Link struct {
	Target    string   `json:"target"`     // Real file/directory
	LinkMount string   `json:"link_mount"` // Where the symlink is
	Type      LinkType `json:"link_type"`
}

var ErrDeclinedOverwrite = errors.New("user declined overwriting existing file, no action taken")

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
func Construct(state *state.TrovlState, targetPath, symlinkPath string) (Link, error) {
	targetFileInfo, err := utils.GetPathInfo(targetPath)
	if !targetFileInfo.Exists || err != nil {
		return Link{}, fmt.Errorf("invalid target path '%v': %v", targetPath, err)
	}

	targetFile, err := os.Open(targetPath)
	if err != nil {
		return Link{}, err
	}

	var linkType LinkType
	if targetFileInfo.IsDir {
		linkType = LinkDirectory
	} else {
		linkType = LinkFile
	}

	targetFile.Close()

	symlinkInfo, err := utils.GetPathInfo(symlinkPath)
	if err != nil {
		return Link{}, fmt.Errorf("could not get symlink info: %v", err)
	}

	// Conflict: existing file at the symlink position
	if symlinkInfo.Exists {
		if state.Options.DryRun {
			state.Logger.Info("[DRY-RUN] conflict with existing file", "link", symlinkPath)
			return Link{}, nil
		}

		if !symlinkInfo.IsSymlink {
			// TODO: consider a backup feature if the existing target is a simple ordinary file
			return Link{}, fmt.Errorf("existing file at symlink is not a symlink, exiting")
		}

		shouldOverwrite := false

		if state.Options.OverwriteYes {
			shouldOverwrite = true
		} else if state.Options.OverwriteNo {
			shouldOverwrite = false
		} else {
			if symlinkInfo.TargetPath == targetPath {
				state.Logger.Warn(fmt.Sprintf("Symlink %v already exists and already points to %v, should it be overwritten? [y/N]", targetPath, symlinkPath))
			} else {
				state.Logger.Warn(fmt.Sprintf("Symlink %v already exists but points to another target (%v), should it be overwritten? [y/N]", symlinkInfo.TargetPath, symlinkPath))
			}
			fmt.Printf("> ")
			var input = 'n'
			if _, err := fmt.Scanf("%c\n", &input); err != nil {
				return Link{}, fmt.Errorf("could not read input, no action taken: %v", err)
			}
			shouldOverwrite = unicode.ToLower(input) == 'y'
		}

		if shouldOverwrite {
			state.Logger.Warn("Overwriting existing file...")
			if err := os.Remove(symlinkPath); err != nil {
				return Link{}, fmt.Errorf("could not delete existing file: %v", err)
			}
		} else {
			state.Logger.Warn("Declined overwriting existing file, no action taken")
			return Link{}, ErrDeclinedOverwrite
		}
	}

	return Link{
		Target:    targetPath,
		LinkMount: symlinkPath,
		Type:      linkType,
	}, nil
}

// Add a symlink specified by the Link class.
// Precondition: there is no existing file where the symlink was specified
func Add(state *state.TrovlState, targetPath, symlinkPath string) error {
	targetPath, err := CleanLink(targetPath, state.Options.UseRelative)
	if err != nil {
		return fmt.Errorf("invalid path (target): %v", err)
	}
	symlinkPath, err = CleanLink(symlinkPath, state.Options.UseRelative)
	if err != nil {
		return fmt.Errorf("invalid path (symlink): %v", err)
	}

	link, err := Construct(state, targetPath, symlinkPath)
	if err != nil {
		return fmt.Errorf("failed to construct link: %v", err)
	}

	if state.Options.DryRun {
		state.Logger.Info("[DRY-RUN] would create symlink", "target", targetPath, "link", symlinkPath)
		return nil
	}

	err = os.Symlink(link.Target, link.LinkMount)
	return err
}

// RemoveByPath takes in the path to a symlink to remove, while keeping the original
// file intact (note: target file is not checked for existence as the symlink is being removed.)
func RemoveByPath(state *state.TrovlState, path string) error {
	path, err := CleanLink(path, true)
	if err != nil {
		return fmt.Errorf("invalid path (symlink): %v", err)
	}
	info, err := utils.GetPathInfo(path)
	if err != nil {
		return fmt.Errorf("could not get symlink info: %v", err)
	}

	if !info.Exists {
		return fmt.Errorf("no symlink exists at %v", path)
	}

	if !info.IsSymlink {
		return fmt.Errorf("invalid symlink: %v", err)
	}

	if state.Options.DryRun {
		state.Logger.Info("[DRY-RUN] would remove symlink", "link", path)
		return nil
	}

	return os.Remove(path)
}
