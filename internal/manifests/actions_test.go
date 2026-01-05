package manifests

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sneha-afk/trovl/internal/state"
	"github.com/sneha-afk/trovl/internal/utils"
)

var differentOS string
var teststate *state.TrovlState
var backupDir string

var (
	validSingleLink              = `{"links":[{"target":"actual_file","link":"test_symlink"}]}`
	validMultipleLinks           = `{"links":[{"target":"actual1","link":"symlink1"},{"target":"actual2","link":"symlink2"}]}`
	manifestWithOverride         = `{"links":[{"target":"actual_file","link":"default_symlink","platforms":["linux", "windows"],"platform_overrides":{"` + runtime.GOOS + `":{"link":"override_symlink"}}}]}`
	emptyManifest                = `{"links":[]}`
	allPlatformsManifest         = `{"links":[{"target":"actual_file","link":"symlink","platforms":["all"]}]}`
	specificPlatformsManifest    = `{"links":[{"target":"actual_file","link":"symlink","platforms":["linux","darwin","windows"]}]}`
	relativeLink                 = `{"links":[{"target":"actual_file","link":"symlink","relative":true}]}`
	absoluteLink                 = `{"links":[{"target":"actual_file","link":"symlink","relative":false}]}`
	multipleOverrides            = `{"links":[{"target":"actual_file","link":"default_symlink","platforms":["all"],"platform_overrides":{"linux":{"link":"linux_symlink"},"darwin":{"link":"darwin_symlink"},"windows":{"link":"windows_symlink"}}}]}`
	complexManifest              = `{"links":[{"target":"actual1","link":"symlink1","platforms":["all"]},{"target":"actual2","link":"symlink2","platforms":["` + runtime.GOOS + `"],"relative":true},{"target":"actual3","link":"symlink3","platform_overrides":{"` + runtime.GOOS + `":{"link":"override_symlink3"}}}]}`
	currentPlatformOnly          = `{"links":[{"target":"actual_file","link":"symlink","platforms":["` + runtime.GOOS + `"]}]}`
	differentPlatformOnly        string
	overrideForDifferentOS       string
	overrideRemovesDefault       = `{"links":[{"target":"actual_file","link":"default_symlink","platforms":["all"],"platform_overrides":{"` + runtime.GOOS + `":{"link":"override_symlink"}}}]}`
	nonexistentSource            = `{"links":[{"target":"nonexistent_file","link":"symlink"}]}`
	multipleOverridesWithCurrent = `{"links":[{"target":"actual_file","link":"default_symlink","platforms":["linux","darwin"],"platform_overrides":{"linux":{"link":"linux_symlink"},"darwin":{"link":"darwin_symlink"},"windows":{"link":"windows_symlink"}}}]}`

	invalidPlatform           = `{"links":[{"target":"actual_file","link":"test_symlink", "platforms":["lolos"]}]}`
	invalidPlatformInOverride = `{"links":[{"target":"actual_file","link":"default_symlink","platforms":["linux", "windows"],"platform_overrides":{"lolos":{"link":"override_symlink"}}}]}`
	invalidJSONSyntax         = `{"links":[{"target":"test","link":}]}`
	invalidJSONStructure      = `["not", "an", "object"]`
	malformedJSON             = `{this is not json}`
	emptyFile                 = ``
)

