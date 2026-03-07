# Command Reference

> **[← Back to README](../README.md)** · [Quick Start](QUICKSTART.md) · [Installation](INSTALLATION.md) · [Usage](USAGE.md)

Every command, subcommand, and flag available in ShellDock — with examples.

---

## `shelldock [command-set-name]`

Run a command set directly.

**Flags:**
- `-l, --local` — Only check local repository (skip bundled repository)
- `--skip <steps>` — Skip specific steps (comma-separated or range)
- `--only <steps>` — Run only specific steps (comma-separated or range)
- `--version <version>` or `--ver <version>` — Run specific version or tag (e.g., v1, v2, certonly, nginx)
- `-a, --yes` — Execute all commands without prompting for confirmation
- `--args <key=value,...>` — Provide dynamic arguments (e.g., `--args name=John,email=john@example.com`)

**Examples:**
```bash
shelldock docker
shelldock docker --local
shelldock docker --skip 1,2,3
shelldock docker --only 1-3
shelldock docker@v1
shelldock docker --version v1
shelldock certbot@certonly
shelldock docker -a
```

## `shelldock run [command-set-name]`

Explicitly run a command set. Same as direct execution but more explicit.

**Flags:** Same as direct execution.

**Examples:**
```bash
shelldock run docker
shelldock run docker --skip 1,2
shelldock run docker --only 3,4,5
shelldock run docker -a
```

## `shelldock show [command-set-name]`

Preview commands without executing them.

**Flags:**
- `-l, --local` — Only check local repository
- `--version <version>` or `--ver <version>` — Show specific version or tag

**Examples:**
```bash
shelldock show docker
shelldock show docker --local
shelldock show docker@v1
shelldock show certbot@certonly
```

## `shelldock echo [command-set-name]`

Output commands in a plain, copyable format (no descriptions or formatting).

**Flags:**
- `-l, --local` — Only check local repository
- `--skip <steps>` — Skip specific steps
- `--only <steps>` — Run only specific steps
- `--version <version>` or `--ver <version>` — Specific version or tag

**Examples:**
```bash
shelldock echo docker
shelldock echo docker --skip 1,2
shelldock echo docker --only 3,4
shelldock echo docker --version v2
shelldock echo certbot@certonly
shelldock echo docker | bash
```

## `shelldock list`

List all available command sets from both bundled and local repositories.

**Example:**
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

## `shelldock manage`

Open interactive terminal UI for full management of local command sets (create, view, edit, delete).

When creating or editing, everything is displayed on a **single scrollable page** — metadata, steps, and expanded step details (platforms, arguments) are all visible together. You never lose sight of your previous entries.

**Example:**
```bash
shelldock manage
```

**List View Controls:**

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate command sets |
| `Enter` | View command set details |
| `n` | Create new command set |
| `e` | Edit selected command set |
| `d` | Delete selected command set |
| `q` or `Ctrl+C` | Quit |

**Detail View Controls:**

| Key | Action |
|-----|--------|
| `↑/↓` | Scroll through steps |
| `e` | Edit this command set |
| `d` | Delete this command set |
| `Esc` | Back to list |

**Form — Metadata Fields** (Name → Description → Version):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance to next field |
| `Shift+Tab` | Go back to previous field |
| `Esc` | Cancel |

**Form — Step List** (after advancing past Version):

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate steps |
| `n` | Add a new step (opens it immediately) |
| `e` / `Enter` | Expand and edit the selected step inline |
| `d` | Remove the selected step |
| `Ctrl+S` | Save and exit |
| `Esc` | Go back to metadata |

**Form — Step Fields** (when a step is expanded inline):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance: Description → Command → Skip on error → Platforms → Arguments |
| `Shift+Tab` | Go back to previous field |
| `Esc` | Close the step, return to step list |

**Form — Platforms** (inline sub-editor under the step):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance: platform key → command, then next entry |
| `Shift+Tab` | Go back to previous platform field |
| `Ctrl+N` | Add a new platform entry |
| `Ctrl+D` | Remove the current platform entry |
| `Esc` | Done with platforms, move to Arguments |

**Form — Arguments** (inline sub-editor under the step):

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Advance: name → prompt → default → required, then next entry |
| `Shift+Tab` | Go back to previous argument field |
| `Ctrl+N` | Add a new argument entry |
| `Ctrl+D` | Remove the current argument entry |
| `Esc` | Done with arguments, close the step |

## `shelldock versions [command-set-name]`

List all available versions for a command set.

**Examples:**
```bash
shelldock versions docker
shelldock versions certbot
```

**Example Output:**
```
Available versions for 'docker':

  - v1
  - v2
  * v3 (latest)

Use 'shelldock docker@<version>' or 'shelldock docker --version <version>' to run a specific version or tag
```

## `shelldock config show`

Show current configuration.

**Example:**
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

## `shelldock config set [platform]`

Set the platform for command execution.

**Examples:**
```bash
shelldock config set auto
shelldock config set ubuntu
shelldock config set centos
shelldock config set darwin
```

**Supported Platforms:**
- **Linux Distributions:** `ubuntu`, `debian`, `centos`, `rhel`, `fedora`, `arch`, `opensuse`, `alpine`, `amazon`, `oracle`
- **Other:** `darwin` (macOS), `windows`
- **Auto:** `auto` (automatically detects your platform)

## `shelldock sync`

Sync command sets from the cloud repository to the bundled repository.

Requires `sudo` if the bundled repository is in a system directory.

**Example:**
```bash
sudo shelldock sync
```
