# Usage Guide

> **[← Back to README](../README.md)** · [Quick Start](QUICKSTART.md) · [Installation](INSTALLATION.md) · [Command Reference](COMMAND_REFERENCE.md)

Detailed guide covering configuration, command execution, the management UI, versioning, dynamic arguments, and platform support.

---

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

## Basic Commands

### Get Help

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
  -l, --local     Only check local repository (skip bundled)
  -v, --version   version for shelldock
```

### List Available Command Sets

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

## Command Execution

### Direct Execution (Recommended)

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
     $ apt-get update

  2. Install Docker
     $ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh

  3. Start Docker service
     $ systemctl start docker

  4. Enable Docker on boot
     $ systemctl enable docker

Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o: y

🚀 Executing commands...

[1/4] Update package index (step 1)
$ apt-get update
Run this step? (y/N): y
✅ Success

[2/4] Install Docker (step 2)
$ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
Run this step? (y/N): y
✅ Success

...
```

The initial prompt accepts three options:
- **`a`** (all) — run every step without further prompts
- **`y`** (yes) — proceed step-by-step, prompting before each step with `y/N`
- **`N`** (no) — cancel execution

### Skip All Prompts

Execute all commands without any prompting:

```bash
shelldock docker -a
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
$ apt-get update
✅ Success

[2/4] Install Docker (step 2)
...
```

### Dynamic Arguments

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
     $ apt-get update && apt-get install -y git

  2. Verify Git installation
     $ git --version

  3. Configure Git (set your name and email)
     $ git config --global user.name "John Doe" && git config --global user.email "john@example.com"

Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o:
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
     $ apt-get update && apt-get install -y git

  2. Verify Git installation
     $ git --version

  3. Configure Git (set your name and email)
     $ git config --global user.name "{{name}}" && git config --global user.email "{{email}}"
     📝 Will prompt for: name, email

Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o: a

🚀 Executing commands...

[1/3] Install Git (step 1)
$ apt-get update && apt-get install -y git
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

Arguments can have default values. If a default is set and you don't provide a value, the default will be used automatically.

### Using the `run` Subcommand

```bash
shelldock run docker
```

This is equivalent to `shelldock docker` but more explicit.

### Prefer Local Repository

```bash
# Using flag
shelldock --local docker
shelldock -l docker

# Using subcommand
shelldock run --local docker
```

### Skip Specific Steps

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
     $ systemctl enable docker

Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o:
```

### Run Only Specific Steps

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
     $ apt-get update

  3. Start Docker service
     $ systemctl start docker

Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o:
```

You cannot use both `--skip` and `--only` flags together.

## Command Management

### Preview Commands (Show Without Executing)

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
     $ apt-get update
     📱 Other platforms:
        centos: yum check-update || true
        fedora: dnf check-update || true
        arch: pacman -Sy
        darwin: brew update

  2. Install Docker
     $ curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
     📱 Other platforms:
        centos: yum install -y docker
        fedora: dnf install -y docker
        arch: pacman -S docker
        darwin: brew install --cask docker

💡 To execute these commands, run: shelldock docker
```

### Echo Commands (Copyable Format)

Output commands in a plain, copyable format (no descriptions or formatting):

```bash
shelldock echo docker
```

**Example Output:**
```
apt-get update
curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
systemctl start docker
systemctl enable docker
```

**With flags:**
```bash
shelldock echo docker --skip 1,2
shelldock echo docker --only 1,3
shelldock echo docker --local
shelldock echo docker --version v2
shelldock echo docker@v2
shelldock echo certbot@certonly
```

This is useful for:
- Copying commands to run manually
- Piping commands to a shell: `shelldock echo docker | bash`
- Scripting and automation
- Generating scripts from command sets

### Interactive Management UI

```bash
shelldock manage
```

This opens a full-featured interactive terminal UI for managing your local command sets. You can create, view, edit, and delete command sets with complete control over every field — including platform-specific commands, dynamic arguments, skip-on-error settings, and versioning.

#### Unified Single-Page Form

When creating or editing a command set, everything is displayed on a **single scrollable page**. You never lose sight of your previous entries — metadata fields (Name, Description, Version), the step list, and expanded step details (including platforms and arguments) are all visible together. The view auto-scrolls to keep the active field in view.

**Example of the unified form view:**

```
 ✏  New Command Set

