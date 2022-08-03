package kruise

import "strings"

type (
	Option struct {
		Arguments   string
		Description string
	}

	Options []Option
)

func NewOption(args []string, desc string) Option {
	return Option{
		Arguments:   strings.Join(args, ", "),
		Description: desc,
	}
}

func (o Options) GetValidArgs() []string {
	var valid []string
	for _, opt := range o {
		valid = append(valid, strings.Split(opt.Arguments, ", ")...)
	}
	return valid
}
