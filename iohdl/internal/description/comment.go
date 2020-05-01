package description

import (
	"strings"

	"github.com/batiazinga/hdl/iohdl/internal/token"
)

// Comment is a documentation associated to an element of the chip
// or the chip itsefl.
type Comment struct {
	start, end token.Position
	lit        string
}

// Start returns the position at which the comment starts.
func (c Comment) Start() token.Position { return c.start }

// End returns the position of the end of the comment.
func (c Comment) End() token.Position { return c.end }

// Literal returns the full text of the comment, including the delimiters.
func (c Comment) Literal() string { return c.lit }

// IsEOL returns true if the comment is a "until end-of-line" comment.
func (c Comment) IsEOL() bool { return strings.HasPrefix(c.lit, "//") }

// IsMultiline returns true if the comment is multi-line comment.
// Both "/*" and "/**" are multiline comments.
func (c Comment) IsMultiline() bool { return strings.HasPrefix(c.lit, "/*") }

// IsAPI returns true if the comment is an API comment,
// i.e. whose left delimiter is "/**".
func (c Comment) IsAPI() bool { return strings.HasPrefix(c.lit, "/**") }
