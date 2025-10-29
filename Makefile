# Makefile for kongtools

.PHONY: build run clean test deps help install dev version

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
	@mkdir -p $(HOME)/bin
	cp $(BUILD_DIR)/$(APP_NAME) $(HOME)/bin/

# 显示版本信息
version:
	@echo "App Name: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

# 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build      - Build the application with version info"
	@echo "  build-dev  - Build the application without version info (faster)"
	@echo "  run        - Build and run the application"
	@echo "  dev        - Run in development mode (go run)"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build files"
	@echo "  install    - Install to ~/bin"
	@echo "  deps       - Install dependencies"
	@echo "  version    - Show version information"
	@echo "  help       - Show this help"

