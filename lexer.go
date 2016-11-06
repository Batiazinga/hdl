package hdl

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type token struct {
	typ          tokenType
	file         string
	line, column int
	val          string
}

func (t token) String() string {
	switch {
	case t.typ == tokenEOF:
		return "EOF"
	case t.typ == tokenError:
		return t.val
	case t.typ > tokenKeywords:
		return fmt.Sprintf("<%s>", t.val)
	case len(t.val) > 10:
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenComma
	tokenSemiCol
	tokenColumn
	tokenLeftDelim
	tokenRightDelim
	tokenLeftPar
	tokenRightPar
	tokenLeftIndex
	tokenRightIndex
	tokenPipe
	tokenNumber
	tokenRange
	tokenTrue
	tokenFalse
	tokenIdentifier
	tokenCommentAPI // API documentation comment
	// Keywords all appear after this token
	tokenKeywords // this is just a separator
	tokenDecl     // the CHIP keyword
	tokenIN
	tokenOUT
	tokenPARTS
	tokenCLOCKED
)

var keywords = map[string]tokenType{
	"CHIP":    tokenDecl,
	"IN":      tokenIN,
	"OUT":     tokenOUT,
	"PARTS":   tokenPARTS,
	"CLOCKED": tokenCLOCKED,
}

const (
	eof        = -1
	spaceChars = " \t\r\n"
)

// stateLex is a state of the lexer.
// It is a function which returns another state (function).
type stateLex func(*lexer) stateLex

type lexer struct {
	file   string // name of the input file/chip
	tokens chan token
	input  string
	state  stateLex // the next state
	pos    int      // current position in the input
	width  int      // width of the last rune
	start  int      // start of the current token
}

//functions for internal use in the lexer

// next returns the next rune in the input string.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	// the error rune cannot be returned here
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// backup steps back one rune.
// Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// ignores runes between the last emitted token and the current position.
func (l *lexer) ignore() {
	l.start = l.pos
}

// lineNumber computes the current line number from the position of the current token.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.start], "\n")
}

// columnNumber computes the column number of the start position of the current token.
func (l *lexer) columnNumber() int {
	// find last '\n' before l.start
	last := strings.LastIndex(l.input[:l.start], "\n")
	// count number of runes up to l.start (included)
	return 1 + utf8.RuneCountInString(l.input[last+1:l.start])
}

// emit sends the current token to the channel.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.file, l.lineNumber(), l.columnNumber(), l.input[l.start:l.pos]}
	l.start = l.pos
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// errorf emits a token error
// and return the nil (end) state.
func (l *lexer) errorf(format string, args ...interface{}) stateLex {
	l.tokens <- token{tokenError, l.file, l.lineNumber(), l.columnNumber(), fmt.Sprintf(format, args...)}
	return nil
}

// run runs the lexing state machine.
func (l *lexer) run() {
	// initialize the lexer state
	// and loop over states until the nil state
	for l.state = lexBase; l.state != nil; {
		l.state = l.state(l)
	}
	// source (lexer) closes the channel when it's done
	close(l.tokens)
}

//functions used by users of the lexer

// nextToken returns the next token from the input.
// After a tokenEOF or tokenError, no more token is received.
// It is called by the parser in another goroutine.
func (l *lexer) nextToken() token {
	t := <-l.tokens
	return t
}

// drain the token output until the lexer goroutine is done and exits.
// This is called by the parser in another goroutine.
func (l *lexer) drain() {
	for range l.tokens {
	}
}

// lex creates a new lexer for the given input string
func lex(file, input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l
}

// states of the lexer

