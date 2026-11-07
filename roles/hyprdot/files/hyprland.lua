---@diagnostic disable: undefined-global
-- Required for global submaps
hl.dsp.submap("global")

-- Sources
require("env")
require("cursor")
require("variables")
require("execs")
require("general")
require("rules")
require("colors")
require("keybinds")

-- nwg-displays support
require("workspaces")
require("monitors")

-- Absent until Noctalia's first palette change
if pcall(require, "noctalia") then
	require("noctalia").apply_theme()
end
