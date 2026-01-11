# ShellDock Manual Testing Guide

This guide walks through testing all ShellDock features manually.

## Prerequisites

1. Build the application:
   ```bash
   make build
   ```

2. Ensure you have a test command set:
   ```bash
   cp test-clean.yaml ~/.shelldock/test.yaml
   ```

## Test Checklist

### 1. Basic Commands

#### Test: List command sets
```bash
./shelldock list
```
**Expected:** Shows available command sets including "test"

#### Test: Show command set (preview)
```bash
./shelldock show test
```
**Expected:** Displays command set details without executing

#### Test: Show with platform info
```bash
./shelldock show test
```
**Expected:** Shows detected platform

### 2. Versioning

#### Test: List versions
```bash
./shelldock versions test
```
**Expected:** Lists v1, v2, v3 (with v3 marked as latest)

#### Test: Show specific version (v1)
```bash
./shelldock show test@v1
```
**Expected:** Shows v1 commands

#### Test: Show specific version (v2)
```bash
./shelldock show test@v2
```
**Expected:** Shows v2 commands with platform support

#### Test: Show latest version (default)
```bash
./shelldock show test
```
**Expected:** Shows v3 (latest) by default

#### Test: Show version using --ver flag
```bash
./shelldock show test --ver v1
```
**Expected:** Shows v1 commands

#### Test: Show version using --version flag
```bash
./shelldock show test --version v1
```
**Expected:** Shows v1 commands (--version is alias for --ver)

#### Test: Show version using tag
```bash
./shelldock show certbot@certonly
./shelldock show certbot --version nginx
```
**Expected:** Shows correct version matching the tag

#### Test: Error handling - non-existent version
```bash
./shelldock show test@v999
```
**Expected:** Shows error message

### 3. Platform Support

#### Test: Show current platform
```bash
./shelldock config show
```
**Expected:** Shows current platform setting and active platform

#### Test: Set platform to ubuntu
```bash
./shelldock config set ubuntu
./shelldock config show
```
**Expected:** Platform set to ubuntu

#### Test: Set platform to centos
```bash
./shelldock config set centos
./shelldock config show
```
**Expected:** Platform set to centos

#### Test: Set platform to auto
```bash
./shelldock config set auto
./shelldock config show
```
**Expected:** Platform set to auto

#### Test: Platform-specific commands
```bash
./shelldock config set ubuntu
./shelldock show test@v2
```
**Expected:** Shows ubuntu-specific commands

```bash
./shelldock config set centos
./shelldock show test@v2
```
**Expected:** Shows centos-specific commands

### 4. Step Filtering

#### Test: Skip steps (preview)
```bash
echo "n" | ./shelldock test --skip 1,2
```
**Expected:** Shows "Skipping steps: 1,2" and only displays remaining steps

#### Test: Skip range of steps
```bash
echo "n" | ./shelldock test --skip 1-3
```
**Expected:** Shows "Skipping steps: 1-3" and only displays remaining steps

#### Test: Only specific steps
```bash
echo "n" | ./shelldock test --only 1,3
```
**Expected:** Shows "Running only steps: 1,3" and only displays those steps

#### Test: Only range of steps
```bash
echo "n" | ./shelldock test --only 1-3
```
**Expected:** Shows "Running only steps: 1-3" and only displays those steps

#### Test: Error - Skip and only together
```bash
./shelldock test --skip 1 --only 2
```
**Expected:** Shows error: "cannot use both --skip and --only flags together"

#### Test: Skip all steps (edge case)
```bash
echo "n" | ./shelldock test --skip 1,2,3,4,5,6,7
```
**Expected:** Shows error: "No commands to execute after filtering"

#### Test: Only non-existent step (edge case)
```bash
./shelldock test --only 999
```
**Expected:** Shows error: "No commands to execute after filtering"

### 5. Flag Combinations

