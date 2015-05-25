package args

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringParsing(t *testing.T) {
	r := strings.NewReader("faggot")
	s := NewScanner(r)
	for {
		tok, lit := s.Scan()
		fmt.Printf("Token: '%s' | Literal: '%v'\n", tok, lit)
		if tok == EOF {
			break
		}
		if tok == ILLEGAL {
			break
		}
	}
}

func TestIntParsing(t *testing.T) {
	r := strings.NewReader("faggot 1235105")
	s := NewScanner(r)
	for {
		tok, lit := s.Scan()
		fmt.Printf("Token:'%s' | Literal: '%v'\n", tok, lit)
		if tok == EOF || tok == ILLEGAL {
			break
		}
	}
}

func TestFlagParsing(t *testing.T) {
	r := strings.NewReader("faggot -flag 1246787654")
	s := NewScanner(r)
	for {
		tok, lit := s.Scan()
		fmt.Printf("Token:'%s' | Literal: '%v'\n", tok, lit)
		if tok == EOF || tok == ILLEGAL {
			break
		}
	}
}
