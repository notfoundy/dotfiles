# Gaming

The `steam` role installs Steam from Flathub (user install) with the `gamescope` and `MangoHud`
layers. The sandbox ships its own Mesa, 32-bit included, so no RPM Fusion nor multilib. The image
adds `steam-devices` for controller and VR udev rules.

Proton-GE: install the `com.valvesoftware.Steam.CompatibilityTool.Proton-GE` flatpak.

Upscaling, in the launch options:

```
gamescope -W 1280 -H 720 -w 1920 -h 1080 -F fsr -f -- %command%
```
