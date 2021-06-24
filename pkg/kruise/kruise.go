package kruise

import (
	"fmt"
	"log"
	"os"

	c "github.com/j2udevelopment/kruise/pkg/config"
	del "github.com/j2udevelopment/kruise/pkg/kruise/delete"
	dep "github.com/j2udevelopment/kruise/pkg/kruise/deploy"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var config c.Config

// NewKruiseCmd represents to root command of the kruise CLI
func NewKruiseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kruise",
		Short: "Kruise streamlines the local development experience",
		Long:  "Kruise is a CLI that aims to streamline the local development experience.\nModern software development can involve an overwhelming number of tools.\nYou can think of Kruise as a CLI wrapper that abstracts (but doesn't hide) the finer details of using many other CLIs that commonly make their way into a software engineers tool kit.",
	}
	cmd.AddCommand(
		dep.NewDeployCmd(dep.NewDeployOpts()),
		del.NewDeleteCmd(del.NewDeleteOpts()),
	)
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.kruise.yaml)")
	return cmd
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kruise" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kruise")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}
}
