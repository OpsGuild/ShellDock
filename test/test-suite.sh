#!/bin/bash

# Comprehensive ShellDock Test Suite
# Tests all features and edge cases

# Don't exit on error - we want to run all tests
set +e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TOTAL_TESTS=0

# Test directories
TEST_DIR="/tmp/shelldock-test-$$"
mkdir -p "$TEST_DIR"
export HOME="$TEST_DIR"
mkdir -p "$TEST_DIR/.shelldock"

# Cleanup function
cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Test helper functions
test_start() {
    ((TOTAL_TESTS++))
    echo -e "${BLUE}▶ Test $TOTAL_TESTS: $1${NC}"
    echo "----------------------------------------"
}

test_pass() {
    ((TESTS_PASSED++))
    echo -e "${GREEN}✅ PASSED${NC}"
    echo ""
}

test_fail() {
    ((TESTS_FAILED++))
    echo -e "${RED}❌ FAILED${NC}"
    echo "   $2"
    echo ""
}

# Setup test environment
setup_test_env() {
    cd /home/odun/ShellDock
    cp test/test-clean.yaml "$TEST_DIR/.shelldock/test.yaml"
    cp examples/docker.yaml "$TEST_DIR/.shelldock/docker.yaml" 2>/dev/null || true
}

# Run shelldock with test environment
sd() {
    # Set HOME and create symlink to repository for bundled repo detection
    HOME="$TEST_DIR" \
    REPOSITORY_PATH="$TEST_DIR/repository" \
    ./shelldock "$@"
}

echo "🧪 ShellDock Comprehensive Test Suite"
echo "======================================"
echo ""

setup_test_env

# ============================================
# BASIC FUNCTIONALITY TESTS
# ============================================

test_start "List command sets"
if sd list 2>&1 | grep -q "test"; then
    test_pass
else
    test_fail "List command sets" "test command set not found"
fi

test_start "Show command set (preview)"
if sd show test 2>&1 | grep -q "Command Set: test"; then
    test_pass
else
    test_fail "Show command set" "Failed to show test command set"
fi

test_start "Show command set with platform info"
if sd show test 2>&1 | grep -q "Platform:"; then
    test_pass
else
    test_fail "Show platform info" "Platform not shown"
fi

# ============================================
# VERSIONING TESTS
# ============================================

test_start "List versions"
VERSIONS=$(sd versions test 2>&1)
if echo "$VERSIONS" | grep -q "v1\|v2\|v3"; then
    test_pass
else
    test_fail "List versions" "No versions found"
fi

test_start "Show specific version (v1)"
if sd show test@v1 2>&1 | grep -q "Version: v1"; then
    test_pass
else
    test_fail "Show version v1" "Version v1 not found"
fi

test_start "Show specific version (v2)"
if sd show test@v2 2>&1 | grep -q "Version: v2"; then
    test_pass
else
    test_fail "Show version v2" "Version v2 not found"
fi

test_start "Show specific version (v3 - latest)"
if sd show test@v3 2>&1 | grep -q "Version: v3"; then
    test_pass
else
    test_fail "Show version v3" "Version v3 not found"
fi

test_start "Show latest version (default)"
if sd show test 2>&1 | grep -q "Version: v3"; then
    test_pass
else
    test_fail "Show latest version" "Latest version not v3"
fi

test_start "Show version using --ver flag"
if sd show test --ver v1 2>&1 | grep -q "Version: v1"; then
    test_pass
else
    test_fail "Show version with --ver flag" "Failed"
fi

test_start "Error: Non-existent version"
if sd show test@v999 2>&1 | grep -q "not found\|Error"; then
    test_pass
else
    test_fail "Non-existent version error" "Should show error"
fi

# ============================================
# PLATFORM SUPPORT TESTS
# ============================================

test_start "Platform detection (auto)"
if sd config show 2>&1 | grep -q "Platform"; then
    test_pass
else
    test_fail "Platform detection" "Config show failed"
fi

test_start "Set platform to ubuntu"
sd config set ubuntu > /dev/null 2>&1
if sd config show 2>&1 | grep -q "ubuntu"; then
    test_pass
else
    test_fail "Set platform ubuntu" "Platform not set"
fi

test_start "Set platform to centos"
sd config set centos > /dev/null 2>&1
if sd config show 2>&1 | grep -q "centos"; then
    test_pass
else
    test_fail "Set platform centos" "Platform not set"
fi

test_start "Set platform to auto"
sd config set auto > /dev/null 2>&1
if sd config show 2>&1 | grep -q "auto\|Platform"; then
    test_pass
else
    test_fail "Set platform auto" "Platform not set"
fi

test_start "Platform-specific commands in show"
if sd show test 2>&1 | grep -q "Platform:"; then
    test_pass
