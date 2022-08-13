package kruise

import "strings"

type (
	// Option represents the arguments and description for a CLI option
	Option struct {
		Arguments   string
		Description string
	}
	//Options represents a slice of Option objects
	Options []Option
)

// newOption is a helper function used to create a new Option object from a
// slice of arguments and a description
func newOption(args []string, desc string) Option {
	return Option{
		Arguments:   strings.Join(args, ", "),
		Description: desc,
	}
}

//getValidArgs is used to get arguments from a slice of Options
func (o Options) getValidArgs() []string {
	var valid []string
	for _, opt := range o {
		valid = append(valid, strings.Split(opt.Arguments, ", ")...)
	}
	return valid
}
