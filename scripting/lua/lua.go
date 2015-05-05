package lua

import (
	"github.com/THUNDERGROOVE/SDETool/scripting"
	"github.com/THUNDERGROOVE/SDETool/scripting/langs"
	"github.com/THUNDERGROOVE/SDETool/util/log"
	"github.com/yuin/gopher-lua"
)

func init() {
	langs.Register("lua", &Lua{})
}

// Lua implements ScriptingLang
type Lua struct {
	state *lua.LState
}

func (l *Lua) RunScript(filename string) error {
	if l.state == nil {
		log.Println("WARNING; lua state is nil")
	}
	return l.state.DoFile(filename)
}

func (l *Lua) RunString(s string) error {
	return l.state.DoString(s)
}

func (l *Lua) Interpreter() error {
	log.Println("[LUA] Not implemented")
	return scripting.InterpreterNotImplemented
}

func (l *Lua) Init() error {
	log.Println("[LUA] Starting up GopherLua VM")
	l.state = lua.NewState()
	if l.state == nil {
		log.Println("[LUA] Error creating VM")
	}
	Loader(l.state)
	return nil
}

func (l *Lua) IsInit() bool {
	if l.state == nil {
		return false
	}
	return true
}
