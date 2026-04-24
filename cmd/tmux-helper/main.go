package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/lotus-creations/tmux-helper/internal/config"
	"github.com/lotus-creations/tmux-helper/internal/tmux"
	"github.com/lotus-creations/tmux-helper/internal/ui"
)

const version = "0.1.0"

func main() {
	// Handle --help and --version flags
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		case "--version", "-v":
			fmt.Printf("tmux-helper version %s\n", version)
			os.Exit(0)
		}
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	client := tmux.NewClient()

	switch command {
	case "picker":
		ui.RunPicker()

	case "help-overlay", "overlay":
		ui.RunHelp()

	case "config":
		handleConfig(os.Args[2:])

	case "apply":
		handleApply()

	case "layout-next":
		if err := client.NextLayout(); err != nil {
			fmt.Fprintf(os.Stderr, "Error cycling layout: %v\n", err)
			notifyError("Failed to cycle layout")
			os.Exit(1)
		}

	case "layout-prev":
		layout := client.GetCurrentLayout()
		if layout == "" {
			fmt.Fprintln(os.Stderr, "Error: No tmux session active")
			notifyError("No tmux session active")
			os.Exit(1)
		}
		fmt.Println("Current layout:", layout)
		fmt.Println("(Use layout-next to cycle)")

	case "layout":
		layout := client.GetCurrentLayout()
		if layout == "" {
			fmt.Fprintln(os.Stderr, "Error: No tmux session active")
			notifyError("No tmux session active")
			os.Exit(1)
		}
		fmt.Println(layout)

	case "sessions":
		listSessions()

	case "version", "ver":
		fmt.Printf("tmux-helper version %s\n", version)

	case "help", "?":
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleConfig(args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		notifyError("Failed to load config")
		os.Exit(1)
	}

	if len(args) == 0 {
		// Show current config
		config.PrintConfig(cfg)
		return
	}

	subcommand := args[0]

	switch subcommand {
	case "show":
		config.PrintConfig(cfg)

	case "path":
		fmt.Println(config.ConfigPath())

	case "get":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: tmux-helper config get <key>")
			os.Exit(1)
		}
		key := args[1]
		switch key {
		case "prefix":
			fmt.Println(cfg.Prefix)
		case "split-vertical-size":
			fmt.Printf("%.2f\n", cfg.SplitVerticalSize)
		case "split-horizontal-size":
			fmt.Printf("%.2f\n", cfg.SplitHorizontalSize)
		case "mouse":
			fmt.Printf("%t\n", cfg.Mouse)
		case "theme":
			fmt.Println(cfg.Theme)
		case "terminal":
			fmt.Println(cfg.Terminal)
		default:
			fmt.Fprintf(os.Stderr, "Unknown key: %s\n", key)
			os.Exit(1)
		}

	case "set":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: tmux-helper config set <key> <value>")
			fmt.Fprintln(os.Stderr, "\nAvailable keys:")
			fmt.Fprintln(os.Stderr, "  prefix               (e.g., C-a, C-t)")
			fmt.Fprintln(os.Stderr, "  split-vertical-size   (0.1 - 0.9)")
			fmt.Fprintln(os.Stderr, "  split-horizontal-size (0.1 - 0.9)")
			fmt.Fprintln(os.Stderr, "  mouse                (true/false)")
			fmt.Fprintln(os.Stderr, "  theme                (purple/green)")
			fmt.Fprintln(os.Stderr, "  terminal             (e.g., screen-256color)")
			os.Exit(1)
		}
		key := args[1]
		value := args[2]

		switch key {
		case "prefix":
			cfg.Prefix = value
		case "split-vertical-size":
			var v float64
			fmt.Sscanf(value, "%f", &v)
			if v < 0.1 || v > 0.9 {
				fmt.Fprintln(os.Stderr, "Error: split-vertical-size must be between 0.1 and 0.9")
				os.Exit(1)
			}
			cfg.SplitVerticalSize = v
		case "split-horizontal-size":
			var v float64
			fmt.Sscanf(value, "%f", &v)
			if v < 0.1 || v > 0.9 {
				fmt.Fprintln(os.Stderr, "Error: split-horizontal-size must be between 0.1 and 0.9")
				os.Exit(1)
			}
			cfg.SplitHorizontalSize = v
		case "mouse":
			cfg.Mouse = value == "true" || value == "1" || value == "yes"
		case "theme":
			cfg.Theme = value
		case "terminal":
			cfg.Terminal = value
		default:
			fmt.Fprintf(os.Stderr, "Unknown key: %s\n", key)
			os.Exit(1)
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			notifyError("Failed to save config")
			os.Exit(1)
		}
		fmt.Printf("Updated %s = %s\n", key, value)
		fmt.Println("Run 'tmux-helper apply' to regenerate ~/.tmux.conf")

	case "edit":
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		path := config.ConfigPath()
		// Create default config if it doesn't exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			config.SaveConfig(cfg)
		}
		execCmd(editor, path)
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "Unknown config command: %s\n", subcommand)
		fmt.Fprintln(os.Stderr, "\nUsage: tmux-helper config [show|set|get|edit]")
		os.Exit(1)
	}
}

