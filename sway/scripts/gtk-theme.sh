#!/usr/bin/env bash

## Copyright (C) 2020-2023 Aditya Shakya <adi1090x@gmail.com>
##
## Set GTK Themes, Icons, Cursor and Fonts
## Edited by @yurihikari

THEME='Catppuccin-Macchiato-Standard-Lavender-Dark'
ICONS='Tela-circle-purple-dark'
FONT='JetBrainsMono 10'
CURSOR='Catppuccin-Macchiato-Mauve'
CURSOR_SIZE='32'

SCHEMA='gsettings set org.gnome.desktop.interface'

apply_themes () {
	${SCHEMA} gtk-theme "$THEME"
	${SCHEMA} icon-theme "$ICONS"
	${SCHEMA} cursor-theme "$CURSOR"
	${SCHEMA} font-name "$FONT"
	${SCHEMA} cursor-size "$CURSOR_SIZE"
}

apply_themes
