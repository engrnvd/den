package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func printDirectorySizes(path string) error {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening directory: %v", err)
	}
	defer dir.Close()

	// Calculate total parent directory size
	var totalParentSize int64
	parentInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting parent directory info: %v", err)
	}
	totalParentSize = parentInfo.Size()

	// Read directory entries
	entries, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("error reading directory contents: %v", err)
	}

	fmt.Printf("Parent Directory: %s (Total Size: %d bytes)\n", path, totalParentSize)
	fmt.Println("Child Directories:")

	// Iterate through entries and print directory details
	for _, entry := range entries {
		if entry.IsDir() {
			childPath := filepath.Join(path, entry.Name())
			size, err := calculateDirectorySize(childPath)
			if err != nil {
				fmt.Printf("Error calculating size for %s: %v\n", entry.Name(), err)
				continue
			}
			fmt.Printf("%s: %s\n", entry.Name(), humanReadableSize(size))
		}
	}

	return nil
}

func calculateDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func match(regex string, str string) bool {
	m, err := regexp.MatchString(regex, str)
	if err != nil {
		printError(err)
	}

	return !!m
}

func printError(err error) {
	log.Fatal(err)
}

func humanReadableSize(bytes int64) string {
	const (
		KB = 1 << 10 // 1024 bytes
		MB = 1 << 20 // 1024 * 1024 bytes
		GB = 1 << 30 // 1024 * 1024 * 1024 bytes
		TB = 1 << 40 // 1024 * 1024 * 1024 * 1024 bytes
		PB = 1 << 50 // 1024 * 1024 * 1024 * 1024 * 1024 bytes
	)

	switch {
	case bytes >= PB:
		return fmt.Sprintf("%.2f PB", float64(bytes)/PB)
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
