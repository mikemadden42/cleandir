// Command cleandir sorts loose files in a directory into Documents/<ext>.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// skipExtensions lists extensions (lowercase, without the dot) that are never moved.
// "9" is kept from the original behavior; the reason for it isn't recorded.
var skipExtensions = map[string]bool{
	"9": true,
}

// errExists means the destination file is already there; the source is left in place.
var errExists = errors.New("destination already exists")

// moveFile moves src to dst without overwriting an existing dst.
// It hard-links then removes src, so an existing dst makes the link fail instead
// of being replaced. Filesystems without hard links fall back to check-then-rename.
func moveFile(src, dst string) error {
	err := os.Link(src, dst)
	if err == nil {
		if err := os.Remove(src); err != nil {
			_ = os.Remove(dst) // Undo the link so the file isn't left in both places
			return err
		}
		return nil
	}
	if errors.Is(err, fs.ErrExist) {
		return errExists
	}

	if _, err := os.Lstat(dst); err == nil {
		return errExists
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.Rename(src, dst)
}

func clean(dir string, dryRun bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
	}

	failed := 0
	for _, f := range files {
		// Skip dotfiles, directories, symlinks, pipes and sockets
		if !f.Type().IsRegular() || strings.HasPrefix(f.Name(), ".") {
			continue
		}
		extension := strings.ToLower(strings.TrimPrefix(filepath.Ext(f.Name()), "."))
		if extension == "" || skipExtensions[extension] {
			continue
		}

		destDir := filepath.Join(dir, "Documents", extension)
		srcPath := filepath.Join(dir, f.Name())
		destPath := filepath.Join(destDir, f.Name())

		if dryRun {
			if _, err := os.Lstat(destPath); err == nil {
				fmt.Println("File", srcPath, "already exists in", destDir)
			} else if !errors.Is(err, fs.ErrNotExist) {
				fmt.Fprintln(os.Stderr, "Error checking destination:", err)
				failed++
			} else {
				fmt.Println("Would move", srcPath, "to", destDir)
			}
			continue
		}

		if err := os.MkdirAll(destDir, 0750); err != nil { // Create directory if it doesn't exist
			fmt.Fprintln(os.Stderr, "Error creating directory:", err)
			failed++
			continue
		}

		err := moveFile(srcPath, destPath)
		switch {
		case errors.Is(err, errExists):
			fmt.Println("File", srcPath, "already exists in", destDir)
		case err != nil:
			fmt.Fprintln(os.Stderr, "Error moving file:", err)
			failed++
		default:
			fmt.Println("Moved", srcPath, "to", destDir)
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d file(s) could not be moved", failed)
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
