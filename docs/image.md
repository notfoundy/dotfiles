# Bootc image

The system is a [bootc](https://bootc-dev.github.io/bootc/) image built from
`quay.io/fedora/fedora-bootc:44`, which ships no desktop. `image/Containerfile` adds only what
is used:

- `image/packages/` — packages from the official Fedora repositories, plus 1Password from its
  own repository, one file per category (`desktop.txt`, `hardware.txt`, `cli.txt`, `dev.txt`…).
  No COPR.
- **Hyprland** and its libraries, `xdg-desktop-portal-hyprland`, `hyprpicker`,
  `hyprland-guiutils` (the "app not responding" and permission dialogs, which Hyprland warns
  about when missing) and `hyprpolkitagent`, both built on `hyprtoolkit`, are not in Fedora 44, and neither **Ghostty** nor **Noctalia Greeter** ever was: each is built from
  source by its own Containerfile (`image/hyprland/`, `image/ghostty/`,
  `image/noctalia-greeter/`). Those builder images keep only the installed files
  (`/out`) and their runtime package list (`/deps`), which the image copies in: no build
  dependency reaches the system.

Upstream versions are pinned by `ARG`s in the builder Containerfiles, one `RUN` per project in
dependency order, so bumping one only rebuilds it and what depends on it. Renovate opens a pull
request for each release, grouping Hyprland and its libraries since they move together. The
runtime libraries of what is built from source are resolved in the builder (`ldd`, then
`rpm -qf`), so a version bump needs no package list update, and the image build fails if a
library is still missing.

`.github/workflows/image.yml` publishes `ghcr.io/notfoundy/dotfiles` when an `img-v*` tag is
pushed (see [Releases](releases.md)), and rebuilds the latest `img-v*` tag daily for Fedora
updates. Merging to `main`, Renovate pull requests included, publishes nothing until the next
tag; pull requests touching `image/` only build. A builder image is tagged with the hash of its inputs and only
rebuilt when they change, so the daily build takes minutes, not the half hour Hyprland needs.
If a Fedora update ever breaks the ABI the builders were compiled against, the library check
fails the build: run the workflow by hand with **rebuild builders** checked. The same daily
run re-enables its own schedule, which GitHub would otherwise turn off after 60 days without
activity.

The machine follows with `bootc upgrade`. The stock update timer runs
`bootc upgrade --apply`, which reboots on its own; it is overridden to only stage the update,
applied at the next reboot. `bootc rollback` returns to the previous image.

## Building locally

Builders first, then the image:

```bash
podman build -f image/hyprland/Containerfile -t localhost/hyprland-build image
podman build -f image/ghostty/Containerfile -t localhost/ghostty-build image
podman build -f image/noctalia-greeter/Containerfile -t localhost/noctalia-greeter-build image
podman build -f image/Containerfile -t localhost/dotfiles image
```

Moving to a new Fedora release means bumping `ARG FEDORA_VERSION` in the four Containerfiles.
