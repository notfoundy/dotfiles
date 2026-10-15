#!/usr/bin/env bash
set -euo pipefail

mkdir -p /deps
find /out -type f -print0 | xargs -0 file | awk -F: '/ELF/ { sub("^/out", "", $1); print $1 }' >/deps/binaries.txt

while read -r elf; do
  ldd "/out$elf" 2>/dev/null || true
done </deps/binaries.txt |
  awk '$2 == "=>" && $3 ~ /^\// { print $3 }' |
  sort -u |
  while read -r lib; do
    rpm -qf --qf '%{NAME}\n' "$lib" 2>/dev/null || true
  done |
  awk '!/not owned/' |
  sort -u >/deps/packages.txt
