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

This playbook includes a shell script located at `bin/dotfiles`. This script makes sure any Ansible dependencies are installed and updated and initializes your environment.
Ansible will take care of everything (or almost). You just need to grab a coffee while it works!

### On a fresh installation

```bash
bash -c "$(curl -fsSL https://raw.githubusercontent.com/notfoundy/dotfiles/main/bin/dotfiles)"
```

### Otherwise

To update your environment run the `dotfiles` command in your shell

```bash
dotfiles
```

## 1Password

SSH authentication and commit signing both go through the 1Password SSH agent. Private keys
stay in the vault — nothing is ever written to `~/.ssh` except public keys — and every use
asks for confirmation.

Two steps the playbook cannot do for you:

1. **Enable the agent.** In the 1Password desktop app: _Settings → Developer → Use the SSH agent_.
2. **Register the signing key on GitHub.** Add the public key a second time, as a **Signing Key**.
   A key registered only as an Authentication Key leaves your commits `Unverified`.

Which key signs is set by `op.git.signing_key` in `group_vars/all.yml`.

---

If you have any questions, suggestions, or just want to chat about life, feel free to open an issue or drop me a message.
