#!/bin/sh

# set the service with systemctl --user enable wallpaper.timer
# set a random wallpaper
backgrounds="${HOME}/.config/sway/backgrounds"
wallpaper=$(ls $backgrounds | shuf -n 1)
path="${backgrounds}/${wallpaper}"
swaymsg output "*" bg $path fill
