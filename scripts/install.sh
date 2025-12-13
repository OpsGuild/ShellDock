#!/bin/bash

set -e

VERSION="1.0.0"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="shelldock"

echo "🚀 Installing ShellDock v${VERSION}..."

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
else
    echo "❌ Unsupported OS: $OSTYPE"
    exit 1
fi

# Detect architecture
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
else
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
fi

# Download URL (update this to your actual release URL)
DOWNLOAD_URL="https://github.com/shelldock/shelldock/releases/download/v${VERSION}/shelldock-${OS}-${ARCH}"

echo "📥 Downloading from ${DOWNLOAD_URL}..."

# Download binary
if command -v curl &> /dev/null; then
    curl -L -o /tmp/${BINARY_NAME} ${DOWNLOAD_URL}
elif command -v wget &> /dev/null; then
    wget -O /tmp/${BINARY_NAME} ${DOWNLOAD_URL}
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Make executable
chmod +x /tmp/${BINARY_NAME}

# Install
echo "📦 Installing to ${INSTALL_DIR}..."
sudo mv /tmp/${BINARY_NAME} ${INSTALL_DIR}/${BINARY_NAME}

# Verify installation
if command -v ${BINARY_NAME} &> /dev/null; then
    echo "✅ ShellDock installed successfully!"
    echo ""
    echo "Run 'shelldock --help' to get started"
    echo "Run 'shelldock manage' to open the interactive UI"
else
    echo "❌ Installation failed. Please check your PATH."
    exit 1
fi



