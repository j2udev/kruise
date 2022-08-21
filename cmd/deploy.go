package cmd

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewDeployKmd creates the Kruise deploy command
func NewDeployKmd() kruise.Kommand {
	return kruise.NewKmd("deploy").
		WithAliases([]string{"dep"}).
		WithArgs(cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs)).
		WithValidArgs(kruise.GetValidDeployArgs()).
		WithShortDescription("Deploy the specified options to your Kubernetes cluster").
		WithOptions(NewDeployOptions()).
		WithRunFunc(NewDeployFunc).
		WithFlags(NewDeployFlags()).
		Build()
}

// NewDeployFunc is used to define the action performed when the deploy command is called
func NewDeployFunc(cmd *cobra.Command, args []string) {
	kruise.Deploy(cmd.Flags(), args)
}

// NewDeployOptions creates options for the Kruise deploy command
//
// Options are dynamically populated from `deploy` config in the kruise manifest
func NewDeployOptions() []kruise.Option {
	return kruise.GetDeployOptions()
}

// NewDeployFlags creates flags for the Kruise deploy command
//
// See the pflag package for more information: https://pkg.go.dev/github.com/spf13/pflag
func NewDeployFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("deploy", pflag.ContinueOnError)
	fs.BoolP("shallow-dry-run", "d", false, "output the command being performed under the hood")
	fs.BoolP("concurrent", "c", false, "deploy the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)")
	fs.BoolP("init", "i", false, "add Helm repositories and create Kubernetes secrets for the specified options")
	fs.StringP("profile", "p", "", "deploy a profile")
	return fs
}
