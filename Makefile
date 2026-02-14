.PHONY: build install clean test package deb rpm test-sandbox test-sandbox-build

# Get version from git tag, or use "dev" if no tag
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")
BUILD_DIR := build
BINARY_NAME := shelldock
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

build:
	@echo "Building shelldock..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing shelldock..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf dist
	@echo "Clean complete!"

test: test-unit test-integration

test-unit:
	@echo "=========================================="
	@echo "🧪 Running Unit Tests"
	@echo "=========================================="
	@go test -v -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "📊 Coverage Report:"
	@go tool cover -func=coverage.out
	@echo ""
	@echo "✅ Unit tests completed"
	@echo ""

test-integration: build
	@echo "=========================================="
	@echo "🧪 Running Integration Tests"
	@echo "=========================================="
	@chmod +x test/test-suite.sh
	@./test/test-suite.sh
	@echo ""
	@echo "✅ Integration tests completed"

# Cross-platform builds
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "Building Windows binary..."
	@if command -v goversioninfo >/dev/null 2>&1; then \
		echo "Adding Windows version info..."; \
		goversioninfo -64 -o resource.syso packaging/chocolatey/versioninfo.json && \
		GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-X main.version=$(VERSION)" -buildmode=exe -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe . && \
		rm -f resource.syso; \
	else \
		GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-X main.version=$(VERSION)" -buildmode=exe -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .; \
	fi
	@echo "Cross-platform builds complete!"
	@echo ""
	@echo "Note: To reduce false positives, consider code signing the Windows binary:"
	@echo "  signtool sign /f certificate.pfx /p password /t http://timestamp.digicert.com $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

# Debian/Ubuntu package
deb: build
	@echo "Creating Debian package..."
	@mkdir -p dist/deb/DEBIAN
	@mkdir -p dist/deb/usr/local/bin
	@mkdir -p dist/deb/usr/share/doc/shelldock
	@mkdir -p dist/deb/usr/share/shelldock/repository
	@cp $(BUILD_DIR)/$(BINARY_NAME) dist/deb/usr/local/bin/
	@cp README.md dist/deb/usr/share/doc/shelldock/
	@cp -r repository/. dist/deb/usr/share/shelldock/repository/ 2>/dev/null || true
	@cat packaging/deb/control | sed "s/VERSION/$(VERSION)/g" > dist/deb/DEBIAN/control
	@cp packaging/deb/postinst dist/deb/DEBIAN/postinst
	@chmod +x dist/deb/DEBIAN/postinst
	@dpkg-deb --build dist/deb dist/shelldock_$(VERSION)_$(GOARCH).deb
	@echo "Debian package created: dist/shelldock_$(VERSION)_$(GOARCH).deb"

# RPM package (RedHat/CentOS/Fedora)
rpm: build
	@echo "Creating RPM package..."
	@mkdir -p dist/rpm/BUILD dist/rpm/BUILDROOT dist/rpm/RPMS dist/rpm/SOURCES dist/rpm/SPECS
	@cp $(BUILD_DIR)/$(BINARY_NAME) dist/rpm/SOURCES/
	@cat packaging/rpm/shelldock.spec | sed "s/VERSION/$(VERSION)/g" > dist/rpm/SPECS/shelldock.spec
	@rpmbuild --define "_topdir $(PWD)/dist/rpm" -bb dist/rpm/SPECS/shelldock.spec
	@echo "RPM package created in dist/rpm/RPMS/"

# Arch Linux package
arch: build
	@echo "Creating Arch Linux package..."
	@mkdir -p dist/arch/pkg dist/arch/src
	@cp $(BUILD_DIR)/$(BINARY_NAME) dist/arch/src/
	@cat packaging/arch/PKGBUILD | sed "s/VERSION/$(VERSION)/g" > dist/arch/PKGBUILD
	@cd dist/arch && makepkg -c
	@echo "Arch package created in dist/arch/"

deps:
	@go mod download
	@go mod tidy

# Cross-platform sandbox tests using Docker
# Usage:
#   make test-sandbox                              # All platforms, all commands
#   make test-sandbox PLATFORM=ubuntu              # Single platform
#   make test-sandbox COMMAND=openssh              # Single command on all platforms
#   make test-sandbox PLATFORM=fedora COMMAND=git  # Single command on single platform
#   make test-sandbox BUILD=0                      # Skip Docker image rebuild
PLATFORM ?=
COMMAND ?=
BUILD ?= 1

test-sandbox-build:
	@echo "Building shelldock for linux/amd64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/shelldock-linux-amd64 .
	@echo "Build complete"

test-sandbox: test-sandbox-build
	@echo "=========================================="
	@echo "🐳 Running Sandbox Tests"
	@echo "=========================================="
	@chmod +x test/sandbox/entrypoint.sh test/sandbox/run-sandbox.sh
	@PLATFORM=$(PLATFORM) COMMAND=$(COMMAND) BUILD=$(BUILD) ./test/sandbox/run-sandbox.sh
	@echo ""
	@echo "✅ Sandbox tests completed"

