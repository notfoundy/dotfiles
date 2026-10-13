---@diagnostic disable: undefined-global

require("variables")

local hyprScripts = "$HOME/.config/hypr/scripts"
local ipc = "noctalia msg "

-- General
hl.bind("SUPER + D", hl.dsp.exec_cmd(ipc .. "panel-toggle launcher"), { description = "Widgets:Toggle launcher" })
hl.bind(
	"SUPER + N",
	hl.dsp.exec_cmd(ipc .. "panel-toggle control-center"),
	{ description = "Widgets:Toggle control center" }
)
hl.bind(
	"CTRL + ALT + Delete",
	hl.dsp.exec_cmd(ipc .. "panel-toggle session"),
	{ description = "Session:Toggle session menu" }
)
hl.bind("SUPER + SHIFT + Space", hl.dsp.exec_cmd(ipc .. "bar-toggle"), { description = "Widgets:Toggle bar" })

hl.bind(
	"XF86MonBrightnessUp",
	hl.dsp.exec_cmd(ipc .. "brightness-up"),
	{ locked = true, repeating = true, description = "Screen:Increase brightness" }
)
hl.bind(
	"XF86MonBrightnessDown",
	hl.dsp.exec_cmd(ipc .. "brightness-down"),
	{ locked = true, repeating = true, description = "Screen:Decrease brightness" }
)
hl.bind(
	"XF86AudioRaiseVolume",
	hl.dsp.exec_cmd(ipc .. "volume-up"),
	{ locked = true, repeating = true, description = "Audio:Raise volume" }
)
hl.bind(
	"XF86AudioLowerVolume",
	hl.dsp.exec_cmd(ipc .. "volume-down"),
	{ locked = true, repeating = true, description = "Audio:Lower volume" }
)

hl.bind("XF86AudioMute", hl.dsp.exec_cmd(ipc .. "volume-mute"), { locked = true })
hl.bind(
	"SUPER + SHIFT + M",
	hl.dsp.exec_cmd(ipc .. "volume-mute"),
	{ locked = true, description = "Audio:Toggle mute" }
)
hl.bind("ALT + XF86AudioMute", hl.dsp.exec_cmd(ipc .. "mic-mute"), { locked = true })
hl.bind("XF86AudioMicMute", hl.dsp.exec_cmd(ipc .. "mic-mute"), { locked = true })
hl.bind(
	"SUPER + ALT + M",
	hl.dsp.exec_cmd(ipc .. "mic-mute"),
	{ locked = true, description = "Audio:Toggle mic" }
)
hl.bind(
	"CTRL + SUPER + T",
	hl.dsp.exec_cmd(ipc .. "panel-toggle wallpaper"),
	{ description = "Wallpaper:Toggle wallpaper selector" }
)
hl.bind(
	"CTRL + SUPER + ALT + T",
	hl.dsp.exec_cmd(ipc .. "wallpaper-random"),
	{ description = "Wallpaper:Select random wallpaper" }
)
hl.bind(
	"CTRL + SUPER + R",
	hl.dsp.exec_cmd("pkill noctalia; noctalia"),
	{ description = "Widgets:Restart shell" }
)

-- Clipboard history
hl.bind(
	"SUPER + V",
	hl.dsp.exec_cmd(ipc .. "panel-toggle clipboard"),
	{ description = "Clipboard:History >> clipboard" }
)
hl.bind(
	"SUPER + Period",
	hl.dsp.exec_cmd("pkill fuzzel || " .. hyprScripts .. "/fuzzel-emoji.sh"),
	{ description = "Clipboard:Emoji >> clipboard" }
)

-- Color picker
hl.bind(
	"SUPER + SHIFT + C",
	hl.dsp.exec_cmd("hyprpicker -a"),
	{ description = "Clipboard:Pick color #RRGGBB >> clipboard" }
)

-- Screenshot
hl.bind(
	"SUPER + SHIFT + S",
	hl.dsp.exec_cmd(ipc .. "screenshot-region"),
	{ description = "Screenshot:Region snip >> clipboard" }
)

-- AI
hl.bind(
	"SUPER + SHIFT + A",
	hl.dsp.exec_cmd("pidof slurp || " .. hyprScripts .. "/snip_to_search.sh"),
	{ description = "AI:Google Lens" }
)

