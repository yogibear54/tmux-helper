package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lotus-creations/tmux-helper/internal/tmux"
	"github.com/lotus-creations/tmux-helper/internal/ui"
)

func main() {
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

	case "layout-next":
		if err := client.NextLayout(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "layout-prev":
		layout := client.GetCurrentLayout()
		fmt.Println("Current layout:", layout)
		fmt.Println("(Use layout-next to cycle)")

	case "layout":
		layout := client.GetCurrentLayout()
		fmt.Println(layout)

	case "sessions":
		listSessions()

	case "help", "?":
		printHelp()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: tmux-helper <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  picker       Open interactive session picker")
	fmt.Println("  help-overlay Show help in popup")
	fmt.Println("  layout       Show current layout")
	fmt.Println("  layout-next  Cycle to next layout")
	fmt.Println("  sessions     List all sessions")
	fmt.Println("  help         Show keybindings (text)")
	fmt.Println("")
	fmt.Println("Keybindings (Prefix: Ctrl-a):")
	fmt.Println("  ?             Help overlay")
	fmt.Println("  F             Session picker")
	fmt.Println("  h/j/k/l       Navigate panes (vim-style)")
	fmt.Println("  |             Split left/right")
	fmt.Println("  -             Split top/bottom")
	fmt.Println("  Space         Cycle layout")
	fmt.Println("  c             New window")
	fmt.Println("  d             Detach")
	fmt.Println("  x             Kill pane")
	fmt.Println("  X             Kill window")
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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