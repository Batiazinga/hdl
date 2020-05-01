package token

import "fmt"

// Position is a position in an input file.
type Position struct {
	line, column int
}

// NewPosition returns a position pointing to (line, column).
// Parameters line and column are 0-indexed.
func NewPosition(line, column int) Position { return Position{line, column} }

// Line returns the line number.
// The first line has index 1.
func (p Position) Line() int { return p.line + 1 }

// Column returns the column number (byte count), i.e. the position in the current line.
// The first column has index 1.
func (p Position) Column() int { return p.column + 1 }

func (p Position) String() string {
	return fmt.Sprintf("line %d column %d", p.Line(), p.Column())
}

// Less returns true if p is strictly less than q.
func (p Position) Less(q Position) bool {
	switch {
	case p.line < q.line:
		return true
	case p.line > q.line:
		return false
	default:
		return p.column < q.column
	}
}
