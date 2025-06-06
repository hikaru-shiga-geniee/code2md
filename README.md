# code2md-go

## 概要

`code2md-go` は、指定されたファイルやディレクトリの内容を読み込み、Markdownのコードブロック形式で標準出力するコマンドラインツールです。Python版 `code2md` のGo言語への移植版です。
特定のディレクトリパターンや、ドットから始まるファイル/ディレクトリをデフォルトで除外する機能を備えています。
主に、複数のコードファイルの内容をまとめてドキュメントやプロンプトに貼り付けたい場合に便利です。

## 機能

* 指定されたファイルの内容をMarkdownコードブロックとして出力します (` ```<lang>:<path> `)。
* 指定されたディレクトリ内を再帰的に探索し、含まれるファイルの内容をMarkdownコードブロックとして出力します。
* 出力されるコードブロックには、実行ディレクトリからの相対パスが付与されます。
* デフォルトで、`.` で始まるファイルやディレクトリ（例: `.env`, `.git`, `.vscode`）は無視されます。
* デフォルトで、特定の **ディレクトリ名** パターン（`__pycache__`, `build*`, `dist*`, `*.egg-info`, `node_modules`）に一致するディレクトリは探索対象から除外されます（ワイルドカード `*`, `?`, `[]` を使用）。
* `-i` または `--ignore` オプションで、探索時に無視する **ディレクトリ名やファイル名** のパターン（ワイルドカード使用可）を指定できます。
    * このオプションで指定されたパターンに一致するディレクトリが見つかった場合、そのディレクトリ以下の探索は行われません。
    * このオプションで指定されたパターンに一致するファイルは出力から除外されます。
    * 複数のパターンをカンマ区切りで一度に指定できます（例: `--ignore "*.md,*.py,*.json"`）。
* `--include-dotfiles` オプションを指定すると、デフォルトで無視される `. ` で始まるファイルやディレクトリも処理対象に含めます。
* `--no-default-ignores` オプションを指定すると、上記のデフォルト無視 **ディレクトリ** パターン (`__pycache__`, `build*` など) を適用しません。
* バイナリファイルなど、UTF-8テキストとして読み込めないファイルは警告メッセージを標準エラー出力に出力してスキップします。

## 動作環境

* **Go:** 1.22 以上

## インストールとビルド

### ソースからビルド

1. **リポジトリのクローン:**
   ```bash
   git clone <リポジトリのURL>
   cd code2md-go
   ```

2. **依存関係の取得:**
   ```bash
   go mod download
   ```

3. **ビルド:**
   ```bash
   make build
   ```
   または
   ```bash
   go build -o bin/code2md ./cmd/code2md
   ```

### バイナリのインストール

ビルドしたバイナリを、PATHが通っている場所に配置します:

```bash
# Linuxの例
sudo cp bin/code2md /usr/local/bin/

# macOSの例
cp bin/code2md /usr/local/bin/

# Windowsの例
# bin\code2md.exe をPATHが通っている場所にコピー
```

## 使い方

```bash
# 単一ファイル
code2md main.go

# 複数ファイル
code2md main.go go.mod README.md

# ディレクトリ (内部を再帰的に処理)
code2md src/

# ファイルとディレクトリの組み合わせ
code2md main.go docs/

# 出力をファイルに保存
code2md src/ > output.md
```

### オプション

* **`-i <パターン>` / `--ignore <パターン>`:** 無視する **ディレクトリ名やファイル名** のパターンを指定します。複数指定可能です。ワイルドカード (`*`, `?`, `[]`) が使えます。
    ```bash
    # node_modules ディレクトリと *.md ファイルを無視してカレントディレクトリを処理
    code2md . -i node_modules -i "*.md"

    # temp で始まるディレクトリを無視
    code2md . --ignore "temp*"
    
    # 複数のファイルタイプをカンマ区切りで一度に指定
    code2md . --ignore "*.md,*.py,*.json"
    
    # __init__.py ファイルを無視
    code2md . --ignore "__init__"
    ```

* **`--include-dotfiles`:** 通常無視される `.git`, `.env` のようなドットから始まるファイルやディレクトリを処理対象に含めます。
    ```bash
    # .env ファイルも出力に含める
    code2md . --include-dotfiles
    ```

* **`--no-default-ignores`:** デフォルトで設定されている無視 **ディレクトリ** パターン (`__pycache__`, `build*`, `dist*`, `*.egg-info`, `node_modules`) を無効化し、これらのディレクトリも探索対象とします。
    ```bash
    # build ディレクトリの中身も処理対象とする
    code2md . --no-default-ignores
    ```

## 開発者向け情報

* **テストの実行:**
    ```bash
    make test
    # または
    go test ./...
    ```

* **クリーンアップ:**
    ```bash
    make clean
    ```

* **ヘルプ:**
    ```bash
    make help
    ```

## ライセンス

[MITライセンス](LICENSE)

## リリース方法

このプロジェクトはGitHub ActionsとGoReleaserを使用して、クロスプラットフォーム（Windows、macOS、Linux）向けのバイナリを自動的にビルドしリリースします。

リリースするには、以下の手順に従ってください：

1. コードの変更をコミットしてプッシュします
   ```bash
   git add .
   git commit -m "リリース準備"
   git push origin main
   ```

2. 新しいバージョンタグを作成してプッシュします
   ```bash
   git tag v0.1.0  # バージョン番号は適宜変更してください
   git push origin v0.1.0
   ```

3. GitHub Actionsが自動的に実行され、リリースページに各プラットフォーム向けのバイナリが公開されます

