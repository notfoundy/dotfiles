# 1Password

SSH authentication and commit signing go through the 1Password SSH agent: only public keys are
written to `~/.ssh`.

By hand:

1. In the app: _Settings → Developer → Use the SSH agent_.
2. On GitHub, add the public key a second time as a **Signing Key**, or commits stay
   `Unverified`.

The signing key is `op.git.signing_key` in `group_vars/all.yml`.

## Groups

The app only trusts `op` and the browser helper when the GID of `onepassword-cli` /
`onepassword` is 1000 or above (otherwise `op` fails with `connection reset`). The image pins
them to 1500–1502 in `image/files/usr/lib/sysusers.d/onepassword.conf`.

sysusers never renumbers an existing group: a machine installed from an older image (900–902)
needs, once, then an app restart:

```sh
sudo groupmod -g 1500 onepassword
sudo groupmod -g 1501 onepassword-mcp
sudo groupmod -g 1502 onepassword-cli
```
