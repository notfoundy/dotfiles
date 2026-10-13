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
> For now, only Fedora is supported because it’s the only distribution I use.
> I'm not a distro hopper, so it might stay this way.

Everything runs through the `dotfiles` command, a small Go program living in `cli/`.
It prepares the host, makes sure the Ansible dependencies are installed and updated,
refreshes the checkout in `~/.dotfiles`, then hands the terminal over to
`ansible-playbook`.
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
repository to `~/.dotfiles` and asks you to reboot once it is done.

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

## Desktop shell

The `noctalia` role installs [Noctalia](https://docs.noctalia.dev) from the official Fedora
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
  palette to Hyprland's borders, Ghostty and GTK 3/4.
- **Inter** is the shell font; the `font` role installs it from the Fedora repositories.
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

The `steam` role enables RPM Fusion, installs Steam along with `gamescope`, `mangohud` and
`vulkan-tools`, and pulls the 32-bit Mesa drivers that Proton titles need. It then asserts
that the loader actually sees a hardware Vulkan driver — if only `llvmpipe` answers, games
would silently render on the CPU — and that a 32-bit ICD is present.

Steam Play is enabled by default by recent Steam clients, so nothing else is required.
For titles that need a community Proton build, drop GE-Proton into
`~/.steam/root/compatibilitytools.d/`.

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