-- OCR
hl.bind(
	"SUPER + SHIFT + X",
	hl.dsp.exec_cmd(
		"pidof slurp || grim -g \"$(slurp $SLURP_ARGS)\" \"/tmp/ocr_image.png\" && tesseract \"/tmp/ocr_image.png\" stdout -l $(tesseract --list-langs | awk 'NR>1{print $1}' | tr '\\\\n' '+' | sed 's/\\\\+$/\\\\n/') | wl-copy && rm \"/tmp/ocr_image.png\""
	),
	{ description = "OCR:Character recognition >> clipboard" }
)

-- Fullscreen screenshot
local grimhyprctl = "grim -o \"$(hyprctl activeworkspace -j | jq -r '.monitor')\""
hl.bind(
	"Print",
	hl.dsp.exec_cmd(grimhyprctl .. " - | wl-copy"),
	{ locked = true, description = "Screenshot:Fullscreen >> clipboard" }
)
hl.bind(
	"CTRL + Print",
	hl.dsp.exec_cmd(
		"mkdir -p $(xdg-user-dir PICTURES)/Screenshots && "
			.. grimhyprctl
			.. " $(xdg-user-dir PICTURES)/Screenshots/Screenshot_\"$(date '+%Y-%m-%d_%H.%M.%S')\".png"
	),
	{ locked = true, non_consuming = true, description = "Screenshot:Fullscreen >> clipboard & file" }
)
hl.bind("CTRL + Print", hl.dsp.exec_cmd(grimhyprctl .. " - | wl-copy"), { locked = true, non_consuming = true })

-- Window
-- Focusing
hl.bind("SUPER + mouse:272", hl.dsp.window.drag(), { mouse = true, description = "Windows:Move window (drag)" })
hl.bind("SUPER + mouse:274", hl.dsp.window.drag(), { mouse = true })
hl.bind("SUPER + mouse:273", hl.dsp.window.resize(), { mouse = true, description = "Windows:Resize window (drag)" })
--#/# bind = SUPER + ←/↑/→/↓,, -- Focus in direction
for i = 1, 6 do
	local arrowkey = { "Left", "Right", "Up", "Down", "BracketLeft", "BracketRight" }
	local focusdir = { "l", "r", "u", "d", "l", "r" }
	local focusname = { "left", "right", "up", "down" }
	local options = focusname[i] and { description = "Windows:Focus " .. focusname[i] } or {}
	hl.bind("SUPER + " .. arrowkey[i], hl.dsp.focus({ direction = focusdir[i] }), options)
end
--#/# bind = SUPER + SHIFT, ←/↑/→/↓,, -- Move in direction
for i = 1, 4 do
	local arrowkey = { "Left", "Right", "Up", "Down" }
	local focusdir = { "l", "r", "u", "d" }
	local movename = { "left", "right", "up", "down" }
	hl.bind(
		"SUPER + SHIFT + " .. arrowkey[i],
		hl.dsp.window.move({ direction = focusdir[i] }),
		{ description = "Windows:Move window " .. movename[i] }
	)
end

hl.bind("ALT + F4", function()
	hl.exec_cmd('notify-send "Wrong close keybind" "Super+Q to close. Use Alt+F4 for Windows VMs" -a Hyprland')
end, { non_consuming = true })
hl.bind("SUPER + Q", hl.dsp.window.close(), { description = "Windows:Close window" })
hl.bind("SUPER + SHIFT + ALT + Q", hl.dsp.exec_cmd("hyprctl kill"), { description = "Windows:Force kill window" })

