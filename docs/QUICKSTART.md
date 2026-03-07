# Quick Start Guide

> **[← Back to README](../README.md)** · [Installation](INSTALLATION.md) · [Usage](USAGE.md) · [Command Reference](COMMAND_REFERENCE.md)

Get up and running with ShellDock in under two minutes.

---

## 1. Install ShellDock

**One-line install (Linux & macOS):**

```bash
curl -fsSL https://shelldock.opsguild.tech/install.sh | bash
```

**Alternative (direct from GitHub):**

```bash
curl -sSL https://raw.githubusercontent.com/OpsGuild/ShellDock/master/scripts/install.sh | sudo bash
```

**Or build from source:**

```bash
git clone https://github.com/OpsGuild/ShellDock.git
cd ShellDock
make build
sudo make install
```

**Verify installation:**

```bash
shelldock --version
```

For platform-specific package managers (Homebrew, apt, yum, pacman, Chocolatey, Snap), see the [Installation Guide](INSTALLATION.md).

---

## 2. List Available Command Sets

```bash
shelldock list
```

**Example output:**

```
☁️  Bundled Repository:
  • certbot
  • docker
  • git
  • kubernetes
  • nginx
  • nodejs
  • python
  • rust

💾 Local Repository:
  (empty)
```

---

## 3. Preview a Command Set

```bash
shelldock show docker
```

This shows every step, including platform-specific commands, without executing anything.

---

## 4. Run a Command Set

```bash
shelldock docker
```

ShellDock previews all steps, then prompts:

```
Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o:
```

| Key | Behavior |
|-----|----------|
| `a` | Run every step without further prompts |
| `y` | Proceed step-by-step — confirm each step with `y/N` |
| `N` | Cancel |

To skip all prompts entirely:

```bash
shelldock docker -a
```

---

## 5. Create Your Own Command Set

### Option A: Interactive TUI

```bash
shelldock manage
```

Press `n` to create a new command set. Everything happens on a **single scrollable page** — you fill in the name, description, version, then add steps. When editing a step, its fields (description, command, skip-on-error, platforms, and arguments) expand inline so you can always see your previous entries. Press `Ctrl+S` from the step list to save.

**List view controls:**

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate |
| `Enter` | View details |
| `n` | New command set |
| `e` | Edit |
| `d` | Delete |
| `q` | Quit |

**Form controls (single-page):**

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance to next field |
| `Shift+Tab` | Go back to previous field |
| `n` | Add a new step (from step list) |
| `e` / `Enter` | Expand a step inline |
| `Ctrl+N` | Add platform or argument entry |
| `Ctrl+D` | Remove platform or argument entry |
| `Esc` | Close current section / go back |
| `Ctrl+S` | Save and exit (from step list) |

### Option B: Write YAML manually

Create a file in `~/.shelldock/`:

```bash
mkdir -p ~/.shelldock
cat > ~/.shelldock/my-setup.yaml << 'EOF'
name: my-setup
description: My custom setup commands
version: "v1"
commands:
  - description: Say hello
    command: echo "Hello from ShellDock!"
  - description: Show system info
    command: uname -a
    skip_on_error: true
EOF
```

Then run it:

```bash
shelldock my-setup
```

See the [Command Set Format Guide](COMMAND_SET_FORMAT.md) for the full YAML specification — including multi-version formats, platform-specific commands, dynamic arguments, and version tags.

---

## 6. Useful Flags

| Flag | Example | Description |
|------|---------|-------------|
| `-a, --yes` | `shelldock docker -a` | Run all steps without prompting |
| `-l, --local` | `shelldock -l my-setup` | Only search the local repository |
| `--skip` | `shelldock docker --skip 1,2` | Skip specific steps |
| `--only` | `shelldock docker --only 3-5` | Run only specific steps |
| `--ver` / `@` | `shelldock docker@v2` | Run a specific version or tag |
| `--args` | `shelldock git --args name="Jane"` | Provide dynamic arguments |

---

## 7. Set Your Platform

ShellDock auto-detects your OS, but you can override it:

```bash
# View current platform
shelldock config show

# Set manually
shelldock config set ubuntu
shelldock config set darwin
shelldock config set auto    # restore auto-detection
```

---

## 8. Copy Commands for Scripting

Output raw commands without descriptions or prompts:

```bash
shelldock echo docker
```

Pipe directly to a shell:

```bash
shelldock echo docker | bash
```

---

## What's Next?

| Topic | Link |
|-------|------|
| Full usage guide with examples | [Usage Guide](USAGE.md) |
| Every command and flag | [Command Reference](COMMAND_REFERENCE.md) |
| YAML format deep dive | [Command Set Format](COMMAND_SET_FORMAT.md) |
| All feature details | [Features](FEATURES.md) |
| Shell autocompletion | [Bash Completion](BASH_COMPLETION.md) |
| Contributing to ShellDock | [Contributing](CONTRIBUTING.md) |