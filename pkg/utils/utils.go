package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/j2udevelopment/kruise/pkg/config"
)

// Check is an abstraction for common error debugging
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// ExecuteCommand is used as a repeatable means of calling CLI commands wrapped
// by kruise
func ExecuteCommand(shallowDryRun bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if shallowDryRun {
		fmt.Printf("%s\n", cmd)
	} else {
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			log.Printf("%s", err)
			cmd.Wait()
			return err
		}
		cmdErr, _ := io.ReadAll(stderr)
		cmdOut, _ := io.ReadAll(stdout)
		if len(cmdErr) > 0 {
			log.Printf("%s", cmdErr)
			return errors.New(string(cmdErr))
		}
		fmt.Printf("%s", cmdOut)
		cmd.Wait()
	}
	return nil
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
func CollectValidArgs(opts []config.Option) []string {
	var collector []string
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
func CollectValidArgsMap(opts []config.Option) map[string][]string {
	argsMap := make(map[string][]string)
	for _, opt := range opts {
		optArgs := strings.Split(opt.Arguments, ", ")
		argsMap[optArgs[0]] = optArgs
	}
	return argsMap
}
