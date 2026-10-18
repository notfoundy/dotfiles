# 1Password

SSH authentication and commit signing both go through the 1Password SSH agent. Private keys
stay in the vault — nothing is ever written to `~/.ssh` except public keys — and every use
asks for confirmation.

Two steps the playbook cannot do for you:

1. **Enable the agent.** In the 1Password desktop app: _Settings → Developer → Use the SSH agent_.
2. **Register the signing key on GitHub.** Add the public key a second time, as a **Signing Key**.
   A key registered only as an Authentication Key leaves your commits `Unverified`.

Which key signs is set by `op.git.signing_key` in `group_vars/all.yml`.
