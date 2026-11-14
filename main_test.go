package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func TestCleanMovesRegularFiles(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "a.txt"))

	if err := clean(dir, false); err != nil {
		t.Fatal(err)
	}

	if exists(filepath.Join(dir, "a.txt")) {
		t.Error("a.txt still in source directory")
	}
	if !exists(filepath.Join(dir, "Documents", "txt", "a.txt")) {
		t.Error("a.txt not moved to Documents/txt")
	}
}

func TestCleanSkipsDotfilesAndNoExtension(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{".bashrc", "noext", "trailing.", "notes.9"} {
		touch(t, filepath.Join(dir, name))
	}

	if err := clean(dir, false); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{".bashrc", "noext", "trailing.", "notes.9"} {
		if !exists(filepath.Join(dir, name)) {
			t.Errorf("%s was moved", name)
		}
	}
}

func TestCleanSkipsSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need extra privileges on Windows")
	}
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "realdir"), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("realdir", filepath.Join(dir, "link.d")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../elsewhere.txt", filepath.Join(dir, "rel.txt")); err != nil {
		t.Fatal(err)
	}

	if err := clean(dir, false); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"link.d", "rel.txt"} {
		if !exists(filepath.Join(dir, name)) {
			t.Errorf("symlink %s was moved", name)
		}
	}
}

func TestCleanDryRunMovesNothing(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, "a.txt"))

	if err := clean(dir, true); err != nil {
		t.Fatal(err)
	}

	if !exists(filepath.Join(dir, "a.txt")) {
		t.Error("a.txt was moved during dry run")
	}
	if exists(filepath.Join(dir, "Documents")) {
		t.Error("Documents created during dry run")
	}
}

func TestCleanKeepsExistingDestination(t *testing.T) {
	dir := t.TempDir()
	destDir := filepath.Join(dir, "Documents", "txt")
	if err := os.MkdirAll(destDir, 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "a.txt"), []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("new"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := clean(dir, false); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(filepath.Join(destDir, "a.txt")) //nolint:gosec // path is under t.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "old" {
		t.Errorf("destination overwritten: got %q", got)
	}
	if !exists(filepath.Join(dir, "a.txt")) {
		t.Error("source removed despite existing destination")
	}
}

func TestCleanMissingDirectory(t *testing.T) {
	if err := clean(filepath.Join(t.TempDir(), "missing"), false); err == nil {
		t.Error("expected error for missing directory")
	}
}
