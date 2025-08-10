package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	// Get the actual home directory for comparison
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "tilde only",
			input:    "~",
			expected: homeDir,
			wantErr:  false,
		},
		{
			name:     "tilde with slash",
			input:    "~/",
			expected: homeDir,
			wantErr:  false,
		},
		{
			name:     "tilde with subdirectory",
			input:    "~/Documents",
			expected: filepath.Join(homeDir, "Documents"),
			wantErr:  false,
		},
		{
			name:     "tilde with nested path",
			input:    "~/Documents/Projects/myapp",
			expected: filepath.Join(homeDir, "Documents/Projects/myapp"),
			wantErr:  false,
		},
		{
			name:     "tilde with file",
			input:    "~/.bashrc",
			expected: filepath.Join(homeDir, ".bashrc"),
			wantErr:  false,
		},
		{
			name:     "absolute path",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
			wantErr:  false,
		},
		{
			name:     "relative path",
			input:    "Documents/file.txt",
			expected: "Documents/file.txt",
			wantErr:  false,
		},
		{
			name:     "current directory",
			input:    "./config.yaml",
			expected: "./config.yaml",
			wantErr:  false,
		},
		{
			name:     "parent directory",
			input:    "../config.yaml",
			expected: "../config.yaml",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "tilde not at beginning",
			input:    "/home/user~backup",
			expected: "/home/user~backup",
			wantErr:  false,
		},
		{
			name:     "tilde in middle",
			input:    "path/~/file",
			expected: "path/~/file",
			wantErr:  false,
		},
		{
			name:     "windows-style path with tilde",
			input:    "~\\Documents\\file.txt",
			expected: filepath.Join(homeDir, "\\Documents\\file.txt"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandPath(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ExpandPath() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ExpandPath() unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("ExpandPath() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestExpandPath_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "single character tilde",
			input:   "~",
			wantErr: false,
		},
		{
			name:    "multiple tildes at start",
			input:   "~~",
			wantErr: false,
		},
		{
			name:    "tilde with special characters",
			input:   "~/@#$%^&*()",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandPath(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ExpandPath() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ExpandPath() unexpected error: %v", err)
				return
			}

			// For paths starting with ~, ensure they start with home directory
			if len(tt.input) > 0 && tt.input[0] == '~' {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					t.Fatalf("Failed to get home directory: %v", err)
				}

				if !strings.HasPrefix(result, homeDir) {
					t.Errorf("ExpandPath() result %q should start with home directory %q", result, homeDir)
				}
			}
		})
	}
}

func TestExpandPath_HomeDirectoryError(t *testing.T) {
	// This test is challenging to implement reliably across different systems
	// because we can't easily mock os.UserHomeDir() without changing the function signature
	// or using dependency injection.

	// We'll test the happy path and document that error cases would need
	// integration testing or dependency injection to test properly.

	t.Run("verify function handles home directory retrieval", func(t *testing.T) {
		// Test that the function can successfully get home directory
		_, err := os.UserHomeDir()
		if err != nil {
			t.Skip("Skipping test: cannot get home directory on this system")
		}

		// Test normal expansion works
		result, err := ExpandPath("~/test")
		if err != nil {
			t.Errorf("ExpandPath() unexpected error when home directory is available: %v", err)
		}

		homeDir, _ := os.UserHomeDir()
		expected := filepath.Join(homeDir, "test")
		if result != expected {
			t.Errorf("ExpandPath() = %q, expected %q", result, expected)
		}
	})
}

func TestShortHomePath(t *testing.T) {
	// Get the actual home directory for comparison
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "home directory",
			input:    homeDir,
			expected: "~",
		},
		{
			name:     "home directory with trailing slash",
			input:    homeDir + "/",
			expected: "~/",
		},
		{
			name:     "subdirectory in home",
			input:    filepath.Join(homeDir, "Documents"),
			expected: "~/Documents",
		},
		{
			name:     "nested path in home",
			input:    filepath.Join(homeDir, "Documents/Projects/myapp"),
			expected: "~/Documents/Projects/myapp",
		},
		{
			name:     "hidden file in home",
			input:    filepath.Join(homeDir, ".bashrc"),
			expected: "~/.bashrc",
		},
		{
			name:     "absolute path outside home",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "relative path",
			input:    "Documents/file.txt",
			expected: "Documents/file.txt",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "path containing home as substring",
			input:    "/tmp" + homeDir,
			expected: "/tmp" + homeDir,
		},
		{
			name:     "windows-style path in home",
			input:    filepath.Join(homeDir, "Documents\\file.txt"),
			expected: "~/Documents\\file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShortHomePath(tt.input)
			if result != tt.expected {
				t.Errorf("ShortHomePath() = %q, want %q", result, tt.expected)
			}
		})
	}

	// Test with mocked error from os.UserHomeDir()
	t.Run("error getting home dir", func(t *testing.T) {
		// Save the original function and restore it after the test
		original := osUserHomeDir
		defer func() { osUserHomeDir = original }()

		// Mock os.UserHomeDir to return an error
		osUserHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		// When UserHomeDir fails, ShortHomePath should return the original path
		testPath := "/some/random/path"
		result := ShortHomePath(testPath)
		if result != testPath {
			t.Errorf("ShortHomePath() = %q, want %q", result, testPath)
		}
	})
}
