# Testing Input Handling

## Issue
The command should wait for user to type 'y' or 'n' before executing commands.

## Solution
The code now:
1. Uses `bufio.NewReader(os.Stdin)` to read from terminal
2. Flushes stdout to ensure prompt is visible
3. Blocks on `ReadString('\n')` until user presses Enter
4. Properly handles errors if stdin is unavailable

## Testing

To test interactively:
```bash
./shelldock test
# You should see the prompt and it should wait for your input
# Type 'y' and press Enter to execute
# Type 'n' or just press Enter to cancel
```

The prompt will display:
```
Do you want to execute these commands? (y/N):
```

And it will wait until you press Enter.

## Show Command

You can also preview commands without executing:
```bash
./shelldock show test
./shelldock show docker
```

This shows all commands without asking for confirmation.

