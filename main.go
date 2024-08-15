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
			extension := filepath.Ext(f.Name())
			if extension != "" && extension != "." && extension[1:] != "9" { // Check if extension is not empty before accessing its first character
				destDir := filepath.Join("Documents", extension[1:])

				err := os.MkdirAll(destDir, os.ModePerm) // Create directory if it doesn't exist
				if err != nil {
					fmt.Println("Error creating directory:", err)
				}

				srcPath := f.Name()
				destPath := filepath.Join(destDir, f.Name())

				// Check if the file already exists in the destination directory
				if _, err := os.Stat(destPath); os.IsNotExist(err) {
					err = os.Rename(srcPath, destPath)
					if err != nil {
						fmt.Println("Error moving file:", err)
					} else {
						fmt.Println("Moved", srcPath, "to", destDir)
					}
				} else {
					fmt.Println("File", srcPath, "already exists in", destDir)
				}
			}
		}
	}
}

func main() {
	clean()
}
