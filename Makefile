# Makefile for kongtools

.PHONY: build run clean test deps help install dev version changelog changelog-init

APP_NAME := kongtools
BUILD_DIR := build
BIN_DIR := bin

# 版本信息
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# ldflags for version injection
LDFLAGS := -ldflags "\
	-X 'kongtools/internal/pkg/version.Version=$(VERSION)' \
	-X 'kongtools/internal/pkg/version.Commit=$(COMMIT)' \
	-X 'kongtools/internal/pkg/version.BuildTime=$(BUILD_TIME)' \
	-X 'kongtools/internal/pkg/version.GoVersion=$(GO_VERSION)'"

# 默认目标
all: build

# 构建应用程序（带版本信息）
build:
	@echo "Building $(APP_NAME)..."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) .

# 构建应用程序（不带版本信息，用于快速开发）
build-dev:
	@echo "Building $(APP_NAME) (dev mode)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) .

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 运行应用程序
run: build
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME)

# 开发模式运行（不注入版本信息）
dev:
	@echo "Running in development mode..."
	go run .

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -f $(APP_NAME)

# 安装到系统
install: build
	@echo "Installing $(APP_NAME)..."
	@mkdir -p $(HOME)/.local/bin
	cp $(BUILD_DIR)/$(APP_NAME) $(HOME)/.local/bin/

# 显示版本信息
version:
	@echo "App Name: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

# git-chglog: 初始化配置
changelog-init:
	@echo "Initializing git-chglog configuration..."
	@if ! command -v git-chglog >/dev/null 2>&1; then \
		echo "Error: git-chglog is not installed."; \
		echo "Install it with: go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest"; \
		exit 1; \
	fi
	git-chglog --init

# git-chglog: 生成 CHANGELOG
changelog:
	@echo "Generating CHANGELOG.md..."
	@if ! command -v git-chglog >/dev/null 2>&1; then \
		echo "Error: git-chglog is not installed."; \
		echo "Install it with: go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest"; \
		exit 1; \
	fi
	git-chglog -o CHANGELOG.md

# git-chglog: 生成指定版本的 CHANGELOG
changelog-tag:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make changelog-tag TAG=v1.0.0"; \
		exit 1; \
	fi
	@echo "Generating CHANGELOG.md for tag $(TAG)..."
	git-chglog -o CHANGELOG.md $(TAG)

# 发布新版本：生成 CHANGELOG，提交并打标签
release:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG is required. Usage: make release TAG=v1.0.0"; \
		exit 1; \
	fi
	@if ! command -v git-chglog >/dev/null 2>&1; then \
		echo "Error: git-chglog is not installed."; \
		echo "Install it with: go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest"; \
		exit 1; \
	fi
	@echo "Releasing version $(TAG)..."
	@echo "1. Generating CHANGELOG.md..."
	git-chglog -o CHANGELOG.md --next-tag $(TAG)
	@echo "2. Committing CHANGELOG.md..."
	git add CHANGELOG.md
	git commit -m "chore(changelog): update changelog for $(TAG)"
	@echo "3. Creating tag $(TAG)..."
	git tag $(TAG)
	@echo "Done! Don't forget to push: git push origin main --tags"

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build           - Build the application with version info"
	@echo "  build-dev       - Build the application without version info (faster)"
	@echo "  run             - Build and run the application"
	@echo "  dev             - Run in development mode (go run)"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build files"
	@echo "  install         - Install to ~/bin"
	@echo "  deps            - Install dependencies"
	@echo "  version         - Show version information"
	@echo "  changelog-init  - Initialize git-chglog configuration"
	@echo "  changelog       - Generate CHANGELOG.md from all commits"
	@echo "  changelog-tag   - Generate CHANGELOG.md for specific tag (use TAG=v1.0.0)"
	@echo "  release         - Generate CHANGELOG, commit, and tag a new version (use TAG=v1.0.0)"
	@echo "  help            - Show this help"

