package kruise

import (
	"github.com/j2udev/boa"
)

type (
	// Option represents the arguments and description for a CLI option
	Option boa.Option
	//Options represents a slice of Option objects
	Options []Option
)

// newOption is a helper function used to create a new Option object from a
// slice of arguments and a description
func newOption(args []string, desc string) Option {
	return Option{
		Args: args,
		Desc: desc,
	}
}
