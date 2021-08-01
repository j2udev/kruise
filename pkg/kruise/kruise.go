package kruise

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var kfg Konfig
var deployer Deployer

// Kommand is used to wrap cobra commands to support command options
type Kommand struct {
	Cmd  *cobra.Command
	Opts *[]Option
}

// Initialize is used to initialize kruise configuration and command options
func Initialize() {
	home, err := homedir.Dir()
	checkErr(err)
	file := &kfg.Metadata
	file.Path = home
	file.Name = ".kruise"
	file.Extension = "yaml"
	kfg.Initialize()
	// must initialize konfig before constructing a deployer
	deployer = NewDeployer()
}

// NewCmd represents the kruise command
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kruise",
		Short: "Kruise streamlines the local development experience",
		Long:  "Kruise is a CLI that aims to streamline the local development experience.\nModern software development can involve an overwhelming number of tools.\nYou can think of Kruise as a CLI wrapper that abstracts (but doesn't hide) the finer details of using many other CLIs that commonly make their way into a software engineers tool kit.",
	}
	cmd.AddCommand(
		NewDeployCmd(),
		NewDeleteCmd(),
	)
	cmd.PersistentFlags().StringVarP(&kfg.Metadata.Override, "config", "c", "", "Specify a custom config file (default is ~/.kruise.yaml)")
	return cmd
}

// ExecuteCommand is used as a repeatable means of calling CLI commands
func ExecuteCommand(shallowDryRun bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if shallowDryRun {
		fmt.Printf("%s\n", cmd)
	} else {
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			log.Printf("%s", err)

			waitErr := cmd.Wait()
			checkErr(waitErr)
			return err
		}
		cmdErr, _ := io.ReadAll(stderr)
		cmdOut, _ := io.ReadAll(stdout)
		if len(cmdErr) > 0 {
			log.Printf("%s", cmdErr)
			return errors.New(string(cmdErr))
		}
		fmt.Printf("%s", cmdOut)
		waitErr := cmd.Wait()
		checkErr(waitErr)
	}
	return nil
}
