local sde = require("sde")
print("Hello SDETool")
print("Loading latest SDE file")
sde.loadLatest()
print("Done! Looking up an Amarr Assault ak.0")

local t = sde.getTypeByID(364022)
print(t.Attributes.mDisplayName.." contains "..#t.Attributes.. " attributes")
