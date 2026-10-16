# Dotfiles

Welcome to the messy world of my **dotfiles**!

This is where all my configuration files for my working environments live.
In short, you’ll find configuration files for the terminal, text editor, window manager and other things to make it fun to use for **me**.

Everything is managed with **Ansible**, and yes, it may seem like overkill.
Guess what? It really is!

This repo is heavily influenced by [TechDufus](https://github.com/TechDufus/dotfiles)'s repo. Go check it out!

## Why Ansible?

- I wanted to learn Ansible.
- I wanted something fully automated.

It turns out I really enjoy using it, and this process has been kind of fun
At the moment, it’s probably not fully automated, but I’ll live with it. Maybe it will be someday... who knows?

## Usage

> [!WARNING]
> Only the Fedora bootc image built from `image/` is supported: it carries every system
> package, and the playbook only manages `$HOME`.

Everything runs through the `dotfiles` command, a small Go program living in `cli/`.
It checks the host, updates the Ansible dependencies, refreshes the checkout in
`~/.dotfiles`, then hands the terminal over to `ansible-playbook`. Nothing needs `sudo`.
Ansible will take care of the rest (or almost). You just need to grab a coffee while it works!

### On a fresh installation

Drop the latest release somewhere on your `PATH` and let it do the rest:

```bash
mkdir -p ~/.local/bin
curl -fsSL https://github.com/notfoundy/dotfiles/releases/latest/download/dotfiles_linux_amd64.tar.gz \
  | tar -xz -C ~/.local/bin dotfiles
~/.local/bin/dotfiles
```

Swap `amd64` for `arm64` if that is what you are running. The first run clones the
repository to `~/.dotfiles`; log out once it is done and sign in from the greeter to start
Hyprland.

To build it from a clone instead:

```bash
go build -C cli -o ~/.local/bin/dotfiles .
```

### Otherwise

```bash
dotfiles                        # every default role, as listed by default_roles in group_vars/all.yml
dotfiles -t nvim,go             # only those roles
dotfiles --select               # pick them interactively
dotfiles --skip-update          # leave the checkout alone
dotfiles --repo . -- --check    # run against this clone; anything after -- goes to ansible-playbook
```

`dotfiles --help` lists the rest, and `dotfiles self-update` fetches the latest
release. Every run is appended to `~/.local/state/dotfiles/dotfiles.log`, so a role
that failed can be read back afterwards.

## Bootc image

The system is a [bootc](https://bootc-dev.github.io/bootc/) image built from
`quay.io/fedora/fedora-bootc:44`, which ships no desktop. `image/Containerfile` adds only what
is used:

- `image/packages/` — packages from the official Fedora repositories, plus 1Password from its
  own repository, one file per category (`desktop.txt`, `hardware.txt`, `cli.txt`, `dev.txt`…).
  No COPR.
- **Hyprland** and its libraries, `xdg-desktop-portal-hyprland` and `hyprpicker` are not in
  Fedora 44, and neither **Ghostty** nor **Noctalia Greeter** ever was: each is built from
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

`.github/workflows/image.yml` builds on every change and daily for Fedora updates, and publishes
`ghcr.io/notfoundy/dotfiles`. A builder image is tagged with the hash of its inputs and only
rebuilt when they change, so the daily build takes minutes, not the half hour Hyprland needs.
If a Fedora update ever breaks the ABI the builders were compiled against, the library check
fails the build: run the workflow by hand with **rebuild builders** checked. The same daily
run re-enables its own schedule, which GitHub would otherwise turn off after 60 days without
activity.

The machine follows with `bootc upgrade`. The stock update timer runs
`bootc upgrade --apply`, which reboots on its own; it is overridden to only stage the update,
applied at the next reboot. `bootc rollback` returns to the previous image.

A few things differ from a regular Fedora:

- greetd is the display manager, with [Noctalia Greeter](https://docs.noctalia.dev/greeter/)
  as its greeter. It only offers the `Hyprland` session, which runs `start-hyprland`: Hyprland
  is built without uwsm. The greeter runs as Fedora's `greetd` account, so the tmpfiles rule
  creating its state directory is overridden in `/etc/tmpfiles.d/`, as upstream expects.
  `greetd` gets a fixed UID (903, `image/files/usr/lib/sysusers.d/greetd-id.conf`) so that
  `/var/lib/greetd` keeps its owner across images. `dbus-daemon` ships only for
  `dbus-run-session`, which the greeter session needs; the system bus stays on dbus-broker.
  Syncing Noctalia's wallpaper and palette to the greeter is not set up.
- Noctalia is the polkit agent.
- Swap is zram only, through `zram-generator-defaults`: the base image ships none, and
  without it the system has no swap at all. `systemd-oomd-defaults` kills the offending
  cgroup under sustained memory pressure, before the desktop freezes.
- The base image enables sshd with password login and ships no firewall. sshd is disabled
  (`systemctl enable --now sshd` when needed) and firewalld runs with Fedora's `public`
  zone. libvirt and podman (netavark) both use nftables and keep their own tables, and
  firewalld ships the `libvirt` zones, so VMs and containers keep their networking.
- The base image only knows `C.UTF-8`. `glibc-langpack-en` and `-fr` are added, and
  `/etc/locale.conf` sets `en_US.UTF-8` with French formats (time, currency, paper, units).
  Numbers stay English: a decimal comma breaks `printf '%.2f'` and friends in scripts.
- Anaconda writes `compress=zstd:1` to fstab, but bootc mounts the root from the kernel
  command line, and btrfs takes its options from that first mount. Compression is set back
  through `rootflags` in `/usr/lib/bootc/kargs.d/`.
- The base image ships no `audit`: without auditd, the kernel prints its audit events on the
  console, over the TTY. `audit` is added and `auditd` enabled so they go to its log.
- libvirt is managed by the `wheel` group through a polkit rule: `usermod -aG libvirt` does not
  work, since the group lives in `/usr/lib/group`.
- bootc only copies the image's `/var` at install time. The directories packages leave there
  without a tmpfiles rule (libvirt, flatpak, AccountsService…) are listed by
  `bootc container lint` and written to `/usr/lib/tmpfiles.d/dotfiles-var.conf` at the end of
  the build, so they are recreated if missing, as `rpm-ostree-autovar.conf` does for the base
  image.
- `/` is the composefs image, read-only and always 100 % full. Disk usage is read on `/var`,
  the btrfs partition that also holds `/var/home`: `df -h /var`, and the fastfetch role points
  its disk module there.
- Go, Node, Typst, Starship and Lazygit are installed in `~/.local` by their roles, so they
  stay on the latest release without touching the image. Starship goes through its official
  install script, Lazygit is checked against the release's `checksums.txt`.

To build it locally, builders first:

```bash
podman build -f image/hyprland/Containerfile -t localhost/hyprland-build image
podman build -f image/ghostty/Containerfile -t localhost/ghostty-build image
podman build -f image/noctalia-greeter/Containerfile -t localhost/noctalia-greeter-build image
podman build -f image/Containerfile -t localhost/dotfiles image
```

Moving to a new Fedora release means bumping `ARG FEDORA_VERSION` in the four Containerfiles.

### Installing

The machine is installed by Anaconda from the Fedora Everything network ISO, which pulls
`ghcr.io/notfoundy/dotfiles:latest` through the `ostreecontainer` command of
`image/install.ks`. The package must be public, or the kickstart needs a `%pre` writing
`/etc/ostree/auth.json`. Embed the kickstart in the ISO, so that every boot entry loads it:

```bash
podman run --rm -v "$PWD":/work:z -w /work registry.fedoraproject.org/fedora:44 sh -c \
  'dnf -y install lorax && mkksiso --ks image/install.ks Fedora-Everything-netinst-x86_64-44-*.iso dotfiles-install.iso'
```

Write `dotfiles-install.iso` to a USB stick and boot it; Secure Boot can stay on. The
kickstart only holds what never changes, Anaconda stops on its summary for the rest:

- **Network**: the Wi-Fi, copied to the installed system.
- **Installation destination**: the disk, automatic partitioning with btrfs, **Encrypt my
  data** checked, then the LUKS passphrase.
- **User creation**: with **Make this user administrator** checked (`wheel`). root stays
  locked.

A few kickstart lines undo what Anaconda gets wrong on a bootc system:

- `firewall --use-system-defaults`: otherwise Anaconda runs `firewall-offline-cmd` and opens
  ssh in the default zone.
- Anaconda always rewrites `/etc/locale.conf` with a lone `LANG=`, with or without a `lang`
  directive. The `%post` copies back the image's version from `/usr/etc`: identical to the
  default, the file is no longer a local change and follows the image again.
- The `%post` also drops the `/` entry Anaconda writes to `/etc/fstab`: with composefs, `/`
  is an overlay and `systemd-remount-fs` fails applying btrfs options to it. The root is
  mounted from the kernel arguments anyway, as `bootc install` does.

## Desktop shell

The image installs [Noctalia](https://docs.noctalia.dev) from the official Fedora
repositories; Hyprland starts it and every shell keybind goes through `noctalia msg`.
It also handles notifications, the lock screen and idle, so there is no mako, hyprlock or
hypridle.

Its idle behaviours (lock, screen off, suspend) ship **disabled**: turn them on once in
Noctalia's settings (`SUPER + I`), then commit the result.

The repository is the only place its configuration lives. Noctalia layers two sources, and
the role symlinks both into the checkout that ran the playbook:

| Link                                    | Target                               | Written by       |
| --------------------------------------- | ------------------------------------ | ---------------- |
| `~/.config/noctalia/`                   | `roles/noctalia/files/config/`       | hand             |
| `~/.local/state/noctalia/settings.toml` | `roles/noctalia/files/settings.toml` | the settings GUI |

Noctalia writes through symlinks and hot-reloads both, so anything changed from `SUPER + I`
shows up in `git diff roles/noctalia` straight away; GUI values win over the hand-written
ones. `noctalia config validate` checks the hand-written files.

An uncommitted change leaves the checkout dirty, and `dotfiles` then refuses to update it:
commit it, or pass `--skip-update`. Ansible also refuses to replace a non-empty
`~/.config/noctalia` directory; move its `*.toml` files into `roles/noctalia/files/config/`
first.

The wallpapers in `roles/noctalia/files/wallpapers/` are copied to `~/Pictures/Wallpapers`,
which `config.toml` sets as Noctalia's wallpaper folder.

### Look

The shell aims to be quiet: one typeface, one colour source, no outlines.

- **Colours** come from the current wallpaper (`muted` scheme), so changing the wallpaper
  re-themes the whole system without ever turning loud. Noctalia's built-in templates push the
  palette to Hyprland's borders, Ghostty, GTK 3/4 and Qt.
- **Inter** is the shell font; the image installs it from the Fedora repositories.
- **One cursor everywhere.** The cursor theme is user data, not part of the image: changing it
  never needs a rebuild. `hyprdot_cursor` in `roles/hyprdot/defaults/main.yml` sets the theme,
  its size and the archive it comes from (Bibata Modern Classic, 24 px, from the upstream
  release). The theme is extracted into `~/.local/share/icons`, then set for every toolkit,
  since each one reads its own setting:
  - Hyprland and the apps it starts: `XCURSOR_THEME`/`XCURSOR_SIZE` from the generated
    `~/.config/hypr/cursor.lua`;
  - GTK: `gsettings` (`org.gnome.desktop.interface`) and `settings.ini` for GTK 3 and 4;
  - Qt: `XCURSOR_THEME`/`XCURSOR_SIZE` too, which Qt reads on Wayland;
  - XWayland and other fallbacks: `~/.local/share/icons/default/index.theme`;
  - Flatpak: a user override that exposes `~/.local/share/icons` read-only and sets
    `XCURSOR_*`.

  To change themes, edit `hyprdot_cursor` and run the playbook again. The old theme stays in
  `~/.local/share/icons` until you delete it.
- **No borders.** Panels, OSDs and notifications are separated by contrast and Hyprland's blur,
  not outlines.
- **Wallpaper and lock screen transitions** are drawn at random from all of Noctalia's shader
  effects; the rest of the shell keeps its own animations.
- **The bar only shows state, never actions.** Launcher, clipboard, wallpapers, screenshots and
  the session menu all have a keybind; the control center opens from its icon or a right-click
  on the bar. Empty indicators hide themselves (no unread notification, no Bluetooth device).
- **Geometry follows Hyprland.** The `main` bar is a 32 px full-width band with a transparent
  background and concave bottom corners of 16 px, the most a 32 px bar allows, close to the
  windows' `rounding = 18`. The screen corners are rounded to 18 px too, so bar, windows and screen read
  as one frame.
- **Panels grow out of the bar.** Control center, session, clipboard and wallpapers attach to
  the bar; only the launcher floats in the middle of the screen.

The single `main` bar shows workspaces at the start, the clock in the centre, and
notifications, network, Bluetooth, volume, battery and the control center at the end.
`SUPER + SHIFT + Space` hides or shows it.

The templates would normally edit `~/.config/hypr/hyprland.lua` and `~/.config/ghostty/config`
to include the generated colours. Both files are copies Ansible overwrites, so the repository
already carries the lines the templates look for (`require("noctalia")` and
`theme = noctalia`) and the templates leave them alone. Hyprland loads the theme through
`pcall`, so it starts cleanly before Noctalia has generated it. GTK's `gtk.css` is not managed
here and Noctalia edits it directly.

Qt apps go through `qt6ct` (`QT_QPA_PLATFORMTHEME=qt6ct`), installed by the image. Noctalia's
`qt` template only writes `~/.config/qt6ct/colors/noctalia.conf`, so the `noctalia` role points
`qt6ct.conf` at it, along with the Fusion style and the Papirus-Dark icons. Only those keys
are set: whatever else is changed in the `qt6ct` window is kept.

## 1Password

SSH authentication and commit signing both go through the 1Password SSH agent. Private keys
stay in the vault — nothing is ever written to `~/.ssh` except public keys — and every use
asks for confirmation.

Two steps the playbook cannot do for you:

1. **Enable the agent.** In the 1Password desktop app: _Settings → Developer → Use the SSH agent_.
2. **Register the signing key on GitHub.** Add the public key a second time, as a **Signing Key**.
   A key registered only as an Authentication Key leaves your commits `Unverified`.

Which key signs is set by `op.git.signing_key` in `group_vars/all.yml`.

## Gaming

The `steam` role installs Steam from Flathub (`com.valvesoftware.Steam`, user install) along
with the `gamescope` and `MangoHud` Vulkan layers, on the same branch as Steam's runtime. The
sandbox ships its own Mesa, including the 32-bit drivers Proton needs, so the host needs
neither RPM Fusion nor multilib packages. The image ships `steam-devices`, for the udev rules
controllers and VR headsets need.

Steam Play is enabled by default by recent Steam clients, so nothing else is required.
For titles that need a community Proton build, drop GE-Proton into
`~/.var/app/com.valvesoftware.Steam/data/Steam/compatibilitytools.d/`, or install the
`com.valvesoftware.Steam.CompatibilityTool.Proton-GE` flatpak.

To run a game upscaled instead of at native resolution, pass it through `gamescope` in the
launch options:

```
gamescope -W 1280 -H 720 -w 1920 -h 1080 -F fsr -f -- %command%
```

## Releases

`dotfiles self-update` downloads the latest release published on GitHub, checks its
SHA-256 and swaps the running binary. Until a release exists it stops with
_no release published yet_ — there is nothing to update to.

Cutting one is a tag away; pushing it runs GoReleaser through
`.github/workflows/release.yml`:

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

Assets are published **without a version in their name** (`dotfiles_linux_amd64.tar.gz`,
`checksums.txt`). That is deliberate: it keeps the `/releases/latest/download/` URLs
stable, so `self-update` needs no GitHub API call. Renaming them in `.goreleaser.yaml`
breaks `cli/internal/selfupdate`.

To rehearse a release without publishing anything:

```bash
goreleaser release --snapshot --clean --skip=publish
```

---

If you have any questions, suggestions, or just want to chat about life, feel free to open an issue or drop me a message.
