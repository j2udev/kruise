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
		WithArgs(cobra.MinimumNArgs(1)).
		// TODO: Dynamically populate valid deployment args from config
		// WithValidArgs([]string{"jaeger", "kafka", "mongodb", "mysql"}).
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
	fs.BoolP("parallel", "p", false, "Delete the arguments in parallel")
	return fs
}
