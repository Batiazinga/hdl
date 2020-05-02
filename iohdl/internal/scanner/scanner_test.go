package scanner

import (
	"testing"

	"github.com/batiazinga/hdl/iohdl/internal/token"
)

// TestScanner checks that the scanner returns the expected tokens in the expected order.
func TestScanner(t *testing.T) {

	testcases := []struct {
		label  string
		src    string
		tokens []token.Token
	}{
		{
			label: "empty",
		},
		{
			"keywords",
			"CHIP IN OUT PARTS CLOCKED",
			[]token.Token{
				token.DECL,
				token.IN,
				token.OUT,
				token.PARTS,
				token.CLOCKED,
			},
		},
		{
			"delimiters and separators",
			",;:{}()[]..=",
			[]token.Token{
				token.COMMA,
				token.SEMICOLON,
				token.COLON,
				token.LEFTDELIM,
				token.RIGTDELIM,
				token.LEFTPAR,
				token.RIGHTPAR,
				token.LEFTINDEX,
				token.RIGHTINDEX,
				token.RANGE,
				token.PIPE,
			},
		},
		{
			"identifiers and literals",
			"Chip invalid declaration_1\n142 true false",
			[]token.Token{
				token.IDENT,
				token.IDENT,
				token.IDENT,
				token.NUMBER,
				token.TRUE,
				token.FALSE,
			},
		},
		{
			"non ASCII",
			"γθιπ",
			[]token.Token{
				token.IDENT,
			},
		},
		{
			"indexing and buses",
			"a[16] b[2..8]",
			[]token.Token{
				token.IDENT,
				token.LEFTINDEX,
				token.NUMBER,
				token.RIGHTINDEX,
				token.IDENT,
				token.LEFTINDEX,
				token.NUMBER,
				token.RANGE,
				token.NUMBER,
				token.RIGHTINDEX,
			},
		},
		{
			"simple comments", // finish with a simple comment
			"// comment\nCHIP // until end of line\n// and useless //",
			[]token.Token{
				token.COMMENT,
				token.DECL,
				token.COMMENT,
				token.COMMENT,
			},
		},
		{
			"multiline comments",
			"/* CHIP *//**/ /* a*b */CHIP",
			[]token.Token{
				token.COMMENT,
				token.COMMENT,
				token.COMMENT,
				token.DECL,
			},
		},
		{
			"API comment",
			"/**\n  doc of the CHIP\n*/\nCHIP And",
			[]token.Token{
				token.COMMENT,
				token.DECL,
				token.IDENT,
			},
		},

		// illegal tokens
		{
			"illegal",
			"CHIP *",
			[]token.Token{
				token.DECL,
				token.ILLEGAL,
			},
		},
		{
			"illegal float number",
			"1.2",
			[]token.Token{
				token.NUMBER,
				token.ILLEGAL,
				token.NUMBER,
			},
		},
		{
			"wrong comment",
			"/ this is an error",
			[]token.Token{
				token.ILLEGAL,
				token.IDENT,
				token.IDENT,
				token.IDENT,
				token.IDENT,
			},
		},
		{
			"unclosed comment",
			"/* unclosed comment\nsecond line",
			[]token.Token{
				token.ILLEGAL,
			},
		},
		{
			"unclosed API comment",
			"/** Another\nunclosed\ncomment",
			[]token.Token{
				token.ILLEGAL,
			},
		},
	}

	// run tests
	for _, testcase := range testcases {
		t.Run(
			testcase.label,
			func(t *testing.T) {
				// build lexer for this test
				s := New("filename", []byte(testcase.src))

				// store all tokens in a slice (except EOF)
				var tokens []token.Token
				_, tok, _ := s.Scan()
				for tok != token.EOF {
					tokens = append(tokens, tok)
					_, tok, _ = s.Scan()
				}

				// checks
				if len(tokens) != len(testcase.tokens) {
					t.Fatalf(
						"unexpected number of tokens: %d instead of %d\n  %v\n  %v",
						len(tokens), len(testcase.tokens),
						tokens, testcase.tokens,
					)
				}
				for i := range tokens {
					if tokens[i] != testcase.tokens[i] {
						t.Errorf("unexpected %d-th token: %s instead of %s", i, tokens[i], testcase.tokens[i])
					}
				}
			},
		)
	}
}

func TestPosition(t *testing.T) {
	type position struct {
		line, column int
	}
	type tc struct {
		label     string
		src       string
		positions []position
	}

	testcases := []tc{
		{
			"declaration",
			"CHIP And {\n\tIN a, b;\n\tOUT out;",
			[]position{
				{1, 1},
				{1, 6},
				{1, 10},
				{2, 2},
				{2, 5},
				{2, 6},
				{2, 8},
				{2, 9},
				{3, 2},
				{3, 6},
				{3, 9},
			},
		},
		{
			"invalid line 1",
			"CHIP *",
			[]position{
				{1, 1},
				{1, 6},
			},
		},
		{
			"invalid line 2",
			"CHIP And {\n\t*}",
			[]position{
				{1, 1},
				{1, 6},
				{1, 10},
				{2, 2},
				{2, 3},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(
			testcase.label,
			func(t *testing.T) {
				s := New("test.hdl", []byte(testcase.src))

				// store all positions in a slice (except for EOF)
				var positions []position
				pos, tok, _ := s.Scan()
				for tok != token.EOF {
					positions = append(positions, position{pos.Line(), pos.Column()})
					pos, tok, _ = s.Scan()
				}

				if len(positions) != len(testcase.positions) {
					t.Fatalf("unexpected number of tokens: %d instead of %d", len(positions), len(testcase.positions))
				}
				for i := range positions {
					if positions[i].line != testcase.positions[i].line {
						t.Errorf("wrong line number for %d-th token (%v instead of %v)", i, positions[i].line, testcase.positions[i].line)
					}
					if positions[i].column != testcase.positions[i].column {
						t.Errorf("wrong column number for %d-th token (%v instead of %v)", i, positions[i].column, testcase.positions[i].column)
					}
				}
			},
		)
	}
}
