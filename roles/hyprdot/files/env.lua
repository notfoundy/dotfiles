---@diagnostic disable: undefined-global

-- Wayland
hl.env("ELECTRON_OZONE_PLATFORM_HINT", "auto")

-- Applications
local home = os.getenv("HOME")
hl.env(
	"XDG_DATA_DIRS",
	home .. "/.local/share/flatpak/exports/share:" .. "/var/lib/flatpak/exports/share:/usr/local/share:/usr/share"
)

-- Themes
hl.env("QT_QPA_PLATFORM", "wayland")
hl.env("QT_QPA_PLATFORMTHEME", "kde")
hl.env("XDG_MENU_PREFIX", "plasma-")

-- Virtual environment
hl.env("ILLOGICAL_IMPULSE_VIRTUAL_ENV", os.getenv("HOME") .. "/.local/state/quickshell/.venv")

-- Terminal application
hl.env("TERMINAL", "ghostty")
