package hdl

import "testing"

type typeTest struct {
	label string
	input string
	types []tokenType
}

func TestLex(t *testing.T) {

	tests := []typeTest{
		typeTest{
			label: "keywords",
			input: "CHIP IN OUT PARTS CLOCKED",
			types: []tokenType{
				tokenDecl,
				tokenIN,
				tokenOUT,
				tokenPARTS,
				tokenCLOCKED,
			},
		},
		typeTest{
			label: "decl",
			input: "CHIP And {",
			types: []tokenType{
				tokenDecl,
				tokenIdentifier,
				tokenLeftDelim,
			},
		},
		typeTest{
			label: "indexing",
			input: "a[16] b[2..8]",
			types: []tokenType{
				tokenIdentifier,
				tokenLeftIndex,
				tokenNumber,
				tokenRightIndex,
				tokenIdentifier,
				tokenLeftIndex,
				tokenNumber,
				tokenRange,
				tokenNumber,
				tokenRightIndex,
			},
		},
		typeTest{
			label: "simple comments", // finish with a simple comment
			input: "// comment\nCHIP\n//another one\rNot/**/// end",
			types: []tokenType{
				tokenDecl,
				tokenIdentifier,
			},
		},
		typeTest{
			label: "multline comments",
			input: "/* CHIP *//** This is a test */CHIP",
			types: []tokenType{
				tokenCommentAPI,
				tokenDecl,
			},
		},
		typeTest{
			label: "booleans",
			input: "a=true, b=false",
			types: []tokenType{
				tokenIdentifier,
				tokenPipe,
				tokenTrue,
				tokenComma,
				tokenIdentifier,
				tokenPipe,
				tokenFalse,
			},
		},

		//errors
		typeTest{
			label: "error base",
			input: "CHIP *",
			types: []tokenType{
				tokenDecl,
				tokenError,
			},
		},
		typeTest{
			label: "error dot",
			input: "1.2",
			types: []tokenType{
				tokenNumber,
				tokenError,
			},
		},
		typeTest{
			label: "error /",
			input: "/ this is an error",
			types: []tokenType{tokenError},
		},
		typeTest{
			label: "error unclosed 1",
			input: "/* unclosed comment\nsecond line",
			types: []tokenType{tokenError},
		},
		typeTest{
			label: "errpr unclosed 2",
			input: "/** Another\nunclosed\ncomment",
			types: []tokenType{tokenError},
		},
	}

	// run tests
	for _, test := range tests {

		// build lexer for this test
		l := lex(test.label, test.input)

		// store all tokens in a slice (except EOF)
		var output []token
		for tok := l.nextToken(); tok.typ != tokenEOF; tok = l.nextToken() {
			output = append(output, tok)
			if tok.typ == tokenError {
				// stop here, no more token will be emitted
				break
			}
		}

		//checks
		for i := 0; i < len(output); i++ {
			// too many outputs
			if i >= len(test.types) {
				t.Errorf(
					"%q: too many outputs (%v instead of %v) %s (l.%v)",
					test.label,
					len(output),
					len(test.types),
					output[i].val,
					output[i].line,
				)
				continue
			}

			// correct output?
			if output[i].typ != test.types[i] {
				t.Errorf(
					"%q-%v: wrong type %v instead of %v",
					test.label,
					i,
					output[i].typ,
					test.types[i],
				)
			}
		}

		// too few outputs
		if len(output) < len(test.types) {
			t.Errorf(
				"%q: too few outputs (%v instead of %v)",
				test.label,
				len(output),
				len(test.types),
			)
		}
	}
}
