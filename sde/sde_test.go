package sde

import (
	"fmt"
	"testing"
)

func BenchmarkSDEGetType(b *testing.B) {
	s, err := Load("dust-wl-11.sde")
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	b.ResetTimer()
	_, err = s.GetType(364022)
	if err != nil {
		b.Log(err.Error())
	}

}

func BenchmarkGenerateCache(b *testing.B) {
	s, err := Load("dust-wl-11.sde")
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	b.ResetTimer()
	s.generateCache()
}

func BenchmarkCacheTypeNameLookup(b *testing.B) {
	s, err := Load("dust-wl-11.sde")
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	s.DoCaching(true)

	b.ResetTimer()
	_, err = s.lookupByTypeName("arm_assault_am_pro_ak0")
	if err != nil {
		b.Log(err.Error())
	}
}

func BenchmarkTypeNameLookup(b *testing.B) {

	s, err := Load("dust-wl-11.sde")
	if err != nil {
		b.Error(err)
		b.Fail()
	}

	b.ResetTimer()

	_, err = s.lookupByTypeName("arm_assault_am_pro_ak0")
	if err != nil {
		b.Log(err.Error())
	}
}

func ExampleLookupType() {
	s, err := Load("dust-wl-11.sde")
	if err != nil {
		//...
	}

	t, _ := s.GetType(364022)

	fmt.Println(t.GetName())
	// Output:
	// Assault ak.0
}
