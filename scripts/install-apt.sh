#!/bin/bash

# ShellDock APT Repository Installation Script
# This script adds the ShellDock repository to your system's apt sources

set -e

REPO_URL="https://raw.githubusercontent.com/OpsGuild/ShellDock/master/repo/deb"
GITHUB_REPO="OpsGuild/ShellDock"

echo "🔧 Installing ShellDock APT repository..."

# Detect distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO=$ID
    VERSION_CODENAME=${VERSION_CODENAME:-$VERSION_ID}
else
    echo "❌ Cannot detect Linux distribution"
    exit 1
fi

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

# Install required dependencies
echo "📦 Installing required packages..."
apt-get update
apt-get install -y curl gnupg2 ca-certificates apt-transport-https

# Add repository
echo "➕ Adding ShellDock repository..."
cat > /etc/apt/sources.list.d/shelldock.list << EOF
deb [trusted=yes] $REPO_URL stable main
EOF

# Update package list
echo "🔄 Updating package list..."
apt-get update

echo ""
echo "✅ ShellDock repository added successfully!"
echo ""
echo "To install ShellDock, run:"
echo "  sudo apt install shelldock"
echo ""
echo "To update ShellDock in the future:"
echo "  sudo apt update && sudo apt upgrade shelldock"
echo ""



