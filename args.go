package main

import (
	"flag"
)

// @TODO:  Maybe convert each flagset into a struct
// Don't see a real benefit at this point other than ease of reading and not
// over populating the global space
const (
	sdeFlagConst     = "An SDE file to load"
	tidFlagConst     = "A TypeID to lookup"
	typeNameConst    = "A TypeName to lookup"
	displayNameConst = "A display name to lookup"
)

var (
	lookupFlagset = flag.NewFlagSet("lookup", flag.ContinueOnError)

	lookupSDE  = lookupFlagset.String("sde", "", sdeFlagConst)
	lookupTID  = lookupFlagset.Int("t", 0, tidFlagConst)
	lookupTN   = lookupFlagset.String("tn", "", typeNameConst)
	lookupTD   = lookupFlagset.String("td", "", displayNameConst)
	lookupAttr = lookupFlagset.Bool("attr", false, "Print type attributes")

	searchFlagset = flag.NewFlagSet("search", flag.ContinueOnError)

	searchSDE  = searchFlagset.String("sde", "", sdeFlagConst)
	searchName = searchFlagset.String("t", "", "A string to search for")
	searchAttr = searchFlagset.Bool("attr", false, "Print's type attributes only if one type is returned by the search")

	dumperFlagset = flag.NewFlagSet("dumper", flag.ContinueOnError)

	dumperInFile        = dumperFlagset.String("i", "", "The input SQLite3 database file the SDE resides in")
	dumperOutFile       = dumperFlagset.String("o", "", "The output file that is compatible with SDETool")
	dumperVersionString = dumperFlagset.String("ver", "", "The string encoded into the dump")
	dumperOfficial      = dumperFlagset.Bool("official", false, "Set this flag if the SQLite db is from CCP")
	dumperVerbose       = dumperFlagset.Bool("v", false, "Print more stuffs")
)
