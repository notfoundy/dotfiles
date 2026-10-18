# Desktop shell

The image installs [Noctalia](https://docs.noctalia.dev) from the official Fedora
repositories; Hyprland starts it and every shell keybind goes through `noctalia msg`.
It also handles notifications, the lock screen and idle, so there is no mako, hyprlock or
hypridle.

Its idle behaviours (lock, screen off, suspend) ship **disabled**: turn them on once in
Noctalia's settings (`SUPER + I`), then commit the result.

The repository is the only place its configuration lives. Noctalia layers two sources, and
the role symlinks both into the checkout that ran the playbook:

| Link                                    | Target                               | Written by       |
| --------------------------------------- | ------------------------------------ | ---------------- |
| `~/.config/noctalia/`                   | `roles/noctalia/files/config/`       | hand             |
| `~/.local/state/noctalia/settings.toml` | `roles/noctalia/files/settings.toml` | the settings GUI |

Noctalia writes through symlinks and hot-reloads both, so anything changed from `SUPER + I`
shows up in `git diff roles/noctalia` straight away; GUI values win over the hand-written
ones. `noctalia config validate` checks the hand-written files.

An uncommitted change leaves the checkout dirty, and `dotfiles` then refuses to update it:
commit it, or pass `--skip-update`. Ansible also refuses to replace a non-empty
`~/.config/noctalia` directory; move its `*.toml` files into `roles/noctalia/files/config/`
first.

The wallpapers in `roles/noctalia/files/wallpapers/` are copied to `~/Pictures/Wallpapers`,
which `config.toml` sets as Noctalia's wallpaper folder.

## Look

The shell aims to be quiet: one typeface, one colour source, no outlines.

- **Colours** come from the current wallpaper (`muted` scheme), so changing the wallpaper
  re-themes the whole system without ever turning loud. Noctalia's built-in templates push the
  palette to Hyprland's borders, Ghostty, GTK 3/4 and Qt.
- **Inter** is the shell font; the image installs it from the Fedora repositories.
- **One cursor everywhere.** The cursor theme is user data, not part of the image: changing it
  never needs a rebuild. `hyprdot_cursor` in `roles/hyprdot/defaults/main.yml` sets the theme,
  its size and the archive it comes from (Bibata Modern Classic, 24 px, from the upstream
  release). The theme is extracted into `~/.local/share/icons`, then set for every toolkit,
  since each one reads its own setting:
  - Hyprland and the apps it starts: `XCURSOR_THEME`/`XCURSOR_SIZE` from the generated
    `~/.config/hypr/cursor.lua`;
  - GTK: `gsettings` (`org.gnome.desktop.interface`) and `settings.ini` for GTK 3 and 4;
  - Qt: `XCURSOR_THEME`/`XCURSOR_SIZE` too, which Qt reads on Wayland;
  - XWayland and other fallbacks: `~/.local/share/icons/default/index.theme`;
  - Flatpak: a user override that exposes `~/.local/share/icons` read-only and sets
    `XCURSOR_*`.

  To change themes, edit `hyprdot_cursor` and run the playbook again. The old theme stays in
  `~/.local/share/icons` until you delete it.
- **No borders.** Panels, OSDs and notifications are separated by contrast and Hyprland's blur,
  not outlines.
- **Wallpaper and lock screen transitions** are drawn at random from all of Noctalia's shader
  effects; the rest of the shell keeps its own animations.
- **The bar only shows state, never actions.** Launcher, clipboard, wallpapers, screenshots and
  the session menu all have a keybind; the control center opens from its icon or a right-click
  on the bar. Empty indicators hide themselves (no unread notification, no Bluetooth device).
- **Geometry follows Hyprland.** The `main` bar is a 32 px full-width band with a transparent
  background and concave bottom corners of 16 px, the most a 32 px bar allows, close to the
  windows' `rounding = 18`. The screen corners are rounded to 18 px too, so bar, windows and screen read
  as one frame.
- **Panels grow out of the bar.** Control center, session, clipboard and wallpapers attach to
  the bar; only the launcher floats in the middle of the screen.

The single `main` bar shows workspaces at the start, the clock in the centre, and
notifications, network, Bluetooth, volume, battery and the control center at the end.
`SUPER + SHIFT + Space` hides or shows it.

The templates would normally edit `~/.config/hypr/hyprland.lua` and `~/.config/ghostty/config`
to include the generated colours. Both files are copies Ansible overwrites, so the repository
already carries the lines the templates look for (`require("noctalia")` and
`theme = noctalia`) and the templates leave them alone. Hyprland loads the theme through
`pcall`, so it starts cleanly before Noctalia has generated it. GTK's `gtk.css` is not managed
here and Noctalia edits it directly.

Qt apps go through `qt6ct` (`QT_QPA_PLATFORMTHEME=qt6ct`), installed by the image. Noctalia's
`qt` template only writes `~/.config/qt6ct/colors/noctalia.conf`, so the `noctalia` role points
`qt6ct.conf` at it, along with the Fusion style and the Papirus-Dark icons. Only those keys
are set: whatever else is changed in the `qt6ct` window is kept.

hyprtoolkit, behind hyprpolkitagent and Hyprland's own dialogs, has no built-in template: a
user template in `config.toml` renders `roles/noctalia/files/config/templates/hyprtoolkit.conf`
to `~/.config/hypr/hyprtoolkit.conf`, which the next dialog reads.
