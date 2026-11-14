# cleandir

Clean home directory: sorts loose files into `Documents/<extension>/`.

> **Warning:** cleandir moves files. Run it with `-n` first to see what it would do.

## Usage

```
cleandir [-n] [dir]
```

- `dir`: the directory to clean. Defaults to your home directory.
- `-n`: dry run. Prints what would be moved without moving anything.

For example, in a home directory containing `report.PDF` and `notes.txt`:

```
$ cleandir -n
Would move /Users/me/notes.txt to /Users/me/Documents/txt
Would move /Users/me/report.PDF to /Users/me/Documents/pdf
```

## What it moves

Only regular files with an extension, directly inside `dir`. It skips:

- dotfiles, directories, symlinks, pipes and sockets
- files without an extension
- `.9` files

Extensions are lowercased, so `report.PDF` and `scan.pdf` share `Documents/pdf/`.
A file that already exists at the destination is never overwritten; it stays where it is.

cleandir prints errors to stderr and exits with status 1 if any file can't be moved.

## Install

```
go install github.com/mikemadden42/cleandir@latest
```
