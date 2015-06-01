// @TODO: Stop rolling our own parser.  Switch to flags flagsets with commands.
// Might create redundency with flagsets such as --sde.  We will see

package args

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Token uint8

var eof = rune(0)

const (
	FLAG Token = iota

	WHITESPACE

	STRING
	INT

	EOF
	ILLEGAL
)

type Scanner struct {
	r     *bufio.Reader
	index int
}

type TokLit struct {
	Token   Token
	Literal string
}

type Tokens []TokLit

func (t Tokens) Next(i int) *TokLit {
	if len(t) < i+1 {
		return nil
	}
	if t[i+1].Token == WHITESPACE {
		return t.Next(i + 1) // Skip whitespace
	}
	return &t[i+1]
}

func (t Tokens) Prev(i int) *TokLit {
	if i == 0 {
		return nil
	}
	return &t[i-1]
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (t Token) String() string {
	switch t {
	case FLAG:
		return "Flag"
	case WHITESPACE:
		return "Whitespace"
	case STRING:
		return "String"
	case INT:
		return "Int"
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "Illegal"
	default:
		return "Unknown token '" + strconv.Itoa(int(t)) + "'"
	}
	return ""
}

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isFlag(ch) {
		s.unread()
		return s.scanFlag()
	} else if isInt(ch) {
		s.unread()
		return s.scanInt()
	} else if isString(ch) || ch == '"' {
		return s.scanString(true)
	} else if isString(ch) && ch != '"' {
		s.unread()
		return s.scanString(false)
	} else if ch == eof {
		return EOF, "EOF"
	}
	return ILLEGAL, string(ch)
}

func (s *Scanner) ScanAll() Tokens {
	o := Tokens(make([]TokLit, 0))
	for {
		t, l := s.Scan()
		o = append(o, TokLit{Token: t, Literal: l})
		if t == EOF || t == ILLEGAL {
			break
		}
	}

	return o
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}

	}
	return WHITESPACE, buf.String()
}

func (s *Scanner) scanFlag() (tok Token, lit string) {
	var buf bytes.Buffer

	for {
		ch := s.read()
		if ch == eof {
			return FLAG, buf.String()
		}
		switch {
		case (!isString(ch) && !isFlag(ch)):
			s.unread()
			return FLAG, buf.String()
		case (isString(ch) || isFlag(ch)):
			buf.WriteRune(ch)
		default:
		}
	}
	return ILLEGAL, "Failed to parse flag"
}

func (s *Scanner) scanString(quoteStarted bool) (tok Token, lit string) {
	var buf bytes.Buffer

	var isStr bool
	if quoteStarted {
		isStr = true
	}

	for {
		ch := s.read()
		if ch == eof {
			break
		}
		fmt.Printf("Found char '%v': ", string(ch))
		fmt.Printf("Not in quoted string\n")
		switch {
		case (isWhitespace(ch) && isStr):
			buf.WriteRune(ch)
		case (!isString(ch) && !isWhitespace(ch) && isStr): // We are in a quoted string and character doesn't match isString
			return STRING, buf.String()
		case (isQuote(ch) && isStr): // Our string's end quote was reached
			return STRING, buf.String() // Don't write the quote
		case (isQuote(ch) && !isStr): // Found first string quote
			isStr = true // Don't write the quote
		case (isWhitespace(ch) && !isStr): // Whitespace found while not in a quoted string
			s.unread()
			return STRING, buf.String()
		case isString(ch): // Nothing above was reached but still matches isString
			buf.WriteRune(ch)
		default: // ILLEGAL token
			return ILLEGAL, "Token '" + string(ch) + "' is illegal in the context of a string"
		}
	}
	return STRING, buf.String()
}

func (s *Scanner) scanInt() (tok Token, lit string) {
	var buf bytes.Buffer
	for {
		ch := s.read()
		if ch == eof {
			break
		}
		switch {
		case !isInt(ch):
			s.unread()
			return INT, buf.String()
		case isInt(ch):
			buf.WriteRune(ch)
		default:
		}
	}
	return INT, buf.String()
}

func (s *Scanner) read() rune {
	s.index += 1
	ch, _, err := s.r.ReadRune()
	fmt.Printf("read() -> '%v'\n", string(ch))
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	s.index -= 1
	if err := s.r.UnreadRune(); err != nil {
		fmt.Printf("Error unreading rune: %v\n", err.Error())
	}
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isString(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '-' || ch == '_' || ch == '.' || isInt(ch)
}

func isQuote(ch rune) bool { return ch == '"' }
func isFlag(ch rune) bool  { return ch == '-' }
func isInt(ch rune) bool   { return (ch >= '0' && ch <= '9') }
