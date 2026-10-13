---@diagnostic disable: undefined-global
-- This file sources other source files

-- Required for global submaps
hl.dsp.submap("global")

-- Sources
require("env")
require("variables")
require("execs")
require("general")
require("rules")
require("colors")
require("keybinds")

-- nwg-displays support
require("workspaces")
require("monitors")

-- Noctalia colour template, generated on first palette change
if pcall(require, "noctalia") then
	require("noctalia").apply_theme()
end
