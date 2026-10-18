# Gaming

The `steam` role installs Steam from Flathub (`com.valvesoftware.Steam`, user install) along
with the `gamescope` and `MangoHud` Vulkan layers, on the same branch as Steam's runtime. The
sandbox ships its own Mesa, including the 32-bit drivers Proton needs, so the host needs
neither RPM Fusion nor multilib packages. The image ships `steam-devices`, for the udev rules
controllers and VR headsets need.

Steam Play is enabled by default by recent Steam clients, so nothing else is required.
For titles that need a community Proton build, drop GE-Proton into
`~/.var/app/com.valvesoftware.Steam/data/Steam/compatibilitytools.d/`, or install the
`com.valvesoftware.Steam.CompatibilityTool.Proton-GE` flatpak.

To run a game upscaled instead of at native resolution, pass it through `gamescope` in the
launch options:

```
gamescope -W 1280 -H 720 -w 1920 -h 1080 -F fsr -f -- %command%
```
