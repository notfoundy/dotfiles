ostreecontainer --url=ghcr.io/notfoundy/dotfiles:latest
lang en_US.UTF-8
keyboard us
timezone Europe/Paris --utc
firewall --use-system-defaults
rootpw --lock

%post
cp /usr/etc/locale.conf /etc/locale.conf
%end
