local sde = require("sde")
print("Loading an SDE file")
sde.loadLatest()

print("Looking up an Amarr Assault")
local a = sde.getTypeByID(364022)
print("Original speed: "..a.Attributes["mVICProp.groundSpeed"])

print("Looking up a complex armor plate")
local p = sde.getTypeByID(351671)

print("Applying the plate to the assault")
local n = sde.applyType(a, p)
print("New speed: "..n.Attributes["mVICProp.groundSpeed"])
