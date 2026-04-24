package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lotus-creations/tmux-helper/internal/tmux"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tmux-helper <command>")
		fmt.Println("Commands: picker, layout-next, layout-prev")
		os.Exit(1)
	}

	client := tmux.NewClient()
	command := os.Args[1]

	switch command {
	case "picker":
		sessions, err := client.ListSessions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sessions found:", len(sessions))
		for _, s := range sessions {
			attached := ""
			if s.Attached {
				attached = " (attached)"
			}
			fmt.Printf("  %s%s - %d windows\n", s.Name, attached, s.Windows)
		}

	case "layout-next":
		// tmux next-layout is handled via keybinding
		// This could be expanded for integration
		layout := client.GetCurrentLayout()
		fmt.Println("Current layout:", layout)

	case "layout-prev":
		layout := client.GetCurrentLayout()
		fmt.Println("Current layout:", layout)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// argsToString converts os.Args starting from index to a string
func argsToString() string {
	if len(os.Args) < 3 {
		return ""
	}
	return strings.Join(os.Args[2:], " ")
}