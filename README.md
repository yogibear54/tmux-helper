# TMUXIFY

```
╺┳╸┏┳┓╻ ╻╻ ╻╻┏━╸╻ ╻
 ┃ ┃┃┃┃ ┃┏╋┛┃┣╸ ┗┳┛
 ╹ ╹ ╹┗━┛╹ ╹╹╹   ╹ 
```

**An i3-inspired tmux experience with intuitive keybindings and TUI overlays.**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-0.1.0-green.svg)]()

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
- ⚙️ **Configurable** - Edit `~/.tmux-helper.conf` and regenerate

## Quick Start

### One-line Install

```bash
curl -fsSL https://raw.githubusercontent.com/yogibear54/tmux-helper/main/install.sh | bash
```

### Manual Install

```bash
# Clone the repo
git clone git@github.com:yogibear54/tmux-helper.git
cd tmux-helper

# Run the installer (installs binary + config)
./install.sh

# Or do it manually:
go build -o tmux-helper ./cmd/tmux-helper
cp tmux-helper ~/.local/bin/
./tmux-helper apply
tmux source-file ~/.tmux.conf
```

### Restart tmux

```bash
tmux kill-server && tmux
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
tmux-helper --version          # Show version (v0.1.0)
tmux-helper --help             # Show help
tmux-helper picker             # Open session picker TUI
tmux-helper help-overlay       # Open help overlay TUI
tmux-helper sessions           # List all sessions (text)
tmux-helper layout             # Show current layout
tmux-helper layout-next        # Cycle to next layout
tmux-helper config show        # Show current configuration
tmux-helper config set <k> <v> # Set a config value
tmux-helper apply              # Regenerate ~/.tmux.conf
tmux-helper help               # Show keybindings (text)
```

## Configuration

Customize tmux-helper by editing `~/.tmux-helper.conf`:

```bash
# ~/.tmux-helper.conf
prefix=C-a
split-vertical-size=0.50
split-horizontal-size=0.50
mouse=true
theme=purple
terminal=screen-256color
```

### Configurable Options

| Option | Default | Description |
|--------|---------|-------------|
| `prefix` | `C-a` | Command prefix key |
| `split-vertical-size` | `0.50` | Vertical split size (0.1-0.9) |
| `split-horizontal-size` | `0.50` | Horizontal split size (0.1-0.9) |
| `mouse` | `true` | Enable mouse support |
| `theme` | `purple` | UI theme (purple/green) |
| `terminal` | `screen-256color` | Terminal type |

### Config Commands

```bash
tmux-helper config show           # View current config
tmux-helper config get <key>     # Get a specific value
tmux-helper config set <key> <value>  # Set a value
tmux-helper config edit          # Edit in $EDITOR
tmux-helper apply                # Regenerate ~/.tmux.conf
```

### Workflow

```bash
1. Edit ~/.tmux-helper.conf (or use `config set`)
2. Run `tmux-helper apply`
3. Run `tmux source-file ~/.tmux.conf`
```

## Architecture

```
tmux-helper/
├── cmd/tmux-helper/
│   └── main.go                  # CLI entry point (v0.1.0)
├── internal/
│   ├── tmux/
│   │   ├── client.go           # tmux CLI wrapper
│   │   ├── model.go            # Session/Window/Pane structs
│   │   └── errors.go           # Notifications & error handling
│   ├── config/
│   │   ├── config.go           # Config loading/saving/validation
│   │   └── generate.go         # tmux.conf template generator
│   └── ui/
│       ├── common.go           # Shared UI styles
│       ├── picker.go           # Session picker (Purple TUI)
│       └── help.go             # Help overlay (Green TUI)
├── install.sh                   # Installation script
├── build-release.sh             # Release builder
├── .github/workflows/           # CI/CD
│   ├── ci.yml                  # Continuous integration
│   └── release.yml             # Auto-release on tags
├── .tmux.conf                  # Generated tmux config
├── go.mod
└── README.md
```

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

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 | ✅ | Foundation (Go project, tmux client, .tmux.conf) |
| Phase 2 | ✅ | Session Picker (Bubble Tea TUI + display-popup) |
| Phase 3 | ✅ | Layout Cycling (built into tmux) |
| Phase 4 | ✅ | Help Overlay (Bubble Tea TUI + display-popup) |
| Phase 5 | ✅ | Configuration System (config file, generate tmux.conf) |
| Phase 6 | ✅ | Polish & Error Handling (notifications, edge cases, version) |
| Phase 7 | ✅ | Installer & Distribution (install.sh, release workflow) |

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
export PATH="$PATH:$HOME/.local/bin"
```

### tmux not responding to keybindings
```bash
tmux source-file ~/.tmux.conf
```

### TUI popups won't open
Make sure your tmux version supports `display-popup`:

```bash
tmux -V  # Should be 3.3 or newer
```

### Mouse not working
```bash
tmux set -g mouse on
```

## Dependencies

- **Go 1.21+** - Required to build
- **tmux 3.3+** - Required for `display-popup` support
- **Bubble Tea** - TUI framework (charmbracelet/bubbletea)
- **Lipgloss** - Styling library (charmbracelet/lipgloss)

## Releases

### Creating a Release

```bash
# Build and package
./build-release.sh 0.1.0

# Create git tag
git tag v0.1.0
git push origin v0.1.0

# GitHub Actions will build and create a draft release
```

### Release Assets

| File | Description |
|------|-------------|
| `tmux-helper-X.Y.Z` | Static binary (4.5MB) |
| `tmux-helper-X.Y.Z.sha256` | SHA256 checksum |
| `install.sh` | Installation script |
| `README.md` | Documentation |

## License

MIT License - See [LICENSE](LICENSE) for details.

## Contributing

Contributions welcome! Feel free to submit issues and pull requests.
