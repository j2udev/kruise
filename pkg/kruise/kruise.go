package kruise

import (
	c "github.com/j2udevelopment/kruise/pkg/config"
	del "github.com/j2udevelopment/kruise/pkg/kruise/delete"
	dep "github.com/j2udevelopment/kruise/pkg/kruise/deploy"
	"github.com/spf13/cobra"
)

var configFile c.ConfigFile

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
	cmd.PersistentFlags().StringVar(&configFile.Override, "config", "", "config file (default is $HOME/.kruise.yaml)")
	return cmd
}

func init() {
	cobra.OnInitialize(c.InitConfig)
}
