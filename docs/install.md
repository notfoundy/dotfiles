# Installing

The installer is the Fedora Everything network ISO with `image/install.ks` embedded: Anaconda
then installs `ghcr.io/notfoundy/dotfiles:latest` instead of a regular Fedora. The package must
stay public, or the kickstart needs a `%pre` writing `/etc/ostree/auth.json`.

## 1. Build the ISO

Download `Fedora-Everything-netinst-x86_64-44-*.iso` from
[fedoraproject.org/everything/download](https://fedoraproject.org/everything/download) into the
repository root, then embed the kickstart so that every boot entry loads it:

```bash
podman run --rm -v "$PWD":/work:z -w /work registry.fedoraproject.org/fedora:44 sh -c \
  'dnf -y install lorax && mkksiso --ks image/install.ks Fedora-Everything-netinst-x86_64-44-*.iso dotfiles-install.iso'
```

## 2. Write it to a USB stick

Find the stick with `lsblk` (here `/dev/sdX`): everything on it is erased.

```bash
sudo dd if=dotfiles-install.iso of=/dev/sdX bs=4M status=progress conv=fsync
```

## 3. Install

Boot the stick; Secure Boot can stay on. Anaconda stops on its summary for what the kickstart
leaves out:

- **Network**: the Wi-Fi, kept on the installed system.
- **Installation destination**: the disk, automatic btrfs partitioning, **Encrypt my data**.
- **User creation**: check **Make this user administrator** (`wheel`). root stays locked.

Reboot, log in, then set up the home

## What the kickstart fixes

- `firewall --use-system-defaults`: otherwise Anaconda opens ssh in the default zone.
- `/etc/locale.conf`: Anaconda always rewrites it with a lone `LANG=`; `%post` restores the
  image's copy so that it keeps following the image.
- `/etc/fstab`: `%post` drops the `/` entry. With composefs, `/` is an overlay and
  `systemd-remount-fs` fails on its btrfs options; the root is mounted from the kernel arguments.
