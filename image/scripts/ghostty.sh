#!/usr/bin/env bash
set -euo pipefail

mkdir -p /src/ghostty /opt/zig
curl -fsSL "https://release.files.ghostty.org/$GHOSTTY_VERSION/ghostty-$GHOSTTY_VERSION.tar.gz" |
  tar -xz -C /src/ghostty --strip-components=1

zig_version=$(sed -n 's/.*\.minimum_zig_version = "\(.*\)".*/\1/p' /src/ghostty/build.zig.zon)
zig_tarball=$(curl -fsSL https://ziglang.org/download/index.json |
  jq -r --arg v "$zig_version" '.[$v]["x86_64-linux"] | "\(.tarball) \(.shasum)"')
read -r zig_url zig_sha256 <<<"$zig_tarball"
curl -fsSL -o /tmp/zig.tar.xz "$zig_url"
echo "$zig_sha256  /tmp/zig.tar.xz" | sha256sum -c -
tar -xJf /tmp/zig.tar.xz -C /opt/zig --strip-components=1

cd /src/ghostty
DESTDIR=/out ZIG_GLOBAL_CACHE_DIR=/tmp/zig-cache /opt/zig/zig build --prefix /usr -Doptimize=ReleaseFast -Dcpu=baseline -Dpie=true -Dversion-string="$GHOSTTY_VERSION"
rm -rf /out/usr/include /out/usr/lib/libghostty-vt* /out/usr/share/pkgconfig
