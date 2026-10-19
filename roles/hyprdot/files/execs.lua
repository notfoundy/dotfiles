---@diagnostic disable: undefined-global

hl.on("hyprland.start", function()
	-- Shell
	hl.exec_cmd("noctalia")

	-- Core components
	hl.exec_cmd("gnome-keyring-daemon --start --components=secrets")
	hl.exec_cmd("dbus-update-activation-environment --all")
	hl.exec_cmd(
		"sleep 1 && dbus-update-activation-environment --systemd WAYLAND_DISPLAY XDG_CURRENT_DESKTOP"
			.. " && systemctl --user start hyprland-session.target hyprpolkitagent.service"
	)
end)

-- Stops graphical-session.target and the services bound to it (portals, polkit agent)
hl.on("hyprland.shutdown", function()
	hl.exec_cmd("systemctl --user stop hyprland-session.target")
end)
