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

	fmt.Println("Getting Assault ak.0 1000 times")
	b.ResetTimer()
	for i := 0; i < 1000; i++ {
		_, err := s.GetType(364022)
		if err != nil {
			b.Log(err.Error)
		}
	}

}
