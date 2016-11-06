package hdl

import "testing"

// test that tokens have the correct token type

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

// test that line and column number are correct

type positionTest struct {
	label     string
	input     string
	positions []struct{ line, column int }
}

func TestPosition(t *testing.T) {

	tests := []positionTest{
		{
			"valid",
			"CHIP And {\n\tIN a, b;\n\tOUT out;",
			[]struct{ line, column int }{
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
			[]struct{ line, column int }{
				{1, 1},
				{1, 6},
			},
		},
		{
			"invalid line 2",
			"CHIP And {\n\t*}",
			[]struct{ line, column int }{
				{1, 1},
				{1, 6},
				{1, 10},
				{2, 2},
			},
		},
	}

	for i, test := range tests {
		// lexer for this test
		l := lex(test.label, test.input)

		var outputs []token
		for tok := l.nextToken(); tok.typ != tokenEOF; tok = l.nextToken() {
			outputs = append(outputs, tok)
			if tok.typ == tokenError {
				// no more token will be emitted
				break
			}
		}

		if len(outputs) != len(test.positions) {
			t.Fatalf("test %v: invalid number of tokens (%v instead of %v)", i, len(outputs), len(test.positions))
		}

		for j, tok := range outputs {
			if tok.line != test.positions[j].line {
				t.Errorf("test %v: wrong line number for token %s (%v instead of %v)", j, tok, tok.line, test.positions[j].line)
			}
			if tok.column != test.positions[j].column {
				t.Errorf("test %v: wrong column number for token %s (%v instead of %v)", j, tok, tok.column, test.positions[j].column)
			}
		}

	}
}
