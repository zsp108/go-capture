# GoCapture Makefile
# 纯 Golang 桌面截图软件编译、测试与 macOS .app 打包

BINARY_NAME=gocapture
APP_NAME=GoCapture.app
BUILD_DIR=bin
APP_BUNDLE=$(BUILD_DIR)/$(APP_NAME)
CONTENTS_DIR=$(APP_BUNDLE)/Contents
MACOS_DIR=$(CONTENTS_DIR)/MacOS
RESOURCES_DIR=$(CONTENTS_DIR)/Resources
GO_FILES=$(shell find . -name '*.go')

.PHONY: all build test clean run app run-app kill help

all: test build

help: ## 显示帮助信息
	@echo "GoCapture 桌面截图软件构建指南:"
	@echo "  make build       编译命令行原生二进制 (bin/gocapture)"
	@echo "  make app         编译并打包为 macOS 原生应用 (bin/GoCapture.app)"
	@echo "  make run-app     关闭旧进程并启动 macOS GoCapture.app"
	@echo "  make run         编译并在当前终端直接前台运行 (方便查看日志)"
	@echo "  make kill        关闭正在运行的 GoCapture 进程"
	@echo "  make test        运行所有核心模块单元测试"
	@echo "  make clean       清理编译产物"

kill: ## 关闭已在运行的 GoCapture 后台进程
	@echo "==> 终止后台已有 GoCapture 实例..."
	@killall $(BINARY_NAME) 2>/dev/null || true

build: ## 编译 Go 原生二进制文件
	@echo "==> 正在编译 GoCapture 原生二进制..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gocapture
	@echo "==> 二进制编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

app: kill build ## 编译并打包为标准的 macOS .app 应用程序包
	@echo "==> 正在构建 macOS 原生应用程序包 ($(APP_BUNDLE))..."
	@rm -rf $(APP_BUNDLE)
	@mkdir -p $(MACOS_DIR)
	@mkdir -p $(RESOURCES_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(MACOS_DIR)/$(BINARY_NAME)
	@chmod +x $(MACOS_DIR)/$(BINARY_NAME)
	@cp build/macos/Info.plist $(CONTENTS_DIR)/Info.plist
	@cp build/macos/PkgInfo $(CONTENTS_DIR)/PkgInfo
	@which codesign >/dev/null 2>&1 && codesign --force --deep --sign - $(APP_BUNDLE) 2>/dev/null || true
	@echo "==> 成功生成 macOS 应用程序: $(APP_BUNDLE)"

run-app: app ## 终止旧进程并启动全新 macOS .app 应用程序
	@echo "==> 启动 $(APP_BUNDLE)..."
	open -n $(APP_BUNDLE)

test: ## 执行所有核心单元测试
	@echo "==> 正在运行 GoCapture 核心算法与模块单元测试..."
	go test -v ./pkg/ocr ./pkg/loupe ./pkg/annotation ./pkg/pin
	@echo "==> 所有单元测试已通过！"

run: build ## 编译并在当前终端前台运行 (实时查看输出)
	@echo "==> 终止旧进程并前台启动 GoCapture..."
	@killall $(BINARY_NAME) 2>/dev/null || true
	./$(BUILD_DIR)/$(BINARY_NAME)

clean: ## 清理构建产物
	@echo "==> 清理构建目录..."
	@killall $(BINARY_NAME) 2>/dev/null || true
	@rm -rf $(BUILD_DIR)
	@echo "==> 清理完成"
