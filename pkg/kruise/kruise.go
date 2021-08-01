package kruise

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

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
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
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
		Short: "Kruise is a black-box CLI",
		Long: `Kruise is a configurable CLI. It has a set of core commands whose
options are determined by a config file.`,
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
			cobra.CheckErr(cmd.Wait())
			return err
		}
		cmdErr, _ := io.ReadAll(stderr)
		cmdOut, _ := io.ReadAll(stdout)
		if len(cmdErr) > 0 {
			log.Printf("%s", cmdErr)
			return errors.New(string(cmdErr))
		}
		fmt.Printf("%s", cmdOut)
		cobra.CheckErr(cmd.Wait())
	}
	return nil
}
