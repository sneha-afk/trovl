package utils_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sneha-afk/trovl/internal/utils"
)

// helper to temporarily override goos
func withGOOS(osname string) func() {
	orig := utils.GOOS
	utils.GOOS = osname
	return func() { utils.GOOS = orig }
}

func getAbs(p string) string {
	path, err := filepath.Abs(p)
	if err != nil {
		os.Exit(1)
	}
	return path
}

func TestNormalizeWindowsEnvVars_RegexComprehensive(t *testing.T) {
	withGOOS("windows")

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "single percent var (cmd.exe)",
			in:   `%FOO%`,
			want: `${FOO}`,
		},
		{
			name: "multiple percent vars (cmd.exe)",
			in:   `%A%%B%%C%`,
			want: `${A}${B}${C}`,
		},
		{
			name: "percent with surrounding text",
			in:   `C:\%FOO%\bar`,
			want: `C:\${FOO}\bar`,
		},
		{
			name: "invalid percent var ignored",
			in:   `%1BAD%`,
			want: `%1BAD%`,
		},
		{
			name: "single powershell env",
			in:   `$env:FOO`,
			want: `${FOO}`,
		},
		{
			name: "mixed casing $Env:",
			in:   `$Env:FOO`,
			want: `${FOO}`,
		},
		{
			name: "powershell env with path",
			in:   `$env:FOO\bar`,
			want: `${FOO}\bar`,
		},
		{
			name: "brace env (powershell)",
			in:   `${env:FOO}`,
			want: `${FOO}`,
		},
		{
			name: "brace env with path",
			in:   `${env:FOO}\bar`,
			want: `${FOO}\bar`,
		},
		{
			name: "brace mixed casing",
			in:   `${ENV:FOO}`,
			want: `${FOO}`,
		},
		{
			name: "mixed percent + $env + ${env:}",
			in:   `%A%\$env:B\${env:C}`,
			want: `${A}\${B}\${C}`,
		},
		{
			name: "percent followed by literal $env",
			in:   `%A%\$env:B`,
			want: `${A}\${B}`,
		},
		{
			name: "text with no envs",
			in:   `C:\path\file.txt`,
			want: `C:\path\file.txt`,
		},
		{
			name: "dollar sign not env",
			in:   `price: $100`,
			want: `price: $100`,
		},
		{
			name: "percent signs not env",
			in:   `50% complete`,
			want: `50% complete`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.NormalizeWindowsEnvVars(tt.in)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCleanPath_AllCases(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("cannot get home dir: %v", err)
	}

	// Set some env vars for testing
	t.Setenv("A", "foo")
	if runtime.GOOS != "windows" {
		// other platforms are case sensitive
		t.Setenv("a", "foo")
	}
	t.Setenv("B", "bar")
	t.Setenv("C", "baz")
	t.Setenv("TRAVERASL", "../../")

	tests := []struct {
		name         string
		in           string
		useRelative  bool
		want         string
		overrideGOOS string // optional: "" means don't override
	}{
		// --- Windows normalization tests ---
		{
			name:         "cmd percent vars",
			in:           `%A%\%B%`,
			want:         "foo\\bar",
			useRelative:  true,
			overrideGOOS: "windows",
		},
		{
			name:         "unclosed percent",
			in:           "%A",
			want:         "%A",
			useRelative:  true,
			overrideGOOS: "windows",
		},
		{
			name:         "powershell env",
			in:           `$env:A\$env:B`,
			useRelative:  true,
			want:         "foo\\bar",
			overrideGOOS: "windows",
		},
		{
			name:         "powershell case insensitivity",
			in:           "${ENV:a}/",
			want:         "foo",
			useRelative:  true,
			overrideGOOS: "windows",
		},
		{
			name:         "brace powershell env",
			in:           `${env:A}\${env:B}`,
			useRelative:  true,
			want:         "foo\\bar",
			overrideGOOS: "windows",
		},
		{
			name:         "mixed syntax",
			in:           `%A%\$env:B\${env:C}`,
			useRelative:  true,
			want:         "foo\\bar\\baz",
			overrideGOOS: "windows",
		},
		{
			name:         "multiple vars no separators",
			in:           `%A%%B%%C%`,
			useRelative:  true,
			want:         "foobarbaz",
			overrideGOOS: "windows",
		},
		{
			name:         "multiple vars with separators",
			in:           `%A%\%B%\%C%`,
			useRelative:  true,
			want:         "foo\\bar\\baz", // UNC path
			overrideGOOS: "windows",
		},
		{
			name:         "undefined vars become empty",
			in:           `%NOPE%\foo`,
			useRelative:  true,
			want:         "\\foo", // UNC path
			overrideGOOS: "windows",
		},
		{
			name: "tilde alone",
			in:   "~",
			want: home,
		},
		{
			name: "tilde with path",
			in:   "~/foo",
			want: filepath.Join(home, "foo"),
		},
		{
			name:        "tilde with cleaning",
			in:          "~/foo/../bar",
			useRelative: true,
			want:        filepath.Join(home, "bar"),
		},
		{
			name:        "mid-path tilde no-op",
			in:          "foo/~/bar",
			want:        filepath.Join("foo", "~", "bar"),
			useRelative: true,
		},
		{
			name: "tilde with expanded var",
			in:   "~/%A%",
			want: filepath.Join(home, "foo"),
		},
		{
			name:        "relative path clean",
			in:          "foo/../bar",
			useRelative: true,
			want:        "bar",
		},
		{
			name:        "absolute path resolve",
			in:          "foo/../bar",
			useRelative: false,
			want:        getAbs("bar"),
		},
		{
			name:         "non-windows noop",
			in:           `%A%`,
			useRelative:  true,
			want:         filepath.Clean(`%A%`),
			overrideGOOS: "linux",
		},
		{
			name:         "pathological empty string",
			in:           "",
			want:         ".",
			useRelative:  true,
			overrideGOOS: "windows",
		},
		{
			name:         "single dot",
			in:           ".",
			want:         ".",
			useRelative:  true,
			overrideGOOS: "windows",
		},
		{
			name:         "double dot",
			in:           "..",
			want:         "..",
			useRelative:  true,
			overrideGOOS: "windows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// optionally override GOOS
			var cleanup func()
			if tt.overrideGOOS != "" {
				cleanup = withGOOS(tt.overrideGOOS)
				defer cleanup()
			}

			got, err := utils.CleanPath(tt.in, tt.useRelative)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

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
