package cmd

import (
	"github.com/j2udevelopment/kruise/pkg/kruise"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewKruiseKmd creates the root Kruise command
func NewKruiseKmd() kruise.Kommand {
	return kruise.NewKmd("kruise").
		WithLongDescription(`Kruise is a configurable CLI. It has a set of core commands whose options are determined by a config file.`).
		WithSubKommands(
			NewDeployKmd(),
			NewDeleteKmd(),
		).
		WithPersistentFlags(NewKruisePersistentFlags()).
		WithPersistentPreRunFunc(persistentPreRun).
		Version("0.1.0").
		Build()
}

// NewKruisePersistentFlags creates flags for the kruise command
//
// See the pflag package for more information: https://pkg.go.dev/github.com/spf13/pflag
func NewKruisePersistentFlags() *pflag.FlagSet {
	pfs := pflag.NewFlagSet("kruise", pflag.ContinueOnError)
	pfs.StringP("verbosity", "V", "error", "specify the log level to be used (trace, debug, info, warn, error)")
	return pfs
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	setLogLevel(cmd)
}

func setLogLevel(cmd *cobra.Command) {
	logger := kruise.Logger
	lvl, err := cmd.Flags().GetString("verbosity")
	kruise.Error(err)
	switch lvl {
	case "trace":
		kruise.Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		kruise.Logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.Fatalf("Invalid verbosity level: %s", lvl)
	}
}
