package lang

import (
	"path/filepath"
	"strings"
)

// 拡張子から言語を推測するためのマッピング
var extMap = map[string]string{
	".py":         "python",
	".pyw":        "python",
	".js":         "javascript",
	".mjs":        "javascript",
	".cjs":        "javascript",
	".html":       "html",
	".htm":        "html",
	".css":        "css",
	".md":         "markdown",
	".json":       "json",
	".yaml":       "yaml",
	".yml":        "yaml",
	".toml":       "toml",
	".sh":         "bash",
	".bash":       "bash",
	".zsh":        "zsh",
	".java":       "java",
	".c":          "c",
	".h":          "c",
	".cpp":        "cpp",
	".hpp":        "cpp",
	".cc":         "cpp",
	".hh":         "cpp",
	".cs":         "csharp",
	".go":         "go",
	".rs":         "rust",
	".php":        "php",
	".rb":         "ruby",
	".swift":      "swift",
	".kt":         "kotlin",
	".scala":      "scala",
	".R":          "r",
	".lock":       "toml",
	".sql":        "sql",
	".xml":        "xml",
	".dockerfile": "dockerfile",
}

// ファイル名から言語を推測するためのマッピング
var fileMap = map[string]string{
	"dockerfile": "dockerfile",
	".gitignore": "gitignore",
	"makefile":   "makefile",
	"Makefile":   "makefile",
}

// Detect は、ファイルパスから言語を推測します
func Detect(path string) string {
	// ファイル名で検索
	filename := filepath.Base(path)
	if lang, ok := fileMap[filename]; ok {
		return lang
	}

	// 拡張子で検索
	ext := filepath.Ext(path)
	if ext != "" {
		if lang, ok := extMap[strings.ToLower(ext)]; ok {
			return lang
		}
	}

	// 一致するものがなければ空文字を返す
	return ""
}