-- Window split ratio
-- bind = SUPER, ;/',, -- Adjust split ratio
hl.bind(
	"SUPER + Semicolon",
	hl.dsp.layout("splitratio -0.1"),
	{ repeating = true, description = "Layout:Decrease split ratio" }
)
hl.bind(
	"SUPER + Apostrophe",
	hl.dsp.layout("splitratio +0.1"),
	{ repeating = true, description = "Layout:Increase split ratio" }
)
-- Positioning mode
hl.bind("SUPER + T", hl.dsp.window.float({ action = "toggle" }), { description = "Layout:Toggle floating/tiling" })
-- hl.bind("SUPER + D", hl.dsp.window.fullscreen({ "maximized" }), { description = "Layout:Maximize" })
hl.bind("SUPER + F", hl.dsp.window.fullscreen({ "fullscreen" }), { description = "Layout:Fullscreen" })
hl.bind(
	"SUPER + ALT + F",
	hl.dsp.window.fullscreen_state({ internal = 0, client = 3 }),
	{ description = "Layout:Fullscreen spoof (client only)" }
)
hl.bind("SUPER + P", hl.dsp.window.pseudo(), { description = "Layout:Pseudo window" })
hl.bind("SUPER + J", hl.dsp.layout("togglesplit"), { description = "Layout:Toggle window split" })
hl.bind(
	"SUPER + O",
	hl.dsp.exec_cmd("$HOME/.config/hypr/scripts/window_pop.sh"),
	{ description = "Layout:Pop window out (float & pin)" }
)
hl.bind(
	"SUPER + L",
	hl.dsp.exec_cmd("$HOME/.config/hypr/scripts/workspace_layout.sh"),
	{ description = "Layout:Toggle workspace layout (dwindle/scrolling)" }
)

-- Groups
hl.bind("SUPER + G", hl.dsp.group.toggle(), { description = "Groups:Toggle window grouping" })
hl.bind(
	"SUPER + ALT + G",
	hl.dsp.window.move({ out_of_group = true }),
	{ description = "Groups:Move window out of group" }
)
-- bind = SUPER+ALT, ←/↑/→/↓,, -- Move into the neighbouring group
for i = 1, 4 do
	local arrowkey = { "Left", "Right", "Up", "Down" }
	local groupdir = { "l", "r", "u", "d" }
	local groupname = { "left", "right", "up", "down" }
	hl.bind(
		"SUPER + ALT + " .. arrowkey[i],
		hl.dsp.window.move({ into_group = groupdir[i] }),
		{ description = "Groups:Move window into the group " .. groupname[i] }
	)
end
hl.bind("SUPER + ALT + Tab", hl.dsp.group.next(), { description = "Groups:Next window in group" })
hl.bind("SUPER + SHIFT + ALT + Tab", hl.dsp.group.prev(), { description = "Groups:Previous window in group" })
hl.bind("CTRL + SUPER + Right", hl.dsp.group.next())
hl.bind("CTRL + SUPER + Left", hl.dsp.group.prev())
hl.bind("SUPER + ALT + mouse_down", hl.dsp.group.next())
hl.bind("SUPER + ALT + mouse_up", hl.dsp.group.prev())

-- bind = SUPER+ALT, Hash,, -- Send to workspace -- (1, 2, 3,...)
-- We use raw keycodes because some keyboard layouts register number keys as different chars. The codes can be verified with `wev`
for i = 1, 10 do
	local numberkey = { 10, 11, 12, 13, 14, 15, 16, 17, 18, 19 }
	local options = i == 1 and { description = "Workspace:Send window to workspace 1-10" } or {}
	hl.bind("SUPER + ALT + code:" .. numberkey[i], hl.dsp.window.move({ workspace = i, follow = false }), options)
end
-- keypad numbers
for i = 1, 10 do
	local numpadkey = { 87, 88, 89, 83, 84, 85, 79, 80, 81, 90 }
	hl.bind("SUPER + ALT + code:" .. numpadkey[i], hl.dsp.window.move({ workspace = i, follow = false }))
end

-- bind = SUPER+SHIFT, Hash,, -- Send to workspace and follow -- (1, 2, 3,...)
for i = 1, 10 do
	local numberkey = { 10, 11, 12, 13, 14, 15, 16, 17, 18, 19 }
	local options = i == 1 and { description = "Workspace:Send window to workspace 1-10 and follow" } or {}
	hl.bind("SUPER + SHIFT + code:" .. numberkey[i], hl.dsp.window.move({ workspace = i }), options)
end
-- keypad numbers
for i = 1, 10 do
	local numpadkey = { 87, 88, 89, 83, 84, 85, 79, 80, 81, 90 }
	hl.bind("SUPER + SHIFT + code:" .. numpadkey[i], hl.dsp.window.move({ workspace = i }))
end

hl.bind(
	"SUPER + ALT + S",
	hl.dsp.window.move({ workspace = "special:scratchpad", follow = false }),
	{ description = "Workspace:Send window to scratchpad" }
)
hl.bind("CTRL + SUPER + S", hl.dsp.workspace.toggle_special("scratchpad"))

