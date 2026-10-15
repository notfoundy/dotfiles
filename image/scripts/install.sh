#!/usr/bin/env bash
set -euo pipefail

systemd-sysusers

mapfile -t packages < <(cat /ctx/packages/*.txt /ctx/deps/*/packages.txt | sort -u)
dnf -y install "${packages[@]}"
ldconfig

broken=0
while read -r elf; do
  if ldd "$elf" 2>/dev/null | grep 'not found'; then
    echo "missing libraries for $elf" >&2
    broken=1
  fi
done < <(cat /ctx/deps/*/binaries.txt)
((broken == 0))

ln -sf ../usr/share/zoneinfo/Europe/Paris /etc/localtime
rm /etc/greetd/config.toml.rpmnew
systemctl enable greetd.service
systemctl set-default graphical.target

dnf clean all
rm -rf /var/cache/* /var/log/* /var/lib/dnf /tmp/*
rm -rf /run/dnf /run/gluster /run/selinux-policy /run/tuned
