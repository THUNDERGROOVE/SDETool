package langs

import (
	"errors"
	"log"

	"github.com/THUNDERGROOVE/SDETool/scripting"
)

// ScriptingLangs is a global map that contains all of out registered scripting
// languages
var ScriptingLangs = make(map[string]scripting.ScriptingLang)

// ErrNoSuchLang is called when a language interface function is called on a
// language that isn't registered or doesn't exist.
var ErrNoSuchLang = errors.New("No such language registered")

// Register registers a scripting language for use
func Register(name string, lang scripting.ScriptingLang) {
	log.Println("[LANGS] Registering", name)
	ScriptingLangs[name] = lang
}

// RunString runs the string for the scripting language supplied
func RunString(lang, s string) error {
	if v, ok := ScriptingLangs[lang]; ok {
		if !v.IsInit() {
			v.Init()
		}
		return v.RunString(s)
	}
	return ErrNoSuchLang
}

// RunScript runs the script for the given file
func RunScript(lang, filename string) error {
	if v, ok := ScriptingLangs[lang]; ok {
		if !v.IsInit() {
			v.Init()
		}
		return v.RunScript(filename)
	}
	return errors.New(ErrNoSuchLang.Error() + ": '" + lang + "'")
}

// Interpreter should open an interpreter for the provided language
func Interpreter(lang string) error {
	if v, ok := ScriptingLangs[lang]; ok {
		if !v.IsInit() {
			v.Init()
		}
		return v.Interpreter()
	}
	return ErrNoSuchLang
}
