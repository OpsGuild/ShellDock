# Installation

> **[← Back to README](../README.md)** · [Quick Start](QUICKSTART.md) · [Usage](USAGE.md) · [Command Reference](COMMAND_REFERENCE.md)

All the ways to install, update, and uninstall ShellDock.

---

## Quick Install (All Platforms)

**One-command installation for Linux and macOS:**

```bash
curl -fsSL https://shelldock.opsguild.tech/install.sh | bash
```

**Alternative (direct from GitHub):**

```bash
curl -sSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install.sh | sudo bash
```

This script automatically:
- Detects your OS and architecture
- Downloads the appropriate binary
- Installs ShellDock to `/usr/local/bin/shelldock`
- Makes it executable

**Example Output:**
```
🔍 Checking for latest release...
🚀 Installing ShellDock v1.0.0...
📥 Downloading from https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-linux-amd64...
📦 Installing to /usr/local/bin...
📦 Installing repository files...
✅ ShellDock installed successfully!

Run 'shelldock --help' to get started
Run 'shelldock manage' to open the interactive UI
```

## Platform-Specific Installation

### Debian/Ubuntu (apt)

**Option 1: One-Line Install Script (Easiest)**

```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install-apt.sh | sudo bash
```

This script automatically:
- Detects your architecture (amd64 or arm64)
- Downloads the latest .deb package
- Installs ShellDock
- Fixes any dependencies

**Option 2: Manual .deb Installation**

```bash
# Download the latest .deb package for your architecture
# For amd64:
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock_*_amd64.deb

# For arm64:
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock_*_arm64.deb

# Install
sudo dpkg -i shelldock_*_*.deb

# Fix dependencies if needed
sudo apt-get install -f
```

**Option 3: Direct Binary (No Package Manager)**

```bash
# Download binary
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-linux-amd64
# or for ARM64
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-linux-arm64

# Install
chmod +x shelldock-linux-*
sudo mv shelldock-linux-* /usr/local/bin/shelldock
```

### RedHat/CentOS/Fedora (yum/dnf)

**Option 1: Direct .rpm Package**

```bash
# Download the latest .rpm package
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-*-1.x86_64.rpm

# Install
sudo rpm -i shelldock-*-1.x86_64.rpm
# or for Fedora
sudo dnf install shelldock-*-1.x86_64.rpm
```

**Option 2: Install Script**

```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install-yum.sh | sudo bash
sudo yum install shelldock
# or
sudo dnf install shelldock
```

### Arch Linux (pacman)

**Option 1: AUR (Arch User Repository) — Recommended**

```bash
# Using yay
yay -S shelldock

# Using paru
paru -S shelldock

# Manual installation
git clone https://aur.archlinux.org/shelldock.git
cd shelldock
makepkg -si
```

**Option 2: Install Script**

```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install-arch.sh | bash
```

### macOS

**Option 1: Homebrew (Recommended)**

```bash
brew install OpsGuild/tap/shelldock
```