func TestMain(m *testing.M) {
	for os := range allSupportedPlatforms.Iter() {
		if runtime.GOOS != os {
			differentOS = os
			differentPlatformOnly = `{"links":[{"target":"actual_file","link":"symlink","platforms":["` + differentOS + `"]}]}`
			overrideForDifferentOS = `{"links":[{"target":"actual_file","link":"default_symlink","platforms":["all"],"platform_overrides":{"` + differentOS + `":{"link":"override_symlink"}}}]}`
			break
		}
	}
	if differentOS == "" {
		panic("couldn't get a different os for manifests tests")
	}
	teststate = state.DefaultState()
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
		validate    func(*testing.T, *Manifest)
	}{
		{
			name:    "valid manifest with single link",
			content: validSingleLink,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if len(m.Links) != 1 {
					t.Errorf("expected 1 link, got %d", len(m.Links))
				}
				if m.Links[0].Target != "actual_file" {
					t.Errorf("expected target actual_file, got %s", m.Links[0].Target)
				}
				if len(m.Links[0].Platforms) != 1 || m.Links[0].Platforms[0] != "all" {
					t.Errorf("expected platforms [all], got %v", m.Links[0].Platforms)
				}
			},
		},
		{
			name:    "valid manifest with multiple links",
			content: validMultipleLinks,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if len(m.Links) != 2 {
					t.Errorf("expected 2 links, got %d", len(m.Links))
				}
				for i := range m.Links {
					if len(m.Links[i].Platforms) != 1 || m.Links[i].Platforms[0] != "all" {
						t.Errorf("link %d: expected platforms [all], got %v", i, m.Links[i].Platforms)
					}
				}
			},
		},
		{
			name:    "manifest with platform overrides",
			content: manifestWithOverride,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if _, exists := m.Links[0].PlatformOverrides[runtime.GOOS]; !exists {
					t.Error("expected platform override for current OS")
				}
				if len(m.Links[0].Platforms) != 2 {
					t.Errorf("expected 2 platforms, got %d", len(m.Links[0].Platforms))
				}
			},
		},
		{
			name:    "empty manifest",
			content: emptyManifest,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if len(m.Links) != 0 {
					t.Errorf("expected 0 links, got %d", len(m.Links))
				}
			},
		},
		{
			name:    "manifest with all platforms specified",
			content: allPlatformsManifest,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if len(m.Links[0].Platforms) != 1 || m.Links[0].Platforms[0] != "all" {
					t.Errorf("expected platforms [all], got %v", m.Links[0].Platforms)
				}
			},
		},
		{
			name:    "manifest with specific platforms",
			content: specificPlatformsManifest,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if len(m.Links[0].Platforms) != 3 {
					t.Errorf("expected 3 platforms, got %d", len(m.Links[0].Platforms))
				}
			},
		},
		{
			name:    "manifest with relative link",
			content: relativeLink,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if !m.Links[0].Relative {
					t.Error("expected relative to be true")
				}
			},
		},
		{
			name:    "manifest with multiple overrides",
			content: multipleOverrides,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if len(m.Links[0].PlatformOverrides) != 3 {
					t.Errorf("expected 3 platform overrides, got %d", len(m.Links[0].PlatformOverrides))
				}
			},
		},
		{
			name:    "complex manifest with mixed configurations",
			content: complexManifest,
			wantErr: false,
			validate: func(t *testing.T, m *Manifest) {
				if len(m.Links) != 3 {
					t.Errorf("expected 3 links, got %d", len(m.Links))
				}
			},
		},
		{
			name:        "invalid platform (json unmarshal checks)",
			content:     invalidPlatform,
			wantErr:     true,
			errContains: "unsupported platform",
		},
		{
			name:        "invalid platform in an override (json unmarshal checks)",
			content:     invalidPlatformInOverride,
			wantErr:     true,
			errContains: "unsupported platform",
		},
		{
			name:        "invalid JSON syntax",
			content:     invalidJSONSyntax,
			wantErr:     true,
			errContains: "unmarshal",
		},
		{
			name:        "invalid JSON structure",
			content:     invalidJSONStructure,
			wantErr:     true,
			errContains: "unmarshal",
		},
		{
			name:        "malformed JSON",
			content:     malformedJSON,
			wantErr:     true,
			errContains: "unmarshal",
		},
		{
			name:        "empty file",
			content:     emptyFile,
			wantErr:     true,
			errContains: "unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			manifestPath := filepath.Join(tmpDir, "manifest.json")

			if err := os.WriteFile(manifestPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create manifest file: %v", err)
			}

			m, err := New(manifestPath)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error from New(), got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error from New(): %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, m)
			}
		})
	}
}

func TestNew_FileErrors(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		_, err := New("/non/existent/path/manifest.json")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
		if !strings.Contains(err.Error(), "could not read manifest file") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("directory instead of file", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := New(tmpDir)
		if err == nil {
			t.Error("expected error when reading directory as file")
		}
	})
}