func lexBase(l *lexer) stateLex {
	//ignore spaces
	l.acceptRun(spaceChars)
	l.ignore()

	//read next rune and find next state
	switch r := l.next(); {

	// Comments
	case r == '/':
		l.backup()
		return lexComment

	// simple symbols
	case r == '{':
		l.emit(tokenLeftDelim)
	case r == '}':
		l.emit(tokenRightDelim)
	case r == '(':
		l.emit(tokenLeftPar)
	case r == ')':
		l.emit(tokenRightPar)
	case r == '[':
		l.emit(tokenLeftIndex)
	case r == ']':
		l.emit(tokenRightIndex)
	case r == ':':
		l.emit(tokenColumn)
	case r == ',':
		l.emit(tokenComma)
	case r == ';':
		l.emit(tokenSemiCol)
	case r == '=':
		l.emit(tokenPipe)

	case r == '.':
		// Is it a range? i.e. '..'
		if s := l.next(); s == '.' {
			l.emit(tokenRange)
		} else {
			return l.errorf("Unexpected character %q after '.'", s)
		}

	case unicode.IsDigit(r):
		l.backup()
		return lexInteger
	case unicode.IsLetter(r):
		l.backup()
		return lexIdentifier

	case r == eof:
		l.emit(tokenEOF)
		return nil

	default:
		return l.errorf("Unexpected character %q", r)
	}
	return lexBase
}

func lexIdentifier(l *lexer) stateLex {
LOOP:
	for {
		switch r := l.next(); {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			// if alphanumeric, absorb and continue

		default:
			l.backup()
			// find the word we just read
			word := l.input[l.start:l.pos]
			switch {
			case keywords[word] > tokenKeywords:
				l.emit(keywords[word])
			case word == "true":
				l.emit(tokenTrue)
			case word == "false":
				l.emit(tokenFalse)
			default:
				l.emit(tokenIdentifier)
			}
			break LOOP
		}
	}
	return lexBase
}

func lexInteger(l *lexer) stateLex {
	l.acceptRun("0123456789")
	l.emit(tokenNumber)
	return lexBase
}

func lexComment(l *lexer) stateLex {
	// find comment's left delimiter
	_ = l.next() // this is a '/' rune
	switch second := l.next(); second {

	case '/':
		// nothing more to do -> left delimiter is //

	case '*':
		// multi line comment or API comment?
		if third := l.next(); third != '*' {
			// just a multi line comment -> left delim is /*
			l.backup()
		} else {
			// /** or /**/ ?
			// read the fourth rune to know
			if fourth := l.next(); fourth == '/' {
				// this is a /**/ i.e. useless multiline comment
				l.ignore()
				return lexBase
			} // else, this is a /**
			l.backup()
		}

	default:
		return l.errorf("Unexpected character %q after '/'", second)
	}

	// read left delimiter
	leftDelim := l.input[l.start:l.pos]

	// At this point I know the comment's left delimiter
	// and my current position is the first character following it
	if leftDelim == "//" {
		//single line comment: ignore all runes up to the next eol
		// this infinite loop stops when there is an EOF or EOL
	LOOP1:
		for {
			switch r := l.next(); r {
			case eof:
				// EOF: emit EOF and return the nil state
				l.emit(tokenEOF)
				return nil
			case '\r', '\n':
				// EOL: ignore comment and back to normal lexing
				break LOOP1
			}
		}
		l.ignore() // ignore the comment

	} else { // /* or /** ?
		// if /* ... */ ignore everything
		// else if /** ... */, keep all comment,
		// including delimiters and emit token
		//
		// find the next */ delimiter
		// This infinite loop stops when we meet the EOF character
		// or we find the right delimiter.
	LOOP2:
		for {
			switch r := l.next(); r {
			case eof:
				return l.errorf("Unclosed comment")
			case '*':
				// test next rune
				if r = l.next(); r == '/' {
					//end of comment
					// simple multiline or API comment?
					if leftDelim == "/*" {
						l.ignore()
					} else {
						l.emit(tokenCommentAPI)
					}
					break LOOP2
				}
			}
		}
	}
	return lexBase
}

// Temporary

// Lex is a temporary function
func Lex(file, input string) {
	l := lex(file, input)

	for t := l.nextToken(); t.typ != tokenEOF && t.typ != tokenError; t = l.nextToken() {
		fmt.Println(t)
	}
}