#### Test: Version with skip
```bash
echo "n" | ./shelldock test@v1 --skip 1
```
**Expected:** Shows v1 version with step 1 skipped

#### Test: Version with only
```bash
echo "n" | ./shelldock test@v2 --only 1,2
```
**Expected:** Shows v2 version with only steps 1 and 2

#### Test: Local flag with version
```bash
echo "n" | ./shelldock --local test@v1
```
**Expected:** Shows v1 version from local repository

### 6. Command Execution

#### Test: Execute with --yes flag
```bash
./shelldock test@v1 --yes
```
**Expected:** Executes commands without prompting

#### Test: Execute with skip and yes
```bash
./shelldock test@v1 --skip 2,3,4 --yes
```
**Expected:** Executes only non-skipped steps without prompting

#### Test: Execute with only and yes
```bash
./shelldock test@v1 --only 1 --yes
```
**Expected:** Executes only step 1 without prompting

#### Test: Interactive prompt (without --yes)
```bash
./shelldock test@v1
```
**Expected:** 
- Shows commands
- Prompts "Do you want to execute these commands? (y/N): "
- Waits for input
- Type 'y' to execute, 'n' to cancel

### 7. Direct Execution

#### Test: Direct execution (no subcommand)
```bash
echo "n" | ./shelldock test
```
**Expected:** Same as `./shelldock run test`

#### Test: Run subcommand
```bash
echo "n" | ./shelldock run test
```
**Expected:** Executes test command set

### 8. Error Handling

#### Test: Non-existent command set
```bash
./shelldock nonexistent
```
**Expected:** Shows error: "command set 'nonexistent' not found"

#### Test: Invalid skip format
```bash
./shelldock test --skip invalid
```
**Expected:** Shows error about invalid format

#### Test: Invalid only format
```bash
./shelldock test --only invalid
```
**Expected:** Shows error about invalid format

#### Test: Invalid version format
```bash
./shelldock test@invalid-version
```
**Expected:** Shows error: version not found

### 9. Help and Version

#### Test: Help command
```bash
./shelldock --help
```
**Expected:** Shows help text

#### Test: Version command
```bash
./shelldock --version
```
**Expected:** Shows version number

### 10. Edge Cases

#### Test: Empty command set
Create an empty command set and test:
```bash
cat > ~/.shelldock/empty.yaml << 'EOF'
name: empty
description: Empty command set
version: "v1"
commands: []
EOF

./shelldock show empty
```
**Expected:** Shows empty command set without errors

#### Test: Platform-only commands (no fallback)
```bash
cat > ~/.shelldock/platform-only.yaml << 'EOF'
name: platform-only
description: Platform only commands
version: "v1"
commands:
  - description: Platform command
    platforms:
      ubuntu: echo "ubuntu only"
EOF

./shelldock config set ubuntu
./shelldock show platform-only
```
**Expected:** Shows ubuntu command

### 11. TUI Management

#### Test: Manage command
```bash
./shelldock manage
```
**Expected:** Opens interactive TUI for managing command sets

**In TUI:**
- Navigate with arrow keys
- Press Enter to view
- Press 'a' to add
- Press 'd' to delete
- Press 'q' to quit

## Test Results Template

```
Manual Testing Results
======================

Date: ___________
Tester: ___________

1. Basic Commands: [ ] Pass [ ] Fail
2. Versioning: [ ] Pass [ ] Fail
3. Platform Support: [ ] Pass [ ] Fail
4. Step Filtering: [ ] Pass [ ] Fail
5. Flag Combinations: [ ] Pass [ ] Fail
6. Command Execution: [ ] Pass [ ] Fail
7. Direct Execution: [ ] Pass [ ] Fail
8. Error Handling: [ ] Pass [ ] Fail
9. Help and Version: [ ] Pass [ ] Fail
10. Edge Cases: [ ] Pass [ ] Fail
11. TUI Management: [ ] Pass [ ] Fail

Notes:
_______
_______
_______
```