func TestApply(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		options  *state.TrovlOptions
		setup    func(string)
		validate func(*testing.T, string)
	}{
		{
			name:    "single link with all platforms",
			content: validSingleLink,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "test_symlink")
				info, err := os.Lstat(symlinkPath)
				if err != nil {
					t.Errorf("symlink not created: %v", err)
					return
				}
				if info.Mode()&os.ModeSymlink == 0 {
					t.Error("link is not a symlink")
				}
			},
		},
		{
			name:    "multiple links",
			content: validMultipleLinks,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual1"), []byte("c1"), 0644)
				os.WriteFile(filepath.Join(tmpDir, "actual2"), []byte("c2"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				for _, symlink := range []string{"symlink1", "symlink2"} {
					symlinkPath := filepath.Join(tmpDir, symlink)
					if _, err := os.Lstat(symlinkPath); err != nil {
						t.Errorf("symlink %s not created: %v", symlink, err)
					}
				}
			},
		},
		{
			name:    "platform-specific link matching current OS",
			content: currentPlatformOnly,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "symlink")
				if _, err := os.Lstat(symlinkPath); err != nil {
					t.Errorf("symlink not created: %v", err)
				}
			},
		},
		{
			name:    "platform-specific link NOT matching current OS",
			content: differentPlatformOnly,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "symlink")
				if _, err := os.Lstat(symlinkPath); err == nil {
					t.Errorf("symlink was created when it shouldn't have been: %v", err)
				}
			},
		},
		{
			name:    "platform override for current OS",
			content: manifestWithOverride,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "override_symlink")
				linkDest, err := os.Readlink(symlinkPath)
				if err != nil {
					t.Errorf("failed to read link: %v", err)
					return
				}
				if !strings.Contains(linkDest, "actual_file") {
					t.Errorf("expected link to actual_file, got %s", linkDest)
				}
			},
		},
		{
			name:    "platform override for different OS - uses default",
			content: overrideForDifferentOS,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "default_symlink")
				linkDest, err := os.Readlink(symlinkPath)
				if err != nil {
					t.Errorf("failed to read link: %v", err)
					return
				}
				if !strings.Contains(linkDest, "actual_file") {
					t.Errorf("expected link to actual_file, got %s", linkDest)
				}
			},
		},
		{
			name:    "relative link",
			content: relativeLink,
			wantErr: false,
			options: &state.TrovlOptions{
				UseRelative: true,
			},
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "symlink")
				linkDest, err := os.Readlink(symlinkPath)
				if err != nil {
					t.Errorf("failed to read link: %v", err)
					return
				}
				if filepath.IsAbs(linkDest) {
					t.Errorf("expected relative link, got absolute: %s", linkDest)
				}
			},
		},
		{
			name:    "absolute link (default)",
			content: absoluteLink,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "symlink")
				linkDest, err := os.Readlink(symlinkPath)
				if err != nil {
					t.Errorf("failed to read link: %v", err)
					return
				}
				if !filepath.IsAbs(linkDest) {
					t.Errorf("expected absolute link, got relative: %s", linkDest)
				}
			},
		},
		{
			name:    "empty manifest",
			content: emptyManifest,
			wantErr: false,
			setup:   func(tmpDir string) {},
		},
		{
			name:    "multiple platforms including current OS",
			content: specificPlatformsManifest,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "symlink")
				if _, err := os.Lstat(symlinkPath); err != nil {
					t.Errorf("symlink not created: %v", err)
				}
			},
		},
		{
			name:    "platform override removes platform from default set",
			content: overrideRemovesDefault,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "override_symlink")
				linkDest, err := os.Readlink(symlinkPath)
				if err != nil {
					t.Errorf("failed to read link: %v", err)
					return
				}
				if !strings.Contains(linkDest, "actual_file") {
					t.Errorf("expected override link, got %s", linkDest)
				}
			},
		},
		{
			name:    "source file doesn't exist",
			content: nonexistentSource,
			wantErr: true,
			setup:   func(tmpDir string) {},
		},
		{
			name:    "multiple overrides with current OS",
			content: multipleOverridesWithCurrent,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				expectedSymlink := runtime.GOOS + "_symlink"
				symlinkPath := filepath.Join(tmpDir, expectedSymlink)
				linkDest, err := os.Readlink(symlinkPath)
				if err != nil {
					t.Errorf("failed to read link: %v", err)
					return
				}
				if !strings.Contains(linkDest, "actual_file") {
					t.Errorf("expected link to actual_file, got %s", linkDest)
				}
			},
		},
		{
			name:    "complex manifest with mixed configurations",
			content: complexManifest,
			wantErr: false,
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual1"), []byte("c1"), 0644)
				os.WriteFile(filepath.Join(tmpDir, "actual2"), []byte("c2"), 0644)
				os.WriteFile(filepath.Join(tmpDir, "actual3"), []byte("c3"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				// symlink1 and symlink2 should be the default
				for _, symlink := range []string{"symlink1", "symlink2"} {
					symlinkPath := filepath.Join(tmpDir, symlink)
					if _, err := os.Lstat(symlinkPath); err != nil {
						t.Errorf("symlink %s not created: %v", symlink, err)
					}
				}
				// symlink3 should be an override
				symlinkPath := filepath.Join(tmpDir, "override_symlink3")
				linkDest, err := os.Readlink(symlinkPath)
				if err != nil {
					t.Errorf("failed to read override symlink: %v", err)
					return
				}
				if !strings.Contains(linkDest, "actual3") {
					t.Errorf("expected symlink3 to point to actual3, got %s", linkDest)
				}
				// 1 and 2 overrides and 3 default should not exist
				for _, name := range []string{
					"override_symlink1",
					"override_symlink2",
					"symlink3",
				} {
					path := filepath.Join(tmpDir, name)

					_, err := os.Lstat(path)
					if err == nil {
						t.Errorf("symlink %s should not exist", name)
					} else if !os.IsNotExist(err) {
						t.Errorf("unexpected error checking %s: %v", name, err)
					}
				}
			},
		},
		{
			name:    "BackupYes - backs up existing regular file",
			content: validSingleLink,
			wantErr: false,
			options: &state.TrovlOptions{
				BackupYes: true,
			},
			setup: func(tmpDir string) {
				var err error
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
				// ordinary file
				os.WriteFile(filepath.Join(tmpDir, "test_symlink"), []byte("existing"), 0644)

				backupDir, err = utils.GetCacheDir()
				if err != nil {
					t.Errorf("could not setup backup directory: %v", err)
					return
				}
				backupDir = filepath.Join(backupDir, "backups")
				if err := os.MkdirAll(backupDir, 0755); err != nil {
					t.Errorf("could not create backup parent directory: %v", err)
					return
				}
			},
			validate: func(t *testing.T, tmpDir string) {
				symlinkPath := filepath.Join(tmpDir, "test_symlink")
				info, err := os.Lstat(symlinkPath)
				if err != nil {
					t.Errorf("symlink not created: %v", err)
					return
				}
				if info.Mode()&os.ModeSymlink == 0 {
					t.Error("link is not a symlink")
				}

				entries, err := os.ReadDir(backupDir)
				if err != nil {
					t.Errorf("backup directory not accessible: %v", err)
					return
				}
				if len(entries) == 0 {
					t.Error("no backup files found")
				}
			},
		},
		{
			name:    "BackupYes - errors on existing directory",
			content: validSingleLink,
			wantErr: true,
			options: &state.TrovlOptions{
				BackupYes: true,
			},
			setup: func(tmpDir string) {
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
				// directory
				os.Mkdir(filepath.Join(tmpDir, "test_symlink"), 0755)
			},
		},
		{
			name:    "BackupNo - no error on existing regular file (continues on manifest)",
			content: validSingleLink,
			wantErr: false,
			options: &state.TrovlOptions{
				BackupNo: true,
			},
			setup: func(tmpDir string) {
				// Create source file
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
				// Create existing file at symlink location
				os.WriteFile(filepath.Join(tmpDir, "test_symlink"), []byte("existing"), 0644)
			},
			validate: func(t *testing.T, tmpDir string) {
				data, err := os.ReadFile(filepath.Join(tmpDir, "test_symlink"))
				if err != nil {
					t.Fatalf("expected file to exist: %v", err)
				}
				if string(data) != "existing" {
					t.Fatalf("expected file contents to remain unchanged, got: %q", string(data))
				}
			},
		},
		{
			name:    "BackupNo - errors on existing directory",
			content: validSingleLink,
			wantErr: true,
			options: &state.TrovlOptions{
				BackupNo: true,
			},
			setup: func(tmpDir string) {
				// Create source file
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
				// Create existing directory at symlink location
				os.Mkdir(filepath.Join(tmpDir, "test_symlink"), 0755)
			},
		},
		{
			name:    "no backup/overwrite option set (eof during tests) - errors on existing file",
			content: validSingleLink,
			wantErr: true,
			options: &state.TrovlOptions{},
			setup: func(tmpDir string) {
				// Create source file
				os.WriteFile(filepath.Join(tmpDir, "actual_file"), []byte("content"), 0644)
				// Create existing file at symlink location
				os.WriteFile(filepath.Join(tmpDir, "test_symlink"), []byte("existing"), 0644)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Setenv("XDG_CONFIG_HOME", tmpDir)
			t.Setenv("XDG_CACHE_HOME", tmpDir)
			manifestPath := filepath.Join(tmpDir, "manifest.json")

			if err := os.WriteFile(manifestPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create manifest file: %v", err)
			}

			m, err := New(manifestPath)
			if err != nil {
				t.Fatalf("unexpected error from New(): %v", err)
			}

			if tt.setup != nil {
				tt.setup(tmpDir)
			}

			oldWd, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(oldWd)

			if tt.options == nil {
				err = m.Apply(teststate)
			} else {
				newstate := state.New(tt.options)
				err = m.Apply(newstate)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error from Apply(), got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error from Apply(): %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, tmpDir)
			}
		})
	}
}

func TestIsSupportedPlatform(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		want     bool
	}{
		{
			name:     "valid platform",
			platform: "linux",
			want:     true,
		},
		{
			name:     "invalid platform",
			platform: "lolos",
			want:     false,
		},
		{
			name:     "empty string",
			platform: "",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSupportedPlatform(tt.platform)
			if got != tt.want {
				t.Errorf("platform '%v': got %v, want %v", tt.platform, got, tt.want)
			}
		})
	}
}
