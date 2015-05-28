
-- RoF in rounds/minute
local function GetRoF(t)
	if t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_FullAuto" then
		return t.Attributes["mFireMode0.fireInterval"] * 10000
	elseif t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_ChargeToFire" then
		print(t.Attributes["mFireMode0.fireInterval"])
		print(t.Attributes["m_ChargeInfo.m_fChargeUpTime"])
		return (t.Attributes["mFireMode0.fireInterval"] - t.Attributes["m_ChargeInfo.m_fChargeUpTime"]) * 10000
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

print("Downloading SDE file")
sde.load("dust-wl-11.sde")

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
	print("Damage: "..types[v].Attributes["mFireMode0.instantHitDamage"])
	print("DPS: "..dps)
end
