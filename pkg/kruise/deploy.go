package kruise

import (
	"github.com/spf13/cobra"
)

// NewDeployCmd represents the deploy command
// options are dynamically populated from `deploy` config in the kruise manifest
func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy the specified options to your Kubernetes cluster",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.ValidArgs = deployer.ValidDeployArgs()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cobra.OnlyValidArgs(cmd, args); err != nil {
				return err
			}
			flags := cmd.Flags()
			parallel, err := flags.GetBool("parallel")
			cobra.CheckErr(err)
			if parallel {
				deployer.DeployP(flags, args)
			} else {
				deployer.Deploy(flags, args)
			}
			return nil
		},
	}
	kmd := &Kommand{
		Cmd:  cmd,
		Opts: &deployer.DeployOptions,
	}
	cmd.SetUsageTemplate(UsageTemplate())
	cmd.SetHelpTemplate(UsageTemplate())
	cmd.SetUsageFunc(UsageFunc(*kmd))
	cmd.SetHelpFunc(HelpFunc(*kmd))
	cmd.PersistentFlags().BoolP("shallow-dry-run", "d", false, "Output the command being performed under the hood")
	cmd.PersistentFlags().BoolP("parallel", "p", false, "Delete the arguments in parallel")
	cmd.Flags().BoolP("help", "h", false, "show help for the deploy command")
	return cmd
}
