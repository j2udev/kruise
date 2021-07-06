package kruise

import (
	"strings"

	"github.com/j2udevelopment/kruise/pkg/config"
	"github.com/j2udevelopment/kruise/pkg/helm"
	"github.com/j2udevelopment/kruise/pkg/utils"
	"github.com/j2udevelopment/kruise/tpl"
	"github.com/spf13/cobra"
)

var helmDep []config.HelmDeployment
var deployOpts []config.Option
var validDeployOpts []string

// NewDeployOpts sets deployer and valid option slices
func NewDeployOpts() {
	config.Decode("deploy.helm", &helmDep)
	for _, dep := range helmDep {
		deployOpts = append(deployOpts, dep.Option)
	}
	validDeployOpts = utils.CollectValidArgs(deployOpts)
}

// NewDeployCmd represents the deploy command
func NewDeployCmd() *cobra.Command {
	//TODO: Set this with a flag
	shallowDryRun := true
	cmd := &cobra.Command{
		Use:       "deploy",
		Short:     "Deploy the specified options to your Kubernetes cluster",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: validDeployOpts,
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				for _, dep := range helmDep {
					if utils.Contains(strings.Split(dep.Option.Arguments, ", "), arg) {
						helm.Install(shallowDryRun, &dep.HelmCommand)
					}
				}
			}
		},
	}
	wrapper := config.CommandWrapper{
		Cmd:  cmd,
		Opts: &deployOpts,
	}
	cmd.SetUsageTemplate(tpl.UsageTemplate())
	cmd.SetUsageFunc(tpl.UsageFunc(wrapper))
	return cmd
}
