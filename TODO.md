# TODO

## Behavior (`main.go`)

- [x] **Run against the home directory, not the current directory.** `os.ReadDir(".")` contradicts the README ("Clean home directory"); running from any other directory moves its files into `./Documents/...`. Use `os.UserHomeDir()` or accept a path argument.
- [x] **Add a dry-run flag (`-n`).** There is currently no confirmation or undo.
- [x] **Only move regular files.** The `!f.IsDir()` filter lets symlinks, FIFOs and sockets through. Verified: a symlink to a directory was moved to `Documents/d/`, a FIFO to `Documents/fifo/`, and a relative symlink (`rel.txt -> ../elsewhere.txt`) was moved and broken. Use `if !f.Type().IsRegular() { continue }`.
- [ ] **Normalize extension case.** `b.PDF` and `c.pdf` both went into `Documents/PDF/` on macOS, while the output said `Documents/pdf`; on Linux they would be split into two folders. Use `strings.ToLower` for the folder name.
- [ ] **Skip the file when `MkdirAll` fails.** The error is printed, then the rename is attempted anyway. Add `continue`.
- [ ] **Fix the "already exists" check.** Any `Stat` error other than not-found (e.g. permission denied) is reported as "already exists". The check-then-rename is also racy, and `os.Rename` overwrites silently on Unix. Consider `os.Link` + `os.Remove`, which fails instead of overwriting.
- [ ] **Document or remove the `.9` exception.** `extension[1:] != "9"` has no explanation; make it a named skip list if it's intentional.
- [ ] **Report errors properly.** Write errors to stderr and exit non-zero when anything fails.
- [ ] **Flatten `clean`.** Use early `continue`s to cut the four levels of nesting.

## Tests

- [x] Refactor to `clean(dir string, dryRun bool) error` so it can be tested.
- [ ] Add tests using `t.TempDir()`: symlinks, mixed-case extensions, existing destination file, dotfiles, files without an extension.

## Tooling and config

- [ ] **Migrate `.golangci.yml` to the v2 format** (`golangci-lint migrate`). `config verify` currently rejects `disable-all`, top-level `linters-settings` and `issues.exclude-dirs`, and wants `version: "2"` as a string, so the revive/goimports settings are ignored. The v2 layout:
  - `linters.default: none`
  - `linters.settings.revive`
  - `formatters.settings.goimports.local-prefixes: [github.com/mikemadden42]`
  - `linters.exclusions.paths`
- [x] Add a package comment to `main.go` (the one current `golangci-lint run` finding).
- [ ] **Fix `.goreleaser.yaml`** so `goreleaser check` passes: add `version: 2`, rename `snapshot.name_template` to `version_template`.
- [ ] Add `format_overrides` so Windows archives ship as zip.
- [ ] Drop the no-op `go generate` hook.
- [ ] **Align Go versions:** `go.mod` says `go 1.21` but CI tests only 1.23–1.25. Add 1.21 to the matrix or raise `go.mod` to 1.23.
- [ ] **Run the linter in CI** with `golangci/golangci-lint-action`.

## Docs

- [ ] Expand the README: what the tool does, an example, and a warning that it moves files.
