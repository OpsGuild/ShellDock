#!/bin/bash
set +e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

TESTS_PASSED=0
TESTS_FAILED=0
TESTS_SKIPPED=0
TOTAL_TESTS=0
FAILED_DETAILS=""
TEST_HOME=""

PLATFORM="${SHELLDOCK_PLATFORM:-auto}"
FILTER_COMMAND="${SHELLDOCK_COMMAND:-}"
BINARY="/tmp/shelldock"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." 2>/dev/null && pwd || echo /opt/shelldock)"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  GOARCH="amd64" ;;
    aarch64|arm64) GOARCH="arm64" ;;
    *)       GOARCH="amd64" ;;
esac

SRC_BINARY="${SHELLDOCK_BINARY:-}"
if [ -z "$SRC_BINARY" ] || [ ! -f "$SRC_BINARY" ]; then
    for candidate in \
        "$PROJECT_ROOT/build/shelldock-${OS}-${GOARCH}" \
        "/opt/shelldock/build/shelldock-linux-amd64" \
        "$PROJECT_ROOT/build/shelldock" \
        "$PROJECT_ROOT/.bin/shelldock"; do
        if [ -f "$candidate" ]; then
            SRC_BINARY="$candidate"
            break
        fi
    done
fi

if [ -z "$SRC_BINARY" ] || [ ! -f "$SRC_BINARY" ]; then
    echo -e "${RED}ERROR: Binary not found${NC}"
    echo "Searched in: $PROJECT_ROOT/build/, /opt/shelldock/build/"
    echo "Run 'make test-sandbox-build' or set SHELLDOCK_BINARY=/path/to/binary"
    exit 1
fi

cp "$SRC_BINARY" "$BINARY"
chmod +x "$BINARY"

cleanup() {
    if [ -n "$TEST_HOME" ] && [ -d "$TEST_HOME" ]; then
        rm -rf "$TEST_HOME"
    fi
}
trap cleanup EXIT

TEST_HOME="$(mktemp -d -t shelldock-sandbox-XXXXXX)"
export HOME="$TEST_HOME"
mkdir -p "$HOME/.shelldock"
if [ -d "$PROJECT_ROOT/repository" ]; then
    cp -R "$PROJECT_ROOT/repository/." "$HOME/.shelldock/"
fi

sd() {
    HOME="$HOME" "$BINARY" "$@"
}

test_start() {
    ((TOTAL_TESTS++))
    echo -e "${BLUE}  ▶ Test $TOTAL_TESTS: $1${NC}"
}

test_pass() {
    ((TESTS_PASSED++))
    echo -e "${GREEN}    ✅ PASSED${NC}"
}

test_fail() {
    ((TESTS_FAILED++))
    echo -e "${RED}    ❌ FAILED: $1${NC}"
    FAILED_DETAILS="${FAILED_DETAILS}\n  - $1"
}

test_skip() {
    ((TESTS_SKIPPED++))
    echo -e "${YELLOW}    ⏭  SKIPPED: $1${NC}"
}

sd config set "$PLATFORM" > /dev/null 2>&1

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  ShellDock Sandbox Test${NC}"
echo -e "${CYAN}  Platform: $PLATFORM${NC}"
echo -e "${CYAN}  Container: $(cat /etc/os-release 2>/dev/null | grep ^PRETTY_NAME | cut -d= -f2 | tr -d '\"')${NC}"
if [ -n "$FILTER_COMMAND" ]; then
    echo -e "${CYAN}  Filter: $FILTER_COMMAND${NC}"
fi
echo -e "${CYAN}========================================${NC}"
echo ""

get_all_commands() {
    sd list 2>&1 | grep -E '^\s+\w+:' | sed 's/^[^:]*: //' | tr ',' '\n' | sed 's/ *(also in local)//g' | tr -d ' ' | grep -v '^$' | sort -u
}

get_versions() {
    local cmd="$1"
    sd versions "$cmd" 2>&1 | grep -oE 'v[0-9]+' | sort -u
}

should_test_command() {
    local cmd="$1"
    if [ -z "$FILTER_COMMAND" ]; then
        return 0
    fi
    if [ "$cmd" = "$FILTER_COMMAND" ]; then
        return 0
    fi
    return 1
}

EXEC_TIMEOUT="${SHELLDOCK_EXEC_TIMEOUT:-30}"

