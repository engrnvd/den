package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func printSize(dir string) {
	var size int64 = 0
	children := make([]string, 0)
	dir, _ = filepath.Abs(dir)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			printError(err)
			return err
		}

		size += info.Size()

		shortPath := strings.Replace(path, dir, "", 1)

		if path != dir && !match("/", shortPath) {
			_ = append(children, path)
		}
		return nil
	})

	if err != nil {
		printError(err)
	}

	fmt.Printf("%v: %v\n", dir, humanReadableSize(size))

	for _, p := range children {
		printSize(p)
	}
}

func rootCmdHandler(cmd *cobra.Command, args []string) {
	dir := args[0]
	if dir == "" {
		dir = "."
	}

	err := printDirectorySizes(dir)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
