package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type WordCountOps struct {
	lines int
	bytes int
	words int
	chars int
}

func counter(fileContent []byte) *WordCountOps {
	lineCount := 0
	byteCount := len(fileContent)
	wordCount := len(bytes.Fields(fileContent))
	charCount := len(bytes.Runes(fileContent))

	for i := 0; i < len(fileContent); i++ {
		if string(fileContent[i]) == "\n" {
			lineCount++
		}
	}

	return &WordCountOps{
		lines: lineCount,
		bytes: byteCount,
		words: wordCount,
		chars: charCount,
	}
}

func evalWithFlags(cmd *cobra.Command, file []byte) string {
	flags := cmd.Flags()

	bytes := flags.Changed("bytes")
	lines := flags.Changed("lines")
	words := flags.Changed("words")
	characters := flags.Changed("multibytes")

	result := ""

	flagSet := false
	counts := counter(file)

	if bytes {
		flagSet = true
		result = fmt.Sprintf("%d\t", counts.bytes) + result
	}

	if lines {
		flagSet = true
		result = fmt.Sprintf("%d\t", counts.lines) + result
	}

	if words {
		flagSet = true
		result = fmt.Sprintf("%d\t", counts.words) + result
	}

	if characters {
		flagSet = true
		result = fmt.Sprintf("%d\t", counts.chars) + result
	}

	if !flagSet {
		result = fmt.Sprintf("%d\t%d\t%d\t", counts.bytes, counts.lines, counts.words) + result
	}

	return result
}

func wc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		file, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error opening file")
			os.Exit(1)
		}

		fmt.Print(evalWithFlags(cmd, file))
		os.Exit(0)
	}

	file, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}

	result := fmt.Sprintf("%s\n", args[0])

	result = evalWithFlags(cmd, file) + result

	fmt.Print(result)
	os.Exit(0)
}

var rootCmd = &cobra.Command{
	Use:   "ccwc",
	Short: "Coding Challenges implementation of wc",
	Run:   wc,
}

func init() {
	rootCmd.PersistentFlags().BoolP("bytes", "c", false, "count bytes")
	rootCmd.PersistentFlags().BoolP("lines", "l", false, "count lines")
	rootCmd.PersistentFlags().BoolP("words", "w", false, "count words")
	rootCmd.PersistentFlags().BoolP("multibytes", "m", false, "count characters")
}

func Execute() error {
	return rootCmd.Execute()
}
