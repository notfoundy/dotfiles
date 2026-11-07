# The dotfiles CLI

Roles run in the order of `default_roles`, the others last and alphabetically, whatever the
selection: ansible-core treats `-t` as a set.

Every run is appended to `~/.local/state/dotfiles/dotfiles.log`.

Building from a clone:

```bash
go build -C cli -o ~/.local/bin/dotfiles .
```
