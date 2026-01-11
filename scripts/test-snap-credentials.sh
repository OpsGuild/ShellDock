#!/bin/bash
# Test Snap Store credentials locally
# Usage: ./scripts/test-snap-credentials.sh [path-to-credentials.txt]

set -e

CREDENTIALS_FILE="${1:-credentials.txt}"

if [ ! -f "$CREDENTIALS_FILE" ]; then
    echo "❌ Credentials file not found: $CREDENTIALS_FILE"
    echo ""
    echo "Usage: ./scripts/test-snap-credentials.sh [path-to-credentials.txt]"
    echo ""
    echo "If you have the credentials as a string, save them to a file first:"
    echo "  echo 'YOUR_CREDENTIALS_CONTENT' > credentials.txt"
    exit 1
fi

echo "🔍 Testing Snap Store credentials..."
echo ""

# Check if snapcraft is installed
if ! command -v snapcraft &> /dev/null; then
    echo "❌ snapcraft is not installed"
    echo "Install it with: sudo snap install snapcraft --classic"
    exit 1
fi

# Modern snapcraft uses base64-encoded JSON credentials
FIRST_LINE=$(head -n 1 "$CREDENTIALS_FILE" | tr -d '[:space:]')
echo "📄 Checking credentials file format..."

if echo "$FIRST_LINE" | grep -q "^eyJ"; then
    echo "✅ Credentials appear to be in the correct format (base64-encoded JSON)"
    echo "🔐 Testing credentials with environment variable..."
    # Export credentials to environment variable (snapcraft automatically reads this)
    export SNAPCRAFT_STORE_CREDENTIALS="$FIRST_LINE"
else
    echo "❌ Invalid credentials format"
    echo "Modern snapcraft expects base64-encoded JSON credentials (should start with 'eyJ')"
    echo ""
    echo "To get the correct format, run:"
    echo "  snapcraft export-login --snaps=shelldock --acls=package_upload,package_release credentials.txt"
    echo ""
    echo "The output should be a single line of base64-encoded JSON"
    exit 1
fi

# Test authentication
echo ""
echo "✅ Testing authentication with 'snapcraft whoami'..."
if snapcraft whoami 2>&1; then
    echo ""
    echo "✅ Authentication successful!"
    echo ""
    echo "Your credentials are valid and working."
    echo ""
    echo "To use in GitHub Actions:"
    echo "1. Copy the ENTIRE contents of $CREDENTIALS_FILE (base64-encoded JSON)"
    echo "2. Go to: https://github.com/OpsGuild/ShellDock/settings/secrets/actions"
    echo "3. Create/update secret: SNAP_STORE_CREDENTIALS"
    echo "4. Paste the content as-is (it's already base64-encoded JSON)"
    echo ""
else
    echo ""
    echo "❌ Authentication failed!"
    echo ""
    echo "Possible issues:"
    echo "1. Credentials file is corrupted or incomplete"
    echo "2. Credentials have expired"
    echo "3. Credentials don't have the right permissions"
    echo ""
    echo "To regenerate credentials:"
    echo "  snapcraft login"
    echo "  snapcraft export-login --snaps=shelldock --acls=package_upload,package_release credentials.txt"
    echo ""
    exit 1
fi

# Clean up
unset SNAPCRAFT_STORE_CREDENTIALS

echo "🧪 Test complete!"

