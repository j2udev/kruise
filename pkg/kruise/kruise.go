package kruise

import (
	"github.com/j2udevelopment/kruise/pkg/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var configFile config.File
var cfg config.Manifest

// Initialize initializes kruise configuration and commands options
func Initialize() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	configFile.Path = home
	configFile.Extension = "yaml"
	configFile.FileName = ".kruise"
	config.Initialize(configFile, cfg)
	NewDeployOpts()
	NewDeleteOpts()
}

// NewKruiseCmd represents the kruise command
func NewKruiseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kruise",
		Short: "Kruise streamlines the local development experience",
		Long:  "Kruise is a CLI that aims to streamline the local development experience.\nModern software development can involve an overwhelming number of tools.\nYou can think of Kruise as a CLI wrapper that abstracts (but doesn't hide) the finer details of using many other CLIs that commonly make their way into a software engineers tool kit.",
	}
	cmd.AddCommand(
		NewDeployCmd(),
		NewDeleteCmd(),
	)
	cmd.PersistentFlags().StringVarP(&configFile.Override, "config", "c", "", "Specify a custom config file (default is ~/.kruise.yaml)")
	return cmd
}
