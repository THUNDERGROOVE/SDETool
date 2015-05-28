package commands

import (
	"fmt"

	"github.com/THUNDERGROOVE/SDETool/args"
)

const halp = `SDETool.  Your one stop shop for everything DUST SDE related.. Hopefully...

Commands:
  search <Name/TypeName/TypeID>:
	Prints some all matches we can find with their typeIDs
  lookup <TypeID> {attrs}:
	Prints the TypeID, TypeName and display name of a given typeID.
	If given the subcommand "attrs" we will pretty print the attributes
	for the type.

--help(-h):
	You're looking at it.
--sde(-s) <Filename>:
	Loads an SDE version from file bypassing the internal version system
--sde-info(-i):
	Prints various info about the currently loaded SDE file
`

func RegisterCommands(tokens args.Tokens) {
	RegisterSDE(tokens)
	args.FlagReg.Register("--help", "-h", func(tok args.TokLit, index int) {
		fmt.Printf(halp)
	})
}