--##! Workspace
--# Switching
--#/# bind = SUPER, Hash,, -- Focus workspace -- (1, 2, 3,...)
--# We use raw keycodes because some keyboard layouts register number keys as different chars. The codes can be verified with `wev`
for i = 1, 10 do
	local numberkey = { 10, 11, 12, 13, 14, 15, 16, 17, 18, 19 }
	local options = i == 1 and { description = "Workspace:Focus workspace 1-10" } or {}
	hl.bind("SUPER + code:" .. numberkey[i], hl.dsp.focus({ workspace = i }), options)
end
--# keypad numbers
for i = 1, 10 do
	local numpadkey = { 87, 88, 89, 83, 84, 85, 79, 80, 81, 90 }
	hl.bind("SUPER + code:" .. numpadkey[i], hl.dsp.focus({ workspace = i }))
end

--## Special
hl.bind("SUPER + S", hl.dsp.workspace.toggle_special("scratchpad"), { description = "Workspace:Toggle scratchpad" })
hl.bind("SUPER + mouse:275", hl.dsp.workspace.toggle_special("scratchpad"))

--##! Virtual machines
hl.define_submap("virtual-machine", function()
	hl.bind("SUPER + ALT + F1", function()
		local currentsubmap = hl.get_current_submap()
		if currentsubmap == "virtual-machine" then
			hl.dispatch(
				hl.dsp.exec_cmd("notify-send 'Exited Virtual Machine submap' 'Keybinds re-enabled' -a 'Hyprland'")
			)
			hl.dispatch(hl.dsp.submap("reset"))
		elseif currentsubmap == "" then
			hl.dispatch(
				hl.dsp.exec_cmd(
					"notify-send 'Entered Virtual Machine submap' 'Keybinds disabled. hit SUPER+ALT+F1 to escape' -a 'Hyprland'"
				)
			)
			hl.dispatch(hl.dsp.submap("virtual-machine"))
		end
	end, { submap_universal = true, description = "Session:Toggle virtual machine mode" })
end)

--#!
--# Testing
hl.bind(
	"SUPER + ALT + F11",
	hl.dsp.exec_cmd(
		'bash -c \'RANDOM_IMAGE=$(find ~/Pictures -type f | shuf -n 1); ACTION=$(notify-send "Test notification with body image" "This notification should contain your user account <b>image</b> and <a href=\\"https://discord.com/app\\">Discord</a> <b>icon</b>. Oh and here is a random image in your Pictures folder: <img src=\\"$RANDOM_IMAGE\\" alt=\\"Testing image\\"/>" -a "Hyprland" -p -h "string:image-path:/var/lib/AccountsService/icons/$USER" -t 6000 -i "discord" -A "openImage=Profile image" -A "action2=Open the random image" -A "action3=Useless button"); [[ $ACTION == *openImage ]] && xdg-open "/var/lib/AccountsService/icons/$USER"; [[ $ACTION == *action2 ]] && xdg-open "$RANDOM_IMAGE"\''
	)
) -- # [hidden]
hl.bind(
	"SUPER + ALT + F12",
	hl.dsp.exec_cmd(
		'bash -c \'RANDOM_IMAGE=$(find ~/Pictures -type f | shuf -n 1); ACTION=$(notify-send "Test notification" "This notification should contain a random image in your <b>Pictures</b> folder and <a href=\\"https://discord.com/app\\">Discord</a> <b>icon</b>.\n<i>Flick right to dismiss!</i>" -a "Discord (fake)" -p -h "string:image-path:$RANDOM_IMAGE" -t 6000 -i "discord" -A "openImage=Profile image" -A "action2=Useless button"); [[ $ACTION == *openImage ]] && xdg-open "/var/lib/AccountsService/icons/$USER"\''
	)
) -- # [hidden]
hl.bind(
	"SUPER + ALT + Equal",
	hl.dsp.exec_cmd("notify-send 'Urgent notification' 'Ah hell no' -u critical -a 'Hyprland keybind'")
) -- # [hidden]

--##! Session
hl.bind("CTRL + SUPER + L", hl.dsp.exec_cmd(ipc .. "session lock"), { description = "Session:Lock" })
hl.bind(
	"SUPER + SHIFT + L",
	hl.dsp.exec_cmd(ipc .. "session lock-and-suspend"),
	{ locked = true, description = "Session:Suspend system" }
) -- Sleep
-- hl.bind("switch:on:Lid Switch", hl.dsp.exec_cmd("systemctl suspend || loginctl suspend"), {locked = true} ) -- # [hidden] Suspend when laptop lid is closed, uncomment if for whatever reason it's not the default behavior

