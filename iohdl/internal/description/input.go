package description

import "github.com/batiazinga/hdl/iohdl/internal/token"

// Input is the description of an input.
type Input struct {
	start, end token.Position
	name       string
}

// Start returns the position at which the input starts.
// Comments are taken into account
// so this is not necessarily the position of the input's name.
func (in Input) Start() token.Position {}

// End returns the line at which the input ends.
// Comments are taken into account
// so this is not necessarily the end of the input's name.
func (in Input) End() token.Position { return in.end }

// Name returns the name of the input.
func (in Input) Name() string { return in.name }

// Size returns the number of bits in the input.
// An input is made of one or more bits.
// A bus has more than one bit.
func (in Input) Size() int { return 1 }
