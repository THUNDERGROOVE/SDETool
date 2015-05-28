package version

import (
	"fmt"
	"testing"
)

func TestLoadVersion(t *testing.T) {
	vers, err := LoadVersions()
	if err != nil {
		t.Errorf("Unable to load versions from remote: [%v]\n", err.Error())
		t.Fail()
	}

	for k, v := range vers {
		fmt.Printf("[%v]: %v\n", k, v)
	}
}