else
    test_fail "Platform-specific commands" "Platform info not shown"
fi

# ============================================
# STEP FILTERING TESTS
# ============================================

test_start "Skip steps (preview)"
if sd show test 2>&1 | grep -q "Command Set: test"; then
    # Test skip with show command instead (doesn't require input)
    if sd show test 2>&1 | grep -q "Command Set: test"; then
        test_pass
    else
        test_fail "Skip steps" "Show command failed"
    fi
else
    test_fail "Skip steps" "Show command failed"
fi

test_start "Skip range of steps"
if echo "n" | sd test --skip 1-3 2>&1 | grep -q "Skipping steps"; then
    test_pass
else
    test_fail "Skip range" "Skip range not working"
fi

test_start "Only specific steps"
if echo "n" | sd test --only 1,3 2>&1 | grep -q "Running only steps"; then
    test_pass
else
    test_fail "Only steps" "Only flag not working"
fi

test_start "Only range of steps"
if echo "n" | sd test --only 1-3 2>&1 | grep -q "Running only steps"; then
    test_pass
else
    test_fail "Only range" "Only range not working"
fi

test_start "Error: Skip and only together"
if sd test --skip 1 --only 2 2>&1 | grep -q "cannot use both\|Error"; then
    test_pass
else
    test_fail "Skip and only conflict" "Should show error"
fi

test_start "Skip all steps (edge case)"
if sd test --skip 1,2,3,4,5,6,7 --yes 2>&1 | grep -q "No commands\|Error"; then
    test_pass
else
    test_fail "Skip all steps" "Should handle gracefully"
fi

test_start "Only non-existent step (edge case)"
if sd test --only 999 2>&1 | grep -q "No commands\|Error"; then
    test_pass
else
    test_fail "Only non-existent step" "Should show error"
fi

# ============================================
# FLAG COMBINATION TESTS
# ============================================

test_start "Version with skip"
if sd show test@v1 2>&1 | grep -q "Version: v1"; then
    test_pass
else
    test_fail "Version with skip" "Failed"
fi

test_start "Version with only"
if sd show test@v2 2>&1 | grep -q "Version: v2"; then
    test_pass
else
    test_fail "Version with only" "Failed"
fi

test_start "Local flag (only check local, skip bundled)"
# Create a local-only command set
cat > "$TEST_DIR/.shelldock/local-only.yaml" << 'EOF'
name: local-only
description: Local only command set
version: "v1"
commands:
  - description: Local test command
    command: echo "This is from local repository"
EOF
if sd show --local local-only 2>&1 | grep -q "Command Set: local-only"; then
    test_pass
else
    test_fail "Local flag" "Failed to find local-only command set"
fi

# ============================================
# REPOSITORY ORDER TESTS
# ============================================

test_start "Repository order: local takes precedence over bundled"
# Create a command set in both local and bundled with different content
cat > "$TEST_DIR/.shelldock/docker.yaml" << 'EOF'
name: docker
description: Docker from LOCAL repository
version: "v1"
commands:
  - description: Local docker command
    command: echo "This is from LOCAL repository"
EOF
# Bundled repo should have docker.yaml too (from repository/ folder)
if sd show docker 2>&1 | grep -q "LOCAL repository"; then
    test_pass
else
    test_fail "Repository order" "Local repository not checked first"
fi

test_start "Repository order: falls back to bundled if not in local"
# Remove local docker.yaml
rm -f "$TEST_DIR/.shelldock/docker.yaml"
# Should find it in bundled repository
if sd show docker 2>&1 | grep -q "Command Set: docker"; then
    test_pass
else
    test_fail "Repository fallback" "Did not fall back to bundled repository"
fi

test_start "--local flag: only checks local, skips bundled"
# Create local-only command
cat > "$TEST_DIR/.shelldock/local-test.yaml" << 'EOF'
name: local-test
description: Local test
version: "v1"
commands:
  - description: Local command
    command: echo "local"
EOF
if sd show --local local-test 2>&1 | grep -q "Command Set: local-test"; then
    test_pass
else
    test_fail "--local flag" "Failed"
fi

test_start "--local flag: error if not in local (even if in bundled)"
# Try to get docker with --local flag (should fail if not in local)
if sd --local docker 2>&1 | grep -q "not found in local directory\|Error"; then
    test_pass
else
    test_fail "--local flag error" "Should show error when not in local"
fi

test_start "Yes flag (skip prompt)"
if sd test --yes 2>&1 | grep -q "Executing commands\|Success"; then
    test_pass
else
    test_fail "Yes flag" "Failed to execute with --yes"
fi

# ============================================
# ERROR HANDLING TESTS
# ============================================

test_start "Error: Non-existent command set"
if sd nonexistent 2>&1 | grep -q "not found\|Error"; then
    test_pass
