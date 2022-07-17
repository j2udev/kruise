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
		WithArgs(cobra.MinimumNArgs(1)).
		// TODO: Dynamically populate valid deployment args from config
		// WithValidArgs([]string{"jaeger", "kafka", "mongodb", "mysql"}).
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
	fs.BoolP("parallel", "p", false, "Delete the arguments in parallel")
	return fs
}