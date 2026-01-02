# Testing Guide

This document describes the testing strategy and how to run tests for ShellDock.

## Test Types

ShellDock has two types of tests:

1. **Unit Tests** - Go unit tests that test individual functions and packages
2. **Integration Tests** - Shell scripts that test the full CLI end-to-end

## Running Tests

### Run All Tests

```bash
make test
```

This runs both unit tests and integration tests.

### Run Unit Tests Only

```bash
make test-unit
# or
go test ./...
```

**With coverage:**
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**With HTML coverage report:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Run Integration Tests Only

```bash
make test-integration
# or
./test/test-suite.sh
```

**Run specific integration test:**
```bash
./test/test-all-features.sh
```

## Unit Tests

Unit tests are located alongside the source code in `*_test.go` files:

- `internal/repo/repository_test.go` - Repository operations
- `internal/repo/manager_test.go` - Repository manager
- `internal/config/config_test.go` - Configuration management
- `internal/cli/run_test.go` - CLI command execution logic

### Current Coverage

- `internal/cli`: 23.4% coverage
- `internal/config`: 42.5% coverage
- `internal/repo`: 56.2% coverage

### Running Specific Unit Tests

```bash
# Test a specific package
go test ./internal/repo -v

# Test a specific function
go test ./internal/repo -v -run TestGetCommandSet

# Run with race detector
go test -race ./...

# Run tests multiple times (to catch flaky tests)
go test -count=10 ./...
```

## Integration Tests

Integration tests are shell scripts in the `test/` directory that test the full CLI:

- `test/test-suite.sh` - Comprehensive test suite
- `test/test-all-features.sh` - Feature-specific tests
- `test/test-*.yaml` - Test command set files

### What Integration Tests Cover

- ✅ Basic functionality (list, show, run)
- ✅ Versioning (v1, v2, v3, latest detection)
- ✅ Platform support (detection, configuration, platform-specific commands)
- ✅ Step filtering (--skip, --only, ranges)
- ✅ Flag combinations (version + skip, version + only, etc.)
- ✅ Error handling (non-existent sets, invalid formats, conflicts)
- ✅ Command execution (with --yes flag)
- ✅ Dynamic arguments (--args flag, interactive prompting)
- ✅ Edge cases (empty sets, platform-only commands, etc.)

### Running Integration Tests

```bash
# Full test suite
./test/test-suite.sh

# Feature tests
./test/test-all-features.sh

# With verbose output
bash -x ./test/test-suite.sh
```

## CI/CD Testing

Tests run automatically in GitHub Actions:

### Test Workflow (`.github/workflows/test.yml`)

Runs on:
- Push to main/master/develop branches
- Pull requests
- Manual trigger

Tests:
- Unit tests on multiple OS (Ubuntu, macOS, Windows)
- Multiple Go versions (1.21, 1.22)
- Linting with golangci-lint
- Integration tests
- Build verification

### Release Workflow (`.github/workflows/release.yml`)

**Tests must pass before any release can happen.**

The release workflow:
1. Runs unit tests first (`test` job)
2. Only proceeds with builds if tests pass
3. All package jobs depend on the test job

## Writing Tests

### Writing Unit Tests

Create a `*_test.go` file in the same package:

```go
package repo

import "testing"

func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test"
    
    // Act
    result := MyFunction(input)
    
    // Assert
    if result != "expected" {
        t.Errorf("Expected 'expected', got %q", result)
    }
}
```

### Writing Integration Tests

Add test cases to `test/test-suite.sh`:

```bash
test_start "My new feature test"
if sd my-command 2>&1 | grep -q "expected output"; then
    test_pass
else
    test_fail "Expected output not found"
fi
```

## Test Files

### Unit Test Files

- `internal/repo/repository_test.go` - 11 tests
- `internal/repo/manager_test.go` - 3 tests
- `internal/config/config_test.go` - 7 tests
- `internal/cli/run_test.go` - 6 tests

### Integration Test Files

- `test/test-suite.sh` - Main integration test suite
- `test/test-all-features.sh` - Feature-specific tests
- `test/test-clean.yaml` - Clean test command set
- `test/test-commands.yaml` - Command test set
- `test/test-multi-version.yaml` - Multi-version test set

## Coverage Goals

Current coverage:
- Overall: ~33.6%
- Target: 70%+ for critical packages

### Priority Areas for More Tests

1. **TUI package** (0% coverage) - Terminal UI components
2. **CLI package** (23.4% coverage) - More command execution scenarios
3. **Config package** (42.5% coverage) - Edge cases in platform detection
4. **Repo package** (56.2% coverage) - Error handling and edge cases

## Debugging Tests

### Verbose Output

```bash
go test -v ./...
```

### Run Single Test

```bash
go test -v -run TestGetCommandSet ./internal/repo
```

### Debug with Delve

```bash
dlv test ./internal/repo
```

### Test with Race Detector

```bash
go test -race ./...
```

### Test Timeout

```bash
go test -timeout 30s ./...
```

## Best Practices

1. **Write tests before fixing bugs** - Reproduce the bug in a test first
2. **Test edge cases** - Empty inputs, nil values, invalid data
3. **Use table-driven tests** - For multiple test cases
4. **Keep tests fast** - Unit tests should run in milliseconds
5. **Test error paths** - Don't just test happy paths
6. **Use meaningful test names** - `TestGetCommandSet_NotFound` is better than `Test1`
7. **Clean up resources** - Use `t.TempDir()` for temporary files
8. **Test in isolation** - Tests shouldn't depend on each other

## Troubleshooting

### Tests Pass Locally But Fail in CI

- Check for hardcoded paths
- Verify environment variables
- Check for race conditions (use `-race` flag)
- Ensure test data is included in repository

### Integration Tests Fail

- Ensure `shelldock` binary is built: `go build -o shelldock .`
- Check test directory permissions
- Verify test YAML files are valid
- Check for platform-specific issues

### Coverage Not Updating

- Run with `-coverprofile=coverage.out`
- Ensure you're testing the right packages
- Check that test files are in the same package

## Continuous Improvement

- Add tests for new features
- Increase coverage for existing code
- Refactor to make code more testable
- Add integration tests for new CLI features
- Document test scenarios in this file

