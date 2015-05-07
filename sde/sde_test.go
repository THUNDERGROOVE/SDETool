package sde

import (
	"os"
	"testing"
)

const (
	openVer = "1.0-WL"
)

func init() {
	os.Chdir(os.TempDir())
}

func TestOpen(T *testing.T) {
	if s, err := Open(openVer); err == nil {
		if s.Version != openVer {
			T.Fatal("SDE version did not match after opening file")
		}
		if t, err := s.GetType(364022); err == nil {
			if t.GetName() != "Assault ak.0" {
				T.Fatal("Type mDisplayName missmatch")
			}
		} else {
			T.Fatal("Error getting type for test", err.Error())
		}
		s.DB.Close()
	} else {
		T.Fatal("Got error opening SDE file", err.Error())
	}
}
