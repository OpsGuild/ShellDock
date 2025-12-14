# Package Publishing Guide

This document explains how ShellDock packages are automatically built and published to different package managers when you push a tag.

## Overview

When you push a tag (e.g., `v1.0.0`), GitHub Actions automatically:
1. Builds binaries for all platforms (Linux, macOS, Windows)
2. Creates packages for different package managers (.deb, .rpm, .pkg.tar.zst)
3. Creates package repositories
4. Publishes to GitHub Releases
5. Makes packages available for installation via package managers

## Workflow

The workflow is defined in `.github/workflows/publish-packages.yml` and triggers on:
- Push of tags matching `v*` (e.g., `v1.0.0`)
- Manual workflow dispatch (for testing)

## Package Managers Supported

### 1. Debian/Ubuntu (apt)

**Repository Location:** `repo/deb/`

**Installation:**
```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install-apt.sh | sudo bash
sudo apt update
sudo apt install shelldock
```

**Repository Structure:**
```
repo/deb/
├── dists/
│   └── stable/
│       ├── Release
│       └── main/
│           └── binary-amd64/
│               └── Packages.gz
└── pool/
    └── main/
        └── s/
            └── shelldock/
                └── shelldock_VERSION_amd64.deb
```

### 2. RedHat/CentOS/Fedora (yum/dnf)

**Repository Location:** `repo/rpm/`

**Installation:**
```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install-yum.sh | sudo bash
sudo yum install shelldock
# or
sudo dnf install shelldock
```

**Repository Structure:**
```
repo/rpm/
├── repodata/
│   ├── repomd.xml
│   └── ...
└── shelldock-VERSION-1.x86_64.rpm
```

### 3. Arch Linux

**Installation:**
```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install-arch.sh | bash
```

**Note:** Arch Linux packages can be published to AUR (Arch User Repository) or a custom repository.

## Publishing a New Version

### Step 1: Update Version

Update the version in:
- `Makefile` (VERSION variable)
- `go.mod` (if needed)
- Any other version references

### Step 2: Commit and Tag

```bash
git add .
git commit -m "Release version 1.0.0"
git tag v1.0.0
git push origin master
git push origin v1.0.0
```

### Step 3: Wait for Workflow

The GitHub Actions workflow will:
1. Build all binaries
2. Create packages
3. Create repositories
4. Publish to GitHub Releases

### Step 4: Verify

Check the GitHub Actions tab to ensure the workflow completed successfully.

## Manual Publishing

You can also trigger the workflow manually:

1. Go to GitHub Actions
2. Select "Publish Packages" workflow
3. Click "Run workflow"
4. Enter the version (e.g., `1.0.0`)
5. Click "Run workflow"

## Repository Hosting

Currently, repositories are stored in the GitHub repository under `repo/` directory. For production use, you may want to:

1. **Host on GitHub Pages:**
   - Enable GitHub Pages
   - Point to `repo/` directory
   - Update repository URLs in installation scripts

2. **Use a CDN:**
   - Upload `repo/` to a CDN
   - Update repository URLs in installation scripts

3. **Use a Package Hosting Service:**
   - PackageCloud
   - Bintray (deprecated)
   - Custom server

## Installation Scripts

Installation scripts are located in `scripts/`:
- `install-apt.sh` - For Debian/Ubuntu
- `install-yum.sh` - For RedHat/CentOS/Fedora
- `install-arch.sh` - For Arch Linux

These scripts:
1. Detect the distribution
2. Install required dependencies
3. Add the repository
4. Update package lists
5. Provide installation instructions

## Troubleshooting

### Packages not appearing

1. Check GitHub Actions workflow status
2. Verify tag format (must start with `v`)
3. Check repository structure in `repo/` directory

### Installation fails

1. Ensure you're running as root (use `sudo`)
2. Check internet connectivity
3. Verify repository URLs are accessible
4. Check distribution compatibility

### Repository not updating

1. Clear package manager cache:
   - `apt`: `sudo apt clean && sudo apt update`
   - `yum/dnf`: `sudo yum clean all && sudo yum makecache`
2. Verify repository files are updated
3. Check repository metadata files

## Security Considerations

### GPG Signing

For production use, you should:
1. Generate a GPG key
2. Sign packages and repository metadata
3. Add GPG key to installation scripts
4. Configure package managers to verify signatures

### Example GPG Setup

```bash
# Generate key
gpg --full-generate-key

# Export public key
gpg --armor --export YOUR_KEY_ID > public.key

# Sign repository
gpg --detach-sign --armor repo/deb/dists/stable/Release
```

## Future Enhancements

- [ ] GPG signing for packages
- [ ] Support for more architectures (arm64, etc.)
- [ ] Automatic AUR package updates
- [ ] Homebrew formula for macOS
- [ ] Chocolatey package for Windows
- [ ] Snap package
- [ ] Flatpak package



