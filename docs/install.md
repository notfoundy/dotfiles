# Installing

The machine is installed by Anaconda from the Fedora Everything network ISO, which pulls
`ghcr.io/notfoundy/dotfiles:latest` through the `ostreecontainer` command of
`image/install.ks`. The package must be public, or the kickstart needs a `%pre` writing
`/etc/ostree/auth.json`. Embed the kickstart in the ISO, so that every boot entry loads it:

```bash
podman run --rm -v "$PWD":/work:z -w /work registry.fedoraproject.org/fedora:44 sh -c \
  'dnf -y install lorax && mkksiso --ks image/install.ks Fedora-Everything-netinst-x86_64-44-*.iso dotfiles-install.iso'
```

Write `dotfiles-install.iso` to a USB stick and boot it; Secure Boot can stay on. The
kickstart only holds what never changes, Anaconda stops on its summary for the rest:

- **Network**: the Wi-Fi, copied to the installed system.
- **Installation destination**: the disk, automatic partitioning with btrfs, **Encrypt my
  data** checked, then the LUKS passphrase.
- **User creation**: with **Make this user administrator** checked (`wheel`). root stays
  locked.

A few kickstart lines undo what Anaconda gets wrong on a bootc system:

- `firewall --use-system-defaults`: otherwise Anaconda runs `firewall-offline-cmd` and opens
  ssh in the default zone.
- Anaconda always rewrites `/etc/locale.conf` with a lone `LANG=`, with or without a `lang`
  directive. The `%post` copies back the image's version from `/usr/etc`: identical to the
  default, the file is no longer a local change and follows the image again.
- The `%post` also drops the `/` entry Anaconda writes to `/etc/fstab`: with composefs, `/`
  is an overlay and `systemd-remount-fs` fails applying btrfs options to it. The root is
  mounted from the kernel arguments anyway, as `bootc install` does.
