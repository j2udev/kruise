package cmd

import (
	"github.com/j2udev/boa"
	"github.com/j2udev/kruise/internal/kruise"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewKruiseCmd() *cobra.Command {
	return boa.NewCmd("kruise").
		WithLongDescription("Kruise is a configurable CLI. It has a set of core commands whose options are determined by a config file.").
		WithSubCommands(
			NewDeployCmd(),
			NewDeleteCmd(),
		).
		WithPersistentPreRunFunc(persistentPreRun).
		WithStringPPersistentFlag("verbosity", "V", "error", "specify the log level to be used (trace, debug, info, warn, error)").
		WithVersion("0.1.0").
		Build()
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