This automatically taps the `OpsGuild/tap` repository and installs ShellDock. The formula is maintained in a dedicated [Homebrew tap repository](https://github.com/OpsGuild/homebrew-tap) and is automatically updated with each release.

**Option 2: Direct Binary**

```bash
# For Intel Macs
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-amd64
chmod +x shelldock-darwin-amd64
sudo mv shelldock-darwin-amd64 /usr/local/bin/shelldock

# For Apple Silicon (M1/M2/M3)
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-arm64
chmod +x shelldock-darwin-arm64
sudo mv shelldock-darwin-arm64 /usr/local/bin/shelldock
```

### Windows

**Option 1: Chocolatey**

```powershell
choco install shelldock
```

**Option 2: Direct Binary**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-windows-amd64.exe" -OutFile "shelldock.exe"

# Add to PATH or move to a directory in your PATH
New-Item -ItemType Directory -Path "C:\Program Files\shelldock" -Force
Move-Item shelldock.exe "C:\Program Files\shelldock\"
```

### Snap (Linux)

```bash
sudo snap install shelldock
```

### Flatpak (Linux)

```bash
# Install from Flathub (if published)
flatpak install flathub com.github.opsguild.shelldock

# Or install from downloaded .flatpak file
flatpak install shelldock-*.flatpak
```

### Local Installation (Development)

```bash
go build -o shelldock .
sudo cp shelldock /usr/local/bin/
```

## Package Manager Installation (Detailed)

### Debian/Ubuntu (apt) — Using APT Repository

> **Note:** The APT repository is included in each release. For production use, consider hosting it on GitHub Pages or your own server for better performance.

```bash
# 1. Get latest version from GitHub API
LATEST_VERSION=$(curl -s https://api.github.com/repos/OpsGuild/ShellDock/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "v1.0.0")

# 2. Download and add GPG key from the latest release
curl -fsSL https://github.com/OpsGuild/ShellDock/releases/download/${LATEST_VERSION}/gpg.key | sudo gpg --dearmor -o /usr/share/keyrings/shelldock-archive-keyring.gpg

# 3. Add repository
echo "deb [signed-by=/usr/share/keyrings/shelldock-archive-keyring.gpg] https://github.com/OpsGuild/ShellDock/releases/download/${LATEST_VERSION}/apt-repo stable main" | sudo tee /etc/apt/sources.list.d/shelldock.list

# 4. Update and install
sudo apt update
sudo apt install shelldock
```

**Alternative: Using GitHub Pages (if configured):**

```bash
curl -fsSL https://opsguild.github.io/ShellDock/gpg.key | sudo gpg --dearmor -o /usr/share/keyrings/shelldock-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/shelldock-archive-keyring.gpg] https://opsguild.github.io/ShellDock/apt-repo stable main" | sudo tee /etc/apt/sources.list.d/shelldock.list
sudo apt update
sudo apt install shelldock
```

### RedHat/CentOS/Fedora — Using .rpm Package

```bash
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-*-1.x86_64.rpm
sudo rpm -i shelldock-*-1.x86_64.rpm
# or for Fedora
sudo dnf install shelldock-*-1.x86_64.rpm
```

### Arch Linux — Using AUR

```bash
# Using yay (recommended)
yay -S shelldock

# Using paru
paru -S shelldock

# Manual installation
git clone https://aur.archlinux.org/shelldock.git
cd shelldock
makepkg -si
```

### macOS — Using Homebrew

```bash
brew install OpsGuild/tap/shelldock
```

### Windows — Using Chocolatey

```powershell
choco install shelldock
```

## Updating ShellDock

The update method depends on how you installed ShellDock.

### Quick Install Script

Re-run the installation script:

```bash
curl -fsSL https://shelldock.opsguild.tech/install.sh | bash
```

### APT Repository

```bash
sudo apt update
sudo apt upgrade shelldock
```

### .deb Package

```bash
# For amd64
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock_*_amd64.deb
sudo dpkg -i shelldock_*_amd64.deb
sudo apt-get install -f

# For arm64
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock_*_arm64.deb
sudo dpkg -i shelldock_*_arm64.deb
sudo apt-get install -f
```

Or use the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install-apt.sh | sudo bash
```

### .rpm Package

```bash
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-*-1.x86_64.rpm
sudo rpm -Uvh shelldock-*-1.x86_64.rpm
# or for Fedora
sudo dnf upgrade shelldock-*-1.x86_64.rpm
```

### AUR (Arch Linux)

```bash
# Using yay
yay -Syu shelldock

# Using paru
paru -Syu shelldock

# Manual update
cd ~/shelldock
git pull
makepkg -si
```

### Homebrew (macOS)

```bash
brew upgrade shelldock
```

### Snap

```bash
sudo snap refresh shelldock
```

### Flatpak

```bash
flatpak update com.github.opsguild.shelldock
```

### Chocolatey (Windows)

```powershell
choco upgrade shelldock
```

### Direct Binary

Re-download and replace the binary:

**Linux:**
```bash
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-linux-amd64
chmod +x shelldock-linux-*
sudo mv shelldock-linux-* /usr/local/bin/shelldock
```

**macOS:**
```bash
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-amd64
# or for Apple Silicon
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-arm64
chmod +x shelldock-darwin-*
sudo mv shelldock-darwin-* /usr/local/bin/shelldock
```

**Windows:**
```powershell
Invoke-WebRequest -Uri "https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-windows-amd64.exe" -OutFile "shelldock.exe"
Move-Item -Force shelldock.exe "C:\Program Files\shelldock\shelldock.exe"
```

### Check Current Version

```bash
shelldock --version
```

## Uninstalling ShellDock

### Quick Install Script

```bash
sudo rm -f /usr/local/bin/shelldock
rm -rf ~/.shelldock
```

### APT Repository

```bash
sudo apt remove shelldock
sudo apt purge shelldock
```

### .deb Package

```bash
sudo dpkg -r shelldock
# Or to also remove configuration files:
sudo dpkg --purge shelldock
```

### .rpm Package

```bash
sudo rpm -e shelldock
# Or for Fedora/DNF
sudo dnf remove shelldock
```

### AUR (Arch Linux)

```bash
# Using yay
yay -Rns shelldock

# Using paru
paru -Rns shelldock

# Using pacman
sudo pacman -Rns shelldock
```

### Homebrew (macOS)

```bash
brew uninstall shelldock

# To also remove the tap:
brew untap OpsGuild/tap
```

### Snap

```bash
sudo snap remove shelldock
```

### Flatpak

```bash
flatpak uninstall com.github.opsguild.shelldock
```

### Chocolatey (Windows)

```powershell
choco uninstall shelldock
```

### Direct Binary

**Linux/macOS:**
```bash
sudo rm -f /usr/local/bin/shelldock
rm -rf ~/.shelldock
```

**Windows:**
```powershell
Remove-Item "C:\Program Files\shelldock\shelldock.exe" -Force
Remove-Item "C:\Program Files\shelldock" -Force -ErrorAction SilentlyContinue
Remove-Item "$env:USERPROFILE\.shelldock" -Recurse -Force -ErrorAction SilentlyContinue
```

### Remove Configuration Files

After uninstalling, you may want to remove configuration files:

**Linux/macOS:**
```bash
rm -rf ~/.shelldock
```

**Windows:**
```powershell
Remove-Item "$env:USERPROFILE\.shelldock" -Recurse -Force
```

Configuration files are stored in `~/.shelldock/` (or `%USERPROFILE%\.shelldock\` on Windows) and contain:
- Platform settings (`.sdrc`)
- Local command sets (`*.yaml` files)

## Build from Source

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)

### Build Steps

1. **Clone the repository:**
   ```bash
   git clone https://github.com/OpsGuild/ShellDock.git
   cd shelldock
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Build the binary:**
   ```bash
   # Simple build
   go build -o shelldock .

   # Or use Makefile
   make build
   ```

4. **Install (optional):**
   ```bash
   # Manual installation
   sudo cp shelldock /usr/local/bin/

   # Or use Makefile
   sudo make install
   ```

### Create Packages

```bash
# Debian package
make deb

# RPM package
make rpm

# Arch package
make arch
```
