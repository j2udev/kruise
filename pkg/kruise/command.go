package kruise

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
)

type (
	Command struct {
		Name   string
		Args   []string
		DryRun bool
	}

	CommandBuilder struct {
		Command
	}

	ICommand interface {
		Execute() error
	}

	ICommandBuilder interface {
		WithArgs(a []string) ICommandBuilder
		WithDryRun() ICommandBuilder
		Build() ICommand
	}
)

// NewCmd creates a new CommandBuilder for executing CLI commands
func NewCmd(name string) ICommandBuilder {
	cmd := Command{Name: name, DryRun: false}
	return CommandBuilder{cmd}
}

// WithArgs adds arguments to the CLI command
func (c CommandBuilder) WithArgs(args []string) ICommandBuilder {
	c.Args = args
	return c
}

// WithDryRun prints the command instead of executing it when Execute() is called
func (c CommandBuilder) WithDryRun() ICommandBuilder {
	c.DryRun = true
	return c
}

// Build builds the CommandBuilder into a Command
func (c CommandBuilder) Build() ICommand {
	return Command{
		Name:   c.Name,
		Args:   c.Args,
		DryRun: c.DryRun,
	}
}

// Execute either executes the CLI command or prints it
func (c Command) Execute() error {
	cmd := exec.Command(c.Name, c.Args...)
	if c.DryRun {
		fmt.Printf("%s\n", cmd)
	} else {
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			log.Printf("%s", err)
			Fatal(cmd.Wait())
			return err
		}
		cmdErr, _ := io.ReadAll(stderr)
		cmdOut, _ := io.ReadAll(stdout)
		if len(cmdErr) > 0 {
			log.Printf("%s", cmdErr)
			return errors.New(string(cmdErr))
		}
		fmt.Printf("%s", cmdOut)
		Fatal(cmd.Wait())
	}
	return nil
}
