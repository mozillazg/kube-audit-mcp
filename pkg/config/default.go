package config

import (
	"os"
	"path/filepath"
	"strings"
)

var defaultConfigFile = "~/.config/kube-audit-mcp/config.yaml"

func DefaultConfigFile() string {
	p, err := ExpandPath(defaultConfigFile)
	if err != nil {
		return defaultConfigFile
	}
	return p
}

func ExpandPath(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path, err
		}
		path = filepath.Join(home, path[1:])
	}
	return path, nil
}

func ShortHomePath(path string) string {
	home, err := osUserHomeDir()
	if err != nil {
		return path
	}

	// Only replace the home directory if the path starts with it
	if strings.HasPrefix(path, home) {
		return strings.Replace(path, home, "~", 1)
	}
	return path
}

// Variable for mocking in tests
var osUserHomeDir = os.UserHomeDir
