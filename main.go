package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func clean() {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, f := range files {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") { // Check for dotfiles
			extension := filepath.Ext(f.Name())[1:]
			if extension != "" && extension != "9" {
				destDir := filepath.Join("Documents", extension)
				err := os.MkdirAll(destDir, os.ModePerm) // Create directory if it doesn't exist
				if err != nil {
					fmt.Println("Error creating directory:", err)
				}

				srcPath := f.Name()
				destPath := filepath.Join(destDir, f.Name())

				err = os.Rename(srcPath, destPath)
				if err != nil {
					fmt.Println("Error moving file:", err)
				} else {
					fmt.Println("Moved", srcPath, "to", destDir)
				}
			}
		}
	}
}

func main() {
	clean()
}
