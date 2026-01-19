/*
Package links deals with core actions of handling symlinks.
*/
package links

import (
	"errors"
	"fmt"
	"os"
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

var ErrDryRun = errors.New("no-op: running dry-run")
var ErrDeclinedOverwrite = errors.New("user declined overwriting existing file, no action taken")
var ErrDeclinedBackup = errors.New("user declined backing up exisitng file to place new symlink, no action taken")

// Construct a Link type and validate the target file exists.
func Construct(s *state.TrovlState, targetPath, symlinkPath string) (Link, error) {
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
		s.Logger.Warn("Conflict with existing file", "link", symlinkPath, "existing_is_symlink", symlinkInfo.IsSymlink, "existing_is_dir", symlinkInfo.IsDir)

		if symlinkInfo.IsSymlink {
			if symlinkInfo.TargetPath == targetPath {
				s.Logger.Info("Conflicting symlink already points to target", "path", symlinkPath, "target", targetPath)
			} else {
				s.Logger.Warn("Conflicting symlink points elsewhere", "path", symlinkPath, "current", symlinkInfo.TargetPath, "desired", targetPath)
			}

			if s.Options.DryRun {
				return Link{}, nil
			}

			shouldOverwrite := false
			if s.Options.OverwriteYes {
				shouldOverwrite = true
			} else if s.Options.OverwriteNo {
				shouldOverwrite = false
			} else {
				if symlinkInfo.TargetPath == targetPath {
					fmt.Print("Overwrite anyway? [y/N] > ")
				} else {
					fmt.Print("Overwrite? [y/N] > ")
				}

				var input = 'n'
				if _, err := fmt.Scanf("%c\n", &input); err != nil {
					return Link{}, fmt.Errorf("could not read input, no action taken: %v", err)
				}
				shouldOverwrite = unicode.ToLower(input) == 'y'
			}

			s.Logger.Info("User's decision for overwriting", "overwrite", shouldOverwrite)

			if shouldOverwrite {
				s.LogOverwrite("Overwriting existing file", "existing_path", symlinkPath)
				if err := os.Remove(symlinkPath); err != nil {
					return Link{}, fmt.Errorf("could not delete existing file: %v", err)
				}
			} else {
				s.Logger.Warn("Declined overwriting existing file, no action taken")
				return Link{}, ErrDeclinedOverwrite
			}
		} else {
			if symlinkInfo.IsDir {
				return Link{}, fmt.Errorf("existing file at conflicting symlink path is a directory, exiting")
			}

			s.Logger.Warn("Conflicting file is an ordinary (non-link) file", "existing_path", symlinkPath)

			if s.Options.DryRun {
				return Link{}, nil
			}

			shouldBackup := false
			if s.Options.BackupYes {
				shouldBackup = true
			} else if s.Options.BackupNo {
				shouldBackup = false
			} else {
				fmt.Print("Backup existing file before placing the symlink? [y/N] > ")

				var input = 'n'
				if _, err := fmt.Scanf("%c\n", &input); err != nil {
					return Link{}, fmt.Errorf("could not read input, no action taken: %v", err)
				}
				shouldBackup = unicode.ToLower(input) == 'y'
			}

			s.Logger.Info("User's decision for backup", "backup", shouldBackup)

			if shouldBackup {
				backupDir := s.Options.BackupDir
				if s.Options.BackupDir == "" {
					cacheDir, err := utils.GetCacheDir()
					if err != nil {
						return Link{}, fmt.Errorf("could not get cache directory: %v", err)
					}
					if err := os.MkdirAll(cacheDir, 0o755); err != nil {
						return Link{}, fmt.Errorf("could not create cache directory: %v", err)
					}
					backupDir = filepath.Join(cacheDir, "backups")
				}

				backupPath, err := utils.BackupFile(symlinkPath, backupDir, utils.FileTimeFormat)
				if err != nil {
					return Link{}, fmt.Errorf("could not backup file: %v", err)
				}
				s.LogSuccess("Backed up file", "backup", backupPath, "original", symlinkPath)

				if err := os.Remove(symlinkPath); err != nil {
					return Link{}, fmt.Errorf("could not delete existing file: %v", err)
				}
			} else {
				s.Logger.Warn("Declined backing up existing file, no action taken")
				return Link{}, ErrDeclinedBackup
			}
		}
	}

	s.LogSuccess("Constructed symlink before operation", "target", targetPath, "link", symlinkPath)

	return Link{
		Target:    targetPath,
		LinkMount: symlinkPath,
		Type:      linkType,
	}, nil
}

// Add a symlink specified by the Link class.
// Precondition: there is no existing file where the symlink was specified
func Add(s *state.TrovlState, targetPath, symlinkPath string) error {
	targetPath, err := utils.CleanPath(targetPath, s.Options.UseRelative)
	if err != nil {
		return fmt.Errorf("invalid path (target): %v", err)
	}
	symlinkPath, err = utils.CleanPath(symlinkPath, s.Options.UseRelative)
	if err != nil {
		return fmt.Errorf("invalid path (symlink): %v", err)
	}

	link, err := Construct(s, targetPath, symlinkPath)
	if err != nil && err != ErrDryRun {
		if errors.Is(err, ErrDeclinedOverwrite) || errors.Is(err, ErrDeclinedBackup) {
			return err
		}
		return fmt.Errorf("failed to construct link: %v", err)
	}

	if s.Options.DryRun {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(link.LinkMount), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}
	return os.Symlink(link.Target, link.LinkMount)
}

// RemoveByPath takes in the path to a symlink to remove, while keeping the original
// file intact (note: target file is not checked for existence as the symlink is being removed.)
func RemoveByPath(s *state.TrovlState, path string) error {
	path, err := utils.CleanPath(path, true)
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

	if s.Options.DryRun {
		return nil
	}
	return os.Remove(path)
}
