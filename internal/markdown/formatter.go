package markdown

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/your-org/code2md/internal/lang"
)

// isBinary は、データがバイナリファイルかどうかを判断します
// UTF-8として有効でないか、NUL文字を含む場合はバイナリとみなします
func isBinary(data []byte) bool {
	// UTF-8として有効でない場合はバイナリ
	if !utf8.Valid(data) {
		return true
	}

	// NUL文字を含む場合もバイナリとみなす
	if bytes.IndexByte(data, 0) != -1 {
		return true
	}

	return false
}

// Print は、ファイルリストの内容をMarkdownコードブロック形式で出力します
func Print(w io.Writer, files []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get current directory: %w", err)
	}

	var totalWords, totalChars, totalLines int

	for _, filePath := range files {
		// カレントディレクトリからの相対パスを取得
		relPath, err := filepath.Rel(cwd, filePath)
		if err != nil {
			// 相対パス取得に失敗した場合は絶対パスを使用
			relPath = filePath
		}

		// ファイル内容を読み込み
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error reading file '%s': %v. Skipping.\n", relPath, err)
			continue
		}

		// バイナリファイルのチェック
		if isBinary(data) {
			fmt.Fprintf(os.Stderr, "Warning: File '%s' could not be read as UTF-8 text. Skipping.\n", relPath)
			continue
		}

		// 言語タグを取得
		langTag := lang.Detect(filePath)

		// ファイルの統計情報を計算
		content := string(data)
		lines := strings.Split(content, "\n")
		words := len(strings.Fields(content))
		chars := utf8.RuneCountInString(content)

		// 統計を加算
		totalLines += len(lines)
		totalWords += words
		totalChars += chars

		// Markdownコードブロックとして出力
		fmt.Fprintf(w, "```%s:%s\n%s\n```\n\n", langTag, relPath, content)
	}

	// 最終的な統計情報を標準エラー出力に出力
	fmt.Fprintf(os.Stderr, "Total: %d lines, %d words, %d characters\n", totalLines, totalWords, totalChars)

	return nil
}
