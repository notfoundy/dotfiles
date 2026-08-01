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

-- Shell Overides
require("shellOverides/main")
