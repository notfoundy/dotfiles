# Bootc image

Built on `quay.io/fedora/fedora-bootc:44`, which ships no desktop:

- `image/packages/*.txt`: official Fedora repositories plus 1Password's. No COPR.
- `image/files/etc/dnf/libdnf5.conf.d/90-excludes.conf`: packages dnf must never install, at
  build time and on the machine. A package that hard-requires one fails the build.
- Built from source, each by its own Containerfile, because Fedora 44 lacks them:
  Hyprland and its libraries (`image/hyprland/`), Ghostty (`image/ghostty/`), Noctalia Greeter
  (`image/noctalia-greeter/`). The image only copies their installed files and runtime package
  list: no build dependency reaches the system.

Upstream versions are `ARG`s in the builder Containerfiles, bumped by Renovate (Hyprland and its
libraries grouped). Runtime libraries are resolved in the builder (`ldd`, `rpm -qf`), so a bump
needs no package list change, and a missing library fails the build.

## CI

`.github/workflows/image.yml` publishes `ghcr.io/notfoundy/dotfiles` on an `img-v*` tag (see
[Releases](releases.md)) and rebuilds the latest one daily for Fedora updates. Pull requests
touching `image/` only build.

Builder images are tagged with the hash of their inputs, so the daily build skips Hyprland's half
hour. If a Fedora update breaks their ABI, run the workflow by hand with **rebuild builders**.

The stock update timer runs `bootc upgrade --apply`, which reboots on its own; it is overridden to
only stage.

## Building locally

```bash
podman build -f image/hyprland/Containerfile -t localhost/hyprland-build image
podman build -f image/ghostty/Containerfile -t localhost/ghostty-build image
podman build -f image/noctalia-greeter/Containerfile -t localhost/noctalia-greeter-build image
podman build -f image/Containerfile -t localhost/dotfiles image
```

A new Fedora release: bump `ARG FEDORA_VERSION` in the four Containerfiles.
