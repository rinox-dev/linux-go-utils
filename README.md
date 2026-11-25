# Linux Go Utilities

A collection of 100 Linux monitoring commands written in Go, designed to build for both `amd64` and `riscv64`. Each command is an independent `main` package under `cmd/<name>` and reads Linux interfaces such as `/proc`, `/sys`, and standard files directly, with no architecture-specific assembly and no calls to existing system utilities.

## Build all 100 commands

Install Go, then run:

```sh
./scripts/build-all.sh
```

The script compiles one independent executable per utility, plus the `monitor` dispatcher, for both architectures:

```text
dist/linux-amd64/bin/
dist/linux-riscv64/bin/
```

Compressed release archives are also written to `dist/linux-amd64.tar.xz` and `dist/linux-riscv64.tar.xz`. Run an executable directly or use the dispatcher:

```sh
./dist/linux-amd64/bin/hostinfo
./dist/linux-amd64/bin/monitor list
./dist/linux-amd64/bin/monitor psfind go
```

For local development, use Go directly:

```sh
go run ./cmd/hostinfo
go run ./cmd/psfind go
go test ./...
```

Hardware sensors, systemd, cgroups, logs, and some kernel counters vary by device and permissions; unavailable sources are reported as such. Network probes, webhook delivery, alert policies, and full recursive scans are identified in the command descriptions as requiring additional configuration or integration.

Shared Linux parsing helpers live in `internal/sys`; the command catalog is in `internal/catalog`.
