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

---

If you have any questions, suggestions, or just want to chat about life, feel free to open an issue or drop me a message.
