package utils

import (
	"strings"

	"github.com/spf13/cobra"
)

// CommandWrapper is used to wrap the Command struct to support command options
type CommandWrapper struct {
	Cmd  *cobra.Command
	Opts []Option
}

// Option is used to map an argument to a description and is primarily used for
// usage templating
type Option struct {
	Arguments   string
	Description string
}

// Check is an abstraction for common error handling
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Contains determines whether a string is contained within a slice of strings
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// CollectValidArgs is used to filter a slice of human-readable options into a
// slice of strings to be used with the Cobra Command ValidArgs slice
func CollectValidArgs(opts []Option) []string {
	var collector = []string{}
	for _, opt := range opts {
		collector = append(collector, strings.Split(opt.Arguments, ", ")...)
	}
	return collector
}

// CollectValidArgsDict is used as a set to facilitate quickly looking up valid
// arguments
func CollectValidArgsDict(args []string) map[string]bool {
	argsDict := make(map[string]bool)
	for _, arg := range args {
		argsDict[arg] = true
	}
	return argsDict
}

// CollectValidArgsMap is used as a set to facilitate quickly looking up valid
// arguments
func CollectValidArgsMap(opts []Option) map[string][]string {
	argsMap := make(map[string][]string)
	for _, opt := range opts {
		optArgs := strings.Split(opt.Arguments, ", ")
		argsMap[optArgs[0]] = optArgs
	}
	return argsMap
}
