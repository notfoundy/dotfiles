# Differences from a regular Fedora

- **greetd** with [Noctalia Greeter](https://docs.noctalia.dev/greeter/), offering only
  `start-hyprland` (Hyprland is built without uwsm). `greetd` has a fixed UID (903) so that
  `/var/lib/greetd` keeps its owner across images. `dbus-daemon` is only there for
  `dbus-run-session`, which the greeter needs; the system bus stays dbus-broker.
- **gnome-keyring** is unlocked at login by `pam_gnome_keyring` (`gnome-keyring-pam`): without
  it, PAM skips the line silently and every app prompts for its secret.
- **`hyprland-session.target`**: without uwsm nothing starts `graphical-session.target`, which
  `xdg-desktop-portal` needs (no dark mode in Flatpak and Electron apps otherwise). `execs.lua`
  starts and stops it.
- **hyprpolkitagent** is the polkit agent; Noctalia's stays off, a session has only one.
- **zram swap** (`zram-generator-defaults`): the base image has no swap at all.
  `systemd-oomd-defaults` kills the offending cgroup before the desktop freezes.
- **sshd disabled, firewalld on** (`public` zone): the base image enables sshd with password
  login and ships no firewall. libvirt and podman keep their own nftables tables.
- **Locale**: `en_US.UTF-8` with French formats, but English numbers: a decimal comma breaks
  `printf '%.2f'` in scripts.
- **Kernel arguments** (`/usr/lib/bootc/kargs.d/`): btrfs takes its options from the first
  mount, done by bootc from the command line, so `compress=zstd:1` goes in `rootflags`. `quiet`
  too, which Anaconda omits for bootc.
- **`podman.socket`** for every user, for tools that only speak the Docker API
  (Testcontainers).
- **auditd**: without it, audit events are printed over the TTY.
- **libvirt** is granted to `wheel` by a polkit rule: `usermod -aG libvirt` cannot work, the
  group lives in `/usr/lib/group`.
- **`/var` directories** without a tmpfiles rule are listed by `bootc container lint` and written
  to `/usr/lib/tmpfiles.d/dotfiles-var.conf`: bootc only copies `/var` at install time.
- **Disk usage** is read on `/var` (`df -h /var`): `/` is the composefs image, always 100 %.
- Go, Node, Typst, Starship and Lazygit are installed in `~/.local` by their roles, to follow
  upstream releases without rebuilding the image.
- **No tablet mode** (`/usr/lib/modprobe.d/90-no-tablet-mode.conf`): `lenovo_ymc` and `amd_sfh`
  are blacklisted. On the IdeaPad Flex 5 14ARE05 the touchpad stops answering on I2C and
  closed-lid sleep wakes up by itself; tablet mode handling is the suspect. Delete the file and
  rebuild to undo.
