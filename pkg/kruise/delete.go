package kruise

import (
	"github.com/spf13/cobra"
)

// NewDeleteCmd represents the delete command
// options are dynamically populated from `delete` config in the kruise manifest
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "delete",
		Short:     "Delete the specified options from your Kubernetes cluster",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: collectValidArgs(deployer.DeleteOptions),
		Run: func(cmd *cobra.Command, args []string) {
			deployer.Delete(cmd.Flags(), args)
		},
	}
	kmd := &Kommand{
		Cmd:  cmd,
		Opts: &deployer.DeleteOptions,
	}
	cmd.SetUsageTemplate(UsageTemplate())
	cmd.SetHelpTemplate(UsageTemplate())
	cmd.SetUsageFunc(UsageFunc(*kmd))
	cmd.SetHelpFunc(HelpFunc(*kmd))
	cmd.PersistentFlags().BoolP("shallow-dry-run", "d", false, "Output the command being performed under the hood")
	//TODO: Cobra doesn't call initializers before the help flag attempts to
	// render the usage. Try to find a way around this later, but for now rely on
	// the help command instead of the flag for commands that can pass multiple
	// options.
	cmd.Flags().BoolP("help", "h", false, "show help for the deploy command")
	err := cmd.Flags().MarkHidden("help")
	checkErr(err)
	return cmd
}
