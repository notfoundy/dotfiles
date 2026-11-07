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

- **System**: a [bootc](https://bootc-dev.github.io/bootc/) image built from `image/`,
  published to `ghcr.io/notfoundy/dotfiles`.
- **Home**: the Ansible playbook, one role per tool in `roles/`.
- **`dotfiles`**: a Go CLI (`cli/`) that updates `~/.dotfiles` and runs the playbook, without
  `sudo`.

## Installation

1. Build the ISO, then install from it, see [Installing](docs/install.md).
2. Set up the home:

   ```bash
   mkdir -p ~/.local/bin
   curl -fsSLO https://github.com/notfoundy/dotfiles/releases/latest/download/dotfiles_linux_amd64.tar.gz
   tar -xzf dotfiles_linux_amd64.tar.gz -C ~/.local/bin dotfiles
   ~/.local/bin/dotfiles
   ```

   Then log out and sign in from the greeter.
3. By hand:
   - [1Password](docs/1password.md): enable the SSH agent, register the signing key.
   - Idle (lock, screen off, suspend) ships disabled: enable it in Noctalia (`SUPER + I`), see
     [Desktop shell](docs/desktop.md).

## Usage

```bash
dotfiles                        # default_roles from group_vars/all.yml
dotfiles -t nvim,go
dotfiles --select               # interactive picker
dotfiles --skip-update          # leave the checkout alone
dotfiles --repo . -- --check    # this clone; after -- goes to ansible-playbook
dotfiles self-update
```

The update timer stages new images, applied at reboot (`bootc rollback` to go back).

## Documentation

| Page                             | Covers                                      |
| -------------------------------- | ------------------------------------------- |
| [The dotfiles CLI](docs/cli.md)  | Role order, logs, building from source      |
| [Bootc image](docs/image.md)     | Packages, builders, CI, building locally    |
| [System](docs/system.md)         | What differs from a regular Fedora, and why |
| [Installing](docs/install.md)    | ISO, Anaconda, the kickstart                |
| [Desktop shell](docs/desktop.md) | Noctalia and the look of the desktop        |
| [1Password](docs/1password.md)   | SSH authentication and commit signing       |
| [Gaming](docs/gaming.md)         | Steam, gamescope, Proton-GE                 |
| [Releases](docs/releases.md)     | `cli-v*` and `img-v*` tags                  |

---

If you have any questions, suggestions, or just want to chat about life, feel free to open an issue or drop me a message.
