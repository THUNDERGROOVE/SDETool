package lua

import (
	"fmt"
	"log"

	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/THUNDERGROOVE/SDETool/sde/version"
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// SDE is the global SDE object used by functions in the lua runtime
var SDE *sde.SDE

// @TODO:  Reimplement version api using new system

func SDELoader(l *lua.LState)int {
	// We want to take advantage of the magic of luar for functions that
	// take or provide complex go types
	_luaApplyType := luar.New(l, applyType)
	_luaGetTypeByID := luar.New(l, getTypeByID)
	
	var exports = map[string]lua.LGFunction{
		"load":_luaLoad,
		"loadLatest":_luaLoadLatest,
		"search":_luaSearch,
	}

	mod := l.NewTable() 
	l.SetFuncs(mod, exports, lua.LString("value"))

	// Maybe use a map if we get too many of these?
	l.SetField(mod, "applyType", _luaApplyType)
	l.SetField(mod, "getTypeByID", _luaGetTypeByID)

	l.Push(mod)
	return 1
}

func _luaLoadLatest(l *lua.LState) int {
	sde, err := version.LoadLatest()
	if err == nil {
		SDE = sde
		l.Push(lua.LTrue)
		return 1
	}
	l.Push(lua.LFalse)
	return 1
}

func _luaLoad(l *lua.LState) int {
	filename := l.ToString(1)
	s, err := sde.Load(filename)
	if err != nil {
		fmt.Println("couldn't open SDE file", err.Error())
		l.Push(lua.LFalse)
		return 1
	}
	SDE = s
	l.Push(lua.LTrue)
	return 1
}

func applyType(original *sde.SDEType, newType *sde.SDEType) *sde.SDEType {
	out, err := sde.ApplyTypeToType(*original, *newType)
	if err != nil {
		fmt.Println("ERROR applying type", err.Error())
	}
	if out == nil {
		fmt.Println("ERROR applyType got a nil SDEType")
	}
	return out
}

func getTypeByID(ID int) *sde.SDEType {
	var t *sde.SDEType
	var err error
	if t, err = SDE.GetType(ID); err == nil {
		return t
	}
	log.Println("[LUA][MOD] sde.getTypeByID; returned SDEType had error", err.Error())
	return nil
}

func _luaSearch(l *lua.LState) int {
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

