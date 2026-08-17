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

.PHONY: all build test clean run app run-app help

all: test build

help: ## 显示帮助信息
	@echo "GoCapture 桌面截图软件构建指南:"
	@echo "  make build       编译命令行原生二进制 (bin/gocapture)"
	@echo "  make app         编译并打包为 macOS 原生应用 (bin/GoCapture.app)"
	@echo "  make run-app     编译并启动 macOS GoCapture.app"
	@echo "  make test        运行所有核心模块单元测试 (OCR/放大镜/矢量/贴图)"
	@echo "  make run         编译并运行二进制"
	@echo "  make clean       清理编译产物"

build: ## 编译 Go 原生二进制文件
	@echo "==> 正在编译 GoCapture 原生二进制..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gocapture
	@echo "==> 二进制编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

app: build ## 编译并打包为标准的 macOS .app 应用程序包
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
	@echo "    可以通过双击打开或运行 'make run-app' 启动！"

run-app: app ## 启动编译好的 macOS .app 应用程序
	@echo "==> 正在启动 $(APP_BUNDLE)..."
	open $(APP_BUNDLE)

test: ## 执行所有核心单元测试
	@echo "==> 正在运行 GoCapture 核心算法与模块单元测试..."
	go test -v ./pkg/ocr ./pkg/loupe ./pkg/annotation ./pkg/pin
	@echo "==> 所有单元测试已通过！"

run: build ## 编译并直接运行二进制
	@echo "==> 启动 GoCapture..."
	./$(BUILD_DIR)/$(BINARY_NAME)

clean: ## 清理构建产物
	@echo "==> 清理构建目录..."
	@rm -rf $(BUILD_DIR)
	@echo "==> 清理完成"