hl.bind(
	"CTRL + SHIFT + ALT + SUPER + Delete",
	hl.dsp.exec_cmd("systemctl poweroff || loginctl poweroff"),
	{ description = "Session:Shutdown" }
) -- # [hidden] Power off

--##! Screen
--# Zoom
local function zoomfunction(value)
	local zoomvalue = hl.get_config("cursor:zoom_factor")
	if (zoomvalue + value) > 3.0 then
		hl.config({ cursor = { zoom_factor = 3.0 } })
	elseif (zoomvalue + value) < 1.0 then
		hl.config({ cursor = { zoom_factor = 1.0 } })
	else
		hl.config({ cursor = { zoom_factor = zoomvalue + value } })
	end
end
hl.bind("SUPER + Minus", function()
	zoomfunction(-0.3)
end, { repeating = true, description = "Screen:Zoom out" })
hl.bind("SUPER + Equal", function()
	zoomfunction(0.3)
end, { repeating = true, description = "Screen:Zoom in" })

--# Zoom with keypad
hl.bind("SUPER + code:82", function()
	zoomfunction(-0.3)
end, { repeating = true })
hl.bind("SUPER + code:86", function()
	zoomfunction(0.3)
end, { repeating = true })

--##! Media
local mediaNextCommand =
	'playerctl next || playerctl position `bc <<< "100 * $(playerctl metadata mpris:length) / 1000000 / 100"`'
hl.bind("SUPER + SHIFT + N", hl.dsp.exec_cmd(mediaNextCommand), { locked = true, description = "Media:Next track" })
hl.bind("XF86AudioNext", hl.dsp.exec_cmd(mediaNextCommand), { locked = true })
hl.bind("XF86AudioPrev", hl.dsp.exec_cmd("playerctl previous"), { locked = true })
hl.bind("SUPER + SHIFT + ALT + mouse:275", hl.dsp.exec_cmd("playerctl previous"))
hl.bind("SUPER + SHIFT + ALT + mouse:276", hl.dsp.exec_cmd(mediaNextCommand))
hl.bind(
	"SUPER + SHIFT + B",
	hl.dsp.exec_cmd("playerctl previous"),
	{ locked = true, description = "Media:Previous track" }
)
hl.bind(
	"SUPER + SHIFT + P",
	hl.dsp.exec_cmd("playerctl play-pause"),
	{ locked = true, description = "Media:Play/pause" }
)
hl.bind("XF86AudioPlay", hl.dsp.exec_cmd("playerctl play-pause"), { locked = true })
hl.bind("XF86AudioPause", hl.dsp.exec_cmd("playerctl play-pause"), { locked = true })
--##! Apps
hl.bind("SUPER + Return", hl.dsp.exec_cmd(terminal), { description = "Apps:Terminal" })
hl.bind("CTRL + ALT + T", hl.dsp.exec_cmd(terminal))
hl.bind("SUPER + E", hl.dsp.exec_cmd(fileManager), { description = "Apps:File manager" })
hl.bind("SUPER + W", hl.dsp.exec_cmd(browser), { description = "Apps:Browser" })
hl.bind("SUPER + C", hl.dsp.exec_cmd(codeEditor), { description = "Apps:Code editor" })
hl.bind("CTRL + SUPER + SHIFT + ALT + W", hl.dsp.exec_cmd(officeSoftware), { description = "Apps:Office software" })
hl.bind("SUPER + X", hl.dsp.exec_cmd(textEditor), { description = "Apps:Text editor" })
hl.bind("CTRL + SUPER + V", hl.dsp.exec_cmd(volumeMixer), { description = "Apps:Volume mixer" })
hl.bind("SUPER + I", hl.dsp.exec_cmd(settingsApp), { description = "Apps:Settings app" })
hl.bind("CTRL + SHIFT + Escape", hl.dsp.exec_cmd(taskManager), { description = "Apps:Task manager" })

--# Cursed stuff
--## Make window not amogus large
hl.bind(
	"CTRL + SUPER + Backslash",
	hl.dsp.window.resize({ x = 640, y = 480, "exact" }),
	{ description = "Windows:Resize window to 640x480" }
)
