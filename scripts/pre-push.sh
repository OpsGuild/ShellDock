#!/bin/bash
set -e

echo "🚀 Running pre-push checks..."

echo "🔍 Running golangci-lint..."
if ! golangci-lint run; then
    echo "❌ Linting failed. Push aborted."
    exit 1
fi
echo "✅ Linting passed."

echo "🧪 Running unit tests..."
if ! go test ./...; then
    echo "❌ Unit tests failed. Push aborted."
    exit 1
fi
echo "✅ Unit tests passed."

echo "🎉 All checks passed. Proceeding with push."
exit 0
