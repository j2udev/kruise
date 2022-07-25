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
		WithExample("kruise delete kafka mongodb").
		WithArgs(cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs)).
		WithValidArgs(kruise.GetValidDeleteArgs()).
		WithShortDescription("Delete the specified options from your Kubernetes cluster").
		WithOptions(NewDeleteOptions()).
		WithRunFunc(kruise.Delete).
		WithFlags(NewDeleteFlags()).
		Build()
}

// NewDeleteOptions creates options for the Kruise delete command
//
// Options are dynamically populated from `delete` config in the kruise manifest
func NewDeleteOptions() []kruise.Option {
	return kruise.GetDeleteOptions()
}

// NewDeleteFlags creates flags for the Kruise delete command
//
// See the pflag package for more information: https://pkg.go.dev/github.com/spf13/pflag
func NewDeleteFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("delete", pflag.ContinueOnError)
	fs.BoolP("shallow-dry-run", "d", false, "Output the command being performed under the hood")
	fs.BoolP("concurrent", "c", false, "Delete the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)")
	fs.BoolP("parallel", "p", false, "Delete the arguments in parallel (will be ignored if '--concurrent' is passed)")
	return fs
}
