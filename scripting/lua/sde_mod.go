package lua

import (
	"fmt"
	"log"
	"runtime"

	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

var SDE sde.SDE

func Loader(l *lua.LState) {
	var exports = map[string]lua.LValue{
		"getVersions": luar.New(l, getVersions),
		"loadVersion": luar.New(l, loadVersion),
		"getTypeByID": luar.New(l, getTypeByID),
		"applyType":   luar.New(l, applyType),
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
	depreciated("sde package no longer uses a version system")

	out := make([]string, 0)
	/*
		for k, _ := range sde.Versions {
			out = append(out, k)
		}
	*/
	return out
}

func applyType(original sde.SDEType, newType sde.SDEType) sde.SDEType {
	out, err := sde.ApplyTypeToType(original, newType)
	if err != nil {
		fmt.Println("ERROR applying type", err.Error())
	}
	if out == nil {
		fmt.Println("ERROR applyType got a nil SDEType")
	}
	return *out
}

func loadVersion(version string) error {
	depreciated("sde package no longer uses a version system")
	return nil
	/*
		var err error
		SDE, err = sde.Open(version)
		return err
	*/
}

func getTypeByID(ID int) *sde.SDEType {

	if t, err := SDE.GetType(ID); err == nil {
		return t
	} else {
		log.Println("[LUA][MOD] sde.getTypeByID; returned SDEType had error", err.Error())
		return nil
	}

	return nil
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

func depreciated(message string) {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fmt.Printf("Function: %v is depreciated\n%v\n", f.Name(), message)
}
