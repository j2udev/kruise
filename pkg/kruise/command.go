package kruise

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type (
	Command struct {
		Name   string
		Args   []string
		DryRun bool
		StdOut bool
	}

	CommandBuilder struct {
		Command
	}

	ICommand interface {
		Execute() error
	}

	ICommandBuilder interface {
		WithArgs(a []string) ICommandBuilder
		WithDryRun(dr bool) ICommandBuilder
		WithNoStdOut() ICommandBuilder
		Build() ICommand
	}
)

func NewCmd(name string) ICommandBuilder {
	cmd := Command{Name: name, DryRun: false, StdOut: true}
	return CommandBuilder{cmd}.WithDryRun(false)
}

func (c CommandBuilder) WithArgs(args []string) ICommandBuilder {
	c.Args = args
	return c
}

func (c CommandBuilder) WithDryRun(dr bool) ICommandBuilder {
	c.DryRun = dr
	return c
}

func (c CommandBuilder) WithNoStdOut() ICommandBuilder {
	c.StdOut = false
	return c
}

func (c CommandBuilder) Build() ICommand {
	return Command{
		Name:   c.Name,
		Args:   c.Args,
		DryRun: c.DryRun,
		StdOut: c.StdOut,
	}
}

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
