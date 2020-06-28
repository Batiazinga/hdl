package description

import "github.com/batiazinga/hdl/iohdl/internal/token"

// Chip is the description of a chip
// with its comments, header and body.
type Chip struct {
	// sorted by position
	comments []Comment

	// declaration
	start, end token.Position
	name       string

	// header
	inputs  InputList
	outputs OutputList

	// body
	parts PartList
}

// Start returns the position at which the chip starts.
// All comments are taken into account
// so this is not necessarily the position of the main comment
// or where the chip is declared.
func (c Chip) Start() token.Position {
	if len(c.comments) == 0 {
		return c.start
	}
	return c.comments[0].Start()
}

// End returns the line at which the chip ends.
func (c Chip) End() token.Position { return c.end }

// NumComments returns the number of comments relative to the chip.
// In the hdl format, these are all comments before the chip declaration.
func (c Chip) NumComments() int { return len(c.comments) }

// Comment returns the i-th comment relative to the chip.
// This panics if i is out of bounds.
func (c Chip) Comment(i int) Comment { return c.comments[i] }

// Name returns the name of the chip, i.e. the name in the chip declaration.
func (c Chip) Name() string { return c.name }

// Inputs returns the list of input pins.
func (c Chip) Inputs() InputList { return c.inputs }

// Outputs returns the list of output pins.
func (c Chip) Outputs(i int) OutputList { return c.outputs }

// Parts returns the list of part.
func (c Chip) Parts() PartList { return c.parts }

// ChipBuilder can build a description of a chip on-the-fly.
type ChipBuilder struct{}

// AppendComment appends a comment.
// Even if comments relative to the chip are before the chip declaration,
// it can be called after the declaration.
func (b *ChipBuilder) AppendComment(comment Comment) {}

// Declare starts the declaration of the chip and set its name.
func (b *ChipBuilder) Declare(line int, name string) {}

// DeclareInputs declare the list of inputs.
// Even if inputs should be declared between the chip and outputs declarations,
// it can be called at any moment.
func (b *ChipBuilder) DeclareInputs(inputs InputList) {}

// DeclareOutputs declares the list of outputs.
// Even if outputs should be declared between the inputs and the body declarationw,
// it can be called at any moment.
func (b *ChipBuilder) DeclareOutputs(outputs OutputList) {}

// DeclareParts declares the list of parts.
// Even if the body of the chip should be defined after the outputs declaration,
// it can be called at any moment.
func (b *ChipBuilder) DeclareParts(part Part) {}

// Build return the chip.
// The builder is reset so it can reused without any side effect on previously built chips.
func (b *ChipBuilder) Build() Chip { return Chip{} }
