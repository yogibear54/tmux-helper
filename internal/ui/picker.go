package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/lotus-creations/tmux-helper/internal/tmux"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#4F46E5")).
			Padding(0, 1).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D4D4D4")).
			Padding(0, 1)

	attachedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E")).
			Padding(0, 1)

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Background(lipgloss.Color("#1F2937"))
)

// PickerItem represents a session or window in the list
type PickerItem struct {
	Type     string // "session" or "window"
	Name     string
	Index    int
	Parent   string
	Attached bool
}

func (i PickerItem) FilterValue() string {
	return i.Name + " " + i.Parent
}

// Picker is the main session picker model
type Picker struct {
	sessions   []tmux.Session
	windows    map[string][]tmux.Window
	items      []PickerItem
	filtered   []PickerItem
	cursor     int
	filter     textinput.Model
	viewport   viewport.Model
	ready      bool
	height     int
	width      int
}

// NewPicker creates a new session picker
func NewPicker() *Picker {
	ti := textinput.New()
	ti.Placeholder = "Filter sessions..."
	ti.Focus()
	ti.Prompt = "🔍 "

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.Style{}

	return &Picker{
		sessions: []tmux.Session{},
		windows:  make(map[string][]tmux.Window),
		items:    []PickerItem{},
		filtered: []PickerItem{},
		cursor:   0,
		filter:   ti,
		viewport: vp,
	}
}

// Init initializes the picker
func (p *Picker) Init() tea.Cmd {
	client := tmux.NewClient()
	sessions, err := client.ListSessions()
	if err != nil {
		// No tmux running
		p.sessions = []tmux.Session{}
	} else {
		p.sessions = sessions
		for _, s := range sessions {
			windows, err := client.ListWindows(s.Name)
			if err == nil {
				p.windows[s.Name] = windows
			}
		}
	}
	p.rebuildItems()
	return nil
}

// rebuildItems rebuilds the item list from sessions
func (p *Picker) rebuildItems() {
	p.items = []PickerItem{}

	for _, s := range p.sessions {
		// Add session header
		p.items = append(p.items, PickerItem{
			Type:     "session",
			Name:     s.Name,
			Index:    0,
			Parent:   "",
			Attached: s.Attached,
		})

		// Add windows under session
		for _, w := range p.windows[s.Name] {
			p.items = append(p.items, PickerItem{
				Type:     "window",
				Name:     w.Name,
				Index:    w.Index,
				Parent:   s.Name,
				Attached: w.Active,
			})
		}
	}

	p.applyFilter()
}

// applyFilter filters items based on search query
func (p *Picker) applyFilter() {
	query := strings.ToLower(p.filter.Value())
	if query == "" {
		p.filtered = p.items
	} else {
		p.filtered = []PickerItem{}
		for _, item := range p.items {
			if strings.Contains(strings.ToLower(item.Name), query) ||
				strings.Contains(strings.ToLower(item.Parent), query) {
				p.filtered = append(p.filtered, item)
			}
		}
	}

	// Reset cursor if out of bounds
	if p.cursor >= len(p.filtered) {
		p.cursor = max(0, len(p.filtered)-1)
	}
}

// Update handles messages
func (p *Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return p, tea.Quit

		case tea.KeyUp, tea.KeyShiftTab:
			if p.cursor > 0 {
				p.cursor--
			}

		case tea.KeyDown, tea.KeyTab:
			if p.cursor < len(p.filtered)-1 {
				p.cursor++
			}

		case tea.KeyEnter:
			if len(p.filtered) > 0 && p.cursor < len(p.filtered) {
				item := p.filtered[p.cursor]
				p.selectItem(item)
				return p, tea.Quit
			}

		case tea.KeyCtrlU:
			p.filter.SetValue("")
			p.applyFilter()
		}

	case tea.WindowSizeMsg:
		p.height = msg.Height
		p.width = msg.Width
		p.viewport.Width = msg.Width
		p.viewport.Height = msg.Height - 5 // Space for header, filter, help

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// Calculate which item was clicked
			row := msg.Y - 2 // Account for header and filter
			if row >= 0 && row < len(p.filtered) {
				p.cursor = row
				item := p.filtered[p.cursor]
				p.selectItem(item)
				return p, tea.Quit
			}
		}
	}

	// Update filter input
	newFilter, cmd := p.filter.Update(msg)
	p.filter = newFilter
	cmds = append(cmds, cmd)

	// Check if filter changed
	if p.filter.Value() != "" || len(p.filtered) != len(p.items) {
		p.applyFilter()
	}

	// Scroll viewport
	p.viewport.SetContent(p.renderList())

	return p, tea.Batch(cmds...)
}

// selectItem handles item selection
func (p *Picker) selectItem(item PickerItem) {
	if item.Type == "session" {
		// Switch to session
		exec.Command("tmux", "switch-client", "-t", item.Name).Run()
	} else {
		// Switch to window in session
		exec.Command("tmux", "switch-client", "-t", item.Parent).Run()
		exec.Command("tmux", "select-window", "-t", fmt.Sprintf("%d", item.Index)).Run()
	}
}

// renderList renders the item list
func (p *Picker) renderList() string {
	var lines []string

	// Check for no tmux running
	if len(p.sessions) == 0 {
		warningStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Bold(true)
		lines = append(lines, warningStyle.Render("  ⚠ No tmux sessions running"))
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Render("  Start tmux with: tmux new -s <name>"))
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Render("  Press Esc to exit"))
	} else if len(p.filtered) == 0 {
		// Sessions exist but filtered results are empty
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Render("  No matching sessions"))
		lines = append(lines, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Render("  Try a different filter or press Ctrl+U to clear"))
	} else {
		for i, item := range p.filtered {
			var prefix, name, suffix string

			if item.Type == "session" {
				if item.Attached {
					prefix = "● "
					suffix = " (attached)"
				} else {
					prefix = "○ "
					suffix = ""
				}
				name = item.Name + suffix
			} else {
				prefix = "  └─ "
				name = fmt.Sprintf("%d: %s", item.Index, item.Name)
			}

			var style lipgloss.Style
			if i == p.cursor {
				style = selectedStyle
			} else if item.Attached && item.Type == "session" {
				style = attachedStyle
			} else {
				style = normalStyle
			}

			lines = append(lines, style.Render(prefix+name))
		}
	}

	return strings.Join(lines, "\n")
}

// View renders the picker
func (p *Picker) View() string {
	if !p.ready {
		p.ready = true
	}

	header := headerStyle.Render(" tmux-helper session picker ")

	filterView := p.filter.View()

	helpBar := helpStyle.Render(" ↑↓ Navigate  |  Enter Select  |  Ctrl+U Clear  |  Esc Quit ")

	content := lipgloss.JoinVertical(lipgloss.Left, header, filterView, "", p.renderList())

	// Wrap content with border
	borderBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6B7280")).
		Padding(1).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top, borderBox, helpBar)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RunPicker starts the session picker
func RunPicker() {
	p := tea.NewProgram(NewPicker())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running picker: %v\n", err)
		os.Exit(1)
	}
}