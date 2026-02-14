#!/bin/bash
set +e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yaml"

PLATFORM="${PLATFORM:-}"
COMMAND="${COMMAND:-}"
BUILD_FLAG="${BUILD:-1}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Run ShellDock sandbox tests across Linux platforms using Docker."
    echo ""
    echo "Options:"
    echo "  PLATFORM=<name>    Test a single platform (ubuntu, debian, rockylinux, fedora, archlinux)"
    echo "  COMMAND=<name>     Test a single command (e.g., openssh, docker, git)"
    echo "  BUILD=0            Skip rebuilding Docker images"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Test all platforms, all commands"
    echo "  PLATFORM=ubuntu $0                    # Test only ubuntu"
    echo "  COMMAND=openssh $0                    # Test only openssh on all platforms"
    echo "  PLATFORM=fedora COMMAND=git $0        # Test only git on fedora"
    echo "  BUILD=0 $0                            # Skip Docker image rebuild"
    echo ""
    echo "Makefile targets:"
    echo "  make test-sandbox                     # Test all platforms, all commands"
    echo "  make test-sandbox PLATFORM=ubuntu     # Test only ubuntu"
    echo "  make test-sandbox COMMAND=openssh     # Test only openssh"
    echo "  make test-sandbox PLATFORM=debian COMMAND=git"
}

if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    usage
    exit 0
fi

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  ShellDock Sandbox Test Runner${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

ARCH=$(uname -m)
case "$ARCH" in
    x86_64)  GOARCH="amd64" ;;
    aarch64) GOARCH="arm64" ;;
    *)       GOARCH="amd64" ;;
esac

BINARY="$PROJECT_ROOT/build/shelldock-linux-${GOARCH}"

if [ ! -f "$BINARY" ]; then
    echo -e "${YELLOW}Binary not found at $BINARY${NC}"
    echo -e "${BLUE}Building shelldock for linux/${GOARCH}...${NC}"
    cd "$PROJECT_ROOT"
    GOOS=linux GOARCH="$GOARCH" go build -o "$BINARY" .
    if [ $? -ne 0 ]; then
        echo -e "${RED}Build failed${NC}"
        exit 1
    fi
    echo -e "${GREEN}Build complete${NC}"
    echo ""
fi

AVAILABLE_PLATFORMS="ubuntu debian rockylinux fedora archlinux windows"

if [ -n "$PLATFORM" ]; then
    valid=false
    for p in $AVAILABLE_PLATFORMS; do
        if [ "$p" = "$PLATFORM" ]; then
            valid=true
            break
        fi
    done
    if [ "$valid" = false ]; then
        echo -e "${RED}Invalid platform: $PLATFORM${NC}"
        echo "Available: $AVAILABLE_PLATFORMS"
        exit 1
    fi
    PLATFORMS="$PLATFORM"
else
    PLATFORMS="$AVAILABLE_PLATFORMS"
fi

echo -e "${BLUE}Platforms: ${PLATFORMS}${NC}"
if [ -n "$COMMAND" ]; then
    echo -e "${BLUE}Command filter: ${COMMAND}${NC}"
fi
echo ""

TOTAL_PLATFORMS=0
PASSED_PLATFORMS=0
FAILED_PLATFORMS=0
PLATFORM_RESULTS=""

for plat in $PLATFORMS; do
    ((TOTAL_PLATFORMS++))
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}  Platform: ${plat}${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    BUILD_ARGS=""
    if [ "$BUILD_FLAG" = "1" ]; then
        BUILD_ARGS="--build"
    fi

    COMMAND="$COMMAND" docker compose -f "$COMPOSE_FILE" run --rm $BUILD_ARGS "$plat" 2>&1
    exit_code=$?

    if [ $exit_code -eq 0 ]; then
        ((PASSED_PLATFORMS++))
        PLATFORM_RESULTS="${PLATFORM_RESULTS}\n${GREEN}  ✅ ${plat}${NC}"
    else
        ((FAILED_PLATFORMS++))
        PLATFORM_RESULTS="${PLATFORM_RESULTS}\n${RED}  ❌ ${plat}${NC}"
    fi

    echo ""
done

docker compose -f "$COMPOSE_FILE" down --remove-orphans > /dev/null 2>&1

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  Sandbox Test Summary${NC}"
echo -e "${CYAN}========================================${NC}"
echo -e "$PLATFORM_RESULTS"
echo ""
echo -e "  Platforms: ${TOTAL_PLATFORMS}"
echo -e "  ${GREEN}Passed: ${PASSED_PLATFORMS}${NC}"
echo -e "  ${RED}Failed: ${FAILED_PLATFORMS}${NC}"
echo ""

if [ $FAILED_PLATFORMS -eq 0 ]; then
    echo -e "${GREEN}All platform tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some platform tests failed!${NC}"
    exit 1
fi
