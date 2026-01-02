# Bash Completion for ShellDock

ShellDock supports bash completion for improved command-line experience.

## Installation

### Option 1: Temporary (Current Session Only)

```bash
source <(shelldock completion bash)
```

### Option 2: Permanent (User-specific)

Add to your `~/.bashrc` or `~/.bash_profile`:

```bash
# If shelldock is in your PATH
echo 'source <(shelldock completion bash)' >> ~/.bashrc

# Or if you need to specify the full path
echo 'source <(/usr/local/bin/shelldock completion bash)' >> ~/.bashrc
```

Then reload your shell:
```bash
source ~/.bashrc
```

### Option 3: System-wide Installation

```bash
# Generate completion script
shelldock completion bash > /etc/bash_completion.d/shelldock

# Or if using bash-completion package
shelldock completion bash > /usr/share/bash-completion/completions/shelldock

# Reload (if needed)
source /etc/bash_completion.d/shelldock
```

## Usage

After installation, you can use tab completion:

```bash
# Complete commands
shelldock <TAB>
# Shows: run, show, list, manage, sync, versions, config, echo, completion, help

# Complete command sets
shelldock run <TAB>
# Shows available command sets: docker, git, nodejs, python, kubernetes, etc.

# Complete flags
shelldock git --<TAB>
# Shows: --args, --local, --skip, --only, --ver, --yes, --help

# Complete subcommands
shelldock config <TAB>
# Shows: show, set

# Complete arguments (for --args flag)
shelldock git --args <TAB>
# May show argument names if supported
```

## Features

- Command completion
- Command set name completion
- Flag completion
- Subcommand completion
- Argument name hints (where applicable)

## Troubleshooting

If completion doesn't work:

1. **Check if bash-completion is installed:**
   ```bash
   which bash_completion
   ```

2. **Verify the completion script is loaded:**
   ```bash
   type __start_shelldock
   ```

3. **Check for errors:**
   ```bash
   source <(shelldock completion bash) 2>&1
   ```

4. **Reload your shell:**
   ```bash
   exec bash
   ```

## Other Shells

ShellDock also supports completion for other shells:

- **Zsh:** `shelldock completion zsh`
- **Fish:** `shelldock completion fish`
- **PowerShell:** `shelldock completion powershell`

Install them similarly to bash completion.

