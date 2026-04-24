#!/bin/bash
# Launcher script for tmux-helper picker
# Handles TTY properly for Bubble Tea TUI

# Ensure we have a proper terminal
if [ -z "$TMUX" ]; then
    # Not in tmux, just run normally
    tmux-helper picker
else
    # In tmux, use new-window to get a proper TTY
    tmux new-window -n "tmux-helper-picker" "tmux-helper picker"
fi