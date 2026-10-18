#!/usr/bin/env bash
set -euo pipefail

name=$1
url=$2
target=$3
shift 3

src=/src/$name
mkdir -p "$src"
curl -fsSL "$url" | tar -xz -C "$src" --strip-components=1

if [[ -f "$src/CMakeLists.txt" ]]; then
  cmake -S "$src" -B "$src/build" -G Ninja \
    -DCMAKE_INSTALL_PREFIX=/usr -DCMAKE_INSTALL_LIBDIR=lib64 -DCMAKE_BUILD_TYPE=Release \
    -DBUILD_TESTING=OFF "$@"
  cmake --build "$src/build"
  install_into() { DESTDIR=$1 cmake --install "$src/build"; }
else
  meson setup --prefix=/usr --libdir=lib64 --buildtype=release "$src/build" "$src" "$@"
  meson compile -C "$src/build"
  install_into() { DESTDIR=$1 meson install -C "$src/build"; }
fi

install_into /
if [[ "$target" == image ]]; then
  install_into /out
fi
rm -rf "$src"
