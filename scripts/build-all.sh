#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"
go_cmd="${GO:-go}"

for target in amd64 riscv64; do
  out="dist/linux-$target/bin"
  mkdir -p "$out"
  echo "Building Linux/$target utilities..."
  CGO_ENABLED=0 GOOS=linux GOARCH="$target" "$go_cmd" build -trimpath -o "$out/monitor" ./cmd/monitor
  while read -r name _; do
    [[ -n "$name" ]] && ln -sf monitor "$out/$name"
  done < <("$go_cmd" run ./cmd/monitor list)
done

echo "Built 100 command names for Linux/amd64 and Linux/riscv64 under dist/."
