# SDETool

Okay I'm done writing the plain, go only versions.

This time around SDETool will include built-in scripting capabilities.  
For now it just has a lua interpreter and to see all of the exposed methods 
check the [docs](https://github.com/THUNDERGROOVE/SDETool/blob/master/scripting/lua/lua.md)

### Downloads

I have a build server setup that builds SDETool everytime I push code.  
You can download windows/amd64 and linux/amd64 builds

[Jenkins](http://ci.maximumtwang.com)

### Contributing

Want to work on SDETool?  

- Create an issue describing what you want to change.
- Fork the repo
- *MAKE A NEW BRANCH* with the issue number (This is importent if I have 
  changes that haven't hit the repo)
- Make your changes
- Pull request.

I'll accept it and merge it into the development branch and your changes will 
likely get merged into master and released shortly after.

### Building on your own

Windows instructions:

- Install Go from [here](http://golang.org/dl/)
- Install gcc.  I prefer [msys2](http://msys2.github.io)
- Install make. If using msys2 `pacman -Syu make`
- `./scripts/bootstrap.sh` Downloads our dependencies and check that `go` is 
  setup correctly
- `make`

Should have a properly built executable file

### Adding your own language.

Do you know Golang and have a scripting language that idealy doesn't require 
cgo?  

Create an issue with a repo that contains a package similar to the 
scripting.lua package I've already made and I'll submodule it.

If it requires cgo I will likely add it to a seperate file with a 
`//+build linux` tag.  I've tried plenty of times to deal with cgo dependencies
in Windows and it's not fun one bit.  I've tried mingw/msys2/cygwin(ew) all are
trash when it comes to actually using pkg-config and actually getting gcc to 
link for you.

If you do require a dependency that uses cgo, try to use something like 
go-sqlite3 that keeps the C source and headers within the repository.

This is less relevent as Go isn't shipped with it's own C compiler anymore.

### LICENSE
SDETool is released under the MIT license.  
Check [LICENSE](http://github.com/THUNDERGROOVE/SDETool) for more details
