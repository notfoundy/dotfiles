---@diagnostic disable: undefined-global
---@diagnostic disable: lowercase-global

-- The folder within ~/.config/quickshell containing the config
hl.env("qsConfig", "ii")

-- Apps
terminal = "~/.config/hypr/scripts/launch_first_available.sh 'ghostty' 'kitty -1' "
fileManager = "~/.config/hypr/scripts/launch_first_available.sh 'dolphin' 'nautilus' 'nemo' 'thunar' "
browser = "~/.config/hypr/scripts/launch_first_available.sh 'zen-browser' 'firefox' 'brave' 'chromium' 'opera'"
codeEditor = "~/.config/hypr/scripts/launch_first_available.sh 'command -v nvim && ghostty -e nvim' 'zed'"
officeSoftware = "~/.config/hypr/scripts/launch_first_available.sh 'wps' 'libreoffice'"
textEditor = "~/.config/hypr/scripts/launch_first_available.sh 'kate' 'gnome-text-editor'"
volumeMixer = "~/.config/hypr/scripts/launch_first_available.sh 'pavucontrol-qt' 'pavucontrol'"
settingsApp =
	"XDG_CURRENT_DESKTOP=gnome ~/.config/hypr/scripts/launch_first_available.sh 'qs -p ~/.config/quickshell/$qsConfig/settings.qml' 'systemsettings' 'gnome-control-center' 'better-control'"
taskManager =
	"~/.config/hypr/scripts/launch_first_available.sh 'gnome-system-monitor' 'plasma-systemmonitor --page-name Processes' 'command -v btop && ghostty -1 bash -c btop'"
