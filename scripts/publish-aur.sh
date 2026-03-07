#!/bin/bash
# Usage: AUR_SSH_KEY="..." ./scripts/publish-aur.sh <version>

set -e

VERSION="${1:-}"

if [ -z "$AUR_SSH_KEY" ]; then
    echo "AUR_SSH_KEY not set, skipping AUR publish"
    echo ""
    echo "=== HOW TO SET UP AUR PUBLISHING ==="
    echo "1. Create an AUR account at https://aur.archlinux.org/"
    echo "2. Generate an SSH key: ssh-keygen -t ed25519 -C 'aur@shelldock'"
    echo "3. Add the public key to your AUR account"
    echo "4. Add the private key to GitHub secret AUR_SSH_KEY"
    echo "5. Submit the package first manually or create the repo"
    echo "====================================="
    exit 0
fi

mkdir -p ~/.ssh
echo "$AUR_SSH_KEY" > ~/.ssh/aur_key
chmod 600 ~/.ssh/aur_key
ssh-keyscan aur.archlinux.org >> ~/.ssh/known_hosts 2>/dev/null

git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"

echo "Cloning AUR repository..."
GIT_SSH_COMMAND="ssh -i ~/.ssh/aur_key" git clone ssh://aur@aur.archlinux.org/shelldock.git aur-repo

cd aur-repo

if [ ! -f PKGBUILD ]; then
    echo "This appears to be a new AUR package repository"
fi

cp ../packaging/arch/PKGBUILD .
cp ../packaging/arch/.SRCINFO .

if ! grep -q "^pkgname = " .SRCINFO; then
    echo "ERROR: .SRCINFO is missing pkgname entry"
    cat .SRCINFO
    exit 1
fi

echo "=== PKGBUILD ==="
cat PKGBUILD
echo ""
echo "=== .SRCINFO ==="
cat .SRCINFO
echo ""

git add PKGBUILD .SRCINFO
if git diff --staged --quiet; then
    echo "No changes to commit"
else
    git commit -m "Update to v${VERSION}"
    echo "Pushing to AUR..."
    GIT_SSH_COMMAND="ssh -i ~/.ssh/aur_key" git push origin master
    echo "Successfully published to AUR!"
fi
