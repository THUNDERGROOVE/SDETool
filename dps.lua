
local function GetRoF(t)
	if t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_FullAuto" or t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_ChargeToFire" then
		return t.Attributes["mFireMode0.fireInterval"] * 10000
	elseif t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_SingleBurst" then
		return t.Attributes["m_BurstInfo.m_fBurstInterval"] * 10000
	else
		print("Unknown firemode"..t.Attributes["mFireMode0.m_eFireMode"])
		return 0
	end
end

print("--------------")
print("DPS calculator")
print("--------------")
print("\nPress ctrl+c to exist at any time.\n")

print("Pick a verison:")
local v = sde.getVersions()

for i = 1, #v do
	print("  "..i..") "..v[i])
end

while true do
	io.write("> ")
	ver = tonumber(io.read())
	if ver == nil then print("Invalid input") else break end
end

sde.loadVersion(v[ver])

while true do
	print("Search:")
	io.write("> ")
	local t = io.read()
	local types = sde.search(t)
	print("Search results:")
	for i = 1, #types do
		print("  "..i..") "..types[i].Attributes.mDisplayName)
	end
	io.write("> ")
	local v = tonumber(io.read())
	if v == nil then print("Invalid input") break end 
	if types[v] == nil then print("Invalid input") break end
	print("You picked: "..types[v].Attributes.mDisplayName)
	print("FireMode: "..types[v].Attributes["mFireMode0.m_eFireMode"])
	print("Burst: "..types[v].Attributes["m_BurstInfo.m_fBurstInterval"] .."  Interval: ".. types[v].Attributes["mFireMode0.fireInterval"])
	local rof = GetRoF(types[v])
	local dps = (types[v].Attributes["mFireMode0.instantHitDamage"] * rof) / 60 
	print("RoF: "..rof)
	print("DPS: "..dps)
end