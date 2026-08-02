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

## 1Password

SSH authentication and commit signing both go through the 1Password SSH agent. Private keys
stay in the vault — nothing is ever written to `~/.ssh` except public keys — and every use
asks for confirmation.

Two steps the playbook cannot do for you:

1. **Enable the agent.** In the 1Password desktop app: _Settings → Developer → Use the SSH agent_.
2. **Register the signing key on GitHub.** Add the public key a second time, as a **Signing Key**.
   A key registered only as an Authentication Key leaves your commits `Unverified`.

Which key signs is set by `op.git.signing_key` in `group_vars/all.yml`.

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
