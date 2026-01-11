# ShellDock Quick Start Guide

## Installation

### Build from Source

```bash
cd ShellDock
make build
sudo make install
```

Or manually:
```bash
go build -o shelldock .
sudo cp shelldock /usr/local/bin/
```

## Basic Usage

### 1. Run a Command Set

```bash
# Run from bundled or local repository
shelldock docker

# Explicitly use only local repository
shelldock --local docker
shelldock -l docker
```

### 2. List Available Commands

```bash
shelldock list
```

### 3. Manage Commands (Interactive UI)

```bash
shelldock manage
```

In the UI:
- **↑/↓** - Navigate command sets
- **Enter** - View command set details
- **a** - Add new command set
- **e** - Edit selected command set
- **d** - Delete selected command set
- **q** - Quit

### 4. Add a Command Set Manually

Create a YAML file in `~/.shelldock/`:

```bash
mkdir -p ~/.shelldock
cat > ~/.shelldock/my-commands.yaml << 'EOF'
name: my-commands
description: My custom commands
version: "1.0.0"
commands:
  - description: Example command
    command: echo "Hello World"
  - description: Another command
    command: ls -la
    skip_on_error: false
EOF
```

## Example: Adding Docker Setup Commands

1. Copy the example:
```bash
cp examples/docker.yaml ~/.shelldock/docker.yaml
```

2. Run it:
```bash
shelldock docker
```

## Package Manager Installation

### Debian/Ubuntu (.deb)

```bash
make deb
sudo dpkg -i dist/shelldock_1.0.0_amd64.deb
```

### RedHat/CentOS/Fedora (.rpm)

```bash
make rpm
sudo rpm -i dist/rpm/RPMS/x86_64/shelldock-1.0.0-1.x86_64.rpm
```

### Arch Linux

```bash
make arch
cd dist/arch
makepkg -si
```

## Repository Locations

- **Local Repository**: `~/.shelldock/` - Your custom command sets
- **Bundled Repository**: `/usr/share/shelldock/repository/` - Pre-installed command sets

### Bundled Repository Structure

The bundled repository is organized into subdirectories:

```
repository/
├── devops/      # docker, kubernetes, pm2
├── editors/     # nvim
├── languages/   # go, nodejs, python, rust
├── security/    # openssh, ufw
├── system/      # swap, sysinfo
├── vcs/         # git
└── web/         # nginx, certbot
```

**Note:** Subdirectories are transparent - just use `shelldock run docker`, not the full path.

## Command Set Format

Command sets are YAML files with this structure:

```yaml
name: command-set-name
description: Brief description
version: "1.0.0"
commands:
  - description: What this command does
    command: actual-shell-command
    skip_on_error: false  # optional, default false
```

## Tips

1. **Review Before Running**: Always review commands before executing them
2. **Use Local Repo for Personal Commands**: Keep your custom commands in `~/.shelldock/`
3. **Version Control**: You can version control your `~/.shelldock/` directory
4. **Skip on Error**: Use `skip_on_error: true` for optional commands that shouldn't fail the entire set

