.PHONY: build build-local build-cross run dev clean clean-dist help test fmt lint verify

# 版本信息
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建参数
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)
BINARY_NAME := api-server

# 跨平台构建矩阵与产物目录
PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64
DIST_DIR := dist
CGO_ENABLED ?= 0
# 额外随包文件（存在才会被复制）
PACKAGE_FILES ?= README.md config.yaml.example

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 构建生产版本（CROSS=1 启用跨平台打包）
	@if [ "$(CROSS)" = "1" ]; then \
	  $(MAKE) build-cross; \
	else \
	  $(MAKE) build-local; \
	fi

build-local: ## 本地构建（默认）
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) .
	@echo "✓ 构建完成: bin/$(BINARY_NAME)"

build-cross: clean-dist ## 跨平台构建并打包到 dist/
	@echo "开始跨平台构建与打包: $(PLATFORMS)"
	@mkdir -p $(DIST_DIR)
	@set -e; \
	for PLATFORM in $(PLATFORMS); do \
	  OS=$${PLATFORM%/*}; ARCH=$${PLATFORM#*/}; \
	  EXT=""; [ "$$OS" = "windows" ] && EXT=".exe"; \
	  OUT_DIR="$(DIST_DIR)/$(BINARY_NAME)_$(VERSION)_$${OS}_$${ARCH}"; \
	  echo "- 构建 $$OS/$$ARCH"; \
	  mkdir -p "$$OUT_DIR"; \
	  GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=$(CGO_ENABLED) \
	    go build -trimpath -ldflags "$(LDFLAGS) -s -w" -o "$$OUT_DIR/$(BINARY_NAME)$$EXT" .; \
	  for f in $(PACKAGE_FILES); do \
	    [ -f "$$f" ] && cp "$$f" "$$OUT_DIR/" || true; \
	  done; \
	  if command -v zip >/dev/null 2>&1 && [ "$$OS" = "windows" ]; then \
	    (cd "$(DIST_DIR)" && zip -q -r "$(BINARY_NAME)_$(VERSION)_$${OS}_$${ARCH}.zip" "$$(basename "$$OUT_DIR")"); \
	  else \
	    (cd "$(DIST_DIR)" && tar -czf "$(BINARY_NAME)_$(VERSION)_$${OS}_$${ARCH}.tar.gz" "$$(basename "$$OUT_DIR")"); \
	  fi; \
	  rm -rf "$$OUT_DIR"; \
	done; \
	echo "✓ 打包完成，产物位于 $(DIST_DIR)/"

run: build-local ## 构建并运行（生产模式）
	@echo "启动服务（生产模式）..."
	@./bin/$(BINARY_NAME)

dev: build-local ## 构建并运行（开发模式）
	@echo "启动服务（开发模式）..."
	@./bin/$(BINARY_NAME) --dev

clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -f $(BINARY_NAME)
	@echo "✓ 清理完成"

clean-dist: ## 清理打包产物（dist/）
	@rm -rf $(DIST_DIR)
	@echo "✓ 已清理 $(DIST_DIR)/"

version: ## 显示版本信息
	@echo "Version:    $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

test: ## 运行测试并显示覆盖率
	@echo "运行测试（带覆盖率）..."
	@set -o pipefail; \
	OUT=$$(mktemp -t go-test-XXXXXX); \
	trap 'rm -f "$$OUT"' EXIT; \
	if go test -v -coverprofile=coverage.out -covermode=atomic ./... | tee "$$OUT"; then \
	  STATUS=0; \
	else \
	  STATUS=$$?; \
	fi; \
	PASS_PKGS=$$(grep -c '^ok[[:space:]]' "$$OUT" || true); \
	FAIL_PKGS=$$(grep -c '^FAIL[[:space:]]' "$$OUT" || true); \
	TOTAL_PKGS=$$((PASS_PKGS+FAIL_PKGS)); \
	PASS_TESTS=$$(grep -c '^--- PASS:' "$$OUT" || true); \
	FAIL_TESTS=$$(grep -c '^--- FAIL:' "$$OUT" || true); \
	SKIP_TESTS=$$(grep -c '^--- SKIP:' "$$OUT" || true); \
	if [ -f coverage.out ]; then \
	  TOTAL_COV=$$(go tool cover -func=coverage.out | awk '/^total:/ {print $$3}'); \
	else \
	  TOTAL_COV="N/A"; \
	fi; \
	echo "测试汇总：包 总数=$$TOTAL_PKGS 通过=$$PASS_PKGS 失败=$$FAIL_PKGS"; \
	echo "用例汇总：通过=$$PASS_TESTS 失败=$$FAIL_TESTS 跳过=$$SKIP_TESTS"; \
	echo "总覆盖率：$$TOTAL_COV"; \
	exit $$STATUS

fmt: ## 格式化代码
	@echo "格式化代码..."
	@gofmt -w $$(find . -name "*.go" -not -path "./vendor/*")
	@echo "✓ 格式化完成"

lint: ## 运行 go vet
	@echo "运行 go vet..."
	@go vet ./...
	@echo "✓ 检查完成"

verify: ## 顺序执行 fmt -> lint -> test
	@$(MAKE) fmt
	@$(MAKE) lint
	@$(MAKE) test

.DEFAULT_GOAL := help
