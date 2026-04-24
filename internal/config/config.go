package config

import (
	"os"
	"os/user"
	"path/filepath"
)

// Config holds user preferences
type Config struct {
	// Prefix key (default: C-a)
	Prefix string

	// Split sizes (0.0 - 1.0)
	SplitVerticalSize   float64
	SplitHorizontalSize float64

	// Theme
	Theme string

	// Mouse enabled
	Mouse bool
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Prefix:              "C-a",
		SplitVerticalSize:   0.5,
		SplitHorizontalSize: 0.5,
		Theme:               "default",
		Mouse:               true,
	}
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	// Look for config in standard locations
	paths := []string{
		".tmux-helper.conf",
		".config/tmux-helper.conf",
		filepath.Join(os.Getenv("HOME"), ".tmux-helper.conf"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			// File exists, could parse it here
			// For now, just return defaults
			break
		}
	}

	return cfg, nil
}

// ConfigPath returns the default config file path
func ConfigPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".tmux-helper.conf")
}