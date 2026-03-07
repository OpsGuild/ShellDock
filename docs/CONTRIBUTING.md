# Contributing

> **[← Back to README](../README.md)** · [Testing](TESTING.md) · [Manual Testing](MANUAL_TESTING.md) · [Command Set Format](COMMAND_SET_FORMAT.md)

Contributions are welcome! Please feel free to submit a Pull Request.

---

## Contribution Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Adding Command Sets

To contribute command sets to the bundled repository:

1. Create a well-documented YAML file (see [Command Set Format](COMMAND_SET_FORMAT.md))
2. Include platform-specific commands where applicable
3. Test on multiple platforms
4. Place it in the appropriate subdirectory under `repository/`
5. Submit via Pull Request

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional)

### Setup

```bash
git clone https://github.com/OpsGuild/ShellDock.git
cd shelldock
go mod download
```

### Build

```bash
# Simple build
go build -o shelldock .

# Using Makefile
make build

# Build for all platforms
make build-all
```

### Test

```bash
# Run all unit tests
go test ./... -v

# Run integration test suite
./test/test-suite.sh

# Run feature tests
./test/test-all-features.sh

# Using Makefile
make test
```

See [Testing Guide](TESTING.md) for detailed testing documentation and [Manual Testing Guide](MANUAL_TESTING.md) for manual testing procedures.

### Create Packages

```bash
# Debian package
make deb

# RPM package
make rpm

# Arch package
make arch
```

### Project Structure

```
shelldock/
├── internal/
│   ├── cli/            # Command-line interface
│   ├── config/         # Configuration management
│   ├── repo/           # Repository management
│   └── tui/            # Terminal UI
├── examples/           # Example command sets
├── repository/         # Bundled command sets
├── packaging/          # Package configurations
├── scripts/            # Installation scripts
├── test/               # Integration tests
├── docs/               # Documentation
├── static/             # Static assets
├── go.mod
├── go.sum
├── main.go
├── Makefile
└── README.md
```

### Repository Structure

The bundled repository is organized into subdirectories by category:

```
repository/
├── devops/      # docker, kubernetes, pm2
├── editors/     # nvim
├── languages/   # go, nodejs, python, rust
├── security/    # openssh, ufw
├── system/      # swap, sysinfo
├── vcs/         # git
└── web/         # nginx, certbot
```

Subdirectories are transparent to users — commands are accessed by name only (e.g. `shelldock docker`), not by path.

### Local Repository

Your custom command sets are stored in `~/.shelldock/`. When you run a command, ShellDock checks:

1. **Local repository** (`~/.shelldock/`) — checked first
2. **Bundled repository** — checked if not found locally

This allows you to override bundled commands with your own versions.