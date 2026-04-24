#!/bin/bash
#
# tmux-helper installer
# Installs tmux-helper and configures tmux
#

set -e

VERSION="0.1.0"
INSTALL_DIR="${HOME}/.local/bin"
CONFIG_FILE="${HOME}/.tmux-helper.conf"
TMUX_CONF="${HOME}/.tmux.conf"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[✗]${NC} $1"
}

banner() {
    cat << 'EOF'
    _                                _
   (_) ___ _ __ __ _ _ __   __ _  ___| |__
   | |/ _ \ '__/ _` | '_ \ / _` |/ __| '_ \
   | |  __/ | | (_| | | | | (_| | (__| | | |
   |_|\___|_|  \__,_|_| |_|\__,_|\___|_| |_|

   i3-inspired tmux configuration with TUI overlays

EOF
}

check_requirements() {
    info "Checking requirements..."

    # Check for Go (try PATH first, then known locations)
    GO_CMD=""
    if command -v go &> /dev/null; then
        GO_CMD="go"
    elif [ -x "/usr/local/go/bin/go" ]; then
        GO_CMD="/usr/local/go/bin/go"
    fi
    
    if [ -z "$GO_CMD" ]; then
        error "Go not found. Please install Go 1.21+"
        exit 1
    fi
    GO_VERSION=$($GO_CMD version | grep -oP '\d+\.\d+' | head -1)
    success "Go ${GO_VERSION} found"

    # Check for tmux
    if command -v tmux &> /dev/null; then
        TMUX_VERSION=$(tmux -V | grep -oP '\d+\.\d+')
        success "tmux ${TMUX_VERSION} found"
    else
        error "tmux not found. Please install tmux 3.3+"
        exit 1
    fi

    # Check tmux version for display-popup support
    if (( $(echo "$TMUX_VERSION < 3.3" | bc -l) )); then
        warn "tmux 3.3+ recommended for display-popup support"
        warn "Your version: $TMUX_VERSION"
    fi
}

install_binary() {
    info "Installing tmux-helper binary..."

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Determine script location
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

    # Check if we're running from source (look in multiple places)
    if [ -f "${SCRIPT_DIR}/tmux-helper" ]; then
        # Binary in root directory
        info "Installing from root directory..."
        cp "${SCRIPT_DIR}/tmux-helper" "${INSTALL_DIR}/tmux-helper"
        success "Installed from root directory"
    elif [ -f "${SCRIPT_DIR}/release/tmux-helper-${VERSION}" ]; then
        # Binary in release directory
        info "Installing from release directory..."
        cp "${SCRIPT_DIR}/release/tmux-helper-${VERSION}" "${INSTALL_DIR}/tmux-helper"
        success "Installed from release"
    elif [ -f "${SCRIPT_DIR}/release/tmux-helper" ]; then
        # Binary in release directory (without version)
        info "Installing from release directory..."
        cp "${SCRIPT_DIR}/release/tmux-helper" "${INSTALL_DIR}/tmux-helper"
        success "Installed from release"
    else
        # Build from source
        info "Building from source..."
        (cd "${SCRIPT_DIR}" && /usr/local/go/bin/go build -o tmux-helper ./cmd/tmux-helper)
        cp "${SCRIPT_DIR}/tmux-helper" "${INSTALL_DIR}/tmux-helper"
        success "Built and installed from source"
    fi

    # Make executable
    chmod +x "${INSTALL_DIR}/tmux-helper"

    # Add to PATH if needed
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        warn "${INSTALL_DIR} is not in your PATH"
        info "Add this to your ~/.bashrc or ~/.zshrc:"
        echo ""
        echo -e "    ${GREEN}export PATH=\"\$PATH:${INSTALL_DIR}\"${NC}"
        echo ""
    fi

    success "Binary installed to ${INSTALL_DIR}/tmux-helper"
}

install_config() {
    info "Setting up configuration..."

    # Backup existing tmux.conf
    if [ -f "${TMUX_CONF}" ] && [ ! -L "${TMUX_CONF}" ]; then
        BACKUP="${TMUX_CONF}.backup.$(date +%Y%m%d%H%M%S)"
        warn "Backing up existing ~/.tmux.conf to ${BACKUP}"
        cp "${TMUX_CONF}" "${BACKUP}"
    fi

    # Generate initial config
    info "Generating ~/.tmux-helper.conf..."
    "${INSTALL_DIR}/tmux-helper" config set prefix "C-a" 2>/dev/null || true

    # Generate tmux.conf
    info "Generating ~/.tmux.conf..."
    "${INSTALL_DIR}/tmux-helper" apply

    success "Configuration installed"
}

install_shell_integration() {
    info "Installing shell integration..."

    # Detect shell
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_RC="${HOME}/.zshrc"
        SHELL_NAME="zsh"
    elif [ -n "$BASH_VERSION" ]; then
        SHELL_RC="${HOME}/.bashrc"
        SHELL_NAME="bash"
    else
        SHELL_RC="${HOME}/.profile"
        SHELL_NAME="sh"
    fi

    # Add tmux-helper completion if available
    COMPLETION_MARKER="# tmux-helper"
    if ! grep -q "${COMPLETION_MARKER}" "${SHELL_RC}" 2>/dev/null; then
        cat >> "${SHELL_RC}" << 'RC_EOF'

# tmux-helper
export PATH="$PATH:${INSTALL_DIR}"
RC_EOF
        success "Added tmux-helper to ${SHELL_RC}"
    fi
}

finalize() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    success "Installation complete!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Next steps:"
    echo ""
    echo "1. Reload your shell or run:"
    echo -e "   ${GREEN}source ~/.bashrc${NC}  # or ~/.zshrc"
    echo ""
    echo "2. Restart tmux with:"
    echo -e "   ${GREEN}tmux kill-server && tmux${NC}"
    echo ""
    echo "   Or if you have sessions you want to keep:"
    echo -e "   ${GREEN}tmux source-file ~/.tmux.conf${NC}"
    echo ""
    echo "3. Test the session picker:"
    echo -e "   ${GREEN}Ctrl-a F${NC}"
    echo ""
    echo "4. View help:"
    echo -e "   ${GREEN}Ctrl-a ?${NC}"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

uninstall() {
    info "Uninstalling tmux-helper..."

    # Remove binary
    if [ -f "${INSTALL_DIR}/tmux-helper" ]; then
        rm "${INSTALL_DIR}/tmux-helper"
        success "Removed binary"
    fi

    # Ask about config files
    read -p "Remove ~/.tmux-helper.conf and ~/.tmux.conf? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        [ -f "${CONFIG_FILE}" ] && rm "${CONFIG_FILE}" && success "Removed config"
        [ -f "${TMUX_CONF}" ] && rm "${TMUX_CONF}" && success "Removed tmux.conf"
    fi

    success "Uninstallation complete"
}

# Main
main() {
    banner

    case "${1:-install}" in
        install)
            check_requirements
            install_binary
            install_config
            install_shell_integration
            finalize
            ;;
        uninstall|remove)
            uninstall
            ;;
        update)
            check_requirements
            install_binary
            info "Run 'tmux source-file ~/.tmux.conf' to apply changes"
            ;;
        help|--help|-h)
            echo "Usage: $0 [install|uninstall|update|help]"
            echo ""
            echo "Commands:"
            echo "  install    Install tmux-helper (default)"
            echo "  uninstall  Remove tmux-helper"
            echo "  update     Update to new version"
            echo "  help       Show this help"
            ;;
        *)
            error "Unknown command: $1"
            echo "Use: $0 help"
            exit 1
            ;;
    esac
}

main "$@"
