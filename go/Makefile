# Makefile for Hakoniwa Go Implementation

.PHONY: help build test test-unit test-compat test-integration coverage lint clean run

# デフォルトターゲット
help:
	@echo "Hakoniwa Go Implementation - Available targets:"
	@echo "  make build              - Build executables"
	@echo "  make test               - Run all tests"
	@echo "  make test-unit          - Run unit tests"
	@echo "  make test-compat        - Run compatibility tests (Golden Testing)"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make coverage           - Generate coverage report"
	@echo "  make lint               - Run linter"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make run-main           - Run main server (development)"
	@echo "  make update-golden      - Update Golden test files"

# ビルド
build:
	@echo "Building hako-main..."
	go build -o bin/hako-main ./cmd/hako-main
	@echo "Building hako-mente..."
	go build -o bin/hako-mente ./cmd/hako-mente
	@echo "Build complete: bin/"

# テスト実行
test: test-unit test-compat test-integration

test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/...

test-compat:
	@echo "Running compatibility tests (Golden Testing)..."
	go test -v ./test/compatibility/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./test/integration/...

# カバレッジ
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Linter
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

# Golden ファイル更新
update-golden:
	@echo "Updating Golden test files..."
	UPDATE_GOLDEN=1 go test ./test/compatibility/...

# クリーンアップ
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# 開発サーバー実行
run-main:
	@echo "Starting hako-main (development server)..."
	go run ./cmd/hako-main

run-mente:
	@echo "Starting hako-mente (maintenance tool)..."
	go run ./cmd/hako-mente

# 依存関係の整理
tidy:
	@echo "Tidying go.mod..."
	go mod tidy

# フォーマット
fmt:
	@echo "Formatting code..."
	go fmt ./...

# ベンチマーク
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# すべてをチェック（CI/CD用）
check: lint test coverage
	@echo "All checks passed!"
