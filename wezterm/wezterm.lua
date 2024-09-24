local wezterm = require("wezterm")
local theme = wezterm.plugin.require("https://github.com/neapsix/wezterm").main

return {
	colors = theme.colors(),
	enable_tab_bar = false,
	font_size = 16.0,
	font = wezterm.font("JetBrains Mono Nerd Font"),
	window_background_opacity = 0.78,
	mouse_bindings = {
		-- Ctrl-click will open the link under the mouse cursor
		{
			event = { Up = { streak = 1, button = "Left" } },
			mods = "CTRL",
			action = wezterm.action.OpenLinkAtMouseCursor,
		},
	},
}
