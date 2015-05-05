package sde

import (
	"time"
)

type Version struct {
	URL    string
	Zipped bool
}

var (
	sde17     = Version{"http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.7_674383.zip", true}
	sde18     = Version{"http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_752135.zip", true} // Old: http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_739147.zip
	sde18D    = Version{"http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.8_851720.zip", true}
	sde19     = Version{"http://cdn1.eveonline.com/community/DUST_SDE/Uprising_1.9_853193.zip", true}
	sdewl10   = Version{"http://cdn1.eveonline.com/community/DUST_SDE/Warlords_1.0_857519.zip", true}
	sdewl10UF = Version{"http://dl.dropboxusercontent.com/u/51472257/sde.WL-1.0-02-12.db", false}
	sdeDebug  = Version{"nil", true}
)

var (
	// Versions is a map of all of the available versions.
	Versions map[string]Version
)

func init() {
	defer Debug(time.Now())
	Versions = make(map[string]Version, 0)
	Versions["1.7"] = sde17
	Versions["1.8"] = sde18
	Versions["1.8-delta"] = sde18D
	Versions["1.9"] = sde19
	Versions["1.0-WL"] = sdewl10
	Versions["debug"] = sdeDebug
	Versions["wl-uf-latest"] = sdewl10UF
}
