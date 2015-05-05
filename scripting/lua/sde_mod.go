package lua

import (
	"fmt"
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
	tbl := l.NewTable()
	l.SetGlobal("sde", tbl)
	for k, v := range exports {
		l.SetField(tbl, k, v)
	}

	l.SetFuncs(tbl, map[string]lua.LGFunction{
		"search": search})
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

func search(l *lua.LState) int {
	v := l.ToString(1)
	t := l.NewTable()
	if res, err := SDE.Search(v); err == nil {
		for _, v := range res {
			t.Append(luar.New(l, *v))
		}
	} else {
		fmt.Println("Error encountered in sde.search:", err.Error())
	}
	l.Push(t)
	return 1
}
