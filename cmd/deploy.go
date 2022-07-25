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
		WithExample("kruise deploy kafka mongodb").
		WithArgs(cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs)).
		WithValidArgs(kruise.GetValidDeployArgs()).
		WithShortDescription("Deploy the specified options to your Kubernetes cluster").
		WithOptions(NewDeployOptions()).
		WithRunFunc(kruise.Deploy).
		WithFlags(NewDeployFlags()).
		Build()
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
	fs.BoolP("shallow-dry-run", "d", false, "Output the command being performed under the hood")
	fs.BoolP("concurrent", "c", false, "Deploy the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)")
	fs.BoolP("parallel", "p", false, "Deploy the arguments in parallel (will be ignored if '--concurrent' is passed)")
	fs.BoolP("init", "i", false, "Add Helm repositories for the specified options")
	return fs
}
