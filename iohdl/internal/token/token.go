package token

import "strconv"

// Token is the set of lexical tokens of this HDL.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Identifiers and literals
	IDENT  // Nand
	NUMBER // 123
	TRUE   // true
	FALSE  // false

	// Delimiters and separators
	COMMA      // ,
	SEMICOLON  // ;
	COLON      // :
	LEFTDELIM  // {
	RIGTDELIM  // }
	LEFTPAR    // (
	RIGHTPAR   // )
	LEFTINDEX  // [
	RIGHTINDEX // ]
	PIPE       // =
	RANGE      // ..

	// Keywords
	DECL    // CHIP
	IN      // IN
	OUT     // OUT
	PARTS   // PARTS
	CLOCKED // CLOCKED
)

var tokenStrings = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",
	TRUE:   "true",
	FALSE:  "false",

	COMMA:      ",",
	SEMICOLON:  ";",
	COLON:      ":",
	LEFTDELIM:  "{",
	RIGTDELIM:  "}",
	LEFTPAR:    "(",
	RIGHTPAR:   ")",
	LEFTINDEX:  "[",
	RIGHTINDEX: "]",
	PIPE:       "=",
	RANGE:      "..",

	DECL:    "CHIP",
	IN:      "IN",
	OUT:     "OUT",
	PARTS:   "PARTS",
	CLOCKED: "CLOCKED",
}

var keywords = map[string]Token{
	"CHIP":    DECL,
	"IN":      IN,
	"OUT":     OUT,
	"PARTS":   PARTS,
	"CLOCKED": CLOCKED,
}

func (t Token) String() string {
	s := ""
	if 0 <= t && int(t) < len(tokenStrings) {
		s = tokenStrings[t]
	}

	// if s is still empty, t is an unknown token
	// for debug purpose, print the integer representation of the token
	if s == "" {
		s = "TOKEN(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

// Lookup maps an identifier to its keyword token, TRUE or FALSE literal or IDENT.
func Lookup(s string) Token {
	if keyword, present := keywords[s]; present {
		return keyword
	}

	switch s {
	case "true":
		return TRUE
	case "false":
		return FALSE
	default:
		return IDENT
	}
}
