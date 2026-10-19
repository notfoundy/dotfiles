# 1Password

SSH authentication and commit signing both go through the 1Password SSH agent. Private keys
stay in the vault — nothing is ever written to `~/.ssh` except public keys — and every use
asks for confirmation.

Two steps the playbook cannot do for you:

1. **Enable the agent.** In the 1Password desktop app: _Settings → Developer → Use the SSH agent_.
2. **Register the signing key on GitHub.** Add the public key a second time, as a **Signing Key**.
   A key registered only as an Authentication Key leaves your commits `Unverified`.

Which key signs is set by `op.git.signing_key` in `group_vars/all.yml`.

## Groups

The desktop app only trusts `op` and the browser helper when they connect with the GID of
`onepassword-cli` / `onepassword`, and only if that GID is 1000 or above: below it, `op`
fails with `connection reset` and the app logs `invalid group attempted to connect`. The
image pins the groups to 1500–1502 in `image/files/usr/lib/sysusers.d/onepassword.conf`.

`/etc/group` survives image updates and sysusers never renumbers a group, so a machine
installed from an older image (900–902) needs a one-time fix after booting the new one,
then a restart of the app:

```sh
sudo groupmod -g 1500 onepassword
sudo groupmod -g 1501 onepassword-mcp
sudo groupmod -g 1502 onepassword-cli
```