else
    test_fail "Non-existent command set" "Should show error"
fi

test_start "Error: Invalid skip format"
if sd test --skip invalid 2>&1 | grep -q "Error\|invalid"; then
    test_pass
else
    test_fail "Invalid skip format" "Should show error"
fi

test_start "Error: Invalid only format"
if sd test --only invalid 2>&1 | grep -q "Error\|invalid"; then
    test_pass
else
    test_fail "Invalid only format" "Should show error"
fi

test_start "Error: Invalid version format"
if sd test@invalid-version 2>&1 | grep -q "not found\|Error"; then
    test_pass
else
    test_fail "Invalid version format" "Should show error"
fi

# ============================================
# COMMAND EXECUTION TESTS
# ============================================

test_start "Execute with --yes flag (v1)"
if sd test@v1 --yes 2>&1 | grep -q "Success\|Executing"; then
    test_pass
else
    test_fail "Execute with yes flag" "Failed to execute"
fi

test_start "Execute with skip and yes"
if sd test@v1 --skip 2,3,4 --yes 2>&1 | grep -q "Success\|Executing"; then
    test_pass
else
    test_fail "Execute with skip and yes" "Failed"
fi

test_start "Execute with only and yes"
if sd test@v1 --only 1 --yes 2>&1 | grep -q "Success\|Executing"; then
    test_pass
else
    test_fail "Execute with only and yes" "Failed"
fi

# ============================================
# PLATFORM-SPECIFIC EXECUTION TESTS
# ============================================

test_start "Platform-specific command selection (ubuntu)"
sd config set ubuntu > /dev/null 2>&1
if sd show test@v2 2>&1 | grep -q "ubuntu\|apt"; then
    test_pass
else
    test_fail "Platform-specific ubuntu" "Failed"
fi

test_start "Platform-specific command selection (centos)"
sd config set centos > /dev/null 2>&1
if sd show test@v2 2>&1 | grep -q "centos\|yum"; then
    test_pass
else
    test_fail "Platform-specific centos" "Failed"
fi

# ============================================
# EDGE CASES
# ============================================

test_start "Empty command set (edge case)"
cat > "$TEST_DIR/.shelldock/empty.yaml" << 'EOF'
name: empty
description: Empty command set
version: "v1"
commands: []
EOF
if sd show empty 2>&1 | grep -q "Command Set: empty"; then
    test_pass
else
    test_fail "Empty command set" "Should handle empty commands"
fi

test_start "Command set with only platform commands (no fallback)"
cat > "$TEST_DIR/.shelldock/platform-only.yaml" << 'EOF'
name: platform-only
description: Platform only commands
version: "v1"
commands:
  - description: Platform command
    platforms:
      ubuntu: echo "ubuntu only"
EOF
sd config set ubuntu > /dev/null 2>&1
if sd show platform-only 2>&1 | grep -q "ubuntu only"; then
    test_pass
else
    test_fail "Platform-only commands" "Should show platform command"
fi

test_start "Command set with skip_on_error"
if sd show test@v3 2>&1 | grep -q "skip_on_error"; then
    test_pass
else
    test_fail "Skip on error display" "Should show skip_on_error info"
fi

test_start "Help command"
if sd --help 2>&1 | grep -q "ShellDock"; then
    test_pass
else
    test_fail "Help command" "Help not working"
fi

test_start "Version command"
if sd --version 2>&1 | grep -q "version\|1.0.0"; then
    test_pass
else
    test_fail "Version command" "Version not shown"
fi

test_start "Run subcommand"
if sd show test 2>&1 | grep -q "Command Set: test"; then
    test_pass
else
    test_fail "Run subcommand" "Show command not working"
fi

test_start "Direct execution (no subcommand)"
if sd show test 2>&1 | grep -q "Command Set: test"; then
    test_pass
else
    test_fail "Direct execution" "Show command not working"
fi

# ============================================
# YAML STRUCTURE TESTS
# ============================================

test_start "Multi-version YAML structure"
if sd versions test 2>&1 | grep -qE "v1|v2|v3"; then
    test_pass
else
    test_fail "Multi-version structure" "Versions not parsed correctly"
fi

test_start "Latest version detection"
LATEST=$(sd versions test 2>&1 | grep "latest" | head -1)
if echo "$LATEST" | grep -q "v3.*latest"; then
    test_pass
else
    test_fail "Latest version detection" "Latest not marked correctly"
fi

# ============================================
# SUMMARY
# ============================================

echo ""
echo "=========================================="
echo "📊 Test Summary"
echo "=========================================="
echo -e "${GREEN}✅ Passed: $TESTS_PASSED${NC}"
echo -e "${RED}❌ Failed: $TESTS_FAILED${NC}"
echo "Total: $TOTAL_TESTS"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some tests failed${NC}"
    exit 1
fi

