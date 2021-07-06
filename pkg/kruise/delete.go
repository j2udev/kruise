package kruise

import (
	"strings"

	"github.com/j2udevelopment/kruise/pkg/config"
	"github.com/j2udevelopment/kruise/pkg/helm"
	"github.com/j2udevelopment/kruise/pkg/utils"
	"github.com/j2udevelopment/kruise/tpl"
	"github.com/spf13/cobra"
)

var helmDel []config.HelmDeployment
var deleteOpts []config.Option
var validDeleteOpts []string

// NewDeleteOpts sets deployer and valid option slices
func NewDeleteOpts() {
	config.Decode("delete.helm", &helmDel)
	for _, dep := range helmDel {
		deleteOpts = append(deleteOpts, dep.Option)
	}
	validDeleteOpts = utils.CollectValidArgs(deleteOpts)
}

// NewDeleteCmd represents the delete command
func NewDeleteCmd() *cobra.Command {
	//TODO: Set this with a flag
	shallowDryRun := true
	cmd := &cobra.Command{
		Use:       "delete",
		Short:     "Delete the specified options from your Kubernetes cluster",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: validDeleteOpts,
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				for _, dep := range helmDel {
					if utils.Contains(strings.Split(dep.Option.Arguments, ", "), arg) {
						helm.Uninstall(shallowDryRun, &dep.HelmCommand)
					}
				}
			}
		},
	}
	wrapper := config.CommandWrapper{
		Cmd:  cmd,
		Opts: &deleteOpts,
	}
	cmd.SetUsageTemplate(tpl.UsageTemplate())
	cmd.SetUsageFunc(tpl.UsageFunc(wrapper))
	return cmd
}
