#!/bin/bash

# ShellDock Installation Script for Debian/Ubuntu
# This script downloads and installs the latest ShellDock .deb package

set -euo pipefail

GITHUB_REPO="OpsGuild/ShellDock"
GITHUB_API="https://api.github.com/repos/${GITHUB_REPO}"

echo "🚀 Installing ShellDock..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

# Remove any existing shelldock repository sources (if they exist)
if [ -f /etc/apt/sources.list.d/shelldock.list ]; then
    echo "🧹 Removing old repository configuration..."
    rm -f /etc/apt/sources.list.d/shelldock.list
fi

# Install required dependencies
echo "📦 Installing required packages..."
# Update package list, ignoring errors from bad repositories
apt-get update -qq -o Acquire::Check-Valid-Until=false 2>/dev/null || apt-get update -qq 2>/dev/null || true
apt-get install -y -qq curl wget

# Detect architecture
ARCH=$(dpkg --print-architecture)
case "$ARCH" in
    amd64)
        DEB_ARCH="amd64"
        ;;
    arm64)
        DEB_ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        echo "   Supported architectures: amd64, arm64"
        exit 1
        ;;
esac

# Get latest release version
echo "🔍 Checking for latest release..."
LATEST_VERSION=$(curl -s "${GITHUB_API}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ Could not determine latest version"
    echo "   Please install manually from: https://github.com/${GITHUB_REPO}/releases"
    exit 1
fi

VERSION_NUMBER=${LATEST_VERSION#v}  # Remove 'v' prefix if present
DEB_FILE="shelldock_${VERSION_NUMBER}_${DEB_ARCH}.deb"
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/${DEB_FILE}"

echo "📥 Downloading ShellDock ${LATEST_VERSION} (${DEB_ARCH})..."
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

if ! wget -q "$DOWNLOAD_URL"; then
    echo "❌ Failed to download ${DEB_FILE}"
    echo "   URL: ${DOWNLOAD_URL}"
    echo "   Please check: https://github.com/${GITHUB_REPO}/releases"
    exit 1
fi

# Install the package
echo "📦 Installing ShellDock..."
if dpkg -i "$DEB_FILE" 2>&1 | grep -q "dependency problems"; then
    echo "🔧 Fixing dependencies..."
    apt-get install -f -y -qq
fi

# Cleanup
cd /
rm -rf "$TMP_DIR"

echo ""
echo "✅ ShellDock installed successfully!"
echo ""
echo "Run 'shelldock --help' to get started"
echo "Run 'shelldock manage' to open the interactive UI"
echo ""



