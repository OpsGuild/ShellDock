# ShellDock

A fast, cross-platform shell command repository manager. Save, organize, and execute shell commands from bundled repository or local directory with platform-specific support and versioning.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Quick Install](#quick-install)
  - [Updating ShellDock](#updating-shelldock)
  - [Build from Source](#build-from-source)
  - [Package Manager Installation](#package-manager-installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Basic Commands](#basic-commands)
  - [Command Execution](#command-execution)
  - [Command Management](#command-management)
  - [Versioning](#versioning)
  - [Platform Support](#platform-support)
- [Command Reference](#command-reference)
- [Command Set Format](#command-set-format)
- [Examples](#examples)
- [Repository Structure](#repository-structure)
- [Development](#development)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Features

- ⚡ **Fast** - Built with Go for optimal performance
- 📦 **Bundled Repository** - Pre-installed command sets included with installation
- 💾 **Local Repository** - Manage your own custom command sets in ~/.shelldock
- 🎨 **Interactive TUI** - Beautiful terminal UI for managing commands
- 📦 **Package Manager Support** - Install via apt, yum, pacman, and more
- 🔄 **Step-by-step Execution** - Run commands with confirmation prompts
- 🛡️ **Safe** - Review commands before execution
- 🏷️ **Versioning** - Support for multiple versions of command sets
- 🖥️ **Platform-Agnostic** - Automatic platform detection and distribution-specific commands
- ⏭️ **Selective Execution** - Skip or run only specific steps
- 📋 **Preview Mode** - View commands without executing
- 📤 **Echo Mode** - Output commands in copyable format for scripting
- 🔧 **Dynamic Arguments** - Pass arguments to commands via `--args` flag or interactive prompts

## Installation

### Quick Install (All Platforms)

**One-command installation for Linux and macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install.sh | bash
```

This script automatically:
- Detects your OS and architecture
- Downloads the appropriate binary
- Installs ShellDock to `/usr/local/bin/shelldock`
- Makes it executable

**Example Output:**
```
🚀 Installing ShellDock v1.0.0...
📥 Downloading from https://github.com/OpsGuild/ShellDock/releases/download/v1.0.0/shelldock-linux-amd64...
📦 Installing to /usr/local/bin...
✅ ShellDock installed successfully!

Run 'shelldock --help' to get started
Run 'shelldock manage' to open the interactive UI
```

### Platform-Specific Installation

Choose the installation method that best fits your system:

#### Debian/Ubuntu (apt)

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

#### RedHat/CentOS/Fedora (yum/dnf)

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

#### Arch Linux (pacman)

**Option 1: AUR (Arch User Repository) - Recommended**

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

#### macOS

**Option 1: Homebrew**

```bash
brew tap OpsGuild/tap
brew install shelldock
```

**Option 2: Direct Binary**

```bash
# Download binary
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-amd64
# or for Apple Silicon
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-arm64

# Make executable and install
chmod +x shelldock-darwin-*
sudo mv shelldock-darwin-* /usr/local/bin/shelldock
```

#### Windows

**Option 1: Chocolatey**

```powershell
choco install shelldock
```

**Option 2: Direct Binary**

```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-windows-amd64.exe" -OutFile "shelldock.exe"

# Add to PATH or move to a directory in your PATH
# Example: Move to C:\Program Files\shelldock\
New-Item -ItemType Directory -Path "C:\Program Files\shelldock" -Force
Move-Item shelldock.exe "C:\Program Files\shelldock\"
```

#### Local Installation (Development)

For installing from the local source code:

```bash
# From the project root directory
go build -o shelldock .
sudo cp shelldock /usr/local/bin/
```

This will build and install ShellDock from the current source code.

### Updating ShellDock

The update method depends on how you installed ShellDock:

#### If Installed via Quick Install Script

Simply re-run the installation script - it will download and install the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install.sh | bash
```

#### If Installed via APT Repository

```bash
sudo apt update
sudo apt upgrade shelldock
```

#### If Installed via .deb Package

Re-download and install the latest .deb package:

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

#### If Installed via .rpm Package

Re-download and install the latest .rpm package:

```bash
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-*-1.x86_64.rpm
sudo rpm -Uvh shelldock-*-1.x86_64.rpm
# or for Fedora
sudo dnf upgrade shelldock-*-1.x86_64.rpm
```

#### If Installed via AUR (Arch Linux)

```bash
# Using yay
yay -Syu shelldock

# Using paru
paru -Syu shelldock

# Manual update
cd ~/shelldock  # or wherever you cloned it
git pull
makepkg -si
```

#### If Installed via Homebrew (macOS)

```bash
brew upgrade shelldock
```

#### If Installed via Snap

```bash
sudo snap refresh shelldock
```

#### If Installed via Flatpak

```bash
flatpak update com.github.opsguild.shelldock
```

#### If Installed via Chocolatey (Windows)

```powershell
choco upgrade shelldock
```

#### If Installed via Direct Binary

Re-download and replace the binary:

**Linux:**
```bash
# Download latest binary
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-linux-amd64
# or for ARM64
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-linux-arm64

# Replace existing binary
chmod +x shelldock-linux-*
sudo mv shelldock-linux-* /usr/local/bin/shelldock
```

**macOS:**
```bash
# Download latest binary
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-amd64
# or for Apple Silicon
curl -LO https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-darwin-arm64

# Replace existing binary
chmod +x shelldock-darwin-*
sudo mv shelldock-darwin-* /usr/local/bin/shelldock
```

**Windows:**
```powershell
# Download latest binary
Invoke-WebRequest -Uri "https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-windows-amd64.exe" -OutFile "shelldock.exe"

# Replace existing binary (adjust path as needed)
Move-Item -Force shelldock.exe "C:\Program Files\shelldock\shelldock.exe"
```

#### Check Current Version

To check your current version:

```bash
shelldock --version
```

### Build from Source

#### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)

#### Build Steps

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

**Example Output:**
```
Building shelldock...
Build complete: build/shelldock
Installing shelldock...
Installation complete!
```

### Package Manager Installation

Detailed installation instructions for each package manager:

#### Debian/Ubuntu (apt)

**Using APT Repository (Recommended for automatic updates):**

> **Note:** The APT repository is included in each release. For production use, consider hosting it on GitHub Pages or your own server for better performance.

```bash
# 1. Download and add GPG key from the latest release
# Replace VERSION with the latest version (e.g., v1.0.0)
VERSION="v1.0.0"
curl -fsSL https://github.com/OpsGuild/ShellDock/releases/download/${VERSION}/gpg.key | sudo gpg --dearmor -o /usr/share/keyrings/shelldock-archive-keyring.gpg

# 2. Add repository
echo "deb [signed-by=/usr/share/keyrings/shelldock-archive-keyring.gpg] https://github.com/OpsGuild/ShellDock/releases/download/${VERSION}/apt-repo stable main" | sudo tee /etc/apt/sources.list.d/shelldock.list

# 3. Update and install
sudo apt update
sudo apt install shelldock
```

**Alternative: Using GitHub Pages (if configured):**

If you host the APT repository on GitHub Pages, use:

```bash
# Add GPG key
curl -fsSL https://opsguild.github.io/ShellDock/gpg.key | sudo gpg --dearmor -o /usr/share/keyrings/shelldock-archive-keyring.gpg

# Add repository
echo "deb [signed-by=/usr/share/keyrings/shelldock-archive-keyring.gpg] https://opsguild.github.io/ShellDock/apt-repo stable main" | sudo tee /etc/apt/sources.list.d/shelldock.list

# Install
sudo apt update
sudo apt install shelldock
```

**Using .deb Package:**

```bash
# Download
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock_*_amd64.deb

# Install
sudo dpkg -i shelldock_*_amd64.deb

# Fix dependencies if needed
sudo apt-get install -f
```

**Example Output:**
```
Selecting previously unselected package shelldock.
(Reading database ... 123456 files and directories currently installed.)
Preparing to unpack shelldock_1.0.0_amd64.deb ...
Unpacking shelldock (1.0.0) ...
Setting up shelldock (1.0.0) ...
ShellDock installed successfully!
```

#### RedHat/CentOS/Fedora (yum/dnf)

**Using .rpm Package:**

```bash
# Download
wget https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-*-1.x86_64.rpm

# Install
sudo rpm -i shelldock-*-1.x86_64.rpm
# or for Fedora
sudo dnf install shelldock-*-1.x86_64.rpm
```

**Example Output:**
```
Preparing...                          ################################# [100%]
Updating / installing...
   1:shelldock-1.0.0-1                ################################# [100%]
ShellDock installed successfully!
```

#### Arch Linux (pacman)

**Using AUR (Arch User Repository):**

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

**Example Output:**
```
==> Making package: shelldock 1.0.0-1
==> Checking runtime dependencies...
==> Installing package shelldock with pacman -U...
loading packages...
resolving dependencies...
looking for conflicting packages...
Packages (1) shelldock-1.0.0-1
Total Installed Size:  5.23 MiB
:: Proceed with installation? [Y/n]
```

#### macOS

**Using Homebrew:**

```bash
brew install shelldock/tap/shelldock
```

**Using Direct Binary:**

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

#### Windows

**Using Chocolatey:**

```powershell
choco install shelldock
```

**Using Direct Binary:**

```powershell
# Download
Invoke-WebRequest -Uri "https://github.com/OpsGuild/ShellDock/releases/latest/download/shelldock-windows-amd64.exe" -OutFile "shelldock.exe"

# Add to PATH (example)
New-Item -ItemType Directory -Path "C:\Program Files\shelldock" -Force
Move-Item shelldock.exe "C:\Program Files\shelldock\"
# Add C:\Program Files\shelldock to your system PATH
```

#### Snap (Linux)

```bash
sudo snap install shelldock
```

#### Flatpak (Linux)

```bash
# Install from Flathub (if published)
flatpak install flathub com.github.opsguild.shelldock

# Or install from downloaded .flatpak file
flatpak install shelldock-*.flatpak
```

## Configuration

ShellDock uses a configuration file at `~/.shelldock/.sdrc` to store platform settings.

### View Current Configuration

```bash
shelldock config show
```

**Example Output:**
```
ShellDock Configuration:
  Platform setting: auto
  Active platform: ubuntu
  Config file: ~/.shelldock/.sdrc
```

### Set Platform

```bash
# Auto-detect platform (recommended)
shelldock config set auto

# Set specific platform
shelldock config set ubuntu
shelldock config set centos
shelldock config set fedora
shelldock config set arch
shelldock config set darwin
shelldock config set windows
```

**Supported Platforms:**
- **Linux Distributions:** `ubuntu`, `debian`, `centos`, `rhel`, `fedora`, `arch`, `opensuse`, `alpine`, `amazon`, `oracle`
- **Other:** `darwin` (macOS), `windows`
- **Auto:** `auto` (automatically detects your platform)

**Example Output:**
```
Platform set to: auto
Auto-detected platform: ubuntu
```

## Usage

### Basic Commands

#### Get Help

```bash
shelldock --help
```

**Example Output:**
```
ShellDock is a fast, cross-platform tool for managing and executing
saved shell commands from bundled repository or local directory.

You can run a command set directly:
  shelldock docker
  shelldock --local docker

Or use subcommands:
  shelldock run docker
  shelldock list
  shelldock manage

Usage:
  shelldock [command-set-name] [flags]
  shelldock [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage ShellDock configuration
  help        Help about any command
  list        List available command sets
  manage      Manage command sets (interactive UI)
  run         Run a saved command set
  show        Show commands in a command set without executing
  sync        Sync command sets from bundled repository
  versions    List all available versions for a command set

Flags:
  -h, --help      help for shelldock
  -l, --local     Only check local repository (skip cloud
  -v, --version   version for shelldock
```

#### List Available Command Sets

```bash
shelldock list
```

**Example Output:**
```
☁️  Bundled Repository:
  • docker
  • nodejs
  • python

💾 Local Repository:
  • my-custom-setup
  • test
```

### Command Execution

#### Direct Execution (Recommended)

The simplest way to run a command set:

```bash
shelldock docker
```

**Example Output:**
```
📦 Command Set: docker
📝 Description: Docker installation and setup commands
🔢 Version: v1
🖥️  Platform: ubuntu
📋 Commands to execute:

  1. Update package index
     $ sudo apt-get update

  2. Install Docker
     $ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh

  3. Start Docker service
     $ sudo systemctl start docker

  4. Enable Docker on boot
     $ sudo systemctl enable docker

Do you want to execute these commands? (y/N): 
```

**Note:** The command will wait for you to type `y` or `n` and press Enter. To skip the prompt, use the `--yes` or `-y` flag.

#### Skip Confirmation Prompt

Execute commands without prompting:

```bash
shelldock docker --yes
shelldock docker -y
```

**Example Output:**
```
📦 Command Set: docker
📝 Description: Docker installation and setup commands
🔢 Version: v1
🖥️  Platform: ubuntu
📋 Commands to execute:
...

🚀 Executing commands...
[1/4] Update package index (step 1)
...
```

#### Dynamic Arguments

Some command sets accept dynamic arguments that can be provided via the `--args` flag or through interactive prompts.

**Using `--args` flag:**

```bash
shelldock git --args name="John Doe",email="john@example.com"
```

**Example Output:**
```
📦 Command Set: git
📝 Description: Git installation and basic configuration
🔢 Version: v1
🖥️  Platform: ubuntu
📋 Commands to execute:

  1. Install Git
     $ sudo apt-get update && sudo apt-get install -y git

  2. Verify Git installation
     $ git --version

  3. Configure Git (set your name and email)
     $ git config --global user.name "John Doe" && git config --global user.email "john@example.com"

Do you want to execute these commands? (y/N):
```

**Interactive prompting (when `--args` not provided):**

```bash
shelldock git
```

**Example Output:**
```
📦 Command Set: git
📝 Description: Git installation and basic configuration
🔢 Version: v1
🖥️  Platform: ubuntu
📋 Commands to execute:

  1. Install Git
     $ sudo apt-get update && sudo apt-get install -y git

  2. Verify Git installation
     $ git --version

  3. Configure Git (set your name and email)
     $ git config --global user.name "{{name}}" && git config --global user.email "{{email}}"
     📝 Will prompt for: name, email

Do you want to execute these commands? (y/N): y

🚀 Executing commands...

[1/3] Install Git (step 1)
$ sudo apt-get update && sudo apt-get install -y git
✅ Success

[2/3] Verify Git installation (step 2)
$ git --version
✅ Success

[3/3] Configure Git (set your name and email) (step 3)
Enter your Git user name: John Doe
Enter your Git email address: john@example.com
$ git config --global user.name "John Doe" && git config --global user.email "john@example.com"
✅ Success

🎉 All commands executed successfully!
```

**Note:** Arguments can have default values. If a default is set and you don't provide a value, the default will be used automatically.

#### Using the `run` Subcommand

```bash
shelldock run docker
```

This is equivalent to `shelldock docker` but more explicit.

#### Prefer Local Repository

```bash
# Using flag
shelldock --local docker
shelldock -l docker

# Using subcommand
shelldock run --local docker
```

#### Skip Specific Steps

Skip steps 1, 2, and 3:

```bash
shelldock docker --skip 1,2,3
```

Skip a range of steps:

```bash
shelldock docker --skip 1-3
```

**Example Output:**
```
📦 Command Set: docker
📝 Description: Docker installation and setup commands
🔢 Version: v1
🖥️  Platform: ubuntu
⏭️  Skipping steps: 1,2,3
📋 Commands to execute:

  4. Enable Docker on boot
     $ sudo systemctl enable docker

Do you want to execute these commands? (y/N):
```

#### Run Only Specific Steps

Run only steps 1, 3, and 5:

```bash
shelldock docker --only 1,3,5
```

Run a range of steps:

```bash
shelldock docker --only 1-3
```

**Example Output:**
```
📦 Command Set: docker
📝 Description: Docker installation and setup commands
🔢 Version: v1
🖥️  Platform: ubuntu
🎯 Running only steps: 1,3
📋 Commands to execute:

  1. Update package index
     $ sudo apt-get update

  3. Start Docker service
     $ sudo systemctl start docker

Do you want to execute these commands? (y/N):
```

**Note:** You cannot use both `--skip` and `--only` flags together.

### Command Management

#### Preview Commands (Show Without Executing)

```bash
shelldock show docker
```

**Example Output:**
```
📦 Command Set: docker
📝 Description: Docker installation and setup commands
🔢 Version: v1
🖥️  Platform: ubuntu
📋 Commands:

  1. Update package index
     $ sudo apt-get update
     📱 Other platforms:
        centos: sudo yum check-update || true
        fedora: sudo dnf check-update || true
        arch: sudo pacman -Sy
        darwin: brew update

  2. Install Docker
     $ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
     📱 Other platforms:
        centos: sudo yum install -y docker
        fedora: sudo dnf install -y docker
        arch: sudo pacman -S docker
        darwin: brew install --cask docker

💡 To execute these commands, run: shelldock docker
```

#### Echo Commands (Copyable Format)

Output commands in a plain, copyable format (no descriptions or formatting):

```bash
shelldock echo docker
```

**Example Output:**
```
sudo apt-get update
curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
sudo systemctl start docker
sudo systemctl enable docker
```

**With --skip flag:**
```bash
shelldock echo docker --skip 1,2
```

**Example Output:**
```
sudo systemctl start docker
sudo systemctl enable docker
```

**With --only flag:**
```bash
shelldock echo docker --only 1,3
```

**Example Output:**
```
sudo apt-get update
sudo systemctl start docker
```

**Other options:**
```bash
shelldock echo docker --local          # Only from local repository
shelldock echo docker --ver v2         # Specific version
shelldock echo docker@v2               # Version using @ syntax
```

This is useful for:
- Copying commands to run manually
- Piping commands to a shell: `shelldock echo docker | bash`
- Scripting and automation
- Generating scripts from command sets

**With Flags:**
```bash
# Skip specific steps
shelldock echo docker --skip 1,2

# Run only specific steps
shelldock echo docker --only 3,4

# Use specific version
shelldock echo docker --ver v1

# From local repository only
shelldock echo docker --local
```

**Example with --skip:**
```bash
$ shelldock echo docker --skip 1
curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
sudo systemctl start docker
sudo systemctl enable docker
```

**Example with --only:**
```bash
$ shelldock echo docker --only 2,3
curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
sudo systemctl start docker
```

#### Interactive Management UI

```bash
shelldock manage
```

This opens an interactive terminal UI where you can:
- **View** all command sets
- **Add** new command sets
- **Edit** existing commands
- **Delete** command sets

**UI Controls:**
- `↑/↓` or `j/k` - Navigate command sets
- `Enter` - View command set details
- `a` - Add new command set
- `e` - Edit selected command set
- `d` - Delete selected command set
- `q` or `Ctrl+C` - Quit

**Example UI:**
```
┌─────────────────────────────────────┐
│  ShellDock - Command Sets           │
├─────────────────────────────────────┤
│  ▶ docker                           │
│    nodejs                           │
│    python                           │
│    test                             │
│                                     │
│  Controls: ↑/↓ Navigate | Enter    │
│  View | a Add | d Delete | q Quit  │
└─────────────────────────────────────┘
```

### Versioning

#### List Available Versions

```bash
shelldock versions docker
```

**Example Output:**
```
Available versions for 'docker':

  - v1
  - v2
  * v3 (latest)

Use 'shelldock docker@<version>' or 'shelldock docker --ver <version>' to run a specific version
```

#### Run Specific Version

Using `@` syntax:

```bash
shelldock docker@v1
shelldock docker@v2
```

Using `--ver` flag:

```bash
shelldock docker --ver v1
shelldock run docker --ver v2
```

**Example Output:**
```
📦 Command Set: docker
📝 Description: Docker installation and setup commands
🔢 Version: v1
🖥️  Platform: ubuntu
📋 Commands to execute:

  1. Update package index
     $ sudo apt-get update
...
```

#### Show Specific Version

```bash
shelldock show docker@v1
shelldock show docker --ver v2
```

### Platform Support

ShellDock automatically detects your platform and uses the appropriate commands. Commands can be defined for multiple platforms in a single command set.

**Example with Platform-Specific Commands:**

When you run `shelldock docker` on Ubuntu, it uses:
```bash
sudo apt-get update
```

On CentOS, it automatically uses:
```bash
sudo yum check-update || true
```

On Fedora:
```bash
sudo dnf check-update || true
```

The platform is detected from:
1. Your `.sdrc` configuration file
2. Auto-detection from `/etc/os-release` (Linux)
3. System detection (macOS/Windows)

**Unsupported Platform Handling:**

If you set a platform that doesn't have commands defined for it, ShellDock will gracefully handle this:

- Commands without platform-specific support will be skipped during execution
- A warning message will show which platforms are available
- You can change your platform with `shelldock config set <platform>`

**Example:**
```bash
$ shelldock config set windows
$ shelldock show git

📦 Command Set: git
📝 Description: Git installation and setup
🔢 Version: v1
🖥️  Platform: windows
📋 Commands:

  1. Install Git
     ⚠️  No command available for platform 'windows'
     Available platforms: ubuntu, debian, centos, rhel, fedora, arch, darwin

  2. Verify Git installation
     $ git --version

⚠️  Warning: Some commands are not available for platform 'windows'
   Consider changing your platform with: shelldock config set <platform>
```

## Command Reference

### `shelldock [command-set-name]`

Run a command set directly.

**Flags:**
- `-l, --local` - Only check local repository (skip bundled repository)
- `--skip <steps>` - Skip specific steps (comma-separated or range)
- `--only <steps>` - Run only specific steps (comma-separated or range)
- `--ver <version>` - Run specific version (e.g., v1, v2)
- `-y, --yes` - Execute commands without prompting for confirmation
- `--args <key=value,...>` - Provide dynamic arguments (e.g., `--args name=John,email=john@example.com`)

**Examples:**
```bash
shelldock docker
shelldock docker --local
shelldock docker --skip 1,2,3
shelldock docker --only 1-3
shelldock docker@v1
shelldock docker --ver v1
shelldock docker --yes
shelldock docker -y
```

### `shelldock run [command-set-name]`

Explicitly run a command set. Same as direct execution but more explicit.

**Flags:** Same as direct execution

**Examples:**
```bash
shelldock run docker
shelldock run docker --skip 1,2
shelldock run docker --only 3,4,5
```

### `shelldock show [command-set-name]`

Preview commands without executing them.

**Flags:**
- `-l, --local` - Only check local repository
- `--ver <version>` - Show specific version

**Examples:**
```bash
shelldock show docker
shelldock show docker --local
shelldock show docker@v1
```

### `shelldock list`

List all available command sets from both cloud and local repositories.

**Example:**
```bash
shelldock list
```

### `shelldock manage`

Open interactive terminal UI for managing command sets.

**Example:**
```bash
shelldock manage
```

### `shelldock versions [command-set-name]`

List all available versions for a command set.

**Example:**
```bash
shelldock versions docker
```

### `shelldock config show`

Show current configuration.

**Example:**
```bash
shelldock config show
```

### `shelldock config set [platform]`

Set the platform for command execution.

**Examples:**
```bash
shelldock config set auto
shelldock config set ubuntu
shelldock config set centos
```

### `shelldock sync`

Sync command sets from bundled repository (placeholder for future implementation).

**Example:**
```bash
shelldock sync
```

## Command Set Format

Command sets are stored as YAML files. They support both single-version and multi-version formats.

### Single Version Format

```yaml
name: docker
description: Docker installation and setup commands
version: "v1"
commands:
  - description: Update package index
    command: sudo apt-get update
  - description: Install Docker
    command: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
    skip_on_error: false
  - description: Start Docker service
    command: sudo systemctl start docker
    skip_on_error: true
```

### Multi-Version Format

```yaml
name: docker
description: Docker installation and setup commands
versions:
  - version: "v1"
    description: Docker installation and setup commands
    commands:
      - description: Install Docker
        command: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
  - version: "v2"
    description: Docker installation with updated script
    commands:
      - description: Install Docker
        command: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
  - version: "v3"
    latest: true
    description: Docker installation with platform support
    commands:
      - description: Install Docker
        platforms:
          ubuntu: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
          centos: sudo yum install -y docker
          fedora: sudo dnf install -y docker
          arch: sudo pacman -S docker
          darwin: brew install --cask docker
```
<｜tool▁calls▁begin｜><｜tool▁call▁begin｜>
read_file

### Platform-Specific Commands

Commands can specify different commands for different platforms:

```yaml
name: docker
description: Docker installation (platform-agnostic)
version: "v1"
commands:
  - description: Update package index
    platforms:
      ubuntu: sudo apt-get update
      debian: sudo apt-get update
      centos: sudo yum check-update || true
      rhel: sudo yum check-update || true
      fedora: sudo dnf check-update || true
      arch: sudo pacman -Sy
      darwin: brew update
      windows: choco upgrade chocolatey
  - description: Install Docker
    platforms:
      ubuntu: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
      centos: sudo yum install -y docker
      fedora: sudo dnf install -y docker
      arch: sudo pacman -S docker
      darwin: brew install --cask docker
      windows: choco install docker-desktop
```

**Field Descriptions:**
- `name` - Unique identifier for the command set
- `description` - Human-readable description
- `version` - Version string (e.g., "v1", "v2")
- `commands` - Array of command objects
  - `description` - What this command does
  - `command` - Single command (backward compatible)
  - `platforms` - Map of platform -> command (preferred for multi-platform)
  - `skip_on_error` - Continue execution if this command fails (default: false)
  - `args` - Array of argument definitions for dynamic command arguments

### Dynamic Arguments

Commands can accept dynamic arguments that are provided at runtime. This allows commands to be more flexible and reusable.

**Using Arguments in Commands:**

Arguments are referenced in commands using `{{argName}}` syntax:

```yaml
commands:
  - description: Configure Git
    command: git config --global user.name "{{name}}" && git config --global user.email "{{email}}"
    args:
      - name: name
        prompt: "Enter your Git user name"
        required: true
      - name: email
        prompt: "Enter your Git email address"
        required: true
```

**Argument Definition Fields:**
- `name` - Variable name used in `{{name}}` placeholders (required)
- `prompt` - Custom prompt question shown to user (optional, defaults to "Enter {name}:")
- `default` - Default value if argument not provided (optional)
- `required` - Whether argument is required (default: false)

**Providing Arguments:**

You can provide arguments in two ways:

1. **Using `--args` flag** (no prompting):
   ```bash
   shelldock git --args name="John Doe",email="john@example.com"
   ```

2. **Interactive prompting** (if `--args` not provided):
   ```bash
   shelldock git
   # Will prompt:
   # Enter your Git user name: John Doe
   # Enter your Git email address: john@example.com
   ```

**Example with Default Values:**

```yaml
commands:
  - description: Install npm packages
    command: npm install -g {{packages}}
    args:
      - name: packages
        prompt: "Enter npm package names (space-separated)"
        default: "yarn pnpm"
        required: false
```

In this example, if no `--args` is provided and the user doesn't enter a value, it will use the default "yarn pnpm".

**Notes:**
- Arguments work with both `command` and `platforms` fields
- If a required argument is missing and not in a terminal, the command will fail
- Arguments are substituted as-is, so ensure proper quoting in your commands
- Multiple arguments can be provided: `--args key1=value1,key2=value2`

## Examples

### Example 1: Docker Installation

**Command Set File:** `~/.shelldock/docker.yaml`

```yaml
name: docker
description: Docker installation and setup
version: "v1"
commands:
  - description: Update package index
    platforms:
      ubuntu: sudo apt-get update
      centos: sudo yum check-update || true
      fedora: sudo dnf check-update || true
  - description: Install Docker
    platforms:
      ubuntu: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
      centos: sudo yum install -y docker
      fedora: sudo dnf install -y docker
  - description: Start Docker service
    command: sudo systemctl start docker
    skip_on_error: true
```

**Usage:**
```bash
shelldock docker
```

**Output:**
```
📦 Command Set: docker
📝 Description: Docker installation and setup
🔢 Version: v1
🖥️  Platform: ubuntu
📋 Commands to execute:

  1. Update package index
     $ sudo apt-get update

  2. Install Docker
     $ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh

  3. Start Docker service
     $ sudo systemctl start docker

Do you want to execute these commands? (y/N): y

🚀 Executing commands...

[1/3] Update package index (step 1)
$ sudo apt-get update
Hit:1 http://archive.ubuntu.com/ubuntu jammy InRelease
...
✅ Success

[2/3] Install Docker (step 2)
$ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
...
✅ Success

[3/3] Start Docker service (step 3)
$ sudo systemctl start docker
✅ Success

🎉 All commands executed successfully!
```

### Example 2: Node.js Setup with Version Selection

```bash
# List versions
shelldock versions nodejs

# Run specific version
shelldock nodejs@v1

# Skip installation step if already installed
shelldock nodejs --skip 1

# Execute without prompt
shelldock nodejs --yes
```

### Example 3: Selective Step Execution

```bash
# Only run steps 3 and 4
shelldock docker --only 3,4

# Skip steps 1-3, run the rest
shelldock docker --skip 1-3
```

## Repository Structure

### Local Repository

Location: `~/.shelldock/`

Structure:
```
~/.shelldock/
├── .sdrc                    # Configuration file
├── docker.yaml              # Command set file
├── nodejs.yaml
└── my-custom-setup.yaml
```

### Cloud Cache

Location: `~/.cache/shelldock/cloud/`

Cached command sets from the bundled repository.

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional)

### Setup

```bash
git clone https://github.com/OpsGuild/ShellDock.git
cd shelldock
go mod download
```

### Build

```bash
# Simple build
go build -o shelldock .

# Using Makefile
make build

# Build for all platforms
make build-all
```

### Test

```bash
# Run comprehensive test suite
./test-suite.sh

# Or use Makefile
make test

# See docs/MANUAL_TESTING.md for manual testing guide
```

**Test Suite Coverage:**
- ✅ Basic functionality (list, show, run)
- ✅ Versioning (v1, v2, v3, latest detection)
- ✅ Platform support (detection, configuration, platform-specific commands)
- ✅ Step filtering (--skip, --only, ranges)
- ✅ Flag combinations (version + skip, version + only, etc.)
- ✅ Error handling (non-existent sets, invalid formats, conflicts)
- ✅ Command execution (with --yes flag)
- ✅ Dynamic arguments (--args flag, interactive prompting)
- ✅ Edge cases (empty sets, platform-only commands, etc.)

**Example Output:**
```
🧪 ShellDock Comprehensive Test Suite
======================================

▶ Test 1: List command sets
----------------------------------------
✅ PASSED

...

==========================================
📊 Test Summary
==========================================
✅ Passed: 44
❌ Failed: 0
Total: 44

🎉 All tests passed!
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

### Project Structure

```
shelldock/
├── cmd/                 # CLI commands
├── internal/
│   ├── cli/            # Command-line interface
│   ├── config/         # Configuration management
│   ├── repo/           # Repository management
│   └── tui/            # Terminal UI
├── examples/           # Example command sets
├── packaging/          # Package configurations
├── scripts/            # Installation scripts
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Documentation

Additional documentation is available in the `docs/` folder:

- **[Testing Guide](docs/TESTING.md)** - How to run unit and integration tests
- **[Manual Testing Guide](docs/MANUAL_TESTING.md)** - Comprehensive guide for manual testing of all features
- **[Features](docs/FEATURES.md)** - Detailed feature list and descriptions
- **[Quick Start](docs/QUICKSTART.md)** - Quick start guide for new users
- **[Bash Completion](docs/BASH_COMPLETION.md)** - Install and use bash completion for ShellDock
- **[Test Input](docs/TEST_INPUT.md)** - Test input examples
- **[Test Results](docs/TEST_RESULTS.md)** - Test results and validation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Contribution Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Adding Command Sets

To contribute command sets to the bundled repository:

1. Create a well-documented YAML file
2. Include platform-specific commands where applicable
3. Test on multiple platforms
4. Submit via Pull Request

## License

MIT License - see LICENSE file for details

---

**Made with ❤️ for developers who love automation**
