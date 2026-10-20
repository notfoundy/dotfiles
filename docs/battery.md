# Battery

On firmware with a charge threshold, the image's udev rule stops charging at 80 % and
nothing else is needed (see [System](system.md)).

IdeaPads have no threshold, only charge modes in `/sys/class/power_supply/BAT0/charge_types`:
`Standard` charges to 100 %, `Long_Life` (Lenovo's Conservation mode) stops at 60 %, a value
the EC fixes. The mode lives in the EC, not in the OS: one set by Lenovo Vantage under Windows
survives a reinstall and even a BIOS reset.

The `battery` role, outside `default_roles`, gets the 80 % back:

```bash
dotfiles -t battery
```

It only installs anything when the `ideapad_acpi` driver is loaded, and the service checks it
again. `ideapad-charge-limit`, a `systemd --user` service, reads the charge every 2 minutes:
`Standard` at 75 % or less, `Long_Life` at 80 % or more. Above 60 %, `Long_Life` stops
charging without draining the battery, so it holds around 80 % on AC. In between it keeps the
current mode, so that it does not flip back and forth.

The service only runs while the session does, so it switches back to `Long_Life` when it
stops: otherwise a laptop left in `Standard` would charge to 100 % once logged out or shut
down.

To charge to 100 % once:

```bash
systemctl --user stop ideapad-charge-limit
echo Standard > /sys/class/power_supply/BAT0/charge_types
```

Stopping the service sets `Long_Life`, hence `Standard` afterwards.
`systemctl --user start ideapad-charge-limit`, or the next login, brings the cap back.
