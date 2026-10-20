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

## How it works

> [!WARNING]
> Only the Fedora bootc image built from `image/` is supported: it carries every system
> package, and the playbook only manages `$HOME`.

- **The system** is a [bootc](https://bootc-dev.github.io/bootc/) image, built from `image/`
  and published to `ghcr.io/notfoundy/dotfiles`.
- **The home** is set up by the Ansible playbook, one role per tool in `roles/`.
- **The `dotfiles` command**, a small Go program living in `cli/`, keeps the checkout in
  `~/.dotfiles` up to date and runs the playbook. Nothing needs `sudo`.

## Installation

### 1. Build the ISO

The machine is installed by Anaconda from the Fedora Everything network ISO, which pulls the
image through the kickstart `image/install.ks`. Embed the kickstart in the ISO:

```bash
podman run --rm -v "$PWD":/work:z -w /work registry.fedoraproject.org/fedora:44 sh -c \
  'dnf -y install lorax && mkksiso --ks image/install.ks Fedora-Everything-netinst-x86_64-44-*.iso dotfiles-install.iso'
```

### 2. Install the system

Write `dotfiles-install.iso` to a USB stick and boot it; Secure Boot can stay on. Anaconda
stops on its summary for:

- **Network**: the Wi-Fi, copied to the installed system.
- **Installation destination**: the disk, automatic partitioning with btrfs, **Encrypt my
  data** checked, then the LUKS passphrase.
- **User creation**: with **Make this user administrator** checked (`wheel`). root stays
  locked.

### 3. Set up the home

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

### 4. Finish by hand

- **1Password**: enable the SSH agent and register the signing key on GitHub, see
  [1Password](docs/1password.md).
- **Idle**: lock, screen off and suspend ship disabled; turn them on in Noctalia's settings
  (`SUPER + I`), then commit the result, see [Desktop shell](docs/desktop.md).
- **Battery**: on an IdeaPad, `dotfiles -t battery` caps the charge at 80 %, see
  [Battery](docs/battery.md).

## Usage

```bash
dotfiles                        # every default role, as listed by default_roles in group_vars/all.yml
dotfiles -t nvim,go             # only those roles
dotfiles --select               # pick them interactively
dotfiles --skip-update          # leave the checkout alone
dotfiles --repo . -- --check    # run against this clone; anything after -- goes to ansible-playbook
```

`dotfiles --help` lists the rest.

## Updating

- **The system**: the update timer stages new images, applied at the next reboot.
  `bootc upgrade` does it by hand, `bootc rollback` returns to the previous image.
- **The CLI**: `dotfiles self-update` fetches the latest release.
- **The home**: run `dotfiles` again.

## Documentation

| Page                                 | Covers                                                        |
| ------------------------------------ | ------------------------------------------------------------- |
| [The dotfiles CLI](docs/cli.md)      | Role order, logs, building from source                        |
| [Bootc image](docs/image.md)         | Packages, builders, pinned versions, CI, building locally     |
| [System](docs/system.md)             | What differs from a regular Fedora, and why                   |
| [Installing](docs/install.md)        | The kickstart and what it fixes in Anaconda                   |
| [Desktop shell](docs/desktop.md)     | Noctalia, its configuration and the look of the whole desktop |
| [1Password](docs/1password.md)       | SSH authentication and commit signing                         |
| [Gaming](docs/gaming.md)             | Steam, gamescope, Proton-GE                                   |
| [Battery](docs/battery.md)           | The 80 % charge cap, and its IdeaPad workaround               |
| [Releases](docs/releases.md)         | `cli-v*` and `img-v*` tags, what each one publishes           |

---

If you have any questions, suggestions, or just want to chat about life, feel free to open an issue or drop me a message.
