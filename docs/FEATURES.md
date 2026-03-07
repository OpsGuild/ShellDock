# Features

> **[← Back to README](../README.md)** · [Quick Start](QUICKSTART.md) · [Usage](USAGE.md) · [Command Reference](COMMAND_REFERENCE.md)

A comprehensive overview of everything ShellDock can do.

---

## Table of Contents

- [Command Execution](#command-execution)
- [Interactive Confirmation](#interactive-confirmation)
- [Step Filtering](#step-filtering)
- [Versioning & Tags](#versioning--tags)
- [Platform-Specific Commands](#platform-specific-commands)
- [Dynamic Arguments](#dynamic-arguments)
- [Preview & Echo Modes](#preview--echo-modes)
- [Interactive Management TUI](#interactive-management-tui)
- [Repository System](#repository-system)
- [Configuration](#configuration)
- [Bash Completion](#bash-completion)
- [Package Manager Support](#package-manager-support)
- [Error Handling](#error-handling)

---

## Command Execution

ShellDock lets you run curated command sets with a single command. No more copy-pasting from wikis or digging through shell history.

```bash
# Run directly by name
shelldock docker

# Explicit run subcommand (equivalent)
shelldock run docker
```

Every execution starts with a preview of all steps so you know exactly what will happen before anything runs.

---

## Interactive Confirmation

All command executions show a full preview and prompt before anything is executed:

```
Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o:
```

| Key | Behavior |
|-----|----------|
| `a` | Run every step without further prompts |
| `y` | Step-by-step mode — confirm each step individually with `y/N` |
| `N` | Cancel execution entirely |

When running step-by-step, each step prompts separately:

```
[2/4] Install Docker (step 2)
$ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
Run this step? (y/N):
```

**Skip all prompts** for scripting and CI:

```bash
shelldock docker -a
```

---

## Step Filtering

Run only the steps you need, or skip the ones you don't.

### Skip Steps

```bash
shelldock docker --skip 1,2,3       # Skip steps 1, 2, and 3
shelldock docker --skip 1-3          # Skip a range
shelldock docker --skip 1,4-6        # Mix individual and ranges
```

### Run Only Specific Steps

```bash
shelldock docker --only 1,3,5        # Run only these steps
shelldock docker --only 1-3          # Run a range
shelldock docker --only 2            # Run a single step
```

> **Note:** You cannot use `--skip` and `--only` together — ShellDock will show an error if you try.

---

## Versioning & Tags

Command sets can contain multiple versions, letting you maintain different workflows under a single name.

### List Versions

```bash
shelldock versions docker
```

```
Available versions for 'docker':

  - v1
  - v2
  * v3 (latest)
```

### Run a Specific Version

```bash
# Using @ syntax
shelldock docker@v1

# Using flags
shelldock docker --version v2
shelldock docker --ver v2
```

### Version Tags

Versions can have human-friendly tags for easier identification:

```bash
shelldock versions certbot
```

```
Available versions for 'certbot':

  - v1 [certonly]
  - v2 [nginx] (latest)
```

```bash
# Reference by tag instead of version number
shelldock certbot@certonly
shelldock certbot@nginx
```

The `latest` flag in YAML determines which version runs when no version is specified.

---

## Platform-Specific Commands

ShellDock automatically detects your operating system and selects the right command for each step. A single command set can work across Ubuntu, CentOS, Fedora, Arch, macOS, and Windows.

**Example:** Running `shelldock docker` on different platforms:

| Platform | Command used |
|----------|-------------|
| Ubuntu | `apt-get update` |
| CentOS | `yum check-update \|\| true` |
| Fedora | `dnf check-update \|\| true` |
| Arch | `pacman -Sy` |
| macOS | `brew update` |
| Windows | `choco upgrade chocolatey` |

**Supported platforms:**

| Category | Platforms |
|----------|----------|
| Linux distributions | `ubuntu`, `debian`, `centos`, `rhel`, `fedora`, `arch`, `opensuse`, `alpine`, `amazon`, `oracle` |
| Generic Linux | `linux` (fallback for any unmatched Linux distro) |
| Other | `darwin` (macOS), `windows` |

**Platform resolution order:**

1. Exact match for the user's platform (e.g., `ubuntu`)
2. Generic `linux` fallback (for Linux distros without an exact match)
3. Generic `command` field (if no `platforms` map is defined)
4. Step skipped (if no matching command is found — Windows skips generic commands)

**Unsupported platform handling:**

If a step doesn't have a command for your platform, ShellDock warns you and shows which platforms are available — it won't silently fail.

---

## Dynamic Arguments

Command sets can accept runtime arguments, making them reusable with different values each time.

### Provide via Flag

```bash
shelldock git --args name="Jane Doe",email="jane@example.com"
```

### Interactive Prompting

If `--args` isn't provided, ShellDock prompts interactively:

```
[3/3] Configure Git (set your name and email) (step 3)
Enter your Git user name: Jane Doe
Enter your Git email address: jane@example.com
$ git config --global user.name "Jane Doe" && git config --global user.email "jane@example.com"
✅ Success
```

### Argument Features

- **Default values** — Arguments can have defaults that are used when no value is provided
- **Required flag** — Mark arguments as required to enforce input
- **Custom prompts** — Define user-friendly prompt text
- **Cross-step persistence** — An argument collected in step 1 is reused in later steps without re-prompting
- **Works everywhere** — Arguments are substituted in both `command` and `platforms` fields

---

## Preview & Echo Modes

### Preview (`show`)

View the full command set with descriptions, platform info, and formatting — without executing:

```bash
shelldock show docker
shelldock show docker@v2
```

Shows each step, its command for your platform, and commands available on other platforms.

### Echo (Copyable Output)

Output raw commands in a plain format — no descriptions, no formatting, no prompts:

```bash
shelldock echo docker
```

```
apt-get update
curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
systemctl start docker
systemctl enable docker
```

**Use cases:**

- Copy commands to run manually
- Pipe to a shell: `shelldock echo docker | bash`
- Generate scripts: `shelldock echo docker > setup.sh`
- CI/CD automation

Echo supports all filtering flags:

```bash
shelldock echo docker --skip 1,2
shelldock echo docker --only 3,4
shelldock echo docker@v2
shelldock echo certbot@certonly
```

---

## Interactive Management TUI

A full-featured terminal UI for managing local command sets without editing YAML by hand.

```bash
shelldock manage
```

### Unified Single-Page Form

When creating or editing a command set, everything is displayed on a **single scrollable page**. Metadata fields (Name, Description, Version), the step list, and expanded step details — including platforms and arguments — are all visible together. You never lose sight of your previous entries. The view auto-scrolls to keep the active field in view.

Steps are shown as a list below the metadata. When you expand a step (press `e` or `Enter`), its fields appear inline — including inline sub-editors for platforms and arguments. Collapsed steps show a compact summary (command preview, platform count, argument count). The `❯` cursor indicator and auto-scrolling keep you oriented.

### List View

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate command sets |
| `Enter` | View command set details |
| `n` | Create new command set |
| `e` | Edit selected command set |
| `d` | Delete selected command set |
| `q` or `Ctrl+C` | Quit |

### Detail View

| Key | Action |
|-----|--------|
| `↑/↓` | Scroll through steps |
| `e` | Edit this command set |
| `d` | Delete this command set |
| `Esc` | Back to list |

### Form Controls (Create / Edit)

**Metadata fields** (Name → Description → Version):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance to next field |
| `Shift+Tab` | Go back to previous field |
| `Esc` | Cancel |

**Step list** (after advancing past Version):

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate steps |
| `n` | Add a new step (opens it immediately) |
| `e` / `Enter` | Expand and edit the selected step inline |
| `d` | Remove the selected step |
| `Ctrl+S` | Save and exit |
| `Esc` | Go back to metadata |

**Step fields** (when a step is expanded):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance: Description → Command → Skip on error → Platforms → Arguments |
| `Shift+Tab` | Go back to previous field |
| `Esc` | Close the step, return to step list |

**Platforms** (inline sub-editor under the step):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance: platform key → command, then next entry |
| `Shift+Tab` | Go back to previous platform field |
| `Ctrl+N` | Add a new platform entry |
| `Ctrl+D` | Remove the current platform entry |
| `Esc` | Done with platforms, move to Arguments |

**Arguments** (inline sub-editor under the step):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance: name → prompt → default → required, then next entry |
| `Shift+Tab` | Go back to previous argument field |
| `Ctrl+N` | Add a new argument entry |
| `Ctrl+D` | Remove the current argument entry |
| `Esc` | Done with arguments, close the step |

**What you can manage per step:**

- Description and command
- Skip-on-error flag (true/false)
- Platform-specific commands (e.g., ubuntu, centos, darwin, windows) — edited inline
- Dynamic arguments with name, prompt, default value, and required flag — edited inline

All sections stay visible on the same scrollable page throughout the entire editing experience.

---

## Repository System

ShellDock uses a two-tier repository system:

### Local Repository (`~/.shelldock/`)

- Your personal command sets
- Always checked first
- Can override bundled command sets
- Managed via the TUI or by editing YAML files directly
- Great candidate for version control

### Bundled Repository (`/usr/share/shelldock/repository/`)

- Ships with ShellDock installation
- Curated command sets organized by category
- Updated via `shelldock sync` or reinstallation

```
repository/
├── devops/      # docker, kubernetes, pm2
├── editors/     # nvim
├── languages/   # go, nodejs, python, rust
├── network/     # networking tools
├── security/    # openssh, ufw
├── system/      # swap, sysinfo
├── vcs/         # git
└── web/         # nginx, certbot
```

**Subdirectories are transparent** — commands are accessed by name only (e.g., `shelldock docker`, not `shelldock devops/docker`).

### Repository Priority

1. **Default:** Local repository first, then bundled
2. **With `--local` / `-l`:** Only check local repository (skip bundled)

---

## Configuration

ShellDock stores configuration in `~/.shelldock/.sdrc`.

```bash
# View current configuration
shelldock config show

# Set platform explicitly
shelldock config set ubuntu

# Restore auto-detection
shelldock config set auto
```

```
ShellDock Configuration:
  Platform setting: auto
  Active platform: ubuntu
  Config file: ~/.shelldock/.sdrc
```

---

## Bash Completion

ShellDock supports shell autocompletion for bash, zsh, fish, and PowerShell.

```bash
# Temporary (current session)
source <(shelldock completion bash)

# Permanent
echo 'source <(shelldock completion bash)' >> ~/.bashrc
```

Completions cover commands, subcommands, command set names, flags, and argument hints.

See the [Bash Completion Guide](BASH_COMPLETION.md) for full setup instructions.

---

## Package Manager Support

ShellDock can be installed via multiple package managers:

| Method | Platform |
|--------|----------|
| Quick install script | Linux, macOS |
| Homebrew | macOS |
| APT / `.deb` | Debian, Ubuntu |
| YUM / DNF / `.rpm` | CentOS, RHEL, Fedora |
| AUR / pacman | Arch Linux |
| Chocolatey | Windows |
| Snap | Linux |
| Flatpak | Linux |
| Direct binary | All platforms |
| Build from source | All platforms |

See the [Installation Guide](INSTALLATION.md) for detailed instructions for each method.

---

## Error Handling

ShellDock provides clear, actionable error messages:

| Scenario | Behavior |
|----------|----------|
| Command set not found | Shows error with available sets |
| All steps filtered out | `No commands to execute after filtering` |
| `--skip` and `--only` combined | `Cannot use both --skip and --only flags together` |
| Invalid step numbers | Descriptive error with valid range |
| Non-existent version | Shows error with available versions |
| Missing required argument | Prompts interactively or fails with message |
| `skip_on_error: true` | Execution continues even if that step fails |
| Unsupported platform for a step | Step skipped with warning showing available platforms |

---

## Summary

| Feature | Description |
|---------|-------------|
| Direct execution | `shelldock <name>` runs a command set |
| Interactive confirmation | Preview and choose: all, step-by-step, or cancel |
| Step filtering | `--skip` and `--only` with ranges |
| Versioning | Multiple versions per command set with `@v2` syntax |
| Tags | Human-friendly names like `@nginx` or `@certonly` |
| Platform detection | Automatic OS detection with per-step overrides |
| Dynamic arguments | `--args` flag or interactive prompts |
| Preview mode | `show` to view without executing |
| Echo mode | `echo` for raw, pipeable output |
| Interactive TUI | Full management UI with `manage` |
| Two-tier repos | Local-first with bundled fallback |
| Shell completion | bash, zsh, fish, PowerShell |
| Skip on error | Graceful failure handling per step |
| Package managers | apt, yum, pacman, brew, choco, snap, and more |