package scanner

import (
	"unicode"
	"unicode/utf8"

	"github.com/batiazinga/hdl/iohdl/internal/token"
)

const (
	eof = -1
)

// Scanner can scan a source text to extract its tokens.
type Scanner struct {
	// input file
	filename string
	src      []byte

	// current rune and its width and position
	// line and column start at zero
	current        rune
	w              int
	pos, line, col int

	// start position of the current token
	start, tokLine, tokCol int

	// invalid src (not illegal tokens)
	errs []Error
}

// New returns a ready-to-use scanner.
func New(file string, src []byte) *Scanner {
	s := &Scanner{
		filename: file,
		src:      src,
	}

	s.next()
	return s
}

// Scan scans the next token and returns it with its position and literal string if any.
// The source ends with the token.EOF token.
//
// The position points to the beginning of the token.
//
// All tokens have a literal string except token.EOF.
func (s *Scanner) Scan() (pos Position, tok token.Token, lit string) {
	s.skipSpace()
	pos = Position{s.filename, s.tokLine, s.tokCol}

	switch current := s.current; {

	case current == eof:
		tok = token.EOF

	case isDigit(current):
		for isDigit(s.current) {
			s.next()
		}
		tok = token.NUMBER
		lit = s.literal()
		s.moveToken()

	case unicode.IsLetter(current):
		for isAlphanumeric(s.current) {
			s.next()
		}
		lit = s.literal()
		tok = token.Lookup(lit)
		s.moveToken()

	case current == '.':
		// this may be a range
		s.next()
		if s.current == '.' {
			tok = token.RANGE
			lit = ".."
			s.next()
		} else {
			tok = token.ILLEGAL
			lit = "."
		}

	case startComment(current):
		tok, lit = s.comment()

	// simple tokens
	default:
		switch current {
		case ',':
			tok = token.COMMA
		case ';':
			tok = token.SEMICOL
		case ':':
			tok = token.COLUMN
		case '{':
			tok = token.LEFTDELIM
		case '}':
			tok = token.RIGTDELIM
		case '(':
			tok = token.LEFTPAR
		case ')':
			tok = token.RIGHTPAR
		case '[':
			tok = token.LEFTINDEX
		case ']':
			tok = token.RIGHTINDEX
		case '=':
			tok = token.PIPE

		default:
			tok = token.ILLEGAL
		}

		// move to the position after the token
		// i.e. position of the next token (or whitespace)
		s.next()

		lit = s.literal()
		s.moveToken()
	}

	return
}

// next moves the scanner to the next rune.
// The next rune may be utf8.RuneError, which means that encoding is invalid.
func (s *Scanner) next() {
	// move to the next rune
	s.pos += s.w
	s.col += s.w
	if s.current == '\n' {
		s.line++
		s.col = 0
	}

	if s.pos >= len(s.src) {
		s.current = eof
		s.w = 0
		return
	}

	// read next rune
	// first try ASCII
	s.current, s.w = rune(s.src[s.pos]), 1
	if s.current >= utf8.RuneSelf {
		// not ASCII
		s.current, s.w = utf8.DecodeLastRune(s.src[s.pos:])
		if s.current == utf8.RuneError {
			// input was not empty
			// so it's an invalid UTF8-encoding (and width is 1)
			// collect error and continue
			s.errs = append(s.errs, Error{Position{s.filename, s.line, s.col}})
		}
	}
}

// moveToken moves the token position to the current position.
func (s *Scanner) moveToken() {
	s.start = s.pos
	s.tokLine = s.line
	s.tokCol = s.col
}

func (s *Scanner) skipSpace() {
	for s.current == ' ' || s.current == '\t' || s.current == '\r' || s.current == '\n' {
		s.next()
	}
	s.moveToken()
}

// literal returns the string between the token start position (included)
// and the current position (excluded).
func (s *Scanner) literal() string {
	return string(s.src[s.start:s.pos])
}

// scan a comment assuming the current rune is '/'.
// So call it just after a call to startComment has returned true.
func (s *Scanner) comment() (tok token.Token, lit string) {
	// first rune is '/'
	s.next()

	// find comment's left delimiter
	switch s.current {

	case '/':
		// comment starter is "//"
		s.next() // move to first rune after "//"

	case '*':
		// multi line comment or API comment?
		// read third rune to know
		s.next()
		if s.current == '*' {
			// API comment "/**"
			// but maybe this is just an empty multiline comment "/**/"
			s.next()
			if s.current == '/' {
				// this is a /**/, i.e. empty multiline comment
				s.next()
				tok = token.COMMENT
				lit = s.literal()
				s.moveToken()
				return
			} // else, this is a "/**", i.e. an API comment
		} // else, comment starter is "/*"

	default:
		tok = token.ILLEGAL
		lit = s.literal()
		s.moveToken()
		return
	}

	// at this point we know the comment left delimiter
	// and the scanner points to the first rune after
	if start := string(s.src[s.start:s.pos]); start == "//" {
		// comment until the end of line
		// loop until we find EOL or EOF
		for s.current != eof && s.current != '\n' && s.current != '\r' {
			s.next()
		}
		// move to rune after EOL or stay on EOF
		s.next()
		tok = token.COMMENT
		lit = s.literal()
		s.moveToken()
		return

	}
	// "/**" or "/*"
	// find the next "*/" delimiter
	// This infinite loop stops when we meet the EOF rune
	// or we find the right delimiter
	for {
		for s.current != '*' && s.current != eof {
			s.next()
		}
		// current rune is '*' or eof
		if s.current == eof {
			tok = token.ILLEGAL
			lit = s.literal()
			return
		}
		// current rune is '*'
		s.next()
		if s.current == '/' {
			// end of comment
			s.next()
			tok = token.COMMENT
			lit = s.literal()
			s.moveToken()
			return
		}
	}

}

// only 0-9
func isDigit(r rune) bool { return '0' <= r && r <= '9' }

// 0-9, other digits, letters and _
func isAlphanumeric(r rune) bool { return isDigit(r) || unicode.IsLetter(r) || r == '_' }

func startComment(r rune) bool { return r == '/' }
