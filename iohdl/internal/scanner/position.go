package scanner

import "fmt"

// Position is a position in an input file.
type Position struct {
	file         string
	line, column int
}

// File returns the name of the file.
func (p Position) File() string { return p.file }

// Line returns the line number.
// The first line has index 1.
func (p Position) Line() int { return p.line + 1 }

// Column returns the column number, i.e. the position in the current line.
// The first column has index 1.
func (p Position) Column() int { return p.column + 1 }

func (p Position) String() string {
	return fmt.Sprintf("%s line %d column %d", p.File(), p.Line(), p.Column())
}
