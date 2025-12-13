# ShellDock Features

## Core Features

### 1. Direct Command Execution
```bash
shelldock docker          # Run docker command set
shelldock --local docker  # Prefer local repository
```

### 2. Subcommands
```bash
shelldock run docker     # Explicit run command
shelldock list           # List all command sets
shelldock show docker    # Preview commands without executing
shelldock manage         # Interactive TUI for managing commands
shelldock sync           # Sync from cloud repository
```

### 3. Step Filtering

#### Skip Steps
Skip specific steps using `--skip`:
```bash
shelldock docker --skip 1,2,3        # Skip steps 1, 2, and 3
shelldock docker --skip 1-3           # Skip steps 1 through 3
shelldock docker --skip 1,3,5         # Skip steps 1, 3, and 5
```

#### Run Only Specific Steps
Run only specific steps using `--only`:
```bash
shelldock docker --only 1,3,5         # Run only steps 1, 3, and 5
shelldock docker --only 1-3            # Run only steps 1 through 3
shelldock docker --only 2              # Run only step 2
```

**Note:** You cannot use both `--skip` and `--only` together.

### 4. Interactive Confirmation
All command executions show a preview and wait for confirmation:
```
Do you want to execute these commands? (y/N):
```

### 5. Repository Priority
- **Default**: Cloud repository first, then local
- **With `--local` flag**: Local repository first, then cloud

### 6. Command Set Format
Command sets are stored as YAML files:
```yaml
name: docker
description: Docker setup commands
version: "1.0.0"
commands:
  - description: Install Docker
    command: sudo apt-get install docker
    skip_on_error: false  # optional
```

## Examples

### Basic Usage
```bash
# List available command sets
shelldock list

# Preview commands
shelldock show docker

# Run all commands
shelldock docker

# Run with step filtering
shelldock docker --skip 1,2
shelldock docker --only 3,4,5
```

### Advanced Usage
```bash
# Use local repository
shelldock --local my-commands

# Skip a range of steps
shelldock docker --skip 1-5

# Run only specific steps
shelldock docker --only 6-8

# Combine with local flag
shelldock --local docker --skip 1,2
```

## Error Handling

- If all commands are filtered out, an error is shown
- If both `--skip` and `--only` are used, an error is shown
- Invalid step numbers show descriptive errors
- Commands with `skip_on_error: true` continue on failure

