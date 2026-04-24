package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Notify displays a message in tmux status bar
func Notify(message string) {
	cmd := exec.Command("tmux", "display-message", "-p", message)
	cmd.Run()
}

// NotifySuccess displays a success message
func NotifySuccess(message string) {
	cmd := exec.Command("tmux", "display-message", "-p", fmt.Sprintf("#[fg=green]✓#[default] %s", message))
	cmd.Run()
}

// NotifyError displays an error message
func NotifyError(message string) {
	cmd := exec.Command("tmux", "display-message", "-p", fmt.Sprintf("#[fg=red]✗#[default] %s", message))
	cmd.Run()
}

// NotifyInfo displays an info message
func NotifyInfo(message string) {
	cmd := exec.Command("tmux", "display-message", "-p", fmt.Sprintf("#[fg=cyan]ℹ#[default] %s", message))
	cmd.Run()
}

// ErrorHandler provides consistent error handling
type ErrorHandler struct {
	Silent bool
}

// NewErrorHandler creates a new error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{Silent: false}
}

// Handle processes an error with optional user notification
func (e *ErrorHandler) Handle(err error, context string, notify bool) error {
	if err == nil {
		return nil
	}

	errMsg := fmt.Sprintf("%s: %v", context, err)
	
	if notify && !e.Silent {
		// Check if we're in tmux
		if os.Getenv("TMUX") != "" {
			NotifyError(errMsg)
		}
	}
	
	return err
}

// MustExecute runs a command and panics on error
func MustExecute(args ...string) {
	cmd := exec.Command("tmux", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("tmux %s: %w (output: %s)", strings.Join(args, " "), err, string(out)))
	}
}
