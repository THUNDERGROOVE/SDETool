package main

import (
	"flag"
)

// @TODO:  Maybe convert each flagset into a struct

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
)
