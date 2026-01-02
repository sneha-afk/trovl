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
	type result struct {
		target string
		mount  string
		err    bool
	}

	tests := []struct {
		name         string
		targetExists bool
		linkExists   bool
		userInput    string
		expected     result
	}{
		{
			name:         "success: brand new symlink",
			targetExists: true,
			linkExists:   false,
			expected:     result{"target.txt", "link.txt", false},
		},
		{
			name:         "error: target file does not exist",
			targetExists: false,
			linkExists:   false,
			expected:     result{err: true},
		},
		{
			name:         "success: existing file, accepted overwrite",
			targetExists: true,
			linkExists:   true,
			userInput:    "y\n",
			expected:     result{"target.txt", "link.txt", false},
		},
		{
			name:         "error: existing file, declined overwrite",
			targetExists: true,
			linkExists:   true,
			userInput:    "n\n",
			expected:     result{err: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			targetPath := filepath.Join(tmp, "target.txt")
			linkPath := filepath.Join(tmp, "link.txt")

			// writefile creates file if not exist
			if tc.targetExists {
				os.WriteFile(targetPath, []byte("target"), 0644)
			}
			if tc.linkExists {
				os.Symlink(targetPath, linkPath)
			}

			state := state.DefaultState()
			switch tc.userInput {
			case "y\n":
				state.Options.OverwriteYes = true
			case "n\n":
				state.Options.OverwriteNo = true
			}

			res, err := links.Construct(state, targetPath, linkPath)

			if (err != nil) != tc.expected.err {
				t.Errorf("(Construct) expected error: %v, got: %v", tc.expected.err, err)
			}

			if err == nil {
				if res.Target != targetPath || res.LinkMount != linkPath {
					t.Errorf("(Construct) returned Link struct mismatch: expected target = %v, got %v; expected symlink = %v, got %v", targetPath, res.Target, linkPath, res.LinkMount)
				}

				// if decided to overwrite, Construct would remove the old file so Add is guaranteed to add something
				if tc.linkExists && tc.userInput == "y\n" {
					if _, err := os.Stat(linkPath); err == nil {
						t.Error("(Construct) expected existing link file to be deleted (yes to overwrit), but it still exists")
					}
				}
			}

			err = links.Add(state, targetPath, linkPath)
			if (err != nil) != tc.expected.err {
				t.Errorf("(Add) expected error: %v, got: %v", tc.expected.err, err)
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