func handleApply() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		notifyError("Failed to load config")
		os.Exit(1)
	}

	if err := config.Validate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
		notifyError("Invalid config")
		os.Exit(1)
	}

	// Generate tmux.conf
	content := cfg.TmuxConfig()
	if content == "" {
		fmt.Fprintln(os.Stderr, "Error: Failed to generate tmux.conf")
		notifyError("Failed to generate tmux.conf")
		os.Exit(1)
	}
	
	// Save to ~/.tmux.conf
	usr, _ := user.Current()
	path := filepath.Join(usr.HomeDir, ".tmux.conf")
	
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing ~/.tmux.conf: %v\n", err)
		notifyError("Failed to write tmux.conf")
		os.Exit(1)
	}

	fmt.Printf("Generated ~/.tmux.conf with prefix=%s\n", cfg.Prefix)
	fmt.Println("Run 'tmux source-file ~/.tmux.conf' to apply changes")
}

func execCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s: %v\n", name, err)
		os.Exit(1)
	}
}

// notifyError shows error in tmux status bar if available
func notifyError(msg string) {
	if os.Getenv("TMUX") != "" {
		tmux.NotifyError(msg)
	}
}

func printUsage() {
	fmt.Printf(`tmux-helper %s - i3-inspired tmux keybindings and TUI tools

Usage: tmux-helper <command>

Commands:
  picker         Open interactive session picker (TUI popup)
  help-overlay   Show help overlay (TUI popup)
  config         Show/set configuration
  apply          Regenerate ~/.tmux.conf from config
  layout         Show current layout
  layout-next    Cycle to next layout
  sessions       List all sessions
  version        Show version
  help           Show this help

Config commands:
  config show         Show current config
  config set <k> <v>  Set a config value
  config get <k>      Get a config value
  config edit         Edit config in $EDITOR

Keybindings (Prefix: Ctrl-a):
  ?             Help overlay
  F             Session picker
  h/j/k/l       Navigate panes (vim-style)
  |             Split left/right
  -             Split top/bottom
  Space         Cycle layout
  c             New window
  d             Detach
  x             Kill pane
  X             Kill window

Run 'tmux-helper help' for keybindings reference.
Run 'tmux-helper --version' for version info.
`, version)
}

func printHelp() {
	help := `
tmux-helper Keybindings
=======================

Prefix: Ctrl-a

PANE NAVIGATION
  h           Move left
  j           Move down
  k           Move up
  l           Move right

SPLITS
  |            Split left/right (vertical)
  -            Split top/bottom (horizontal)

LAYOUTS
  Space        Cycle to next layout

SESSIONS
  F            Open session picker (TUI)
  c            New window
  d            Detach

PANE MANAGEMENT
  x            Kill current pane
  X            Kill current window
  !            Break pane into new window
  Shift+H/J/K/L  Swap with adjacent pane

MOUSE
  Click        Select pane

COPY MODE (vim-style)
  [            Enter copy mode
  v            Begin selection
  y            Copy selection
  Enter        Copy selection
`
	fmt.Println(help)
}

func listSessions() {
	client := tmux.NewClient()
	sessions, err := client.ListSessions()
	if err != nil {
		if strings.Contains(err.Error(), "no server running") {
			fmt.Println("No tmux sessions running")
			return
		}
		fmt.Fprintf(os.Stderr, "Error listing sessions: %v\n", err)
		notifyError("Failed to list sessions")
		os.Exit(1)
	}

	if len(sessions) == 0 {
		fmt.Println("No tmux sessions running")
		return
	}

	fmt.Printf("Sessions (%d):\n", len(sessions))
	for _, s := range sessions {
		attached := ""
		if s.Attached {
			attached = " (attached)"
		}
		fmt.Printf("  %s%s - %d windows\n", s.Name, attached, s.Windows)
	}
}