package langs

import (
	"errors"
	"github.com/THUNDERGROOVE/SDETool/scripting"
	"github.com/THUNDERGROOVE/SDETool/util/log"
)

var ScriptingLangs = make(map[string]scripting.ScriptingLang)

var NoSuchLang = errors.New("No such language registered")

func Register(name string, lang scripting.ScriptingLang) {
	log.Println("[LANGS] Registering", name)
	ScriptingLangs[name] = lang
}

func RunString(lang, s string) error {
	if v, ok := ScriptingLangs[lang]; ok {
		if !v.IsInit() {
			v.Init()
		}
		return v.RunString(s)
	} else {
		return NoSuchLang
	}
	return nil
}

func RunScript(lang, filename string) error {
	if v, ok := ScriptingLangs[lang]; ok {
		if !v.IsInit() {
			v.Init()
		}
		return v.RunScript(filename)
	} else {
		return errors.New(NoSuchLang.Error() + ": '" + lang + "'")
	}
	return nil
}

func Interpreter(lang string) error {
	if v, ok := ScriptingLangs[lang]; ok {
		if !v.IsInit() {
			v.Init()
		}
		return v.Interpreter()
	} else {
		return NoSuchLang
	}
	return nil
}
