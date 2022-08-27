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
		WithProfiles(NewDeployProfiles()).
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
func NewDeployOptions() kruise.Options {
	return kruise.GetDeployOptions()
}

// NewDeployProfiles creates profiles for the Kruise deploy command
//
// Profiles  are dynamically populated from `deploy` config in the kruise manifest
func NewDeployProfiles() kruise.Profiles {
	return kruise.GetDeployProfiles()
}

// NewDeployFlags creates flags for the Kruise deploy command
//
// See the pflag package for more information: https://pkg.go.dev/github.com/spf13/pflag
func NewDeployFlags() *pflag.FlagSet {
	return kruise.GetDeployFlags()
}
