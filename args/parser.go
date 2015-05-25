package args

import (
	"bufio"
	"bytes"
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
	} else if isString(ch) || ch == '"' {
		s.unread()
		return s.scanString()
	} else if isInt(ch) {
		s.unread()
		return s.scanInt()
	} else if isFlag(ch) {
		s.unread()
		return s.scanFlag()
	} else if ch == eof {
		return EOF, "EOF"
	}
	return ILLEGAL, string(ch)
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
			break
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

func (s *Scanner) scanString() (tok Token, lit string) {
	var buf bytes.Buffer

	isStr := false

	for {
		ch := s.read()
		if ch == eof {
			break
		}
		switch {
		case (!isString(ch) && !isWhitespace(ch) && isStr): // We are in a quoted string and character doesn't match isString
			break
		case (isQuote(ch) && isStr): // Our string's end quote was reached
			return STRING, buf.String() // Don't write the quote
		case (isQuote(ch) && !isStr): // Found first string quote
			isStr = true // Don't write the quote
		case (isWhitespace(ch) && !isStr): // Whitespace found while not in a quoted string
			s.unread()
			return STRING, buf.String()
			// Don't behaving this way \/
			return ILLEGAL, "Cannot have whitespace in string literal without quotes"
		case (isWhitespace(ch) && isStr):
			buf.WriteRune(ch)
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
			break
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
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	s.index -= 1
	_ = s.r.UnreadRune()
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isString(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isQuote(ch rune) bool { return ch == '"' }
func isFlag(ch rune) bool  { return ch == '-' }
func isInt(ch rune) bool   { return (ch >= '0' && ch <= '9') }
