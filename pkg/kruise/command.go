package kruise

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type (
	// Command defines how to execute a set or arguments on the command line
	//
	// Used with the exec.Command function: https://pkg.go.dev/os/exec#Command
	Command struct {
		Name   string
		Args   []string
		DryRun bool
		StdOut bool
	}

	// CommandBuilder is used to build a Kruise Command
	CommandBuilder struct {
		Command
	}

	// ICommand defines the functions for a Kruise Command
	ICommand interface {
		Execute() error
	}

	// ICommandBuilder defines the builder functions for the Kruise CommandBuilder
	ICommandBuilder interface {
		WithArgs(a []string) ICommandBuilder
		WithDryRun(dr bool) ICommandBuilder
		WithNoStdOut() ICommandBuilder
		Build() ICommand
	}
)

// NewCmd returns a new Kruise ICommmandBuilder
//
// Not to be confused with the Kruise Kommand which is a wrapper for the cobra
// Command
func NewCmd(name string) ICommandBuilder {
	cmd := Command{Name: name, DryRun: false, StdOut: true}
	return CommandBuilder{cmd}.WithDryRun(false)
}

// WithArgs defines the argument list for a command
func (c CommandBuilder) WithArgs(args []string) ICommandBuilder {
	c.Args = args
	return c
}

// WithDryRun determines whether the command should be printed or executed
func (c CommandBuilder) WithDryRun(dr bool) ICommandBuilder {
	c.DryRun = dr
	return c
}

// WithNoStdOut determines whether the command should show stdout upon being
// executed
func (c CommandBuilder) WithNoStdOut() ICommandBuilder {
	c.StdOut = false
	return c
}

// Build returns an ICommand from a CommandBuilder
func (c CommandBuilder) Build() ICommand {
	return Command{
		Name:   c.Name,
		Args:   c.Args,
		DryRun: c.DryRun,
		StdOut: c.StdOut,
	}
}

// Execute is used to execute the Kruise Command
func (c Command) Execute() error {
	cmd := exec.Command(c.Name, c.Args...)
	if c.DryRun {
		fmt.Printf("%s\n", cmd)
	} else {
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			Fatal(cmd.Wait())
			return err
		}
		cmdErr, _ := io.ReadAll(stderr)
		cmdOut, _ := io.ReadAll(stdout)
		if len(cmdErr) > 0 {
			return errors.New(string(cmdErr))
		}
		if c.StdOut {
			fmt.Printf("%s", cmdOut)
		}
		Fatal(cmd.Wait())
	}
	return nil
}
