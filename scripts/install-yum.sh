#!/bin/bash

# ShellDock YUM/DNF Repository Installation Script
# This script adds the ShellDock repository to your system's yum/dnf sources

set -e

REPO_URL="https://raw.githubusercontent.com/OpsGuild/ShellDock/master/repo/rpm"
GITHUB_REPO="OpsGuild/ShellDock"

echo "🔧 Installing ShellDock YUM/DNF repository..."

# Detect distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO=$ID
    VERSION_ID=${VERSION_ID:-}
else
    echo "❌ Cannot detect Linux distribution"
    exit 1
fi

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

# Detect package manager
if command -v dnf &> /dev/null; then
    PKG_MANAGER="dnf"
elif command -v yum &> /dev/null; then
    PKG_MANAGER="yum"
else
    echo "❌ Neither yum nor dnf found"
    exit 1
fi

# Install required dependencies
echo "📦 Installing required packages..."
$PKG_MANAGER install -y curl

# Add repository
echo "➕ Adding ShellDock repository..."
cat > /etc/yum.repos.d/shelldock.repo << EOF
[shelldock]
name=ShellDock Repository
baseurl=$REPO_URL
enabled=1
gpgcheck=0
EOF

# Update package list
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



