# tmux-helper

**An i3-inspired tmux experience with intuitive keybindings and a TUI session picker.**

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

tmux-helper makes tmux more intuitive by bringing i3-window-manager-style keybindings to tmux. It includes a TUI application for quick session/window management and a carefully designed `.tmux.conf` with vim-style navigation.

## Features

- 🎯 **Session/Window Quick Picker** - Fuzzy search and select sessions via TUI
- 🔀 **Layout Cycling** - Cycle through pane layouts with a single key
- 📐 **Direction-based Splits** - Split panes using Alt+h/j/k/l (i3-style)
- 🖱️ **Mouse Support** - Click to select, drag to resize panes
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

# Install binary (optional)
sudo cp tmux-helper /usr/local/bin/
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
tmux kill-server  # or press Prefix+d to detach

# Start a new tmux session
tmux
```

## Keybindings

### Prefix Key
| Key | Description |
|-----|-------------|
| `Ctrl-a` | Prefix (like screen, i3-friendly alternative to tmux's default `Ctrl-b`) |

### Pane Navigation (Vim-style)
| Key | Action |
|-----|--------|
| `h` | Move to pane on the **left** |
| `j` | Move to pane **below** |
| `k` | Move to pane **above** |
| `l` | Move to pane on the **right** |
| `Ctrl-a + ?` | Show help overlay |

### Direction-based Splits (Alt without prefix)
| Key | Action |
|-----|--------|
| `Alt-h` | Split pane **left** (vertical split) |
| `Alt-j` | Split pane **down** (horizontal split) |
| `Alt-k` | Split pane **up** (horizontal split) |
| `Alt-l` | Split pane **right** (vertical split) |

### Layout Cycling
| Key | Action |
|-----|--------|
| `Ctrl-a + Space` | Cycle to **next** layout |
| `Ctrl-a + Shift+Space` | Cycle to **previous** layout |

### Session/Window Management
| Key | Action |
|-----|--------|
| `Ctrl-a + F` | Open session picker (TUI) |
| `Ctrl-a + c` | Create new window |
| `Ctrl-a + x` | Kill current pane |
| `Ctrl-a + X` | Kill current window |
| `Ctrl-a + d` | Detach from session |

### Pane Operations
| Key | Action |
|-----|--------|
| `Ctrl-a + !` | Break pane into new window |
| `Ctrl-a + :` | Join pane into current window |
| `Ctrl-a + z` | Zoom pane (toggle) |
| `Ctrl-a + H/J/K/L` | Swap pane with adjacent |

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
| Drag pane edge | Resize pane |
| Right-click | Context menu |

## Command Line Interface

```bash
tmux-helper picker      # List sessions
tmux-helper layout      # Show current layout
tmux-helper help        # Show keybindings
```

## Configuration

Create `~/.tmux-helper.conf` to customize behavior:

```bash
# Split sizes (0.0 - 1.0)
split-vertical-size=0.5
split-horizontal-size=0.5

# Prefix key (default: C-a)
prefix=C-a

# Mouse support
mouse=true
```

## Architecture

```
tmux-helper/
├── cmd/tmux-helper/main.go     # CLI entry point
├── internal/
│   ├── tmux/
│   │   ├── client.go          # tmux CLI wrapper
│   │   └── model.go           # Data models
│   ├── config/
│   │   └── config.go          # Configuration
│   └── ui/
│       └── common.go          # UI components
├── .tmux.conf                  # tmux keybindings
└── README.md
```

### Components

- **Client** - Wraps tmux CLI for session/window/pane operations
- **Models** - Session, Window, Pane structs with parsers
- **Config** - User preferences (split sizes, prefix key, etc.)
- **UI** - Bubble Tea TUI for session picker (Phase 2)

## Development

### Build

```bash
go build -o tmux-helper ./cmd/tmux-helper
```

### Test

```bash
# Start a test session
tmux new-session -d -s test

# Run the picker
./tmux-helper picker

# Clean up
tmux kill-session -t test
```

### Phases

- [x] Phase 1: Foundation (Go project, tmux client, .tmux.conf)
- [ ] Phase 2: Session/Window Quick Picker (TUI)
- [ ] Phase 3: Layout Cycling
- [ ] Phase 4: Help Overlay
- [ ] Phase 5: Configuration System
- [ ] Phase 6: Polish & Error Handling
- [ ] Phase 7: Installer & Distribution

## i3 Comparison

| i3 | tmux-helper | Notes |
|----|-------------|-------|
| Mod+Enter | `Ctrl-a + c` | New window |
| Mod+d | `Ctrl-a + d` | Detach |
| Mod+h/j/k/l | `h/j/k/l` | Navigate panes |
| Mod+shift+h/j/k/l | `Alt+h/j/k/l` | Split direction |
| Mod+e | `Ctrl-a + Space` | Cycle layouts |
| Mod+1-9 | `Ctrl-a + 1-9` | Switch windows |
| Mod+w | `Ctrl-a + F` | Session picker (TUI) |

## Troubleshooting

### "tmux-helper: command not found"
```bash
# Make sure the binary is in your PATH
export PATH=$PATH:/path/to/tmux-helper

# Or use full path
/path/to/tmux-helper picker
```

### tmux not responding to keybindings
```bash
# Reload the configuration
tmux source-file ~/.tmux.conf
```

### Mouse not working
```bash
# Enable mouse mode (should be on by default)
tmux set -g mouse on
```

## License

MIT License - See [LICENSE](LICENSE) for details.

## Contributing

Contributions welcome! Feel free to submit issues and pull requests.