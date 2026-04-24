package tmux

import (
	"strconv"
	"strings"
)

// Session represents a tmux session
type Session struct {
	ID        string
	Name      string
	Windows   int
	Created   int64
	Attached  bool
}

// Window represents a tmux window
type Window struct {
	ID      string
	Index   int
	Name    string
	Layout  string
	Active  bool
	Panes   int
}

// Pane represents a tmux pane
type Pane struct {
	ID       string
	Index    int
	Title    string
	Command  string
	Active   bool
}

// Layout represents layout states for cycling
var Layouts = []string{
	"even-horizontal",
	"even-vertical",
	"main-horizontal",
	"main-vertical",
	"tiled",
}

// ParseSession parses a tmux session line
func ParseSession(line string) Session {
	parts := strings.Split(line, "|")
	if len(parts) < 5 {
		return Session{}
	}

	attached := strings.TrimSpace(parts[4]) == "1"
	created, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
	windows, _ := strconv.Atoi(strings.TrimSpace(parts[2]))

	return Session{
		ID:       strings.TrimSpace(parts[0]),
		Name:     strings.TrimSpace(parts[1]),
		Windows:  windows,
		Created:  created,
		Attached: attached,
	}
}

// ParseWindow parses a tmux window line
func ParseWindow(line string) Window {
	parts := strings.Split(line, "|")
	if len(parts) < 6 {
		return Window{}
	}

	index, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	active := strings.TrimSpace(parts[4]) == "1"
	panes, _ := strconv.Atoi(strings.TrimSpace(parts[5]))

	return Window{
		ID:     strings.TrimSpace(parts[0]),
		Index:  index,
		Name:   strings.TrimSpace(parts[2]),
		Layout: strings.TrimSpace(parts[3]),
		Active: active,
		Panes:  panes,
	}
}

// ParsePane parses a tmux pane line
func ParsePane(line string) Pane {
	parts := strings.Split(line, "|")
	if len(parts) < 5 {
		return Pane{}
	}

	index, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	active := strings.TrimSpace(parts[4]) == "1"

	return Pane{
		ID:      strings.TrimSpace(parts[0]),
		Index:   index,
		Title:   strings.TrimSpace(parts[2]),
		Command: strings.TrimSpace(parts[3]),
		Active:  active,
	}
}

// LayoutIndex finds the index of a layout in the cycle
func LayoutIndex(layout string) int {
	for i, l := range Layouts {
		if l == layout {
			return i
		}
	}
	return 0
}

// NextLayoutIndex cycles to the next layout
func NextLayoutIndex(current int) int {
	return (current + 1) % len(Layouts)
}