command_should_skip_execution() {
    local cmd="$1"
    local version="$2"
    local show_output
    show_output=$(sd show "${cmd}@${version}" 2>&1)

    if [ "$PLATFORM" = "windows" ]; then
        echo "Windows commands cannot execute in Linux container (shelldock uses runtime.GOOS=linux)"
        return 0
    fi

    if echo "$show_output" | grep -qE '\{\{[a-zA-Z_]+\}\}'; then
        echo "requires interactive args (has {{placeholders}})"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'apt-get install|dnf install|yum install|pacman -S|brew install|Add-WindowsCapability|winget install|choco install|scoop install'; then
        echo "contains package install commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'Start-Process|Set-ItemProperty.*Run'; then
        echo "contains Windows process/registry commands"
        return 0
    fi


    if echo "$show_output" | grep -qiE 'systemctl (enable|start|restart|stop)|Set-Service|Restart-Service|Start-Service|launchctl'; then
        echo "contains service management commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'ssh-keygen|ufw (enable|allow|deny)|firewall-cmd|New-NetFirewallRule|iptables'; then
        echo "contains security/firewall commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'curl.*\| *(sh|bash)|wget.*\| *(sh|bash)|curl -L[Oo]|wget -O'; then
        echo "contains remote download/install commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'tar -(C|x).*/(opt|usr)|ln -sf.*/usr/local/bin'; then
        echo "contains system-level install commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'npm install|pip install|pip3 install|go install|cargo install'; then
        echo "contains package manager install commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'mkswap|swapon|swapoff|fallocate|dd if=|mkfs'; then
        echo "contains disk/swap commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'netsh|New-NetIPAddress|Remove-NetIPAddress|Set-NetIPInterface|Disable-NetAdapter|Enable-NetAdapter|Get-NetAdapter|Restart-NetAdapter|Reset-NetAdapterAdvancedProperty'; then
        echo "contains network config commands"
        return 0
    fi

    if echo "$show_output" | grep -qiE 'Get-DnsClientCache|Clear-DnsClientCache|ipconfig /flushdns|Get-NetIPConfiguration|Set-DnsClientServerAddress'; then
        echo "contains Windows DNS/network commands"
        return 0
    fi

    return 1
}

run_with_timeout() {
    local duration="$1"
    shift

    # Use GNU timeout if available, otherwise use portable fallback
    if command -v timeout >/dev/null 2>&1; then
        timeout "$duration" "$@" 2>&1
        local exit_code=$?
        if [ $exit_code -eq 124 ]; then
            echo "TIMEOUT"
        fi
        return $exit_code
    else
        # Portable fallback for macOS (no GNU timeout)
        local output_file
        output_file=$(mktemp)
        HOME="$HOME" "$@" > "$output_file" 2>&1 &
        local pid=$!
        local elapsed=0
        while kill -0 "$pid" 2>/dev/null; do
            if [ "$elapsed" -ge "$duration" ]; then
                kill -9 "$pid" 2>/dev/null
                wait "$pid" 2>/dev/null
                cat "$output_file"
                echo "TIMEOUT"
                rm -f "$output_file"
                return 124
            fi
            sleep 1
            elapsed=$((elapsed + 1))
        done
        wait "$pid" 2>/dev/null
        local exit_code=$?
        cat "$output_file"
        rm -f "$output_file"
        return $exit_code
    fi
}

test_yaml_parse() {
    local cmd="$1"
    local version="$2"
    local label="${cmd}@${version}"

    test_start "YAML parse ${label}"
    output=$(sd show "${cmd}@${version}" 2>&1)

    if echo "$output" | grep -qiE 'error parsing|yaml:|cannot unmarshal|did not find expected|found character that cannot'; then
        test_fail "[${PLATFORM}] YAML parse error ${label}: $(echo "$output" | grep -iE 'error parsing|yaml:|cannot unmarshal|did not find expected' | head -1)"
    else
        test_pass
    fi
}

test_show() {
    local cmd="$1"
    local version="$2"
    local label="${cmd}@${version}"

    test_start "show ${label} on ${PLATFORM}"
    output=$(sd show "${cmd}@${version}" 2>&1)

    if echo "$output" | grep -q "Command Set: ${cmd}"; then
        test_pass
    else
        test_fail "[${PLATFORM}] show failed for ${label}"
    fi
}

test_platform_resolution() {
    local cmd="$1"
    local version="$2"
    local label="${cmd}@${version}"

    test_start "platform resolution ${label} on ${PLATFORM}"
    output=$(sd show "${cmd}@${version}" 2>&1)

    if echo "$output" | grep -qi "no commands\|0 commands"; then
        test_fail "[${PLATFORM}] no commands resolved for ${label}"
    else
        test_pass
    fi
}

test_dry_run() {
    local cmd="$1"
    local version="$2"
    local label="${cmd}@${version}"

    test_start "dry-run ${label} on ${PLATFORM}"

    output=$(echo "n" | sd "${cmd}@${version}" 2>&1)

    if echo "$output" | grep -qi "panic\|fatal\|segfault\|runtime error"; then
        test_fail "[${PLATFORM}] crash in dry-run ${label}"
    else
        test_pass
    fi
}

