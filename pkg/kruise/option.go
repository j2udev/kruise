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

// type Option latest.Option

// func NewOption(o latest.Option) Option {
// 	return Option(o)
// }

// func NewOptions(opts []latest.Option) []Option {
// 	var o []Option
// 	for _, opt := range opts {
// 		o = append(o, NewOption(opt))
// 	}
// 	return o
// }

// func GetValidArgs(opts []Option) []string {
// 	var valid []string
// 	for _, opt := range opts {
// 		valid = append(valid, strings.Split(opt.Arguments, ", ")...)
// 	}
// 	return valid
// }
