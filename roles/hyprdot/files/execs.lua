---@diagnostic disable: undefined-global

hl.on("hyprland.start", function()
	-- Shell
	hl.exec_cmd("noctalia")

	-- Core components
	hl.exec_cmd("gnome-keyring-daemon --start --components=secrets")
	hl.exec_cmd("dbus-update-activation-environment --all")
	hl.exec_cmd("sleep 1 && dbus-update-activation-environment --systemd WAYLAND_DISPLAY XDG_CURRENT_DESKTOP")
end)
