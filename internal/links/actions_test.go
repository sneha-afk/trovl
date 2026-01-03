package links_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/sneha-afk/trovl/internal/state"
	"github.com/sneha-afk/trovl/internal/utils"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// Calls Construct and Add to simulate full work through of the add subcommand.
func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		options  *state.TrovlOptions
		setup    func(tmp, targetPath, linkPath string)
		validate func(t *testing.T, tmp, targetPath, linkPath string)
	}{
		{
			name:    "error: target does not exist",
			wantErr: true,
		},
		{
			name: "success: brand new symlink",
			setup: func(tmp, targetPath, linkPath string) {
				_ = os.WriteFile(targetPath, []byte("target"), 0644)
			},
			validate: func(t *testing.T, tmp, targetPath, linkPath string) {
				info, err := os.Lstat(linkPath)
				if err != nil {
					t.Fatalf("expected symlink to exist: %v", err)
				}
				if info.Mode()&os.ModeSymlink == 0 {
					t.Fatalf("expected symlink, got %v", info.Mode())
				}
			},
		},
		{
			name: "success: existing symlink, overwrite yes",
			options: &state.TrovlOptions{
				OverwriteYes: true,
			},
			setup: func(tmp, targetPath, linkPath string) {
				_ = os.WriteFile(targetPath, []byte("target"), 0644)
				_ = os.Symlink(targetPath, linkPath)
			},
			validate: func(t *testing.T, tmp, targetPath, linkPath string) {
				if _, err := os.Lstat(linkPath); err != nil {
					t.Fatalf("expected symlink to exist after overwrite: %v", err)
				}
			},
		},
		{
			name:    "error: existing symlink, overwrite no",
			wantErr: true,
			options: &state.TrovlOptions{
				OverwriteNo: true,
			},
			setup: func(tmp, targetPath, linkPath string) {
				_ = os.WriteFile(targetPath, []byte("target"), 0644)
				_ = os.Symlink(targetPath, linkPath)
			},
		},
		{
			name: "success: ordinary file exists, backup yes",
			options: &state.TrovlOptions{
				BackupYes: true,
			},
			setup: func(tmp, targetPath, linkPath string) {
				_ = os.WriteFile(targetPath, []byte("target"), 0644)
				_ = os.WriteFile(linkPath, []byte("ordinary"), 0644)
			},
			validate: func(t *testing.T, tmp, targetPath, linkPath string) {
				if _, err := os.Lstat(linkPath); err != nil {
					t.Fatalf("expected symlink to exist after backup: %v", err)
				}
			},
		},
		{
			name:    "error: ordinary file exists, backup no",
			wantErr: true,
			options: &state.TrovlOptions{
				BackupNo: true,
			},
			setup: func(tmp, targetPath, linkPath string) {
				_ = os.WriteFile(targetPath, []byte("target"), 0644)
				_ = os.WriteFile(linkPath, []byte("ordinary"), 0644)
			},
		},
		{
			name:    "error: directory exists at symlink path",
			wantErr: true,
			setup: func(tmp, targetPath, linkPath string) {
				_ = os.WriteFile(targetPath, []byte("target"), 0644)
				_ = os.Mkdir(linkPath, 0755)
			},
		},
		{
			name: "success: dry-run conflict does nothing",
			options: &state.TrovlOptions{
				DryRun: true,
			},
			setup: func(tmp, targetPath, linkPath string) {
				_ = os.WriteFile(targetPath, []byte("target"), 0644)
				_ = os.WriteFile(linkPath, []byte("ordinary"), 0644)
			},
			validate: func(t *testing.T, tmp, targetPath, linkPath string) {
				// original file must still exist
				if _, err := os.Stat(linkPath); err != nil {
					t.Fatalf("expected file to remain during dry-run")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			t.Setenv("XDG_CACHE_HOME", tmp)

			targetPath := filepath.Join(tmp, "target.txt")
			linkPath := filepath.Join(tmp, "link.txt")

			if tt.setup != nil {
				tt.setup(tmp, targetPath, linkPath)
			}

			st := state.DefaultState()
			if tt.options != nil {
				st.Options = tt.options
			}

			_, err := links.Construct(st, targetPath, linkPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Construct: wantErr=%v, got %v", tt.wantErr, err)
			}

			err = links.Add(st, targetPath, linkPath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Add: wantErr=%v, got %v", tt.wantErr, err)
			}

			if tt.validate != nil {
				tt.validate(t, tmp, targetPath, linkPath)
			}
		})
	}
}

func TestValidateSymlink(t *testing.T) {
	type result struct {
		valid bool
		err   bool
	}

	tests := []struct {
		name          string
		createSymlink bool
		target        string
		expected      result
	}{
		{
			name:          "success: symlink to a valid file",
			createSymlink: true,
			target:        "valid_file.txt",
			expected:      result{valid: true, err: false},
		},
		{
			name:          "success: symlink points to a directory",
			createSymlink: true,
			target:        "valid_directory",
			expected:      result{valid: true, err: false},
		},
		{
			name:          "error: valid symlink to non-existing target",
			createSymlink: true,
			target:        "non_existing_target.txt",
			expected:      result{valid: false, err: true},
		},
		{
			name:          "error: invalid symlink path",
			createSymlink: false,
			target:        "invalid_symlink.txt",
			expected:      result{valid: false, err: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			targetPath := filepath.Join(tmp, tc.target)
			symlinkPath := filepath.Join(tmp, "symlink.txt")

			switch tc.target {
			case "valid_file.txt", "invalid_symlink.txt":
				os.WriteFile(targetPath, []byte("test content"), 0644)
			case "valid_directory":
				os.Mkdir(targetPath, 0755)
			}

			if tc.createSymlink {
				os.Symlink(targetPath, symlinkPath)
			}

			valid, err := utils.ValidateSymlink(symlinkPath)
			if valid != tc.expected.valid {
				t.Errorf("expected valid: %v, got: %v", tc.expected.valid, valid)
			}
			if (err != nil) != tc.expected.err {
				t.Errorf("expected error: %v, got: %v", tc.expected.err, err)
			}
		})
	}
}

func TestRemoveByPath(t *testing.T) {
	teststate := state.DefaultState()

	tests := []struct {
		name         string
		targetExists bool
		createLink   bool
		expectErr    bool
	}{
		{
			name:         "success: remove existing symlink",
			targetExists: true,
			createLink:   true,
			expectErr:    false,
		},
		{
			name:         "error: symlink does not exist",
			targetExists: false,
			createLink:   false,
			expectErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			targetPath := filepath.Join(tmp, "target.txt")
			linkPath := filepath.Join(tmp, "link.txt")

			if tc.targetExists {
				os.WriteFile(targetPath, []byte("target"), 0644)
			}

			if tc.createLink {
				if err := os.Symlink(targetPath, linkPath); err != nil {
					t.Errorf("error during link setup: %v", err)
				}
			}

			err := links.RemoveByPath(teststate, linkPath)

			if (err != nil) != tc.expectErr {
				t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
			}

			if !tc.expectErr {
				if _, err := os.Lstat(linkPath); err == nil {
					t.Error("expected symlink to be removed, but it still exists")
				}
			}
		})
	}
}
