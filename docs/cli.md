# The dotfiles CLI

Everything runs through the `dotfiles` command, a small Go program living in `cli/`.
It checks the host, updates the Ansible dependencies, refreshes the checkout in
`~/.dotfiles`, then hands the terminal over to `ansible-playbook`. Nothing needs `sudo`.
Ansible will take care of the rest (or almost). You just need to grab a coffee while it works!

## Commands

```bash
dotfiles                        # every default role, as listed by default_roles in group_vars/all.yml
dotfiles -t nvim,go             # only those roles
dotfiles --select               # pick them interactively
dotfiles --skip-update          # leave the checkout alone
dotfiles --repo . -- --check    # run against this clone; anything after -- goes to ansible-playbook
```

Whatever the selection, roles run in the order of `default_roles`, the ones outside it last
and alphabetically: ansible-core keeps `-t` as a set, so its order cannot be relied on.

`dotfiles --help` lists the rest, and `dotfiles self-update` fetches the latest
release. Every run is appended to `~/.local/state/dotfiles/dotfiles.log`, so a role
that failed can be read back afterwards.

## Building from source

Instead of downloading a release, from a clone:

```bash
go build -C cli -o ~/.local/bin/dotfiles .
```
