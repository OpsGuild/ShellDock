#!/bin/bash

# ShellDock Installation Script

set -e

VERSION="1.0.0"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="shelldock"
GITHUB_REPO="OpsGuild/ShellDock"

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
DOWNLOAD_URL="https://github.com/OpsGuild/ShellDock/releases/download/v${VERSION}/shelldock-${OS}-${ARCH}"

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

# Install binary
echo "📦 Installing to ${INSTALL_DIR}..."
sudo mv /tmp/${BINARY_NAME} ${INSTALL_DIR}/${BINARY_NAME}

# Install repository files
echo "📦 Installing repository files..."
REPO_DIR="/usr/share/shelldock/repository"
sudo mkdir -p "${REPO_DIR}"

# Download repository files from GitHub
REPO_FILES=("docker.yaml" "git.yaml" "kubernetes.yaml" "nodejs.yaml" "python.yaml" "swap.yaml")
for file in "${REPO_FILES[@]}"; do
    REPO_URL="https://raw.githubusercontent.com/${GITHUB_REPO}/master/repository/${file}"
    if command -v curl &> /dev/null; then
        sudo curl -sL -o "${REPO_DIR}/${file}" "${REPO_URL}" || echo "⚠️  Warning: Could not download ${file}"
    elif command -v wget &> /dev/null; then
        sudo wget -q -O "${REPO_DIR}/${file}" "${REPO_URL}" || echo "⚠️  Warning: Could not download ${file}"
    fi
done

# Verify installation
if command -v ${BINARY_NAME} &> /dev/null; then
    echo "✅ ShellDock installed successfully!"
    echo ""
    echo "Run 'shelldock --help' to get started"
    echo "Run 'shelldock manage' to open the interactive UI"
    echo "Run 'shelldock list' to see available command sets"
else
    echo "❌ Installation failed. Please check your PATH."
    exit 1
fi



