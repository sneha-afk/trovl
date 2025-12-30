package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sneha-afk/trovl/internal/utils"
)

func TestGetPathInfo(t *testing.T) {
	tmp := t.TempDir()

	filePath := filepath.Join(tmp, "file.txt")
	os.WriteFile(filePath, []byte("hello"), 0644)

	dirPath := filepath.Join(tmp, "subdir")
	os.Mkdir(dirPath, 0755)

	linkPath := filepath.Join(tmp, "link_to_file")
	os.Symlink(filePath, linkPath)

	tests := []struct {
		name    string
		path    string
		want    utils.PathInfo
		wantErr bool
	}{
		{
			name: "existing file",
			path: filePath,
			want: utils.PathInfo{Exists: true, IsDir: false, IsSymlink: false},
		},
		{
			name: "existing directory",
			path: dirPath,
			want: utils.PathInfo{Exists: true, IsDir: true, IsSymlink: false},
		},
		{
			name: "symlink",
			path: linkPath,
			want: utils.PathInfo{Exists: true, IsDir: false, IsSymlink: true, TargetPath: filePath},
		},
		{
			name: "non-existent",
			path: filepath.Join(tmp, "ghost.txt"),
			want: utils.PathInfo{Exists: false},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := utils.GetPathInfo(tc.path)

			if (err != nil) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tc.wantErr)
			}

			if got != tc.want {
				t.Errorf("got:  %+v\nwant: %+v", got, tc.want)
			}
		})
	}
}
