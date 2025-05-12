# Makefile for code2md-go project

.PHONY: all build test clean help

# デフォルトターゲット
all: build

# ヘルプメッセージを表示
help:
	@echo "Makefile targets for code2md-go:"
	@echo "  make build       - Go プログラムをビルドします"
	@echo "  make test        - テストを実行します"
	@echo "  make clean       - ビルド成果物を削除します"
	@echo "  make all         - ビルドを実行します (デフォルト)"
	@echo "  make help        - このヘルプメッセージを表示します"
	@echo ""
	@echo "Notes:"
	@echo "  - このプロジェクトは Go 1.22 以上が必要です"

# ビルド
build:
	@echo "Building code2md..."
	@mkdir -p bin
	go build -o bin/code2md ./code2md
	@echo "Build finished successfully."
	@echo "Executable is located in the 'bin/' directory."

# テスト実行
test:
	@echo "Running tests..."
	go test ./...

# 依存関係の整理
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# クリーンアップ
clean:
	@echo "Cleaning up generated files and directories..."
	rm -rf bin
	@echo "Cleanup finished." 