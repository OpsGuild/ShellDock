#!/bin/bash

# ShellDock Installation Script

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="shelldock"
GITHUB_REPO="OpsGuild/ShellDock"
GITHUB_API="https://api.github.com/repos/${GITHUB_REPO}"

# Get latest release version
echo "🔍 Checking for latest release..."
LATEST_VERSION=$(curl -s "${GITHUB_API}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ Could not determine latest version"
    echo "   Please install manually from: https://github.com/${GITHUB_REPO}/releases"
    exit 1
fi

VERSION_NUMBER=${LATEST_VERSION#v}  # Remove 'v' prefix if present
echo "🚀 Installing ShellDock ${LATEST_VERSION}..."

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

# Download URL
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/shelldock-${OS}-${ARCH}"

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

# Function to recursively download all .yaml files from repository directory
# Uses pure bash - no dependencies on Python, jq, or other tools
download_repo_files() {
    local ref="master"
    
    # Function to process a directory recursively
    process_dir() {
        local dir_path="$1"
        local dir_api_url="${GITHUB_API}/contents/${dir_path}"
        
        # Get directory contents
        local response
        if command -v curl &> /dev/null; then
            response=$(curl -sL "${dir_api_url}?ref=${ref}")
        else
            response=$(wget -q -O - "${dir_api_url}?ref=${ref}")
        fi
        
        if ! echo "${response}" | grep -q '"type"'; then
            return 1
        fi
        
        # Extract all path values (handles both single-line and multi-line JSON)
        local paths
        paths=$(echo "${response}" | sed -n 's/.*"path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
        
        # Extract all type values
        local types
        types=$(echo "${response}" | sed -n 's/.*"type"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
        
        # Convert to arrays (handles empty lines)
        local path_array=()
        local type_array=()
        
        while IFS= read -r line; do
            [ -n "${line}" ] && path_array+=("${line}")
        done <<< "${paths}"
        
        while IFS= read -r line; do
            [ -n "${line}" ] && type_array+=("${line}")
        done <<< "${types}"
        
        # Process entries (assuming path and type arrays are in same order)
        # This works because GitHub API returns them in consistent order
        local i
        local max_len=${#path_array[@]}
        [ ${#type_array[@]} -lt ${max_len} ] && max_len=${#type_array[@]}
        
        for ((i=0; i<max_len; i++)); do
            local item_path="${path_array[$i]}"
            local item_type="${type_array[$i]}"
            
            [ -z "${item_path}" ] || [ -z "${item_type}" ] && continue
            
            # Skip test files
            [[ "${item_path}" == *"test.yaml" ]] && continue
            
            if [ "${item_type}" = "file" ] && [[ "${item_path}" == *.yaml ]]; then
                # Calculate relative path from repository/
                local rel_path="${item_path#repository/}"
                local file_dir
                file_dir=$(dirname "${rel_path}")
                local filename
                filename=$(basename "${rel_path}")
                
                # Create subdirectory if needed
                if [ "${file_dir}" != "." ] && [ "${file_dir}" != "repository" ]; then
                    sudo mkdir -p "${REPO_DIR}/${file_dir}"
                fi
                
                # Determine local path
                local local_path
                if [ "${file_dir}" != "." ] && [ "${file_dir}" != "repository" ]; then
                    local_path="${REPO_DIR}/${file_dir}/${filename}"
                else
                    local_path="${REPO_DIR}/${filename}"
                fi
                
                # Download file
                local raw_url="https://raw.githubusercontent.com/${GITHUB_REPO}/master/${item_path}"
                echo "  📥 Downloading ${rel_path}..."
                
                if command -v curl &> /dev/null; then
                    sudo curl -sL -o "${local_path}" "${raw_url}" || echo "⚠️  Warning: Could not download ${item_path}"
                else
                    sudo wget -q -O "${local_path}" "${raw_url}" || echo "⚠️  Warning: Could not download ${item_path}"
                fi
            elif [ "${item_type}" = "dir" ]; then
                # Recursively process subdirectories
                process_dir "${item_path}"
            fi
        done
    }
    
    # Start processing from repository directory
    process_dir "repository"
}

# Download all repository files dynamically
download_repo_files

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



