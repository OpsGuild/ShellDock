#!/bin/bash

set -e

REPO_URL="https://raw.githubusercontent.com/OpsGuild/ShellDock/master/repo/rpm"
GITHUB_REPO="OpsGuild/ShellDock"

echo "🔧 Installing ShellDock YUM/DNF repository..."

if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO=$ID
    VERSION_ID=${VERSION_ID:-}
else
    echo "❌ Cannot detect Linux distribution"
    exit 1
fi

if [ "$EUID" -ne 0 ]; then
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

if command -v dnf &> /dev/null; then
    PKG_MANAGER="dnf"
elif command -v yum &> /dev/null; then
    PKG_MANAGER="yum"
else
    echo "❌ Neither yum nor dnf found"
    exit 1
fi

echo "📦 Installing required packages..."
$PKG_MANAGER install -y curl

echo "➕ Adding ShellDock repository..."
cat > /etc/yum.repos.d/shelldock.repo << EOF
[shelldock]
name=ShellDock Repository
baseurl=$REPO_URL
enabled=1
gpgcheck=0
EOF

echo "🔄 Updating package list..."
$PKG_MANAGER makecache

echo ""
echo "✅ ShellDock repository added successfully!"
echo ""
echo "To install ShellDock, run:"
echo "  sudo $PKG_MANAGER install shelldock"
echo ""
echo "To update ShellDock in the future:"
echo "  sudo $PKG_MANAGER update shelldock"
echo ""
