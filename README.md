# tmux-helper

**An i3-inspired tmux experience with intuitive keybindings and TUI overlays.**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)

## Overview

tmux-helper makes tmux more intuitive by bringing i3-window-manager-style keybindings to tmux. It includes TUI applications for session picking and help overlays, plus a carefully designed `.tmux.conf` with vim-style navigation.

## Features

- 🎯 **Session/Window Quick Picker** - Interactive TUI with fuzzy search (Purple theme)
- ❓ **Help Overlay** - Beautiful popup showing all keybindings (Green theme)
- 🔀 **Layout Cycling** - Cycle through pane layouts with a single key
- 📐 **Split Panes** - Split using `Ctrl-a + |` or `Ctrl-a + -`
- 🖱️ **Mouse Support** - Click to select panes
- ⌨️ **Vim-style Navigation** - Navigate panes with h/j/k/l
- 🚀 **Single Binary** - No dependencies, just drop and run

## Quick Start

### 1. Install tmux-helper

```bash
# Clone the repo
git clone https://github.com/lotus-creations/tmux-helper.git
cd tmux-helper

# Build
go build -o tmux-helper ./cmd/tmux-helper

# Install binary
cp tmux-helper ~/.local/bin/
```

### 2. Copy the tmux configuration

```bash
# Backup your existing config (optional)
cp ~/.tmux.conf ~/.tmux.conf.backup

# Link the new config
ln -sf $(pwd)/.tmux.conf ~/.tmux.conf
```

### 3. Restart tmux

```bash
# If tmux is running, restart it
tmux kill-server

# Start a new tmux session
tmux
```

## Keybindings

### Prefix Key
| Key | Description |
|-----|-------------|
| `Ctrl-a` | Command prefix (like screen, i3-friendly) |

### Pane Navigation (Vim-style)
| Key | Action |
|-----|--------|
| `h` | Move to pane on the **left** |
| `j` | Move to pane **below** |
| `k` | Move to pane **above** |
| `l` | Move to pane on the **right** |

### Splits (Prefix + key)
| Key | Action |
|-----|--------|
| `Ctrl-a + \|` | Split left/right (vertical) |
| `Ctrl-a + -` | Split top/bottom (horizontal) |

### Layout Cycling
| Key | Action |
|-----|--------|
| `Ctrl-a + Space` | Cycle to **next** layout |

### Session/Window Management
| Key | Action |
|-----|--------|
| `Ctrl-a + F` | Open session picker (TUI popup) |
| `Ctrl-a + ?` | Open help overlay (TUI popup) |
| `Ctrl-a + c` | Create new window |
| `Ctrl-a + x` | Kill current pane |
| `Ctrl-a + X` | Kill current window |
| `Ctrl-a + d` | Detach from session |

### Pane Operations
| Key | Action |
|-----|--------|
| `Ctrl-a + !` | Break pane into new window |
| `Ctrl-a + H/J/K/L` | Swap pane with adjacent (Shift+direction) |

### Copy Mode (Vim-style)
| Key | Action |
|-----|--------|
| `Ctrl-a + [` | Enter copy mode |
| `v` | Begin selection (in copy mode) |
| `y` | Copy selection (in copy mode) |
| `Enter` | Copy selection (in copy mode) |

### Mouse Controls
| Action | Result |
|--------|--------|
| Click pane | Select pane |

## TUI Overlays

### Session Picker (`Ctrl-a + F`)
Opens an interactive TUI popup for quick session/window selection:

```
╭────────────────────────────────╮
│   tmux-helper session picker   │
│ 🔍 Filter sessions...           │
│                                │
│  ○ dev                         │
│    └─ 1: editor                │
│    └─ 2: terminal              │
│  ○ prod                        │
│    └─ 1: server                │
│                                │
╰────────────────────────────────╯
 ↑↓ Navigate | Enter Select | Ctrl+U Clear | Esc Quit
```

**Features:**
- Fuzzy search filtering
- Session hierarchy with windows
- Keyboard navigation (↑↓, Enter, Esc)
- Click to select
- Purple/violet theme

### Help Overlay (`Ctrl-a + ?`)
Opens a popup showing all keybindings:

```
╭────────────────────────────────╮
│   tmux-helper Keybindings       │
│                                │
│ PREFIX                         │
│   Ctrl-a          Command prefix│
│                                │
│ PANE NAVIGATION (vim-style)   │
│   h              Move left      │
│   j              Move down      │
│   k              Move up        │
│   l              Move right     │
│                                │
│ SPLITS                         │
│   Ctrl-a + |      Split vertical│
│   Ctrl-a + -      Split horiz.  │
│                                │
│ Press Esc or Enter to close    │
╰────────────────────────────────╯
```

