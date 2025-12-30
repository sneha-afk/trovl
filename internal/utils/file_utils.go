package utils

import (
	"fmt"
	"io/fs"
	"os"
)

type PathInfo struct {
	Exists     bool
	IsDir      bool
	IsSymlink  bool
	TargetPath string // If this is a symlink, what is is pointing to?
}

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
