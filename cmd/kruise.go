package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/j2udev/boa"
	"github.com/j2udev/kruise/internal/kruise"
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
		WithStringPPersistentFlag("verbosity", "V", kruise.Logger.GetLevel().String(), "specify the log level to be used (debug, info, warn, error)").
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
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	default:
		logger.Fatalf("Invalid verbosity level: %s", lvl)
	}
}
