package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Config holds user preferences
type Config struct {
	// Prefix key (default: C-a)
	Prefix string

	// Split sizes (0.0 - 1.0)
	SplitVerticalSize   float64
	SplitHorizontalSize float64

	// Mouse enabled (default: true)
	Mouse bool

	// Theme (default: purple)
	Theme string

	// Terminal to use
	Terminal string
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Prefix:              "C-a",
		SplitVerticalSize:   0.5,
		SplitHorizontalSize: 0.5,
		Mouse:               true,
		Theme:               "purple",
		Terminal:            "screen-256color",
	}
}

// ConfigPath returns the default config file path
func ConfigPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".tmux-helper.conf")
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No config file, use defaults
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	// Parse key=value lines
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "prefix":
			cfg.Prefix = value
		case "split-vertical-size":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				if v >= 0.1 && v <= 0.9 {
					cfg.SplitVerticalSize = v
				}
			}
		case "split-horizontal-size":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				if v >= 0.1 && v <= 0.9 {
					cfg.SplitHorizontalSize = v
				}
			}
		case "mouse":
			cfg.Mouse = value == "true" || value == "1" || value == "yes"
		case "theme":
			cfg.Theme = value
		case "terminal":
			cfg.Terminal = value
		}
	}

	return cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
	path := ConfigPath()
	
	content := fmt.Sprintf(`# tmux-helper configuration
# https://github.com/lotus-creations/tmux-helper

# Prefix key (default: C-a)
prefix=%s

# Split sizes (0.1 - 0.9)
split-vertical-size=%.2f
split-horizontal-size=%.2f

# Mouse support (true/false)
mouse=%t

# Theme (purple/green)
theme=%s

# Terminal (default: screen-256color)
terminal=%s
`, cfg.Prefix, cfg.SplitVerticalSize, cfg.SplitHorizontalSize, cfg.Mouse, cfg.Theme, cfg.Terminal)

	return os.WriteFile(path, []byte(content), 0644)
}

// PrintConfig prints the current configuration
func PrintConfig(cfg *Config) {
	fmt.Println("tmux-helper Configuration")
	fmt.Println("=========================")
	fmt.Printf("Config file: %s\n\n", ConfigPath())
	
	fmt.Printf("prefix              = %s\n", cfg.Prefix)
	fmt.Printf("split-vertical-size = %.2f\n", cfg.SplitVerticalSize)
	fmt.Printf("split-horizontal-size = %.2f\n", cfg.SplitHorizontalSize)
	fmt.Printf("mouse               = %t\n", cfg.Mouse)
	fmt.Printf("theme               = %s\n", cfg.Theme)
	fmt.Printf("terminal            = %s\n", cfg.Terminal)
	fmt.Println()
	fmt.Println("Run 'tmux-helper apply' after changing config to regenerate ~/.tmux.conf")
}

// Validate checks if the config is valid
func Validate(cfg *Config) error {
	if cfg.Prefix == "" {
		return fmt.Errorf("prefix cannot be empty")
	}
	if cfg.SplitVerticalSize < 0.1 || cfg.SplitVerticalSize > 0.9 {
		return fmt.Errorf("split-vertical-size must be between 0.1 and 0.9")
	}
	if cfg.SplitHorizontalSize < 0.1 || cfg.SplitHorizontalSize > 0.9 {
		return fmt.Errorf("split-horizontal-size must be between 0.1 and 0.9")
	}
	return nil
}