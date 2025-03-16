#!/bin/sh

WALLPAPER_DIR="$HOME/.local/share/wallpapers"
CURRENT_WALLPAPER=$(swww query | awk -F': ' '{print $5}')
NEW_WALLPAPER=$(find "$WALLPAPER_DIR" -type f ! -path "$CURRENT_WALLPAPER" | shuf -n 1)

if [ -z "$NEW_WALLPAPER" ]; then
  echo "Error : New wallpaper not found"
  exit 1
fi

swww img $NEW_WALLPAPER \
    --transition-fps 60 \
    --transition-type any

