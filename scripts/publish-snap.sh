#!/bin/bash
# Publish snap to Snap Store
# Usage: SNAPCRAFT_STORE_CREDENTIALS="..." ./scripts/publish-snap.sh

set -e

# Prevent command echoing to avoid exposing credentials
set +x

if [ -z "$SNAPCRAFT_STORE_CREDENTIALS" ]; then
    echo "SNAP_STORE_CREDENTIALS not set, skipping publish"
    echo ""
    echo "=== HOW TO SET UP SNAP STORE PUBLISHING ==="
    echo "1. First, register the snap name (one-time setup):"
    echo "   snapcraft login"
    echo "   snapcraft register shelldock"
    echo ""
    echo "2. Then, create credentials for CI:"
    echo "   snapcraft export-login --snaps=shelldock --acls=package_upload,package_release credentials.txt"
    echo ""
    echo "3. Copy the ENTIRE contents of credentials.txt to GitHub secret SNAP_STORE_CREDENTIALS"
    echo "   (Do NOT base64 encode it - use the raw content)"
    echo "============================================="
    exit 0
fi

# Verify authentication works
echo "Verifying Snap Store authentication..."
if ! snapcraft whoami 2>/dev/null; then
    echo "Authentication failed."
    echo ""
    echo "=== TROUBLESHOOTING ==="
    echo "Make sure SNAP_STORE_CREDENTIALS contains the raw output of:"
    echo "  snapcraft export-login --snaps=shelldock --acls=package_upload,package_release credentials.txt"
    echo ""
    echo "The credential should start with '[login.ubuntu.com]' and be multi-line."
    echo "Do NOT base64 encode it."
    echo "======================="
    exit 1
fi

echo "Authentication successful!"

# Find snap file
SNAP_FILE=$(ls *.snap 2>/dev/null | head -1)
if [ -z "$SNAP_FILE" ]; then
    echo "No snap file found to upload"
    exit 0
fi

echo "Uploading $SNAP_FILE to Snap Store..."

# Upload snap to edge channel first, then release to stable
if snapcraft upload "$SNAP_FILE" --release=edge,stable; then
    echo "Successfully published to Snap Store!"
else
    echo "Upload failed. This could be because:"
    echo "  1. The snap name 'shelldock' is not registered to your account"
    echo "  2. The credentials don't have the right permissions"
    echo ""
    echo "To register the snap name (if not done):"
    echo "  snapcraft login"
    echo "  snapcraft register shelldock"
    exit 1
fi

