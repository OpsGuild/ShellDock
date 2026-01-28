#!/bin/bash

# Pre-push hook to run lint and tests
# This script is automatically invoked by git before a push

echo "🚀 Running pre-push checks..."

# Run golangci-lint
echo "🔍 Running golangci-lint..."
if ! golangci-lint run; then
    echo "❌ Linting failed. Push aborted."
    exit 1
fi
echo "✅ Linting passed."

# Run unit tests
echo "🧪 Running unit tests..."
if ! go test ./...; then
    echo "❌ Unit tests failed. Push aborted."
    exit 1
fi
echo "✅ Unit tests passed."

echo "🎉 All checks passed. Proceeding with push."
exit 0
