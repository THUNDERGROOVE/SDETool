print("Hello SDETool\nLet's open an SDE version")


print("Pick a verison:")
v = sde.getVersions()

for i = 1, #v do
	print("  "..i..") "..v[i])
end

while true do
	io.write("> ")
	ver = tonumber(io.read())
	if ver == nil then print("Invalid input") else break end
end

sde.loadVersion(v[ver])

t = sde.getTypeByID(364022)
print(t.Attributes.mDisplayName.." contains "..#t.Attributes.. " attributes")