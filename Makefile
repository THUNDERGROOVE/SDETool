folder = SDETool-`git rev-parse --abbrev-ref HEAD`-`git describe --abbrev=0`
all:
	go build -ldflags "-X main.commit=`git rev-parse --short HEAD` -X main.branch=`git rev-parse --abbrev-ref HEAD` -X main.tagVersion=`git describe --abbrev=0`"
packagewin: all
	mkdir $(folder)
	cp SDETool.exe $(folder)
	cp LICENSE $(folder)
	cp README.md $(folder)
	cp scripts/interactive.bat $(folder)
	mkdir $(folder)/docs
	cp scripting/lua/lua.md $(folder)/docs
	zip -r dist/$(folder)-windows.zip $(folder)/
	rm -rf $(folder)
packagelinux: all
	mkdir $(folder)
	cp SDETool $(folder)
	cp LICENSE $(folder)
	cp README.md $(folder)
	mkdir $(folder)/docs
	cp scripting/lua/lua.md $(folder)/docs
	zip -r dist/$(folder).zip $(folder)/
	rm -rf $(folder)
clean:
	rm -f log.txt
	rm -f SDETool
	rm -f SDETool.exe
	rm *.sde
	rm *.db
