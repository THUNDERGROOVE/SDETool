package scripting

import (
	"errors"
)

// InterpreterNotImplemented is returned when SDETool tries using an interpreter
// that doen't exist
var InterpreterNotImplemented = errors.New("interpreter not implemented")

// ScriptingLang is an interface to allow expandability for other scripting languages.
//
// RunScript is called when the user choses to run a specific script
//
// RunString is called when the user inputs a string to be run
//
// Interpreter is called when the user choses to open a interpreter.
//   You can return InterpreterNotImplemented
//   And we'll know it's not implemented and notify the user.
//
// Init is called whenever your language is first fired up.  But only when IsInit returns false
type ScriptingLang interface {
	RunScript(filename string) error
	RunString(s string) error
	Interpreter() error
	Init() error
	IsInit() bool
}
