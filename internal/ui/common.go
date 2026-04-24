package ui

import (
	"github.com/charmbracelet/bubbletea"
)

// Common UI components and utilities

// Item represents a selectable item in lists
type Item interface {
	FilterValue() string
}

// Model is the interface for all UI models
type Model interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}