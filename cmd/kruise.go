package cmd

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/spf13/pflag"
)

// NewKruiseKmd represents the kruise command
func NewKruiseKmd() kruise.Kommand {
	return kruise.NewKmd("kruise").
		WithLongDescription(`Kruise is a configurable CLI. It has a set of core commands whose options are determined by a config file.`).
		WithSubKommands(
			NewDeployKmd(),
			NewDeleteKmd(),
		).
		WithPersistentFlags(NewKruisePersistentFlags()).
		Version("0.1.0").
		Build()
}

// NewKruisePersistentFlags creates flags for the kruise command
//
// See the pflag package for more information: https://pkg.go.dev/github.com/spf13/pflag
func NewKruisePersistentFlags() *pflag.FlagSet {
	pfs := pflag.NewFlagSet("kruise", pflag.ContinueOnError)
	pfs.StringVarP(&kruise.Kfg.Metadata.Override, "config", "c", "", "Specify a custom config file (default is ~/.kruise.yaml)")
	return pfs
}
