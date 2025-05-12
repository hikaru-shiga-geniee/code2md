package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/code2md/internal/markdown"
	"github.com/your-org/code2md/internal/scan"
)

var (
	ignorePatterns   []string
	includeDotfiles  bool
	noDefaultIgnores bool
)

func main() {
	root := &cobra.Command{
		Use:   "code2md [paths...]",
		Short: "指定されたファイルやディレクトリの内容をMarkdownコードブロック形式で出力します",
		Long: `code2md は、指定されたファイルやディレクトリの内容を読み込み、
Markdownのコードブロック形式で標準出力するコマンドラインツールです。
特定のディレクトリパターンや、ドットから始まるファイル/ディレクトリを
デフォルトで除外する機能を備えています。`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := scan.Options{
				UserIgnorePatterns:  ignorePatterns,
				IncludeDotfiles:     includeDotfiles,
				ApplyDefaultIgnores: !noDefaultIgnores,
			}
			files, err := scan.Gather(args, opts)
			if err != nil {
				return err
			}
			return markdown.Print(os.Stdout, files)
		},
	}

	root.Flags().StringSliceVarP(&ignorePatterns, "ignore", "i", nil,
		"無視する **ディレクトリ名** のパターン (ワイルドカード可)")
	root.Flags().BoolVar(&includeDotfiles, "include-dotfiles", false,
		"'.'で始まるファイルやディレクトリを処理対象に含める")
	root.Flags().BoolVar(&noDefaultIgnores, "no-default-ignores", false,
		"デフォルトの無視ディレクトリパターンを適用しない")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
