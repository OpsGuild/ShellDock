#!/bin/bash

set -e

REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || echo ".")

if [ ! -d "$REPO_ROOT/.git" ]; then
    echo "❌ Error: Not a git repository."
    exit 1
fi

echo "🔧 Installing git hooks..."

SOURCE="$REPO_ROOT/scripts/pre-push.sh"
TARGET="$REPO_ROOT/.git/hooks/pre-push"

if [ ! -f "$SOURCE" ]; then
    echo "❌ Error: Source hook script not found at $SOURCE"
    exit 1
fi

chmod +x "$SOURCE"
ln -sf "$SOURCE" "$TARGET"

echo "✅ Pre-push hook installed at $TARGET"
echo "   It will run 'golangci-lint run' and 'go test ./...' before every push."
