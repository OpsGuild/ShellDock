# Command Set Format

> **[← Back to README](../README.md)** · [Quick Start](QUICKSTART.md) · [Usage](USAGE.md) · [Command Reference](COMMAND_REFERENCE.md)

Full YAML specification for ShellDock command sets — including single-version, multi-version, version tags, platform-specific commands, and dynamic arguments.

---

Command sets are stored as YAML files. They support both single-version and multi-version formats.

## Single Version Format

```yaml
name: docker
description: Docker installation and setup commands
version: "v1"
commands:
  - description: Update package index
    command: apt-get update
  - description: Install Docker
    command: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
    skip_on_error: false
  - description: Start Docker service
    command: systemctl start docker
    skip_on_error: true
```

## Multi-Version Format

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
          centos: yum install -y docker
          fedora: dnf install -y docker
          arch: pacman -S docker
          darwin: brew install --cask docker
```

## Version Tags

You can add optional tags to versions for easier identification. Tags allow you to reference versions by meaningful names instead of version numbers.

```yaml
name: certbot
description: Certbot SSL certificate tool installation
versions:
  - version: "v1"
    tag: certonly
    description: Certbot standalone installation
    commands:
      - description: Install Certbot
        command: apt-get install -y certbot
  - version: "v2"
    tag: nginx
    latest: true
    description: Certbot with Nginx plugin installation
    commands:
      - description: Install Certbot with Nginx plugin
        command: apt-get install -y certbot python3-certbot-nginx
```

You can then run using either the version or tag:

```bash
shelldock certbot@v1
shelldock certbot@certonly
shelldock certbot@nginx
shelldock certbot --version certonly
```

## Platform-Specific Commands

Commands can specify different commands for different platforms:

```yaml
name: docker
description: Docker installation (platform-agnostic)
version: "v1"
commands:
  - description: Update package index
    platforms:
      ubuntu: apt-get update
      debian: apt-get update
      centos: yum check-update || true
      rhel: yum check-update || true
      fedora: dnf check-update || true
      arch: pacman -Sy
      darwin: brew update
      windows: choco upgrade chocolatey
  - description: Install Docker
    platforms:
      ubuntu: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh
      centos: yum install -y docker
      fedora: dnf install -y docker
      arch: pacman -S docker
      darwin: brew install --cask docker
      windows: choco install docker-desktop
```

**Supported platform values:**
- **Linux Distributions:** `ubuntu`, `debian`, `centos`, `rhel`, `fedora`, `arch`, `opensuse`, `alpine`, `amazon`, `oracle`
- **Generic Linux:** `linux` (used as fallback for any Linux distro not explicitly listed)
- **Other:** `darwin` (macOS), `windows`

**Platform resolution order:**

1. Exact match for the user's platform (e.g. `ubuntu`)
2. If the user is on a Linux distro with no exact match, falls back to `linux`
3. If no `platforms` map is defined, falls back to the generic `command` field
4. On Windows, if no Windows platform entry exists, the step is skipped (no fallback to generic `command`)

## Dynamic Arguments

Commands can accept dynamic arguments that are provided at runtime. This allows commands to be more flexible and reusable.

### Using Arguments in Commands

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

### Argument Definition Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Variable name used in `{{name}}` placeholders |
| `prompt` | no | Custom prompt question shown to user (defaults to "Enter {name}:") |
| `default` | no | Default value if argument not provided |
| `required` | no | Whether argument is required (default: false) |

### Providing Arguments

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

### Example with Default Values

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

### Notes on Arguments

- Arguments work with both `command` and `platforms` fields
- If a required argument is missing and not in a terminal, the command will fail
- Arguments are substituted as-is, so ensure proper quoting in your commands
- Multiple arguments can be provided: `--args key1=value1,key2=value2`
- Arguments persist across steps — if step 1 collects a value for `name`, step 3 can reuse it without re-prompting

## Field Reference

### Top-Level Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Unique identifier for the command set |
| `description` | no | Human-readable description |
| `version` | yes* | Version string (e.g. "v1", "v2") — required for single-version format |
| `versions` | yes* | Array of version objects — required for multi-version format |
| `commands` | yes* | Array of command objects — required for single-version format |

\* Either `version` + `commands` (single-version) or `versions` (multi-version) must be provided.

### Version Object Fields (Multi-Version Format)

| Field | Required | Description |
|-------|----------|-------------|
| `version` | yes | Version string (e.g. "v1", "v2") |
| `tag` | no | Optional tag for easier identification (e.g. "certonly", "nginx") |
| `description` | no | Description for this specific version |
| `latest` | no | Boolean — marks this version as the default when no version is specified |
| `commands` | yes | Array of command objects for this version |

### Command Object Fields

| Field | Required | Description |
|-------|----------|-------------|
| `description` | yes | What this command does |
| `command` | no | Single shell command (generic / fallback) |
| `platforms` | no | Map of platform → command for platform-specific behavior |
| `skip_on_error` | no | Boolean — if true, execution continues even if this command fails (default: false) |
| `args` | no | Array of argument definitions for dynamic command arguments |

### Argument Definition Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Variable name used in `{{name}}` placeholders |
| `prompt` | no | Custom prompt question shown to user |
| `default` | no | Default value if user provides nothing |
| `required` | no | Boolean — if true, user must provide a value (default: false) |

## Examples

### Minimal Command Set

```yaml
name: hello
description: Hello world
version: "v1"
commands:
  - description: Say hello
    command: echo "Hello, World!"
```

### Full-Featured Command Set

```yaml
name: project-setup
description: Set up a new development project
version: "v1"
commands:
  - description: Create project directory
    command: mkdir -p ~/projects/{{project_name}}
    args:
      - name: project_name
        prompt: "Enter project name"
        required: true

  - description: Initialize git repository
    command: cd ~/projects/{{project_name}} && git init
    skip_on_error: true

  - description: Install system dependencies
    platforms:
      ubuntu: sudo apt-get install -y build-essential
      fedora: sudo dnf groupinstall -y "Development Tools"
      arch: sudo pacman -S --noconfirm base-devel
      darwin: xcode-select --install
    skip_on_error: true

  - description: Configure git identity
    command: git config --global user.name "{{author}}" && git config --global user.email "{{email}}"
    args:
      - name: author
        prompt: "Your name for git commits"
        required: true
      - name: email
        prompt: "Your email for git commits"
        required: true

  - description: Done
    command: echo "Project {{project_name}} is ready at ~/projects/{{project_name}}"
```
