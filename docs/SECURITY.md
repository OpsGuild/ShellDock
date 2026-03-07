# Security

ShellDock executes shell commands on your system. This document explains the trust model, how to verify what will run before it does, and best practices for safe usage.

---

## Trust Model

ShellDock is a **command repository manager**, not a sandboxed execution environment. Every command it runs executes directly in your shell with your current user's privileges. If you run ShellDock with `sudo`, the commands run as root.

There are two sources of commands:

| Source | Location | Trust Level |
|---|---|---|
| **Bundled repository** | Shipped inside the ShellDock binary/package under `repository/` | Maintained by OpsGuild, reviewed in the open-source repo |
| **Local repository** | `~/.shelldock/` on your machine | Created or synced by you — you control these entirely |

Local commands always take priority over bundled ones. If a local command set has the same name as a bundled one, the local version is used.

### Synced Commands

When you run `shelldock sync`, command sets are pulled from the configured remote repository (GitHub by default). This is conceptually identical to `git pull` — you are downloading files from a remote source. The same caution applies.

---

## Inspect Before You Execute

ShellDock provides multiple ways to see exactly what will run **before** anything touches your system:

### Preview with `show`

```bash
shelldock show docker
```

This prints every step, the resolved platform-specific command, available platforms, arguments, and version info. Nothing is executed.

### Copy with `echo`

```bash
shelldock echo docker
```

This outputs raw commands with no formatting — one per line. You can pipe this to a file or read it in a pager:

```bash
shelldock echo docker | less
```

### Step-by-step confirmation

When you run a command set without `-a` (auto-yes), ShellDock shows a full preview and asks:

```
Do you want to execute these commands? [a]ll/[y]es step-by-step/[N]o:
```

Choosing **`y`** lets you approve or skip each step individually. Choosing **`N`** aborts entirely.

### Read the YAML directly

Every command set is a plain YAML file. Bundled ones live in `repository/` inside the project. Local ones live in `~/.shelldock/`. Open them in any editor:

```bash
cat ~/.shelldock/my-commands.yaml
```

There is no obfuscation, compilation, or binary encoding. What you see in the YAML is what runs.

---

## Best Practices

### 1. Always preview unfamiliar command sets

If you just installed ShellDock or synced new commands, run `shelldock show <name>` before executing anything. Treat it like reading a script before piping `curl` output to `bash`.

### 2. Use step-by-step mode for destructive operations

Run without `-a` so you can inspect and confirm each step:

```bash
shelldock openssh
```

Reserve `-a` (auto-yes) for command sets you have already reviewed and trust.

### 3. Avoid running as root unless necessary

Most installation commands need `sudo`, but diagnostic commands (like `sysinfo`) do not. Run with the minimum privilege required. ShellDock itself does not escalate privileges — if a command needs `sudo`, it is written explicitly in the YAML.

### 4. Pin versions when automating

In scripts or CI, specify the exact version to avoid surprises if a command set is updated:

```bash
shelldock docker@v1 -a
```

### 5. Override bundled commands with local copies

If you want to customize a bundled command set, copy it to your local repository and edit it:

```bash
shelldock show docker          # review the bundled version
shelldock echo docker > /tmp/docker-review.txt  # save for comparison
```

Then create your own version in `~/.shelldock/docker.yaml`. Your local copy takes priority automatically.

### 6. Review after syncing

After running `shelldock sync`, review what changed before executing any synced command sets:

```bash
shelldock sync
shelldock list                 # see what's available
shelldock show <new-command>   # inspect before running
```

### 7. Audit the repository

The entire bundled command repository is open source at [github.com/OpsGuild/ShellDock](https://github.com/OpsGuild/ShellDock) under the `repository/` directory. Every command set is a readable YAML file. Pull requests that modify command sets are subject to code review like any other change.

---

## What ShellDock Does NOT Do

- **No sandboxing** — Commands run directly in your shell. There is no container, VM, or restricted execution environment.
- **No privilege escalation** — ShellDock never adds `sudo` on its own. If a command in a YAML file says `sudo`, that is what the author wrote.
- **No network access at runtime** — Executing a command set does not phone home. The only network operation is `shelldock sync`, which you invoke explicitly.
- **No hidden commands** — The YAML format is the single source of truth. There are no pre/post hooks, hidden scripts, or compiled payloads.
- **No telemetry** — ShellDock collects no usage data.

---

## Reporting Security Issues

If you find a security vulnerability in ShellDock itself (the Go binary, the sync mechanism, or a bundled command set that does something unsafe), please report it responsibly:

- **Email:** Open a private security advisory on the [GitHub repository](https://github.com/OpsGuild/ShellDock/security/advisories)
- **Do not** open a public issue for security vulnerabilities

We will acknowledge reports within 48 hours and provide a fix or mitigation as quickly as possible.

---

## Summary

ShellDock is transparent by design. Every command is stored in plain YAML, every execution is previewable with `show`, and every step can be confirmed individually. The security posture is simple: **read before you run**. The tools to do that are built in.