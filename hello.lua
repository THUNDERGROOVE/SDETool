print("Hello SDETool")
print("Downloading an SDE file...")
sde.loadHTTP("https://dl.dropboxusercontent.com/u/51472257/dust-wl-11.sde")
print("Done! Looking up an Amarr Assault ak.0")

t = sde.getTypeByID(364022)
print(t.Attributes.mDisplayName.." contains "..#t.Attributes.. " attributes")
