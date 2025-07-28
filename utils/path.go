package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands the ~ character in file paths to the user's home directory.
// It returns the expanded absolute path.
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return filepath.Abs(path)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:]), nil
	}

	// Path like ~user/... - not supported, return as-is
	return filepath.Abs(path)
}