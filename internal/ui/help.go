package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Help styles
var (
	helpHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#059669")).
				Padding(0, 1)

	helpSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F59E0B")).
				MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			Width(20)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D4D4D4"))

	helpFooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Background(lipgloss.Color("#1F2937")).
			MarginTop(1)

	helpBorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6B7280")).
			Padding(1)
)

// HelpOverlay displays keybindings in a TUI overlay
type HelpOverlay struct {
	width  int
	height int
}

func NewHelpOverlay() *HelpOverlay {
	return &HelpOverlay{
		width:  60,
		height: 30,
	}
}

// Init initializes the help overlay
func (h *HelpOverlay) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (h *HelpOverlay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).Type {
		case tea.KeyCtrlC, tea.KeyEscape, tea.KeyEnter, tea.KeySpace:
			return h, tea.Quit
		}
	case tea.WindowSizeMsg:
		h.width = msg.(tea.WindowSizeMsg).Width
		h.height = msg.(tea.WindowSizeMsg).Height
	}
	return h, nil
}

// View renders the help overlay
func (h *HelpOverlay) View() string {
	// Calculate max width based on content
	maxWidth := 50

	content := lipgloss.JoinVertical(
		lipgloss.Top,
		helpHeaderStyle.Render(" tmux-helper Keybindings "),
		"",
		helpSectionStyle.Render("PREFIX"),
		renderHelpLine("Ctrl-a", "Command prefix"),
		"",
		helpSectionStyle.Render("PANE NAVIGATION (vim-style)"),
		renderHelpLine("h", "Move left"),
		renderHelpLine("j", "Move down"),
		renderHelpLine("k", "Move up"),
		renderHelpLine("l", "Move right"),
		"",
		helpSectionStyle.Render("SPLITS"),
		renderHelpLine("Ctrl-a + |", "Split left/right (vertical)"),
		renderHelpLine("Ctrl-a + -", "Split top/bottom (horizontal)"),
		"",
		helpSectionStyle.Render("LAYOUTS"),
		renderHelpLine("Ctrl-a + Space", "Cycle to next layout"),
		"",
		helpSectionStyle.Render("SESSIONS"),
		renderHelpLine("Ctrl-a + F", "Open session picker"),
		renderHelpLine("Ctrl-a + c", "New window"),
		renderHelpLine("Ctrl-a + d", "Detach"),
		"",
		helpSectionStyle.Render("PANE MANAGEMENT"),
		renderHelpLine("Ctrl-a + x", "Kill current pane"),
		renderHelpLine("Ctrl-a + X", "Kill current window"),
		renderHelpLine("Ctrl-a + !", "Break pane into new window"),
		renderHelpLine("Ctrl-a + H/J/K/L", "Swap with adjacent pane"),
		"",
		helpSectionStyle.Render("MOUSE"),
		renderHelpLine("Click", "Select pane"),
		"",
		helpSectionStyle.Render("COPY MODE"),
		renderHelpLine("Ctrl-a + [", "Enter copy mode"),
		renderHelpLine("v", "Begin selection"),
		renderHelpLine("y", "Copy selection"),
	)

	// Wrap content with border
	borderBox := helpBorderStyle.Width(maxWidth + 4).Render(content)
	footer := helpFooterStyle.Render(" Press Esc or Enter to close ")

	return lipgloss.JoinVertical(lipgloss.Center, borderBox, footer)
}

// renderHelpLine creates a formatted help line
func renderHelpLine(key, description string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		helpKeyStyle.Render(key),
		helpDescStyle.Render(description),
	)
}

// RunHelp displays the help overlay
func RunHelp() {
	p := tea.NewProgram(NewHelpOverlay())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running help: %v\n", err)
	}
}