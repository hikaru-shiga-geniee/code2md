package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsIgnored(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		expected bool
	}{
		{"node_modules", []string{"node_modules"}, true},
		{"build", []string{"build*"}, true},
		{"build-temp", []string{"build*"}, true},
		{"src", []string{"node_modules", "build*"}, false},
		{"temp", []string{"temp"}, true},
		{"tmp", []string{"tmp"}, true},
		{"src", []string{"src"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIgnored(tt.name, tt.patterns)
			if result != tt.expected {
				t.Errorf("isIgnored(%q, %v) = %v, expected %v", tt.name, tt.patterns, result, tt.expected)
			}
		})
	}
}

func TestGatherWithDotfiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "code2md-test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用ファイル構造を作成
	files := map[string]string{
		"normal.txt":            "normal file",
		".dotfile":              "dot file",
		"subdir/file.txt":       "subdir file",
		"subdir/.hidden.txt":    "hidden file",
		".dotdir/file.txt":      "file in dotdir",
		"node_modules/file.txt": "node_modules file",
		"build/output.txt":      "build output",
		"dist/package.txt":      "dist package",
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("ディレクトリ作成に失敗: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("ファイル作成に失敗: %v", err)
		}
	}

	// テストケース
	tests := []struct {
		name             string
		opts             Options
		expectedCount    int
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name: "デフォルト設定",
			opts: Options{
				UserIgnorePatterns:  nil,
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 2, // normal.txt, subdir/file.txt
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "subdir/file.txt"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "subdir/.hidden.txt"),
				filepath.Join(tempDir, ".dotdir/file.txt"),
				filepath.Join(tempDir, "node_modules/file.txt"),
				filepath.Join(tempDir, "build/output.txt"),
				filepath.Join(tempDir, "dist/package.txt"),
			},
		},
		{
			name: "ドットファイル含む",
			opts: Options{
				UserIgnorePatterns:  nil,
				IncludeDotfiles:     true,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 5, // 通常ファイル + ドットファイル (デフォルト無視ディレクトリ以外)
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, "subdir/.hidden.txt"),
				filepath.Join(tempDir, ".dotdir/file.txt"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "node_modules/file.txt"),
				filepath.Join(tempDir, "build/output.txt"),
				filepath.Join(tempDir, "dist/package.txt"),
			},
		},
		{
			name: "デフォルト無視なし",
			opts: Options{
				UserIgnorePatterns:  nil,
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: false,
			},
			expectedCount: 5, // 通常ファイル + デフォルト無視ディレクトリのファイル (ドットファイル以外)
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, "node_modules/file.txt"),
				filepath.Join(tempDir, "build/output.txt"),
				filepath.Join(tempDir, "dist/package.txt"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "subdir/.hidden.txt"),
				filepath.Join(tempDir, ".dotdir/file.txt"),
			},
		},
		{
			name: "カスタム無視パターン",
			opts: Options{
				UserIgnorePatterns:  []string{"subdir"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 1, // normal.txt のみ
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "node_modules/file.txt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := Gather([]string{tempDir}, tt.opts)
			if err != nil {
				t.Fatalf("Gather() エラー: %v", err)
			}

			if len(files) != tt.expectedCount {
				t.Errorf("ファイル数 = %d, 期待値 %d", len(files), tt.expectedCount)
				t.Logf("見つかったファイル: %v", files)
			}

			// 含まれるべきファイルの確認
			for _, expected := range tt.shouldContain {
				found := false
				for _, file := range files {
					if file == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ファイル %s が結果に含まれていません", expected)
				}
			}

			// 含まれるべきでないファイルの確認
			for _, unexpected := range tt.shouldNotContain {
				for _, file := range files {
					if file == unexpected {
						t.Errorf("ファイル %s が結果に含まれていますが、含まれるべきではありません", unexpected)
						break
					}
				}
			}
		})
	}
}
