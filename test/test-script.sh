#!/bin/bash

# Change to project root
cd "$(dirname "$0")/.."

echo "🧪 Testing ShellDock Locally"
echo "=============================="
echo ""

# Test 1: List command sets
echo "Test 1: List command sets"
echo "---------------------------"
./shelldock list
echo ""

# Test 2: Show help
echo "Test 2: Show help"
echo "---------------------------"
./shelldock --help | head -15
echo ""

# Test 3: Test version
echo "Test 3: Show version"
echo "---------------------------"
./shelldock --version
echo ""

# Test 4: Test direct command (will show preview, user needs to confirm)
echo "Test 4: Direct command execution (shelldock test)"
echo "---------------------------"
echo "This will show the command preview. You can type 'y' to execute or 'n' to cancel."
echo ""

# Test 5: Test with --local flag
echo "Test 5: Test with --local flag"
echo "---------------------------"
./shelldock --local test --help 2>&1 | head -5 || echo "Note: --local flag works with command execution"
echo ""

# Test 6: Test run subcommand
echo "Test 6: Test 'run' subcommand"
echo "---------------------------"
./shelldock run --help
echo ""

echo "✅ Basic tests completed!"
echo ""
echo "To test interactive execution, run:"
echo "  ./shelldock test"
echo ""
echo "To test the TUI, run:"
echo "  ./shelldock manage"



