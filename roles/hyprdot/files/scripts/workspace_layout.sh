#!/usr/bin/env bash

active="$(hyprctl activeworkspace -j)"
curr_workspace="$(jq -r ".id" <<<"${active}")"

if [[ "$(jq -r ".tiledLayout" <<<"${active}")" == "dwindle" ]]; then
  layout="scrolling"
else
  layout="dwindle"
fi

hyprctl keyword workspace "${curr_workspace}, layout:${layout}"
notify-send -u low "Workspace layout set to ${layout}" -a Hyprland
