package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PathInfo struct {
	Exists     bool
	IsDir      bool
	IsSymlink  bool
	TargetPath string // If this is a symlink, what is is pointing to?
}

var FileTimeFormat = "2006-01-02_15-04-05"

func GetPathInfo(path string) (PathInfo, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return PathInfo{Exists: false}, nil
		}
		return PathInfo{}, err
	}

	pi := PathInfo{
		Exists:    true,
		IsDir:     info.IsDir(),
		IsSymlink: info.Mode()&fs.ModeSymlink != 0,
	}

	if pi.IsSymlink {
		target, err := os.Readlink(path)
		if err != nil {
			return pi, err
		}
		pi.TargetPath = target
	}

	return pi, nil
}

// ValidateSymlink first ensures the symlink is indeed one at all, and that it is pointing
// to a valid target file that exists.
func ValidateSymlink(symlinkPath string) (bool, error) {
	symlinkInfo, err := GetPathInfo(symlinkPath)
	if err != nil && !(symlinkInfo.Exists && symlinkInfo.IsSymlink) {
		return false, err
	}

	targetPath, err := os.Readlink(symlinkPath)
	if err != nil {
		return false, fmt.Errorf("target file is not readable: %v", err)
	}

	if targetInfo, err := GetPathInfo(targetPath); !targetInfo.Exists || err != nil {
		return false, fmt.Errorf("could not validate target: %v", err)
	}

	return true, nil
}

// GetCacheDir is similar to os.UserCacheDir, but always uses $XDG_CACHE_HOME if it is defined,
// regardless of the OS. If this is not defined, the cache directory is that specified by
// os.UserCacheDir. Note: does NOT guarantee the directory has been created yet.
func GetCacheDir() (string, error) {
	// Prioritize XDG_CACHE_HOME if it is defined
	xdgCache := os.Getenv("XDG_CACHE_HOME")
	if xdgCache != "" {
		return filepath.Join(xdgCache, "trovl"), nil
	}

	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "trovl"), nil
}

// GetConfigDir returns the path to the trovl config directory.
// It prioritizes $XDG_CONFIG_HOME if defined, otherwise falls back to ~/.config (on all OSes)
// Note: this does NOT guarantee that the directory exists yet.
func GetConfigDir() (string, error) {
	// Prioritize XDG_CONFIG_HOME if set
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, "trovl"), nil
	}

	// Fallback to ~/.config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "trovl"), nil
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

// BackupFile copies a file into the cache directory, and returns the path it was stored to.
// Default backup directory: $XDG_CACHE_HOME/trovl/backups
func BackupFile(path, backupDir, timestampFormat string) (string, error) {
	currTimeStr := time.Now().Format(timestampFormat)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	backupFilename := fmt.Sprintf("%s_backup_%s%s", name, currTimeStr, ext)
	backupPath := filepath.Join(backupDir, backupFilename)

	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return "", fmt.Errorf("could not create backup parent directory: %v", err)
	}

	if err := CopyFile(path, backupPath); err != nil {
		return "", fmt.Errorf("could not backup file: %v", err)
	}
	return backupPath, nil
}
