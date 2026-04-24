package tmux

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps tmux CLI interactions
type Client struct{}

// NewClient creates a new tmux client
func NewClient() *Client {
	return &Client{}
}

// Run executes a tmux command and returns output
func (c *Client) Run(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("tmux error: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// RunSilent runs tmux command without output (for side effects)
func (c *Client) RunSilent(args ...string) error {
	cmd := exec.Command("tmux", args...)
	_, err := cmd.Output()
	return err
}

// ListSessions returns all tmux sessions
func (c *Client) ListSessions() ([]Session, error) {
	out, err := c.Run("list-sessions", "-F", "#{session_id}|#{session_name}|#{session_windows}|#{session_created}|#{session_attached}")
	if err != nil {
		return nil, err
	}

	var sessions []Session
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		sessions = append(sessions, ParseSession(line))
	}
	return sessions, nil
}

// ListWindows returns all windows for a session
func (c *Client) ListWindows(session string) ([]Window, error) {
	out, err := c.Run("list-windows", "-t", session, "-F", "#{window_id}|#{window_index}|#{window_name}|#{window_layout}|#{window_active}|#{window panes}")
	if err != nil {
		return nil, err
	}

	var windows []Window
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		windows = append(windows, ParseWindow(line))
	}
	return windows, nil
}

// ListPanes returns all panes for a window
func (c *Client) ListPanes(window string) ([]Pane, error) {
	out, err := c.Run("list-panes", "-t", window, "-F", "#{pane_id}|#{pane_index}|#{pane_title}|#{pane_current_command}|#{pane_active}")
	if err != nil {
		return nil, err
	}

	var panes []Pane
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		panes = append(panes, ParsePane(line))
	}
	return panes, nil
}

// AttachSession switches to a session
func (c *Client) AttachSession(session string) error {
	return c.RunSilent("switch-client", "-t", session)
}

// SelectWindow switches to a window in current session
func (c *Client) SelectWindow(window string) error {
	return c.RunSilent("select-window", "-t", window)
}

// NextLayout cycles to next layout
func (c *Client) NextLayout() error {
	return c.RunSilent("next-layout", "-t", strings.Split(GetCurrentPane(), ".")[0])
}

// GetCurrentPane returns current pane ID
func GetCurrentPane() string {
	cmd := exec.Command("tmux", "display-message", "-p", "#{session_id}:#{window_id}.#{pane_id}")
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}

// GetCurrentLayout returns current window layout name
func (c *Client) GetCurrentLayout() string {
	cmd := exec.Command("tmux", "display-message", "-p", "#{window_layout}")
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}

// NewWindow creates a new window in current session
func (c *Client) NewWindow() error {
	return c.RunSilent("new-window")
}

// SplitVertical splits current pane vertically
func (c *Client) SplitVertical() error {
	return c.RunSilent("split-window", "-v", "-c", "#{pane_current_path}")
}

// SplitHorizontal splits current pane horizontally
func (c *Client) SplitHorizontal() error {
	return c.RunSilent("split-window", "-h", "-c", "#{pane_current_path}")
}

// SelectPane moves to pane in direction
func (c *Client) SelectPane(direction string) error {
	return c.RunSilent("select-pane", "-"+direction)
}