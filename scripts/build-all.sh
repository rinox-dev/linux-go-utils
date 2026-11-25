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
    [[ -n "$name" ]] || continue
    [[ ! -L "$out/$name" ]] || unlink "$out/$name"
    cp "$out/monitor" "$out/$name"
    chmod 0755 "$out/$name"
  done < <("$go_cmd" run ./cmd/monitor list)
  count=$(find "$out" -maxdepth 1 -type f ! -name monitor | wc -l)
  [[ "$count" -eq 100 ]] || { echo "Expected 100 utilities for $target, found $count" >&2; exit 1; }
  tar -C dist -cf - "linux-$target" | xz -T0 -9e > "dist/linux-$target.tar.xz"
  echo "Packaged $count independent executables in dist/linux-$target.tar.xz"
done

echo "Built 100 independent utilities for Linux/amd64 and Linux/riscv64 under dist/."
