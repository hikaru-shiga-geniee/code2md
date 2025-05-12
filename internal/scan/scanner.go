package scan

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// デフォルトで無視するディレクトリ名のパターン
var defaultIgnore = []string{
	"__pycache__",
	"build*",
	"dist*",
	"*.egg-info",
	"node_modules",
}

// Options は、ファイル探索の設定オプション
type Options struct {
	UserIgnorePatterns  []string
	IncludeDotfiles     bool
	ApplyDefaultIgnores bool
}

// isIgnored は、指定された名前がパターンのいずれかに一致するか確認します
func isIgnored(name string, patterns []string) bool {
	for _, p := range patterns {
		// doublestarは完全なパスパターンを期待するため、
		// 単純なファイル名パターンの場合は特殊処理
		if strings.ContainsAny(p, "*?[]") {
			// ワイルドカードを含むパターン
			ok, _ := doublestar.Match(p, name)
			if ok {
				return true
			}
		} else {
			// 完全一致
			if p == name {
				return true
			}
		}
	}
	return false
}

// Gather は、指定されたパスから条件に一致するファイルのリストを収集します
func Gather(paths []string, opt Options) ([]string, error) {
	// 無視パターンの準備
	ignore := append([]string{}, opt.UserIgnorePatterns...)
	if opt.ApplyDefaultIgnores {
		ignore = append(ignore, defaultIgnore...)
	}

	var out []string
	for _, p := range paths {
		// 絶対パスに変換
		absPath, err := filepath.Abs(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: パス '%s' の解決中にエラー: %v。スキップします。\n", p, err)
			continue
		}

		// パスの存在確認
		info, err := os.Stat(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: 指定されたパス '%s' が見つかりません。スキップします。\n", p)
			continue
		}

		// ファイルの場合は直接追加
		if !info.IsDir() {
			// ドットファイルチェック
			name := filepath.Base(absPath)
			if !opt.IncludeDotfiles && len(name) > 0 && name[0] == '.' {
				continue
			}
			out = append(out, absPath)
			continue
		}

		// ディレクトリ自体がパターンに一致するかチェック
		dirName := filepath.Base(absPath)
		if !opt.IncludeDotfiles && len(dirName) > 0 && dirName[0] == '.' {
			fmt.Fprintf(os.Stderr, "無視 (ドット直接指定): %s\n", p)
			continue
		}

		// ディレクトリ自体がパターンに一致するかチェック
		if isIgnored(dirName, ignore) {
			fmt.Fprintf(os.Stderr, "無視 (パターン直接指定): %s\n", p)
			continue
		}

		// ディレクトリを再帰的に探索
		if err := filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "警告: '%s' へのアクセス中にエラー: %v。スキップします。\n", path, err)
				return nil // エラーを無視して続行
			}

			name := d.Name()

			// dotfile / dotdir
			if !opt.IncludeDotfiles && len(name) > 0 && name[0] == '.' {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// ignore pattern
			if d.IsDir() && isIgnored(name, ignore) {
				fmt.Fprintf(os.Stderr, "無視 (ディレクトリ): %s\n", path)
				return filepath.SkipDir
			}

			if !d.IsDir() {
				fmt.Fprintf(os.Stderr, "load %s\n", path)
				out = append(out, path)
			}
			return nil
		}); err != nil {
			fmt.Fprintf(os.Stderr, "警告: ディレクトリ '%s' の探索中にエラー: %v\n", absPath, err)
		}
	}
	return out, nil
}
