#!/bin/bash

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="shelldock"
GITHUB_REPO="OpsGuild/ShellDock"
GITHUB_API="https://api.github.com/repos/${GITHUB_REPO}"

echo "🔍 Checking for latest release..."
LATEST_VERSION=$(curl -s "${GITHUB_API}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ Could not determine latest version"
    echo "   Please install manually from: https://github.com/${GITHUB_REPO}/releases"
    exit 1
fi

VERSION_NUMBER=${LATEST_VERSION#v}
echo "🚀 Installing ShellDock ${LATEST_VERSION}..."

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
else
    echo "❌ Unsupported OS: $OSTYPE"
    exit 1
fi

ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
else
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
fi

DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/shelldock-${OS}-${ARCH}"

echo "📥 Downloading from ${DOWNLOAD_URL}..."

if command -v curl &> /dev/null; then
    curl -L -o /tmp/${BINARY_NAME} ${DOWNLOAD_URL}
elif command -v wget &> /dev/null; then
    wget -O /tmp/${BINARY_NAME} ${DOWNLOAD_URL}
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

chmod +x /tmp/${BINARY_NAME}

echo "📦 Installing to ${INSTALL_DIR}..."
sudo mv /tmp/${BINARY_NAME} ${INSTALL_DIR}/${BINARY_NAME}

echo "📦 Installing repository files..."
REPO_DIR="/usr/share/shelldock/repository"
sudo mkdir -p "${REPO_DIR}"

download_repo_files() {
    local ref="master"

    process_dir() {
        local dir_path="$1"
        local dir_api_url="${GITHUB_API}/contents/${dir_path}"

        local response
        if command -v curl &> /dev/null; then
            response=$(curl -sL "${dir_api_url}?ref=${ref}")
        else
            response=$(wget -q -O - "${dir_api_url}?ref=${ref}")
        fi

        if ! echo "${response}" | grep -q '"type"'; then
            return 1
        fi

        local paths
        paths=$(echo "${response}" | sed -n 's/.*"path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')

        local types
        types=$(echo "${response}" | sed -n 's/.*"type"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')

        local path_array=()
        local type_array=()

        while IFS= read -r line; do
            [ -n "${line}" ] && path_array+=("${line}")
        done <<< "${paths}"

        while IFS= read -r line; do
            [ -n "${line}" ] && type_array+=("${line}")
        done <<< "${types}"

        local i
        local max_len=${#path_array[@]}
        [ ${#type_array[@]} -lt ${max_len} ] && max_len=${#type_array[@]}

        for ((i=0; i<max_len; i++)); do
            local item_path="${path_array[$i]}"
            local item_type="${type_array[$i]}"

            [ -z "${item_path}" ] || [ -z "${item_type}" ] && continue

            [[ "${item_path}" == *"test.yaml" ]] && continue

            if [ "${item_type}" = "file" ] && [[ "${item_path}" == *.yaml ]]; then
                local rel_path="${item_path#repository/}"
                local file_dir
                file_dir=$(dirname "${rel_path}")
                local filename
                filename=$(basename "${rel_path}")

                if [ "${file_dir}" != "." ] && [ "${file_dir}" != "repository" ]; then
                    sudo mkdir -p "${REPO_DIR}/${file_dir}"
                fi

                local local_path
                if [ "${file_dir}" != "." ] && [ "${file_dir}" != "repository" ]; then
                    local_path="${REPO_DIR}/${file_dir}/${filename}"
                else
                    local_path="${REPO_DIR}/${filename}"
                fi

                local raw_url="https://raw.githubusercontent.com/${GITHUB_REPO}/master/${item_path}"
                echo "  📥 Downloading ${rel_path}..."

                if command -v curl &> /dev/null; then
                    sudo curl -sL -o "${local_path}" "${raw_url}" || echo "⚠️  Warning: Could not download ${item_path}"
                else
                    sudo wget -q -O "${local_path}" "${raw_url}" || echo "⚠️  Warning: Could not download ${item_path}"
                fi
            elif [ "${item_type}" = "dir" ]; then
                process_dir "${item_path}"
            fi
        done
    }

    process_dir "repository"
}

download_repo_files

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
