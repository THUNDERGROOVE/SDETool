package util

import (
	"os"
	"runtime"
)

var OldDir string

func init() {
	OldDir, _ = os.Getwd()
}

func ChangeToAppDir() {
	switch runtime.GOOS {
	case "windows":
		path := os.ExpandEnv("$APPDATA") + "\\.SDETool\\"
		os.MkdirAll(path, 0777)
		os.Chdir(path)
	default:
		path := os.ExpandEnv("$HOME") + "\\.SDETool\\"
		os.MkdirAll(path, 0777)
		os.Chdir(path)
	}
}

func ChangeToWorkDir() {
	os.Chdir(OldDir)
}
