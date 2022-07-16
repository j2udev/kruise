package kruise

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewDeployKmd represents the Kruise deploy command
//
// Options are dynamically populated from `deploy` config in the kruise manifest
func NewDeployKmd() Kommand {
	return NewKmd("deploy").
		WithAliases([]string{"dep"}).
		WithExample("kruise deploy kafka mongodb").
		WithArgs(cobra.MinimumNArgs(1)).
		// WithArgs(cobra.OnlyValidArgs).
		// TODO: Dynamically populate valid deployment args from config
		// WithValidArgs([]string{"jaeger", "kafka", "mongodb", "mysql"}).
		WithArgAliases([]string{"mongo"}).
		WithShortDescription("Deploy the specified options to your Kubernetes cluster").
		WithOptions(deployer.DeployOptions).
		WithRunEFunc(deploy).
		WithFlags(NewDeployFlags()).
		Build()
}

func NewDeployFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("deploy", pflag.ContinueOnError)
	fs.BoolP("shallow-dry-run", "d", false, "Output the command being performed under the hood")
	fs.BoolP("parallel", "p", false, "Delete the arguments in parallel")
	return fs
}

func deploy(c *cobra.Command, args []string) error {
	if err := cobra.OnlyValidArgs(c, args); err != nil {
		return err
	}
	flags := c.Flags()
	parallel, err := flags.GetBool("parallel")
	cobra.CheckErr(err)
	if parallel {
		deployer.DeployP(flags, args)
	} else {
		deployer.Deploy(flags, args)
	}
	return nil
}
