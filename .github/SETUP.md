# Setup Guide for Automated Releases

## GPG Signing Setup

To enable GPG signing for packages:

1. **Generate a GPG key:**
   ```bash
   gpg --full-generate-key
   # Choose: RSA and RSA, 4096 bits, no expiration (or set expiration)
   ```

2. **Export the private key:**
   ```bash
   gpg --armor --export-secret-keys YOUR_KEY_ID > private.key
   ```

3. **Add to GitHub Secrets:**
   - Go to repository Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `GPG_PRIVATE_KEY`
   - Value: Contents of `private.key` (including `-----BEGIN PGP PRIVATE KEY BLOCK-----` and `-----END PGP PRIVATE KEY BLOCK-----`)

4. **Get your GPG key ID:**
   ```bash
   gpg --list-secret-keys --keyid-format LONG
   ```

5. **Add GPG public key to repository:**
   ```bash
   gpg --armor --export YOUR_KEY_ID > public.key
   ```
   - Upload `public.key` to repository root or add to README

## AUR Package Setup

To automatically update AUR packages:

1. **Create AUR package repository:**
   ```bash
   git clone ssh://aur@aur.archlinux.org/shelldock.git
   cd shelldock
   ```

2. **Add AUR SSH key to GitHub Secrets:**
   - Name: `AUR_SSH_KEY`
   - Value: Your SSH private key for AUR

3. **The workflow will automatically:**
   - Update PKGBUILD
   - Push to AUR (if AUR_SSH_KEY is set)

## Homebrew Tap Setup

1. **Create Homebrew tap repository:**
   - Create `homebrew-shelldock` repository under your GitHub org
   - The workflow generates the formula automatically

2. **To publish manually:**
   ```bash
   # Download the formula from workflow artifacts
   # Then in your homebrew-shelldock repo:
   cp shelldock.rb Formula/shelldock.rb
   git add Formula/shelldock.rb
   git commit -m "Update shelldock to v1.0.0"
   git push
   ```

## Chocolatey Setup

1. **Get Chocolatey API key:**
   - Sign up at https://chocolatey.org
   - Get your API key from account settings

2. **Add to GitHub Secrets:**
   - Name: `CHOCOLATEY_API_KEY`
   - Value: Your Chocolatey API key

3. **The workflow will automatically publish** (if API key is set)

## Snap Store Setup

1. **Register on Snapcraft:**
   - Go to https://snapcraft.io
   - Register your snap name: `shelldock`

2. **Add to GitHub Secrets:**
   - Name: `SNAP_STORE_CREDENTIALS`
   - Value: Your snapcraft login credentials (base64 encoded)

3. **The workflow will automatically publish** (if credentials are set)

## Flatpak Setup

1. **Create Flathub application:**
   - Follow Flathub submission process
   - Or host your own Flatpak repository

2. **The workflow builds the Flatpak** - you can manually publish it

## Usage

### Creating a Release

Simply push a tag:
```bash
git tag v1.0.0
git push origin v1.0.0
```

Or use GitHub UI:
1. Go to Releases → Draft a new release
2. Choose tag (or create new)
3. The workflow will run automatically

### Manual Trigger

1. Go to Actions → Release workflow
2. Click "Run workflow"
3. Enter version (e.g., `1.0.0`)
4. Click "Run workflow"

## What Gets Built

- ✅ Binaries for all platforms (Linux amd64/arm64, macOS amd64/arm64, Windows amd64)
- ✅ DEB packages (Debian/Ubuntu)
- ✅ RPM packages (RedHat/CentOS/Fedora)
- ✅ Snap packages
- ✅ Flatpak packages
- ✅ Chocolatey packages
- ✅ Homebrew formula
- ✅ AUR PKGBUILD
- ✅ GPG signatures (if configured)

All artifacts are uploaded to GitHub Releases automatically!

