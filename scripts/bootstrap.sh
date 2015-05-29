#!/bin/bash

# This script sets up a nice environment for you to work on SDETool

printf "Checking that go is installed..."
if [ ! -x "which go" ] ; then
	echo "You don't have go installed.  FIGURE IT OUT BUDDY"
	exit 1
fi
if [ ! -x "which gcc" ] ; then
	echo "You don't have GCC installed.  FIGURE IT OUT BUDDY"
	exit 1
else
	printf "✓\n"

	printf "Installing go-sqlite3..."
	go get github.com/mattn/go-sqlite3
	printf "✓\n"

	printf "Installing gopher-lua..."
	go get github.com/yuin/gopher-lua
	printf "✓\n"

	printf "github.com/layeh/gopher-luar"
	go get github.com/layeh/gopher-luar
	printf "✓\n"

	echo "Done!"
	echo "You can now do a make to build"
fi
