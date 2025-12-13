#!/bin/bash

echo "🧪 ShellDock Comprehensive Feature Testing"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

test_step() {
    echo -e "${BLUE}▶ Test: $1${NC}"
    echo "----------------------------------------"
}

test_pass() {
    echo -e "${GREEN}✅ PASSED${NC}"
    ((TESTS_PASSED++))
    echo ""
}

test_fail() {
    echo -e "${YELLOW}❌ FAILED${NC}"
    ((TESTS_FAILED++))
    echo ""
}

# Change to project root
cd "$(dirname "$0")/.."

# Test 1: List command sets
test_step "1. List command sets"
if ./shelldock list | grep -q "test"; then
    test_pass
else
    test_fail
fi

# Test 2: Show command (preview)
test_step "2. Show command (preview without execution)"
if ./shelldock show test 2>&1 | grep -q "Command Set: test"; then
    test_pass
else
    test_fail
fi

# Test 3: List versions
test_step "3. List available versions"
if ./shelldock versions test 2>&1 | grep -q "1.0.0\|1.1.0\|2.0.0"; then
    test_pass
else
    test_fail
fi

# Test 4: Show specific version
test_step "4. Show specific version (1.0.0)"
if ./shelldock show test@1.0.0 2>&1 | grep -q "Version: 1.0.0"; then
    test_pass
else
    test_fail
fi

# Test 5: Show another version
test_step "5. Show specific version (2.0.0)"
if ./shelldock show test@2.0.0 2>&1 | grep -q "Version: 2.0.0"; then
    test_pass
else
    test_fail
fi

# Test 6: Platform detection
test_step "6. Platform detection and configuration"
./shelldock config show
if ./shelldock config show 2>&1 | grep -q "Platform"; then
    test_pass
else
    test_fail
fi

# Test 7: Show with platform-specific commands
test_step "7. Show platform-specific commands"
if ./shelldock show test 2>&1 | grep -q "Platform:"; then
    test_pass
else
    test_fail
fi

# Test 8: Test --skip flag (preview)
test_step "8. Test --skip flag (preview)"
if echo "n" | ./shelldock test --skip 1,2 2>&1 | grep -q "Skipping steps"; then
    test_pass
else
    test_fail
fi

# Test 9: Test --only flag (preview)
test_step "9. Test --only flag (preview)"
if echo "n" | ./shelldock test --only 1,3 2>&1 | grep -q "Running only steps"; then
    test_pass
else
    test_fail
fi

# Test 10: Test --local flag
test_step "10. Test --local flag"
if echo "n" | ./shelldock --local test 2>&1 | grep -q "Command Set: test"; then
    test_pass
else
    test_fail
fi

# Test 11: Test version flag
test_step "11. Test --ver flag"
if echo "n" | ./shelldock test --ver 1.0.0 2>&1 | grep -q "Version: 1.0.0"; then
    test_pass
else
    test_fail
fi

# Test 12: Test @version syntax
test_step "12. Test @version syntax"
if echo "n" | ./shelldock test@1.1.0 2>&1 | grep -q "Version: 1.1.0"; then
    test_pass
else
    test_fail
fi

# Test 13: Test run subcommand
test_step "13. Test 'run' subcommand"
if echo "n" | ./shelldock run test 2>&1 | grep -q "Command Set: test"; then
    test_pass
else
    test_fail
fi

# Test 14: Test help
test_step "14. Test help command"
if ./shelldock --help 2>&1 | grep -q "ShellDock"; then
    test_pass
else
    test_fail
fi

# Test 15: Test config set
test_step "15. Test config set command"
./shelldock config set auto > /dev/null 2>&1
if ./shelldock config show 2>&1 | grep -q "Platform"; then
    test_pass
else
    test_fail
fi

# Test 16: Test platform-specific command selection
test_step "16. Test platform-specific command selection"
PLATFORM=$(./shelldock config show 2>&1 | grep "Active platform" | awk '{print $3}')
if ./shelldock show test 2>&1 | grep -q "Platform: $PLATFORM"; then
    test_pass
else
    test_fail
fi

# Test 17: Test error handling (non-existent command set)
test_step "17. Test error handling (non-existent command set)"
if ./shelldock nonexistent 2>&1 | grep -q "not found"; then
    test_pass
else
    test_fail
fi

# Test 18: Test error handling (non-existent version)
test_step "18. Test error handling (non-existent version)"
if ./shelldock test@999.0.0 2>&1 | grep -q "not found\|Error"; then
    test_pass
else
    test_fail
fi

# Test 19: Test skip and only conflict
test_step "19. Test skip and only conflict detection"
if ./shelldock test --skip 1 --only 2 2>&1 | grep -q "cannot use both\|Error"; then
    test_pass
else
    test_fail
fi

# Test 20: Test version with skip
test_step "20. Test version with skip flag"
if echo "n" | ./shelldock test@1.0.0 --skip 1 2>&1 | grep -q "Version: 1.0.0"; then
    test_pass
else
    test_fail
fi

echo ""
echo "=========================================="
echo "📊 Test Summary"
echo "=========================================="
echo -e "${GREEN}✅ Passed: $TESTS_PASSED${NC}"
echo -e "${YELLOW}❌ Failed: $TESTS_FAILED${NC}"
echo "Total: $((TESTS_PASSED + TESTS_FAILED))"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some tests failed${NC}"
    exit 1
fi



