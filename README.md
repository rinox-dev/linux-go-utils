# Linux Go Utilities

A collection of Linux monitoring commands written in Go, designed to build for both `amd64` and `riscv64`. The collector uses Linux interfaces such as `/proc`, `/sys`, and standard files, with no architecture-specific assembly or third-party dependencies.

## Generate all 100 commands

Install Go, then run:

```sh
./scripts/build-all.sh
```

The script creates architecture-specific binaries and command-name symlinks:

```text
dist/linux-amd64/bin/
dist/linux-riscv64/bin/
```

Each command is a symlink to the shared `monitor` binary. The invoked command name selects the utility. You can run a command directly or use the dispatcher:

```sh
./dist/linux-amd64/bin/hostinfo
./dist/linux-amd64/bin/monitor list
./dist/linux-amd64/bin/monitor psfind go
```

For local development, use Go directly:

```sh
go run ./cmd/monitor list
go run ./cmd/monitor meminfo
go run ./cmd/monitor psfind go
```

The 100 names and their data sources are listed in [`cmd/monitor/main.go`](cmd/monitor/main.go). Collectors read the interfaces available on the running machine. Hardware sensors, systemd, cgroups, logs, and some kernel counters vary by device and permissions; unavailable sources are reported as such. Network probes, webhook delivery, alert policies, and full recursive scans are identified in the command descriptions as requiring additional configuration or integration.

Existing starter commands remain available under `cmd/runutil`, `cmd/sysinfo`, `cmd/diskusage`, and `cmd/psfind`.
