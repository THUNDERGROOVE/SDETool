print("Hello SDETool\nLet's open an SDE version")
sde.loadVersion("1.0-WL")
t = sde.getTypeByID(364022)
print(t.Attributes.mDisplayName.." contains "..#t.Attributes.. " attributes")