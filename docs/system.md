# Differences from a regular Fedora

What the image changes compared to a regular Fedora install:

- greetd is the display manager, with [Noctalia Greeter](https://docs.noctalia.dev/greeter/)
  as its greeter. It only offers the `Hyprland` session, which runs `start-hyprland`: Hyprland
  is built without uwsm. The greeter runs as Fedora's `greetd` account, so the tmpfiles rule
  creating its state directory is overridden in `/etc/tmpfiles.d/`, as upstream expects.
  `greetd` gets a fixed UID (903, `image/files/usr/lib/sysusers.d/greetd-id.conf`) so that
  `/var/lib/greetd` keeps its owner across images. `dbus-daemon` ships only for
  `dbus-run-session`, which the greeter session needs; the system bus stays on dbus-broker.
  Syncing Noctalia's wallpaper and palette to the greeter is not set up.
- gnome-keyring provides the Secret Service. greetd's PAM stack unlocks it with the login
  password through `pam_gnome_keyring`, which `gnome-keyring-pam` ships: without it, the
  `-auth` line is skipped silently and every app asking for a secret prompts for it.
- hyprpolkitagent is the polkit agent. Its unit hangs off `graphical-session.target`, which
  never starts without uwsm, so `execs.lua` starts it once the session environment is imported
  into systemd, as the Hyprland wiki recommends. Noctalia's own agent stays off
  (`polkit_agent = false`): a session has a single agent.
- Swap is zram only, through `zram-generator-defaults`: the base image ships none, and
  without it the system has no swap at all. `systemd-oomd-defaults` kills the offending
  cgroup under sustained memory pressure, before the desktop freezes.
- The base image enables sshd with password login and ships no firewall. sshd is disabled
  (`systemctl enable --now sshd` when needed) and firewalld runs with Fedora's `public`
  zone. libvirt and podman (netavark) both use nftables and keep their own tables, and
  firewalld ships the `libvirt` zones, so VMs and containers keep their networking.
- The base image only knows `C.UTF-8`. `glibc-langpack-en` and `-fr` are added, and
  `/etc/locale.conf` sets `en_US.UTF-8` with French formats (time, currency, paper, units).
  Numbers stay English: a decimal comma breaks `printf '%.2f'` and friends in scripts.
- Anaconda writes `compress=zstd:1` to fstab, but bootc mounts the root from the kernel
  command line, and btrfs takes its options from that first mount. Compression is set back
  through `rootflags` in `/usr/lib/bootc/kargs.d/`.
- Anaconda adds no `quiet` for a bootc install, so the console printed every kernel info
  message, such as the SELinux label bootc writes to test its rights on each run. `quiet` comes
  from `kargs.d` as well; the messages stay in `journalctl -k`.
- `podman.socket` is enabled for every user (`systemctl --global`), rootless, for the tools that
  only speak the Docker API (Testcontainers, a `docker` CLI in a container). It is socket
  activated: the service only starts on the first connection, at
  `$XDG_RUNTIME_DIR/podman/podman.sock`. `podman compose` needs no socket, it runs
  `podman-compose`.
- The base image ships no `audit`: without auditd, the kernel prints its audit events on the
  console, over the TTY. `audit` is added and `auditd` enabled so they go to its log.
- libvirt is managed by the `wheel` group through a polkit rule: `usermod -aG libvirt` does not
  work, since the group lives in `/usr/lib/group`.
- bootc only copies the image's `/var` at install time. The directories packages leave there
  without a tmpfiles rule (libvirt, flatpak, AccountsService…) are listed by
  `bootc container lint` and written to `/usr/lib/tmpfiles.d/dotfiles-var.conf` at the end of
  the build, so they are recreated if missing, as `rpm-ostree-autovar.conf` does for the base
  image.
- `/` is the composefs image, read-only and always 100 % full. Disk usage is read on `/var`,
  the btrfs partition that also holds `/var/home`: `df -h /var`, and the fastfetch role points
  its disk module there.
- Go, Node, Typst, Starship and Lazygit are installed in `~/.local` by their roles, so they
  stay on the latest release without touching the image. Starship goes through its official
  install script, Lazygit is checked against the release's `checksums.txt`.
