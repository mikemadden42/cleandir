// Command cleandir sorts loose files in a directory into Documents/<ext>.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func clean(dir string, dryRun bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
	}

	for _, f := range files {
		if f.Type().IsRegular() && !strings.HasPrefix(f.Name(), ".") { // Skip dotfiles, directories, symlinks, pipes and sockets
			extension := filepath.Ext(f.Name())
			if extension != "" && extension != "." && extension[1:] != "9" { // Check if extension is not empty before accessing its first character
				destDir := filepath.Join(dir, "Documents", extension[1:])

				if !dryRun {
					err := os.MkdirAll(destDir, 0750) // Create directory if it doesn't exist
					if err != nil {
						fmt.Println("Error creating directory:", err)
					}
				}

				srcPath := filepath.Join(dir, f.Name())
				destPath := filepath.Join(destDir, f.Name())

				// Check if the file already exists in the destination directory
				if _, err := os.Stat(destPath); os.IsNotExist(err) {
					if dryRun {
						fmt.Println("Would move", srcPath, "to", destDir)
						continue
					}
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
	return nil
}

func main() {
	dryRun := flag.Bool("n", false, "dry run: print what would be moved without moving anything")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-n] [dir]\n\nSorts files in dir (default: home directory) into dir/Documents/<ext>.\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var dir string
	switch flag.NArg() {
	case 0:
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error finding home directory:", err)
			os.Exit(1)
		}
		dir = home
	case 1:
		dir = flag.Arg(0)
	default:
		flag.Usage()
		os.Exit(2)
	}

	if err := clean(dir, *dryRun); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
