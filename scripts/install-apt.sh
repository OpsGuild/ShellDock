#!/bin/bash

# ShellDock APT Repository Installation Script
# This script adds the ShellDock repository to your system's apt sources

set -euo pipefail

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

# Check if repository exists
echo "🔍 Checking repository availability..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$REPO_URL/dists/stable/Release" || echo "000")
if [[ ! "$HTTP_CODE" =~ ^(200|301|302)$ ]]; then
    echo ""
    echo "⚠️  Warning: ShellDock APT repository is not yet available."
    echo "   The repository at $REPO_URL does not exist yet."
    echo ""
    echo "📦 Alternative installation methods:"
    echo "   1. Build from source:"
    echo "      git clone https://github.com/$GITHUB_REPO.git"
    echo "      cd ShellDock"
    echo "      go build -o shelldock ."
    echo "      sudo cp shelldock /usr/local/bin/"
    echo ""
    echo "   2. Download binary from GitHub Releases (when available):"
    echo "      https://github.com/$GITHUB_REPO/releases"
    echo ""
    echo "   3. Wait for the repository to be published"
    echo ""
    exit 1
fi

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



