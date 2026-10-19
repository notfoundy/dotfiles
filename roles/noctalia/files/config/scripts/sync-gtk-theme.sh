#!/usr/bin/env bash
# Electron apps (1Password) read dark mode from the GTK theme name, not from color-scheme
mode=${NOCTALIA_THEME_MODE:-$(gsettings get org.gnome.desktop.interface color-scheme)}
case $mode in
  *dark*) theme=Adwaita-dark ;;
  *) theme=Adwaita ;;
esac
gsettings set org.gnome.desktop.interface gtk-theme "$theme"