❯ Name           my-setup█
                  unique identifier, e.g. docker, my-setup
  Description   —
  Version       v1
  Tag           —

  ──────────────────────────────────────────────────────

  Steps  1 step(s)

  ▼ 1. Install Docker

     ❯ Description     Install Docker█
                        what this step does, e.g. "Install Docker"
       Command         —
       Skip on error   false

     ❯ Platforms  (2)
       Supported: ubuntu, debian, centos, rhel, ...

       ❯ platform: ubuntu█
                type a platform name, e.g. ubuntu
           command: —

       tab/enter next • shift+tab back • ctrl+n add • ctrl+d remove
       esc done with platforms → move to arguments

       Arguments  none

  ──────────────────────────────────────────────────────

  Editing step — use tab/shift+tab to move between fields
  esc closes current section • ctrl+s saves from step list
```

#### List View Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate command sets |
| `Enter` | View command set details |
| `n` | Create new command set |
| `e` | Edit selected command set |
| `d` | Delete selected command set |
| `q` or `Ctrl+C` | Quit |

#### Detail View Controls

| Key | Action |
|-----|--------|
| `↑/↓` | Scroll through steps |
| `e` | Edit this command set |
| `d` | Delete this command set |
| `Esc` | Back to list |

#### Form Controls (Create / Edit)

**Metadata fields** (Name, Description, Version, Tag):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance to next field |
| `Shift+Tab` | Go back to previous field |
| `Esc` | Cancel |

After Tag, you advance into the **step list**:

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate steps |
| `n` | Add a new step (opens it immediately) |
| `e` / `Enter` | Expand and edit the selected step |
| `d` | Remove the selected step |
| `Ctrl+S` | Save and exit |
| `Esc` | Go back to metadata |

When a step is **expanded** (editing its fields inline):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance through Description → Command → Skip on error → Platforms → Arguments |
| `Shift+Tab` | Go back to previous field |
| `Esc` | Close the step and return to the step list |

When editing **Platforms** (inline under the step):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance through platform key → command, then next platform entry |
| `Shift+Tab` | Go back to previous platform field |
| `Ctrl+N` | Add a new platform entry |
| `Ctrl+D` | Remove the current platform entry |
| `Esc` | Done with platforms, move to Arguments |

When editing **Arguments** (inline under the step):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance through name → prompt → default → required, then next argument |
| `Shift+Tab` | Go back to previous argument field |
| `Ctrl+N` | Add a new argument entry |
| `Ctrl+D` | Remove the current argument entry |
| `Esc` | Done with arguments, close the step |

#### What you can manage per step

- Description and command
- Skip-on-error flag (true/false)
- Platform-specific commands (e.g. ubuntu, centos, darwin, windows) — edited inline
- Dynamic arguments with name, prompt, default value, and required flag — edited inline

All sections stay visible on the same page. The cursor indicator (`❯`) and auto-scrolling keep you oriented as you move through the form.

## Versioning

### Latest vs Default

These are two separate fields on each version entry:

- **`latest: true`** — the entry returned when no version or tag is specified (e.g. `shelldock mysql`).
- **`default: true`** — the entry returned when a version number matches multiple entries (e.g. `shelldock mysql@v1` and both `mysql` and `mariadb` tags are `v1`).

Every version entry should have both fields. When multiple tags share the same version number, exactly one should be `default: true` for that version.

| Command | Resolution |
|---------|------------|
| `shelldock mysql` | Entry marked `latest: true` |
| `shelldock mysql@v1` | Entry marked `default: true` among v1 entries |
| `shelldock mysql@mariadb` | Exact tag match |

### List Available Versions

```bash
shelldock versions docker
```

**Example Output:**
```
Available versions for 'docker':

  - v1
  - v2
  * v3 (default, latest)

Use 'shelldock docker@<version>' or 'shelldock docker --version <version>' to run a specific version or tag
```

**Example Output with Tags (shared version):**
```
Available versions for 'mysql':

  * v1 @mysql (default, latest)
  - v1 @mariadb

Use 'shelldock mysql@<version>' or 'shelldock mysql --version <version>' to run a specific version or tag
```

### Run Specific Version or Tag

Using `@` syntax (supports both version and tag):

```bash
shelldock docker@v1
shelldock docker@v2
shelldock certbot@certonly
shelldock certbot@nginx
shelldock mysql@mariadb
```

Using `--version` flag (or `--ver` alias):

```bash
shelldock docker --version v1
shelldock run docker --version v2
shelldock certbot --version certonly
shelldock run certbot --ver nginx
```

### Show Specific Version or Tag

```bash
shelldock show docker@v1
shelldock show docker --version v2
shelldock show certbot@certonly
shelldock show certbot --ver nginx
```

## Platform Support

ShellDock automatically detects your platform and uses the appropriate commands. Commands can be defined for multiple platforms in a single command set.

**Example with Platform-Specific Commands:**

When you run `shelldock docker` on Ubuntu, it uses:
```bash
apt-get update
```

On CentOS, it automatically uses:
```bash
yum check-update || true
```

On Fedora:
```bash
dnf check-update || true
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
