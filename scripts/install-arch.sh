#!/bin/bash

# ShellDock Arch Linux Installation Script
# This script installs ShellDock from AUR or custom repository

set -e

echo "🔧 Installing ShellDock on Arch Linux..."

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    echo "❌ Please do not run as root (AUR packages should be built as user)"
    exit 1
fi

# Check if yay or paru is installed
if command -v yay &> /dev/null; then
    AUR_HELPER="yay"
elif command -v paru &> /dev/null; then
    AUR_HELPER="paru"
else
    echo "📦 Installing yay (AUR helper)..."
    # Install yay
    cd /tmp
    git clone https://aur.archlinux.org/yay.git
    cd yay
    makepkg -si --noconfirm
    AUR_HELPER="yay"
fi

echo "📦 Installing ShellDock from AUR..."
$AUR_HELPER -S shelldock --noconfirm

echo ""
echo "✅ ShellDock installed successfully!"
echo ""
echo "To update ShellDock in the future:"
echo "  $AUR_HELPER -Syu shelldock"
echo ""