**Features:**
- Green header
- Sectioned layout
- Blue key styling
- Rounded border
- Press Esc/Enter to close

## Command Line Interface

```bash
tmux-helper picker         # Open session picker TUI (display-popup)
tmux-helper help-overlay   # Open help overlay TUI (display-popup)
tmux-helper sessions       # List all sessions (text output)
tmux-helper layout         # Show current layout
tmux-helper layout-next    # Cycle to next layout
tmux-helper help           # Show keybindings (text, no popup)
```

## Architecture

```
tmux-helper/
├── cmd/tmux-helper/
│   └── main.go              # CLI entry point
├── internal/
│   ├── tmux/
│   │   ├── client.go       # tmux CLI wrapper
│   │   └── model.go        # Session/Window/Pane structs
│   ├── config/
│   │   └── config.go       # Configuration system
│   └── ui/
│       ├── common.go       # Shared UI components
│       ├── picker.go       # Session picker TUI (Bubble Tea)
│       └── help.go         # Help overlay TUI (Bubble Tea)
├── .tmux.conf              # tmux keybindings
├── go.mod
└── README.md
```

### Components

| Component | Description |
|-----------|-------------|
| **client.go** | Wraps tmux CLI for session/window/pane operations |
| **model.go** | Session, Window, Pane structs with parsers |
| **config.go** | User preferences (placeholder for future config) |
| **picker.go** | Purple-themed Bubble Tea TUI for session picking |
| **help.go** | Green-themed Bubble Tea TUI for help overlay |
| **.tmux.conf** | i3-inspired keybindings, uses `display-popup` for TUIs |

## Development

### Build

```bash
go build -o tmux-helper ./cmd/tmux-helper
```

### Run

```bash
# Start tmux with test sessions
tmux new-session -d -s dev -n editor
tmux new-session -d -s prod -n server

# Open the session picker
tmux-helper picker

# Or list sessions
tmux-helper sessions
```

### Clean up

```bash
# Kill all sessions
tmux kill-server
```

## Project Status

### Phases

- [x] Phase 1: Foundation (Go project, tmux client, .tmux.conf)
- [x] Phase 2: Session/Window Quick Picker (Bubble Tea TUI + display-popup)
- [x] Phase 3: Layout Cycling (built into tmux)
- [x] Phase 4: Help Overlay (Bubble Tea TUI + display-popup)
- [ ] Phase 5: Configuration System
- [ ] Phase 6: Polish & Error Handling
- [ ] Phase 7: Installer & Distribution

## i3 Comparison

| i3 Action | tmux-helper | Notes |
|-----------|-------------|-------|
| `Mod+Enter` | `Ctrl-a + c` | New window |
| `Mod+d` | `Ctrl-a + d` | Detach |
| `Mod+h/j/k/l` | `h/j/k/l` | Navigate panes |
| `Mod+Shift+h/j/k/l` | `Ctrl-a + \| / -` | Split panes |
| `Mod+e` | `Ctrl-a + Space` | Cycle layouts |
| `Mod+1-9` | `Ctrl-a + 1-9` | Switch windows |
| `Mod+w` | `Ctrl-a + F` | Session picker (TUI) |
| `Mod+?` | `Ctrl-a + ?` | Help overlay (TUI) |

## Troubleshooting

### "tmux-helper: command not found"
```bash
# Make sure the binary is in your PATH
export PATH=$PATH:~/.local/bin

# Or use full path
~/.local/bin/tmux-helper picker
```

### tmux not responding to keybindings
```bash
# Reload the configuration
tmux source-file ~/.tmux.conf
```

### TUI popups won't open
Make sure you're running inside tmux (the TUI requires a TTY). Also verify your tmux version supports `display-popup`:

```bash
tmux -V  # Should be 3.3 or newer for display-popup
```

### Mouse not working
```bash
# Enable mouse mode
tmux set -g mouse on
```

## Dependencies

- **Go 1.21+** - Required to build
- **tmux 3.3+** - Required for `display-popup` support
- **Bubble Tea** - TUI framework (charmbracelet/bubbletea)
- **Lipgloss** - Styling library (charmbracelet/lipgloss)

## License

MIT License - See [LICENSE](LICENSE) for details.

## Contributing

Contributions welcome! Feel free to submit issues and pull requests.