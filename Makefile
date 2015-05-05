all:
	go build -ldflags "-X main.Commit `git rev-parse --short HEAD` -X main.Branch `git rev-parse --abbrev-ref HEAD` -X main.Version `git describe --abbrev=0`"