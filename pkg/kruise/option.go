package kruise

import (
	"strings"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
)

type Option latest.Option

func NewOption(o latest.Option) Option {
	return Option(o)
}

func NewOptions(opts []latest.Option) []Option {
	var o []Option
	for _, opt := range opts {
		o = append(o, NewOption(opt))
	}
	return o
}

func GetValidArgs(opts []Option) []string {
	var valid []string
	for _, opt := range opts {
		valid = append(valid, strings.Split(opt.Arguments, ", ")...)
	}
	return valid
}
