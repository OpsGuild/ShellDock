#!/bin/bash
# Generate .SRCINFO for AUR package
# Usage: ./scripts/generate-srcinfo.sh <version> <sha256>

set -e

VERSION="${1:-1.0.0}"
SHA256="${2:-SKIP}"

cd packaging/arch

echo "Generating .SRCINFO for version $VERSION..."

# Try Docker first for proper generation
if command -v docker &> /dev/null; then
    echo "Attempting to generate .SRCINFO using Docker..."
    docker run --rm -v "$(pwd):/pkg" -w /pkg archlinux:latest bash -c "
        pacman -Sy --noconfirm base-devel &>/dev/null
        useradd -m builder &>/dev/null
        chown -R builder:builder /pkg
        su builder -c 'makepkg --printsrcinfo > .SRCINFO'
    " 2>/dev/null && {
        echo "Generated .SRCINFO using Docker"
        cat .SRCINFO
        exit 0
    }
fi

# Fallback: Generate manually
echo "Generating .SRCINFO manually..."

cat > .SRCINFO << SRCINFO
pkgbase = shelldock
	pkgdesc = A fast, cross-platform shell command repository manager
	pkgver = ${VERSION}
	pkgrel = 1
	url = https://github.com/OpsGuild/ShellDock
	arch = x86_64
	arch = aarch64
	license = MIT
	source_x86_64 = https://github.com/OpsGuild/ShellDock/releases/download/v${VERSION}/shelldock-linux-amd64
	source_aarch64 = https://github.com/OpsGuild/ShellDock/releases/download/v${VERSION}/shelldock-linux-arm64
	sha256sums_x86_64 = ${SHA256}
	sha256sums_aarch64 = SKIP

pkgname = shelldock
SRCINFO

echo "Generated .SRCINFO:"
cat .SRCINFO

