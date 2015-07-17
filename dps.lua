
-- RoF in rounds/minute
local function GetRoF(t)
	if t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_FullAuto" then
		return rpmtorof(t.Attributes["mFireMode0.fireInterval"])
	elseif t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_ChargeToFire" then
		print(t.Attributes["mFireMode0.fireInterval"])
		print(t.Attributes["m_ChargeInfo.m_fChargeUpTime"])
		return rpmtorof(t.Attributes["mFireMode0.fireInterval"] - t.Attributes["m_ChargeInfo.m_fChargeUpTime"])
	elseif t.Attributes["mFireMode0.m_eFireMode"] == "DWFM_SingleBurst" then
		return rpmtorof(t.Attributes["mFireMode0.fireInterval"])
	else
		print("Unknown firemode"..t.Attributes["mFireMode0.m_eFireMode"])
		return 0
	end
end

local function rpmtorof(t)
  t = t or 0
  return (1*t) / 60
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
	local burst = types[v].Attributes["m_BurstInfo.m_iBurstLength"]
	local damage = types[v].Attributes["mFireMode0.instantHitDamage"]
	print("You picked: "..types[v].Attributes.mDisplayName)
	print("FireMode: "..types[v].Attributes["mFireMode0.m_eFireMode"])
	print("Burst: "..types[v].Attributes["m_BurstInfo.m_fBurstInterval"] .."  Interval: ".. types[v].Attributes["mFireMode0.fireInterval"])
	local rof = GetRoF(types[v])
	local dps = ((damage * burst) * rof) / 60 
	print("RoF: "..rof)
	print("Burst length: "..burst)
	print("Damage: "..types[v].Attributes["mFireMode0.instantHitDamage"])
	print("DPS: "..dps)
	print("Math: (("..damage.." * "..burst..") * "..rof..") / 60")
end
