package description

import "github.com/batiazinga/hdl/iohdl/internal/token"

// Chip is the description of a chip
// with its comments, interface and parts.
type Chip struct {
	// sorted by position
	comments []Comment

	// declaration
	start, end token.Position
	name       string

	// interface
	// sorted by position
	inputs  []Input
	outputs []Output

	// body
	// sorted by position
	parts []Part
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

// NumInputs returns the number input pins of the chip.
func (c Chip) NumInputs() int { return len(c.inputs) }

// Input returns the i-th input pin.
// This panics if i is out of bounds.
func (c Chip) Input(i int) Input { return c.inputs[i] }

// NumOutputs returns the number output pins of the chip.
func (c Chip) NumOutputs() int { return len(c.outputs) }

// Output returns the i-th output pin.
// This panics if i is out of bounds.
func (c Chip) Output(i int) Output { return c.outputs[i] }

// NumParts returns the number parts in the chip.
func (c Chip) NumParts() int { return len(c.parts) }

// Part returns the i-th part.
// This panics if i is out of bounds.
func (c Chip) Part(i int) Part { return c.parts[i] }

// ChipBuilder can build a description of a chip on-the-fly.
type ChipBuilder struct{}

// AppendComment appends a comment.
// Even if comments relative to the chip are before the chip declaration,
// it can be called after the declaration.
func (b *ChipBuilder) AppendComment(comment Comment) {}

// Declare starts the declaration of the chip and set its name.
func (b *ChipBuilder) Declare(line int, name string) {}

// AppendInput appends an input.
// Even if inputs are declared after the chip declaration and before output declarations,
// it can be called at any moment.
func (b *ChipBuilder) AppendInput(input Input) {}

// AppendOutput appends an output.
// Even if outputs are declared after input declarations and before the body of the chip,
// it can be called at any moment.
func (b *ChipBuilder) AppendOutput(output Output) {}

// AppendPart appends a part.
// Even if the body of the chip is defined after the output declarations,
// it can be called at any moment.
func (b *ChipBuilder) AppendPart(part Part) {}

// Build return the chip.
// The builder is reset so it can reused without any side effect on previously built chips.
func (b *ChipBuilder) Build() Chip { return Chip{} }
