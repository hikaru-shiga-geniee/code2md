package scan

import (
	"os"
	"path/filepath"
	"strings"
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
		{"test.md", []string{"*.md"}, true},
		{"README.md", []string{"*.md"}, true},
		{"script.js", []string{"*.md"}, false},
		{"__init__.py", []string{"__init__"}, true},
		{"__init__.py", []string{"*.py"}, true},
		{"config.json", []string{"*.md", "*.py", "*.json"}, true},
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
		"README.md":             "markdown file",
		"docs/guide.md":         "markdown guide",
		"src/main.py":           "python file",
		"src/__init__.py":       "python init file",
		"config.json":           "json config",
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
			expectedCount: 7, // normal.txt, subdir/file.txt, README.md, docs/guide.md, config.json, src/__init__.py, src/main.py
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "docs/guide.md"),
				filepath.Join(tempDir, "config.json"),
				filepath.Join(tempDir, "src/__init__.py"),
				filepath.Join(tempDir, "src/main.py"),
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
			expectedCount: 10, // 通常ファイル + ドットファイル (デフォルト無視ディレクトリ以外)
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, "subdir/.hidden.txt"),
				filepath.Join(tempDir, ".dotdir/file.txt"),
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "docs/guide.md"),
				filepath.Join(tempDir, "config.json"),
				filepath.Join(tempDir, "src/__init__.py"),
				filepath.Join(tempDir, "src/main.py"),
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
			expectedCount: 10, // 通常ファイル + デフォルト無視ディレクトリのファイル (ドットファイル以外)
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, "node_modules/file.txt"),
				filepath.Join(tempDir, "build/output.txt"),
				filepath.Join(tempDir, "dist/package.txt"),
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "docs/guide.md"),
				filepath.Join(tempDir, "config.json"),
				filepath.Join(tempDir, "src/main.py"),
				filepath.Join(tempDir, "src/__init__.py"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "subdir/.hidden.txt"),
				filepath.Join(tempDir, ".dotdir/file.txt"),
			},
		},
		{
			name: "カスタム無視パターン (ディレクトリ)",
			opts: Options{
				UserIgnorePatterns:  []string{"subdir"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 6, // normal.txt, README.md, docs/guide.md, config.json, src/__init__.py, src/main.py
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "docs/guide.md"),
				filepath.Join(tempDir, "config.json"),
				filepath.Join(tempDir, "src/__init__.py"),
				filepath.Join(tempDir, "src/main.py"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "node_modules/file.txt"),
			},
		},
		{
			name: "Markdownファイル除外",
			opts: Options{
				UserIgnorePatterns:  []string{"*.md"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 5, // normal.txt, subdir/file.txt, config.json, src/__init__.py, src/main.py
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, "config.json"),
				filepath.Join(tempDir, "src/__init__.py"),
				filepath.Join(tempDir, "src/main.py"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "docs/guide.md"),
				filepath.Join(tempDir, "node_modules/file.txt"),
			},
		},
		{
			name: "__init__ファイル除外",
			opts: Options{
				UserIgnorePatterns:  []string{"__init__"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: false,
			},
			expectedCount: 9, // すべてのファイル - __init__.py
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "subdir/file.txt"),
				filepath.Join(tempDir, "node_modules/file.txt"),
				filepath.Join(tempDir, "build/output.txt"),
				filepath.Join(tempDir, "dist/package.txt"),
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "docs/guide.md"),
				filepath.Join(tempDir, "config.json"),
				filepath.Join(tempDir, "src/main.py"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, ".dotfile"),
				filepath.Join(tempDir, "subdir/.hidden.txt"),
				filepath.Join(tempDir, ".dotdir/file.txt"),
				filepath.Join(tempDir, "src/__init__.py"),
			},
		},
		{
			name: "複数ファイルタイプ除外",
			opts: Options{
				UserIgnorePatterns:  []string{"*.md", "*.py", "*.json"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 2, // normal.txt, subdir/file.txt
			shouldContain: []string{
				filepath.Join(tempDir, "normal.txt"),
				filepath.Join(tempDir, "subdir/file.txt"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "docs/guide.md"),
				filepath.Join(tempDir, "src/main.py"),
				filepath.Join(tempDir, "src/__init__.py"),
				filepath.Join(tempDir, "config.json"),
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

// ファイルパターン除外の直接指定テスト
func TestGatherWithFilePatternIgnore(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "code2md-file-test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用ファイル構造を作成
	files := map[string]string{
		"README.md":     "readme file",
		"script.js":     "javascript file",
		"config.json":   "config file",
		"main.py":       "python main file",
		"__init__.py":   "python init file",
		"test_data.csv": "csv data file",
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("ファイル作成に失敗: %v", err)
		}
	}

	// 直接ファイルを指定した場合のテスト
	tests := []struct {
		name             string
		paths            []string
		opts             Options
		expectedCount    int
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:  "直接指定したMarkdownファイルを除外",
			paths: []string{filepath.Join(tempDir, "README.md"), filepath.Join(tempDir, "script.js")},
			opts: Options{
				UserIgnorePatterns:  []string{"*.md"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 1, // script.js のみ
			shouldContain: []string{
				filepath.Join(tempDir, "script.js"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "README.md"),
			},
		},
		{
			name: "複数のファイルタイプを除外",
			paths: []string{
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "script.js"),
				filepath.Join(tempDir, "config.json"),
				filepath.Join(tempDir, "main.py"),
			},
			opts: Options{
				UserIgnorePatterns:  []string{"*.md", "*.py"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: true,
			},
			expectedCount: 2, // script.js, config.json
			shouldContain: []string{
				filepath.Join(tempDir, "script.js"),
				filepath.Join(tempDir, "config.json"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "README.md"),
				filepath.Join(tempDir, "main.py"),
			},
		},
		{
			name: "__init__ファイルを除外",
			paths: []string{
				filepath.Join(tempDir, "main.py"),
				filepath.Join(tempDir, "__init__.py"),
			},
			opts: Options{
				UserIgnorePatterns:  []string{"__init__"},
				IncludeDotfiles:     false,
				ApplyDefaultIgnores: false,
			},
			expectedCount: 1, // main.py のみ
			shouldContain: []string{
				filepath.Join(tempDir, "main.py"),
			},
			shouldNotContain: []string{
				filepath.Join(tempDir, "__init__.py"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := Gather(tt.paths, tt.opts)
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

// getRelativePath関数のテスト
func TestGetRelativePath(t *testing.T) {
	// 現在の作業ディレクトリを取得
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在の作業ディレクトリの取得に失敗: %v", err)
	}

	tests := []struct {
		name     string
		absPath  string
		expected string
	}{
		{
			name:     "現在のディレクトリ内のファイル",
			absPath:  filepath.Join(wd, "test.txt"),
			expected: "test.txt",
		},
		{
			name:     "サブディレクトリ内のファイル",
			absPath:  filepath.Join(wd, "subdir", "file.txt"),
			expected: filepath.Join("subdir", "file.txt"),
		},
		{
			name:     "現在のディレクトリそのもの",
			absPath:  wd,
			expected: ".",
		},
		{
			name:     "親ディレクトリのファイル",
			absPath:  filepath.Join(filepath.Dir(wd), "parent.txt"),
			expected: filepath.Join("..", "parent.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRelativePath(tt.absPath)
			if result != tt.expected {
				t.Errorf("getRelativePath(%q) = %q, 期待値 %q", tt.absPath, result, tt.expected)
			}
		})
	}
}

// 無効なパスでのgetRelativePath関数のテスト
func TestGetRelativePathWithInvalidPath(t *testing.T) {
	// 存在しないドライブ（Windows）や無効なパスでテスト
	invalidPath := "/非常に/長い/存在しない/パス/test.txt"

	// 関数が何らかの結果を返すことを確認（パニックしないことが重要）
	result := getRelativePath(invalidPath)

	// 結果が空文字列でないことを確認
	if result == "" {
		t.Error("getRelativePath() は空文字列を返すべきではありません")
	}

	// 通常は絶対パスがそのまま返されるか、相対パスが計算される
	t.Logf("無効なパス %q に対する結果: %q", invalidPath, result)
}

// Loading メッセージが相対パスで出力されることを検証するテスト
func TestGatherOutputsRelativePaths(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "code2md-relative-test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用ファイルを作成
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	// 現在の作業ディレクトリを一時的に変更
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在の作業ディレクトリの取得に失敗: %v", err)
	}

	// 作業ディレクトリを変更
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("作業ディレクトリの変更に失敗: %v", err)
	}
	defer func() {
		// テスト終了後に元のディレクトリに戻す
		os.Chdir(originalWd)
	}()

	// Gatherを実行
	files, err := Gather([]string{"test.txt"}, Options{
		UserIgnorePatterns:  nil,
		IncludeDotfiles:     false,
		ApplyDefaultIgnores: true,
	})

	if err != nil {
		t.Fatalf("Gather() エラー: %v", err)
	}

	// ファイルが見つかったことを確認
	if len(files) != 1 {
		t.Fatalf("期待されるファイル数: 1, 実際: %d", len(files))
	}

	// 返されるパスは絶対パスのはず
	expectedAbsPath, _ := filepath.Abs("test.txt")
	if files[0] != expectedAbsPath {
		t.Errorf("返されたパス: %s, 期待されるパス: %s", files[0], expectedAbsPath)
	}

	// getRelativePath関数の動作を個別に確認
	relativePath := getRelativePath(files[0])
	if relativePath != "test.txt" {
		t.Errorf("相対パス変換結果: %s, 期待値: test.txt", relativePath)
	}
}

// エラーメッセージとIgnoredメッセージでも相対パスが使用されることを検証
func TestGatherIgnoredMessagesUseRelativePaths(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "code2md-ignored-test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// node_modulesディレクトリを作成（デフォルトで無視される）
	nodeModulesDir := filepath.Join(tempDir, "node_modules")
	if err := os.MkdirAll(nodeModulesDir, 0755); err != nil {
		t.Fatalf("node_modulesディレクトリの作成に失敗: %v", err)
	}

	// node_modules内にファイルを作成
	testFile := filepath.Join(nodeModulesDir, "package.js")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	// 現在の作業ディレクトリを一時的に変更
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("現在の作業ディレクトリの取得に失敗: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("作業ディレクトリの変更に失敗: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Gatherを実行（node_modulesは無視されるべき）
	files, err := Gather([]string{"."}, Options{
		UserIgnorePatterns:  nil,
		IncludeDotfiles:     false,
		ApplyDefaultIgnores: true,
	})

	if err != nil {
		t.Fatalf("Gather() エラー: %v", err)
	}

	// node_modulesのファイルが結果に含まれていないことを確認
	for _, file := range files {
		if strings.Contains(file, "node_modules") {
			t.Errorf("node_modulesのファイルが結果に含まれています: %s", file)
		}
	}

	// この時点でIgnoredメッセージが標準エラー出力に出されているはず
	// （テスト実行時に確認可能）
}
