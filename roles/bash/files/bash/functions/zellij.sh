#!/usr/bin/env bash

zellij() {
  if [ "$#" -eq 0 ]; then
    _zellij_launcher
  else
    command zellij "$@"
  fi
}

_zellij_launcher() {
  SELECTED=$(fd . ~/projects -d 1 -t d | fzf)

  [ -z "$SELECTED" ] && return

  SESSION_NAME=$(basename "$SELECTED" | tr ' ' '_' | tr . _)

  if zellij list-sessions | rg -q $SESSION_NAME; then
    command zellij attach "$SESSION_NAME"
  else
    cd "$SELECTED" || return
    command zellij -s "$SESSION_NAME"
  fi
}
