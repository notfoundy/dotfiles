#!/usr/bin/env bash

width="${1:-1300}"
height="${2:-900}"

active="$(hyprctl activewindow -j)"

if [[ -z "${active}" || "$(jq -r ".address" <<<"${active}")" == "null" ]]; then
  exit 0
fi

if [[ "$(jq -r ".pinned" <<<"${active}")" == "true" ]]; then
  hyprctl -q --batch "dispatch pin ; dispatch togglefloating"
else
  hyprctl -q --batch "dispatch togglefloating ; \
    dispatch resizeactive exact ${width} ${height} ; \
    dispatch centerwindow ; \
    dispatch pin ; \
    dispatch alterzorder top"
fi
