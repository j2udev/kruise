package cmd

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewDeleteKmd creates the Kruise delete command
func NewDeleteKmd() kruise.Kommand {
	return kruise.NewKmd("delete").
		WithAliases([]string{"del"}).
		WithArgs(cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs)).
		WithValidArgs(kruise.GetValidDeployArgs()).
		WithShortDescription("Delete the specified options from your Kubernetes cluster").
		WithOptions(NewDeleteOptions()).
		WithProfiles(NewDeleteProfiles()).
		WithRunFunc(NewDeleteFunc).
		WithFlags(NewDeleteFlags()).
		Build()
}

// NewDeleteFunc is used to define the action performed when the delete command is called
func NewDeleteFunc(cmd *cobra.Command, args []string) {
	kruise.Delete(cmd.Flags(), args)
}

// NewDeleteOptions creates options for the Kruise delete command
//
// Options are dynamically populated from `delete` config in the kruise manifest
func NewDeleteOptions() []kruise.Option {
	return kruise.GetDeleteOptions()
}

// NewDeleteProfiles creates profiles for the Kruise deploy command
//
// Profiles  are dynamically populated from `deploy` config in the kruise manifest
func NewDeleteProfiles() kruise.Profiles {
	return kruise.GetDeleteProfiles()
}

// NewDeleteFlags creates flags for the Kruise delete command
//
// See the pflag package for more information: https://pkg.go.dev/github.com/spf13/pflag
func NewDeleteFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("delete", pflag.ContinueOnError)
	fs.BoolP("shallow-dry-run", "d", false, "output the command being performed under the hood")
	fs.BoolP("concurrent", "c", false, "delete the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)")
	return fs
}
