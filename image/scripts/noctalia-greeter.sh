#!/usr/bin/env bash
set -euo pipefail

mkdir -p /src/noctalia-greeter
curl -fsSL "https://github.com/noctalia-dev/noctalia-greeter/archive/refs/tags/$NOCTALIA_GREETER_VERSION.tar.gz" |
  tar -xz -C /src/noctalia-greeter --strip-components=1

cd /src/noctalia-greeter
meson setup --prefix=/usr --buildtype=plain -Doptimization=2 build
meson compile -C build
DESTDIR=/out meson install -C build
