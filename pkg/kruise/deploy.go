package kruise

import (
	"strings"

	c "github.com/j2udevelopment/kruise/pkg/config"
	h "github.com/j2udevelopment/kruise/pkg/helm"
	u "github.com/j2udevelopment/kruise/pkg/utils"
	t "github.com/j2udevelopment/kruise/tpl"
	"github.com/spf13/cobra"
)

var helmCfg []c.DynamicHelmConfig
var opts []c.Option
var validOpts []string

// NewDeployOpts returns options for the deploy command
func NewDeployOpts() {
	c.Decode("deploy.helm", &helmCfg)
	for _, dep := range helmCfg {
		opts = append(opts, dep.Option)
	}
	validOpts = u.CollectValidArgs(opts)
}

// NewDeployCmd represents the deploy command
func NewDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "deploy",
		Short:     "Deploy the specified options to your Kubernetes cluster",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: validOpts,
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				for _, helm := range helmCfg {
					if u.Contains(strings.Split(helm.Option.Arguments, ", "), arg) {
						h.Install(true, &helm.HelmConfig)
					}
				}
			}
		},
	}
	wrapper := c.CommandWrapper{
		Cmd:  cmd,
		Opts: &opts,
	}
	cmd.SetUsageTemplate(t.UsageTemplate())
	cmd.SetUsageFunc(t.UsageFunc(wrapper))
	return cmd
}
