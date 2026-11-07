# Desktop shell

[Noctalia](https://docs.noctalia.dev) is the shell, notifications, lock screen and idle: no mako,
hyprlock or hypridle.

## Configuration

| Path                                    | Source                               | Deployed as | Written by       |
| --------------------------------------- | ------------------------------------ | ----------- | ---------------- |
| `~/.config/noctalia/`                   | `roles/noctalia/files/config/`       | symlink     | hand             |
| `~/.local/state/noctalia/settings.toml` | `roles/noctalia/files/settings.toml` | copy        | the settings GUI |

GUI values win over the hand-written ones.

`settings.toml` is only copied when missing: the GUI would otherwise dirty the checkout, which
blocks `dotfiles` updates. Syncing is manual:

```bash
cp ~/.local/state/noctalia/settings.toml roles/noctalia/files/settings.toml   # keep GUI changes
cp roles/noctalia/files/settings.toml ~/.local/state/noctalia/settings.toml   # apply the repo's
```

Idle (lock, screen off, suspend) ships **disabled**: enable it in the settings (`SUPER + I`),
then sync.

The battery charge limit is toggled in the control center (UPower; on the IdeaPad, 60 %).

## Look

Quiet: one typeface (Inter), one colour source, no borders.

- **Colours** come from the wallpaper (`muted` scheme). Built-in templates cover Hyprland,
  Ghostty, GTK and Qt; user templates in `config.toml` cover hyprtoolkit and zathura.
- **Geometry follows Hyprland**: a 32 px bar with 16 px concave corners and 18 px screen corners,
  matching `rounding = 18`. Panels attach to the bar, only the launcher floats.
- **The bar only shows state**: every action has a keybind.

### Templates

- `hyprland.lua` and Ghostty's `config` are overwritten by Ansible, so they already carry the
  lines Noctalia would add (`require("noctalia")`, `theme = noctalia`).
- Qt goes through `qt6ct`: the `noctalia` role only sets its colour scheme, style and icons, the
  rest stays editable in `qt6ct`.
- zathura recolours pages but keeps images (`CTRL + R` to see them as printed). `zathura-cb`
  claims folders and archives, so the role removes those associations.
- GTK3 ignores `color-scheme`: `scripts/sync-gtk-theme.sh` switches between `Adwaita` and
  `Adwaita-dark` on mode change.

### Cursor

User data, not part of the image. `hyprdot_cursor` in `roles/hyprdot/defaults/main.yml` sets the
theme, size and archive, then applies it to Hyprland, GTK, Qt, XWayland and Flatpak. The old
theme stays in `~/.local/share/icons` until deleted.

## Known issues

1Password opens light in dark mode on a cold start
([bug](https://www.1password.community/1password-at-home-31/match-system-theme-setting-not-followed-on-startup-21234)):
switch Noctalia to light and back.
