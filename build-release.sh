#!/bin/bash
#
# tmux-helper release builder
# Builds binaries for distribution
#

set -e

VERSION="${1:-0.1.0}"
BUILD_DIR="release"
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

info() {
    echo "[INFO] $1"
}

success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

# Clean previous builds
info "Cleaning previous builds..."
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

# Build for current platform
info "Building for ${OS}/${ARCH}..."
export CGO_ENABLED=0
/usr/local/go/bin/go build -ldflags="-s -w -X main.version=${VERSION}" -o "${BUILD_DIR}/tmux-helper-${VERSION}" ./cmd/tmux-helper

# Make installer executable
chmod +x install.sh

success "Built: ${BUILD_DIR}/tmux-helper-${VERSION}"

# Create checksums
info "Generating checksums..."
cd "${BUILD_DIR}"
sha256sum "tmux-helper-${VERSION}" > "tmux-helper-${VERSION}.sha256"
cd ..

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
success "Release ${VERSION} built successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Files in ${BUILD_DIR}/:"
ls -lh "${BUILD_DIR}"
echo ""
echo "To create a GitHub release:"
echo "  1. Create a tag: git tag v${VERSION}"
echo "  2. Push: git push origin v${VERSION}"
echo "  3. Upload files from ${BUILD_DIR}/"
echo ""