test_execute_and_validate() {
    local cmd="$1"
    local version="$2"
    local label="${cmd}@${version}"

    local skip_reason
    skip_reason=$(command_should_skip_execution "$cmd" "$version")
    if [ $? -eq 0 ]; then
        test_start "execute ${label} on ${PLATFORM}"
        test_skip "${skip_reason}"
        return
    fi

    test_start "execute ${label} on ${PLATFORM}"

    output=$(run_with_timeout "$EXEC_TIMEOUT" "$BINARY" "${cmd}@${version}" --yes 2>&1)
    overall_exit=$?

    echo -e "${CYAN}    ┌── output ──${NC}"
    echo "$output" | sed -n '/Executing commands/,$ p' | sed 's/^/    │ /'
    echo -e "${CYAN}    └────────────${NC}"

    if echo "$output" | grep -q '^TIMEOUT$'; then
        test_fail "[${PLATFORM}] execution timed out after ${EXEC_TIMEOUT}s for ${label}"
        return
    fi

    if echo "$output" | grep -qi "panic\|fatal\|segfault\|runtime error"; then
        test_fail "[${PLATFORM}] crash during execution of ${label}"
        return
    fi

    success_count=$(echo "$output" | grep -c '✅ Success' || true)
    fail_count_warn=$(echo "$output" | grep -c '⚠️.*Command failed' || true)
    fail_count_hard=$(echo "$output" | grep -c '❌ Command failed' || true)
    fail_count=$((fail_count_warn + fail_count_hard))
    total_steps=$(echo "$output" | grep -cE '^\[' || true)

    if [ "$total_steps" -eq 0 ]; then
        test_fail "[${PLATFORM}] no steps executed for ${label}"
        return
    fi

    if echo "$output" | grep -q '🎉 All commands executed successfully!'; then
        test_pass
    elif [ "$fail_count" -gt 0 ]; then
        test_fail "[${PLATFORM}] ${fail_count}/${total_steps} steps failed in ${label}"
        return
    else
        test_pass
    fi

    local executed_steps
    executed_steps=$(echo "$output" | grep -oE '^\[[0-9]+/' | sed 's|^\[||;s|/$||' | sort -un)

    for step_num in $executed_steps; do
        step_header=$(echo "$output" | grep -E "^\[$step_num/" | head -1)
        if [ -z "$step_header" ]; then
            continue
        fi

        step_desc=$(echo "$step_header" | sed 's/^\[[0-9]\+\/[0-9]\+\] //' | sed 's/ (step [0-9]\+)$//')

        step_block=$(echo "$output" | sed -n "/^\[$step_num\//,/^\[\|^🎉/p" | head -n -1)

        test_start "step $step_num '${step_desc}' in ${label} on ${PLATFORM}"

        if echo "$step_block" | grep -q '✅ Success'; then
            test_pass
        elif echo "$step_block" | grep -q '⚠️.*Skipping.*No command available'; then
            test_pass
        elif echo "$step_block" | grep -q '⚠️.*Command failed.*skip_on_error=true'; then
            test_pass
        elif echo "$step_block" | grep -q '⚠️.*Command failed'; then
            test_fail "[${PLATFORM}] step $step_num failed in ${label}: ${step_desc}"
        elif echo "$step_block" | grep -q '❌ Command failed'; then
            test_fail "[${PLATFORM}] step $step_num hard-failed in ${label}: ${step_desc}"
        else
            test_fail "[${PLATFORM}] step $step_num had no result in ${label}: ${step_desc}"
        fi
    done
}

echo -e "${BLUE}--- Command Discovery ---${NC}"
COMMANDS=$(get_all_commands)

if [ -z "$COMMANDS" ]; then
    echo -e "${RED}ERROR: No commands found${NC}"
    exit 1
fi

echo "Found commands: $(echo $COMMANDS | tr '\n' ' ')"
echo ""

for cmd in $COMMANDS; do
    if ! should_test_command "$cmd"; then
        continue
    fi

    echo -e "${CYAN}━━━ Testing: ${cmd} ━━━${NC}"

    VERSIONS=$(get_versions "$cmd")
    if [ -z "$VERSIONS" ]; then
        VERSIONS="v1"
    fi

    for ver in $VERSIONS; do
        echo -e "${BLUE}  --- ${cmd}@${ver} ---${NC}"

        test_yaml_parse "$cmd" "$ver"
        test_show "$cmd" "$ver"
        test_platform_resolution "$cmd" "$ver"
        test_dry_run "$cmd" "$ver"
        test_execute_and_validate "$cmd" "$ver"
    done

    echo ""
done

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  Test Summary - ${PLATFORM}${NC}"
echo -e "${CYAN}========================================${NC}"
echo -e "${GREEN}  Passed:  $TESTS_PASSED${NC}"
echo -e "${RED}  Failed:  $TESTS_FAILED${NC}"
echo -e "${YELLOW}  Skipped: $TESTS_SKIPPED${NC}"
echo "  Total:   $TOTAL_TESTS"

if [ $TESTS_FAILED -gt 0 ]; then
    echo ""
    echo -e "${RED}  Failed tests:${NC}"
    echo -e "$FAILED_DETAILS"
fi

echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
