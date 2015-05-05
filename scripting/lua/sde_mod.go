package lua

import (
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/THUNDERGROOVE/SDETool/util/log"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

var SDE sde.SDE

func Loader(l *lua.LState) {
	var exports = map[string]lua.LValue{
		"getVersions": luar.New(l, getVersions),
		"loadVersion": luar.New(l, loadVersion),
		"getTypeByID": luar.New(l, getTypeByID),
	}
	l.SetGlobal("sde", l.NewTable())
	tbl := l.GetGlobal("sde")
	for k, v := range exports {
		l.SetField(tbl, k, v)
	}
}

func getVersions() []string {
	out := make([]string, 0)
	for k, _ := range sde.Versions {
		out = append(out, k)
	}
	return out
}

func loadVersion(version string) error {
	var err error
	SDE, err = sde.Open(version)

	return err
}

func getTypeByID(ID int) sde.SDEType {
	if SDE.DB != nil {
		if t, err := SDE.GetType(ID); err == nil {
			return t
		} else {
			log.Println("[LUA][MOD] sde.getTypeByID; returned SDEType had error", err.Error())
			return sde.SDEType{}
		}
	} else {
		log.Println("[LUA][MOD] sde.getTypeByID called with no SDE loaded.")
	}

	return sde.SDEType{}
}
