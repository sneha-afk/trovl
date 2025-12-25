package links_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/sneha-afk/trovl/internal/links"
	"github.com/sneha-afk/trovl/internal/models"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name      string
		param     string
		expected  bool
		expectErr bool
	}{
		{"existing file", "test_exists.go", true, false},
		{"non-existing file", "test_exists_no.go", false, true},
	}

	file, err := os.Create("test_exists.go")
	if err != nil {
		log.Fatalf("TestValidatePath: could not create test file: %v", err)
	}
	file.Close()

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			res, err := links.ValidatePath(testcase.param)

			if res != testcase.expected {
				t.Errorf("expected %v, got %v (error: %v)", testcase.expected, res, err)
			}
			if err != nil && testcase.expectErr == false {
				t.Errorf("%v", err)
			}
		})
	}

	if err := os.Remove("test_exists.go"); err != nil {
		log.Fatalf("TestValidatePath: could not delete test file: %v", err)
	}
}

func TestConstruct(t *testing.T) {
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
				os.WriteFile(linkPath, []byte("existing"), 0644)
			}

			// create a pipe to simulate stdin/out
			// newlines in this input to make sure it goes through
			if tc.userInput != "" {
				oldStdin := os.Stdin
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatal(err)
				}

				os.Stdin = r

				w.WriteString(tc.userInput)
				w.Close()

				t.Cleanup(func() { os.Stdin = oldStdin })
			}

			res, err := links.Construct(targetPath, linkPath, models.LinkFile)

			if (err != nil) != tc.expected.err {
				t.Errorf("expected error: %v, got: %v", tc.expected.err, err)
			}

			if err == nil {
				if res.Target != targetPath || res.LinkMount != linkPath {
					t.Errorf("returned Link struct mismatch: expected target = %v, got %v; expected symlink = %v, got %v", targetPath, res.Target, linkPath, res.LinkMount)
				}

				// if decided to overwrite, Construct would remove the old file so Add is guaranteed to add something
				if tc.linkExists && tc.userInput == "y\n" {
					if _, err := os.Stat(linkPath); err == nil {
						t.Error("expected existing link file to be deleted (yes to overwrit), but it still exists")
					}
				}
			}
		})
	}
}
