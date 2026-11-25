# Linux Go Utilities

A collection of small Linux command-line utilities written in Go. Build commands into binaries; Go source files are not directly executable scripts.

## Utilities

Run any command from this directory with `go run`:

```sh
go run ./cmd/sysinfo
go run ./cmd/diskusage [PATH]
go run ./cmd/psfind SEARCH_TERM
go run ./cmd/runutil -- COMMAND [ARG ...]
```

- `sysinfo` prints hostname, OS, kernel, architecture, and load average where available.
- `diskusage` reports total, used, and available space for a path.
- `psfind` searches readable process command lines under `/proc`.
- `runutil` runs another command directly, forwards standard streams, and propagates its exit status. It does not invoke a shell.

Build all utility binaries:

```sh
mkdir -p bin
for name in sysinfo diskusage psfind runutil; do
  go build -o "bin/$name" "./cmd/$name"
done
```
